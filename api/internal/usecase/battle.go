package usecase

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battle"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type BattleUsecase struct {
	battleRepo  repository.BattleRepository
	monsterRepo repository.MonsterRepository
	weaponRepo  repository.WeaponRepository
	itemRepo    repository.ItemRepository
	userRepo    repository.UserRepository
}

func NewBattleUsecase(
	battleRepo repository.BattleRepository,
	monsterRepo repository.MonsterRepository,
	weaponRepo repository.WeaponRepository,
	itemRepo repository.ItemRepository,
	userRepo repository.UserRepository,
) *BattleUsecase {
	return &BattleUsecase{
		battleRepo:  battleRepo,
		monsterRepo: monsterRepo,
		weaponRepo:  weaponRepo,
		itemRepo:    itemRepo,
		userRepo:    userRepo,
	}
}

// Start はバトルを開始しトークンとシードを返す。未登録モンスターは図鑑に自動登録。
func (u *BattleUsecase) Start(ctx context.Context, userID, monsterID uuid.UUID) (token uuid.UUID, seed int64, err error) {
	// モンスター存在確認
	_, err = u.monsterRepo.FindByID(ctx, monsterID)
	if err != nil {
		return uuid.Nil, 0, entity.ErrNotFound
	}

	// 図鑑未登録なら自動登録
	registered, err := u.monsterRepo.IsEntryRegistered(ctx, userID, monsterID)
	if err != nil {
		return uuid.Nil, 0, err
	}
	if !registered {
		if err = u.monsterRepo.RegisterEntry(ctx, userID, monsterID); err != nil {
			return uuid.Nil, 0, err
		}
	}

	// シード生成
	var seedBuf [8]byte
	if _, err = rand.Read(seedBuf[:]); err != nil {
		return uuid.Nil, 0, err
	}
	seed = int64(binary.LittleEndian.Uint64(seedBuf[:]))

	token = uuid.New()
	b := &entity.Battle{
		ID:        token,
		UserID:    userID,
		MonsterID: &monsterID,
		Seed:      seed,
		Status:    entity.BattleStatusInProgress,
	}
	if err = u.battleRepo.Create(ctx, b); err != nil {
		return uuid.Nil, 0, err
	}
	return token, seed, nil
}

// Finish はバトル結果を検証し報酬を付与する
func (u *BattleUsecase) Finish(ctx context.Context, userID uuid.UUID, token uuid.UUID, actions []battle.Action) (*FinishBattleResult, error) {
	b, err := u.battleRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, entity.ErrNotFound
	}
	if b.UserID != userID {
		return nil, entity.ErrUnauthorized
	}
	if b.Status != entity.BattleStatusInProgress {
		return nil, entity.ErrBattleExpired
	}
	// 30分タイムアウト
	if time.Since(b.CreatedAt) > 30*time.Minute {
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusExpired)
		return nil, entity.ErrBattleExpired
	}

	// モンスター情報取得
	monster, err := u.monsterRepo.FindByID(ctx, *b.MonsterID)
	if err != nil {
		return nil, err
	}

	// ユーザー情報取得
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 使用アイテム・変更武器を事前ロード
	items, weapons, err := u.loadActionResources(ctx, actions)
	if err != nil {
		return nil, err
	}

	// シミュレーション実行
	input := buildSimulatorInput(b.Seed, actions, monster, user, items, weapons, false)
	_, err = battle.Simulate(input)
	if err != nil {
		status := entity.BattleStatusLost
		if err == battle.ErrBattleInvalid {
			status = entity.BattleStatusExpired
		}
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, status)
		return nil, err
	}

	// 報酬付与
	_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusCompleted)
	newExp := user.ExperiencePoint + monster.ExperiencePoint
	newLevel := calcLevel(newExp)
	newHP := user.HitPoint
	if newLevel > user.Level {
		newHP += (newLevel - user.Level) * 8
	}
	_ = u.userRepo.Update(ctx, &entity.UpdateUser{
		ID:              userID,
		ExperiencePoint: &newExp,
		Level:           &newLevel,
		HitPoint:        &newHP,
	})

	res := &FinishBattleResult{
		ExperiencePoint: newExp,
		Level:           newLevel,
		HitPoint:        newHP,
	}

	// ドロップ付与
	if monster.WeaponID != nil {
		owned, _ := u.weaponRepo.IsOwnedByUser(ctx, userID, *monster.WeaponID)
		if !owned {
			_ = u.weaponRepo.GrantToUser(ctx, userID, *monster.WeaponID)
			wid := int(*monster.WeaponID)
			res.DropWeaponID = &wid
		}
	} else if monster.ItemID != nil {
		_ = u.itemRepo.GrantToUser(ctx, userID, *monster.ItemID)
		iid := int(*monster.ItemID)
		res.DropItemID = &iid
	}

	return res, nil
}

