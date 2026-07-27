package usecase

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"math"
	mrand "math/rand/v2"
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

	seed, err = newSeed()
	if err != nil {
		return uuid.Nil, 0, err
	}

	// 開始時の装備武器を固定する。戦闘中の持ち替えで users.equipped_weapon_id が
	// 書き換わるため、終了時にそこを見ると別の武器から再計算してしまう。
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return uuid.Nil, 0, err
	}

	token = uuid.New()
	b := &entity.Battle{
		ID:        token,
		UserID:    userID,
		MonsterID: &monsterID,
		Seed:      seed,
		Status:    entity.BattleStatusInProgress,
	}
	if user.Weapon != nil {
		weaponID := user.Weapon.ID
		b.StartWeaponID = &weaponID
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

	// 使用アイテム・変更武器を事前ロードし、所持を検証する
	items, weapons, usedItems, err := u.loadActionResources(ctx, userID, actions)
	if err != nil {
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusExpired)
		return nil, err
	}

	// 開始時の装備武器。記録がなければ素手。
	startWeapon := entityWeaponToParams(&entity.BareHands)
	if b.StartWeaponID != nil {
		w, err := u.weaponRepo.FindByID(ctx, *b.StartWeaponID)
		if err != nil {
			return nil, err
		}
		startWeapon = entityWeaponToParams(w)
	}

	// シミュレーション実行
	input := buildSimulatorInput(b.Seed, actions, monster, user, startWeapon, items, weapons, false)
	_, err = battle.Simulate(input)
	if err != nil {
		status := entity.BattleStatusLost
		if err == battle.ErrBattleInvalid {
			status = entity.BattleStatusExpired
		}
		// 敗北ならアイテムは使い切っている。手順が不正だった場合だけ消費しない。
		if err != battle.ErrBattleInvalid {
			consumeUsedItems(ctx, userID, usedItems, u.itemRepo)
		}
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, status)
		return nil, err
	}

	consumeUsedItems(ctx, userID, usedItems, u.itemRepo)

	// 報酬付与
	_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusCompleted)
	newExp := user.ExperiencePoint + monster.ExperiencePoint
	newLevel := calcLevel(newExp)
	newHP := user.HitPoint
	if newLevel > user.Level {
		newHP += levelUpHitPointGain(newLevel - user.Level)
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

func (u *BattleUsecase) loadActionResources(ctx context.Context, userID uuid.UUID, actions []battle.Action) (map[int64]battle.ItemParams, map[int64]battle.WeaponParams, map[int64]int, error) {
	items, weapons, used, err := loadOwnedActionResources(ctx, userID, actions, u.itemRepo, u.weaponRepo)
	return items, weapons, used, err
}

// loadOwnedActionResources はアクション列で使われる武器とアイテムを読み込む。
//
// 併せて所持を検証する。ここを通さないと、クライアントは持っていない最強武器への
// 持ち替えや、持っていないデバフアイテムの使用を宣言するだけで再計算を通せてしまう。
// 戻り値の第3引数はアイテムIDごとの使用回数で、呼び出し側が消費に使う。
func loadOwnedActionResources(
	ctx context.Context,
	userID uuid.UUID,
	actions []battle.Action,
	itemRepo repository.ItemRepository,
	weaponRepo repository.WeaponRepository,
) (map[int64]battle.ItemParams, map[int64]battle.WeaponParams, map[int64]int, error) {
	items := map[int64]battle.ItemParams{}
	weapons := map[int64]battle.WeaponParams{}
	usedCount := map[int64]int{}

	// 所持アイテムの在庫を一度だけ引く
	owned, err := itemRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	stock := make(map[int64]int, len(owned))
	for _, it := range owned {
		stock[it.ID] = it.Count
	}

	for _, a := range actions {
		switch {
		case a.Type == battle.ActionUseItem && a.ItemID != nil:
			usedCount[*a.ItemID]++
			// 在庫を超えて使ったと主張していないか
			if usedCount[*a.ItemID] > stock[*a.ItemID] {
				return nil, nil, nil, entity.ErrItemStockEmpty
			}
			if _, ok := items[*a.ItemID]; !ok {
				it, err := itemRepo.FindByID(ctx, *a.ItemID)
				if err != nil {
					return nil, nil, nil, err
				}
				items[*a.ItemID] = entityItemToParams(it)
			}

		case a.Type == battle.ActionChangeWeapon && a.WeaponID != nil:
			if _, ok := weapons[*a.WeaponID]; !ok {
				isOwned, err := weaponRepo.IsOwnedByUser(ctx, userID, *a.WeaponID)
				if err != nil {
					return nil, nil, nil, err
				}
				if !isOwned {
					return nil, nil, nil, entity.ErrWeaponNotOwned
				}
				w, err := weaponRepo.FindByID(ctx, *a.WeaponID)
				if err != nil {
					return nil, nil, nil, err
				}
				weapons[*a.WeaponID] = entityWeaponToParams(w)
			}
		}
	}
	return items, weapons, usedCount, nil
}

// consumeUsedItems は検証済みのアクション列で実際に使われた分を減らす。
// 消費はここに一本化してあり、戦闘中には減らさない。
// 途中で減らすと、終了時の再計算から見て在庫が足りない状態になってしまう。
func consumeUsedItems(ctx context.Context, userID uuid.UUID, used map[int64]int, itemRepo repository.ItemRepository) {
	for itemID, n := range used {
		for range n {
			_ = itemRepo.DecrementUserItem(ctx, userID, itemID)
		}
	}
}

func buildSimulatorInput(seed int64, actions []battle.Action, monster *entity.Monster, user *entity.User, startWeapon battle.WeaponParams, items map[int64]battle.ItemParams, weapons map[int64]battle.WeaponParams, isBoss bool) battle.SimulatorInput {
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

// newSeed はバトル用のシードを作る。
//
// JavaScript の数値は倍精度浮動小数なので、2^53 を超える整数は JSON を経由した時点で
// 丸められる。乱数生成（mulberry32）は下位 32bit しか使わないが、丸めで下位ビットまで
// 変わってしまうため、はじめから 32bit に収めてサーバーとクライアントを一致させる。
func newSeed() (int64, error) {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint32(buf[:])), nil
}

// levelUpHitPointGain は levels 段のレベルアップで増える最大HPを返す。
//
// 1 段ごとに 6〜10 の一様乱数を引いて合計する。リファクタ前はこれをフロントの
// reward() が計算していたが、報酬はサーバーが決めるようになったのでここへ移した。
// 結果は finish のレスポンスで返すだけなので、クライアントと同期する必要はない。
func levelUpHitPointGain(levels int) int {
	gain := 0
	for range levels {
		gain += 6 + mrand.IntN(5)
	}
	return gain
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
