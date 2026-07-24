package usecase

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battle"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type BossBattleUsecase struct {
	battleRepo repository.BattleRepository
	weaponRepo repository.WeaponRepository
	itemRepo   repository.ItemRepository
	userRepo   repository.UserRepository
}

func NewBossBattleUsecase(
	battleRepo repository.BattleRepository,
	weaponRepo repository.WeaponRepository,
	itemRepo repository.ItemRepository,
	userRepo repository.UserRepository,
) *BossBattleUsecase {
	return &BossBattleUsecase{
		battleRepo: battleRepo,
		weaponRepo: weaponRepo,
		itemRepo:   itemRepo,
		userRepo:   userRepo,
	}
}

// Start はボスバトルを開始しトークンとシードを返す（monster_id = NULL）
func (u *BossBattleUsecase) Start(ctx context.Context, userID uuid.UUID) (token uuid.UUID, seed int64, err error) {
	var seedBuf [8]byte
	if _, err = rand.Read(seedBuf[:]); err != nil {
		return uuid.Nil, 0, err
	}
	seed = int64(binary.LittleEndian.Uint64(seedBuf[:]))
	token = uuid.New()

	b := &entity.Battle{
		ID:        token,
		UserID:    userID,
		MonsterID: nil, // ボスバトルは NULL
		Seed:      seed,
		Status:    entity.BattleStatusInProgress,
	}
	if err = u.battleRepo.Create(ctx, b); err != nil {
		return uuid.Nil, 0, err
	}
	return token, seed, nil
}

// Finish はボスバトル結果を検証し、クリアタイムを記録して報酬を付与する
func (u *BossBattleUsecase) Finish(ctx context.Context, userID uuid.UUID, token uuid.UUID, actions []battle.Action) (*BossFinishResult, error) {
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
	if time.Since(b.CreatedAt) > 60*time.Minute {
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusExpired)
		return nil, entity.ErrBattleExpired
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items, weapons, err := loadBossActionResources(ctx, actions, u.itemRepo, u.weaponRepo)
	if err != nil {
		return nil, err
	}

	// ボスパラメータをハードコードから使用
	bossMonster := battle.BossParams
	input := battle.SimulatorInput{
		Seed:    b.Seed,
		Actions: actions,
		Monster: bossMonster,
		IsBoss:  true,
		User:    battle.UserState{Level: user.Level, HitPoint: user.HitPoint, MaxHitPoint: user.HitPoint, ExperiencePoint: user.ExperiencePoint},
		Items:   items,
		Weapons: weapons,
	}
	if user.Weapon != nil {
		input.User.Weapon = entityWeaponToParams(user.Weapon)
	}

	_, err = battle.Simulate(input)
	finishTime := time.Now()
	if err != nil {
		status := entity.BattleStatusLost
		if err == battle.ErrBattleInvalid {
			status = entity.BattleStatusExpired
		}
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, status)
		return nil, err
	}

	clearTimeMs := int(finishTime.Sub(b.CreatedAt).Milliseconds())
	_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusCompleted)

	// クリア記録保存
	_ = u.battleRepo.UpsertBossRecord(ctx, &entity.BossRecord{
		UserID:    userID,
		ClearTime: clearTimeMs,
	})

	// 報酬付与（ボスは経験値のみ）
	newExp := user.ExperiencePoint + battle.BossParams.ExperiencePoint
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

	return &BossFinishResult{
		ClearTime:       clearTimeMs,
		ExperiencePoint: newExp,
		Level:           newLevel,
		HitPoint:        newHP,
	}, nil
}

type BossFinishResult struct {
	ClearTime       int
	ExperiencePoint int
	Level           int
	HitPoint        int
}

func loadBossActionResources(ctx context.Context, actions []battle.Action, itemRepo repository.ItemRepository, weaponRepo repository.WeaponRepository) (map[int64]battle.ItemParams, map[int64]battle.WeaponParams, error) {
	items := map[int64]battle.ItemParams{}
	weapons := map[int64]battle.WeaponParams{}
	for _, a := range actions {
		if a.Type == battle.ActionUseItem && a.ItemID != nil {
			if _, ok := items[*a.ItemID]; !ok {
				it, err := itemRepo.FindByID(ctx, *a.ItemID)
				if err != nil {
					return nil, nil, err
				}
				items[*a.ItemID] = entityItemToParams(it)
			}
		}
		if a.Type == battle.ActionChangeWeapon && a.WeaponID != nil {
			if _, ok := weapons[*a.WeaponID]; !ok {
				w, err := weaponRepo.FindByID(ctx, *a.WeaponID)
				if err != nil {
					return nil, nil, err
				}
				weapons[*a.WeaponID] = entityWeaponToParams(w)
			}
		}
	}
	return items, weapons, nil
}