type FinishBattleResult struct {
	ExperiencePoint int
	Level           int
	HitPoint        int
	DropWeaponID    *int
	DropItemID      *int
}

func (u *BattleUsecase) loadActionResources(ctx context.Context, actions []battle.Action) (map[int64]battle.ItemParams, map[int64]battle.WeaponParams, error) {
	items := map[int64]battle.ItemParams{}
	weapons := map[int64]battle.WeaponParams{}
	for _, a := range actions {
		if a.Type == battle.ActionUseItem && a.ItemID != nil {
			if _, ok := items[*a.ItemID]; !ok {
				it, err := u.itemRepo.FindByID(ctx, *a.ItemID)
				if err != nil {
					return nil, nil, err
				}
				items[*a.ItemID] = entityItemToParams(it)
			}
		}
		if a.Type == battle.ActionChangeWeapon && a.WeaponID != nil {
			if _, ok := weapons[*a.WeaponID]; !ok {
				w, err := u.weaponRepo.FindByID(ctx, *a.WeaponID)
				if err != nil {
					return nil, nil, err
				}
				weapons[*a.WeaponID] = entityWeaponToParams(w)
			}
		}
	}
	return items, weapons, nil
}

func buildSimulatorInput(seed int64, actions []battle.Action, monster *entity.Monster, user *entity.User, items map[int64]battle.ItemParams, weapons map[int64]battle.WeaponParams, isBoss bool) battle.SimulatorInput {
	var startWeapon battle.WeaponParams
	if user.Weapon != nil {
		startWeapon = entityWeaponToParams(user.Weapon)
	}
	return battle.SimulatorInput{
		Seed:    seed,
		Actions: actions,
		Monster: battle.MonsterParams{
			Attack:          monster.Attack,
			HitPoint:        monster.HitPoint,
			MaxHitPoint:     monster.HitPoint,
			ExperiencePoint: monster.ExperiencePoint,
			Slash:           monster.Slash,
			Blow:            monster.Blow,
			Shoot:           monster.Shoot,
			Neutral:         monster.Neutral,
			Flame:           monster.Flame,
			Water:           monster.Water,
			Wood:            monster.Wood,
			Shine:           monster.Shine,
			Dark:            monster.Dark,
		},
		IsBoss:  isBoss,
		User:    battle.UserState{Level: user.Level, HitPoint: user.HitPoint, MaxHitPoint: user.HitPoint, ExperiencePoint: user.ExperiencePoint, Weapon: startWeapon},
		Items:   items,
		Weapons: weapons,
	}
}

func entityWeaponToParams(w *entity.Weapon) battle.WeaponParams {
	return battle.WeaponParams{
		ID:            w.ID,
		PhysicsAttack: w.PhysicsAttack,
		ElementAttack: w.ElementAttack,
		PhysicsType:   w.PhysicsType,
		ElementType:   w.ElementType,
	}
}

func entityItemToParams(it *entity.Item) battle.ItemParams {
	return battle.ItemParams{
		ID:         it.ID,
		EffectType: it.EffectType,
		Amount:     it.Amount,
		Rate:       it.Rate,
		Target:     it.Target,
	}
}

// calcLevel は経験値からレベルを計算する。
// フロントエンドの front/src/utils/calculate-level.ts の calculateLevel と同一の計算式。
func calcLevel(exp int) int {
	if exp <= 0 {
		return 1
	}
	const baseExp = 19.0
	const rateOfIncrease = 0.067
	numerator := math.Log(1 + float64(exp)*rateOfIncrease/baseExp)
	denominator := math.Log(1 + rateOfIncrease)
	calculated := 1 + numerator/denominator
	level := int(math.Floor(calculated))
	if level < 1 {
		return 1
	}
	return level
}
