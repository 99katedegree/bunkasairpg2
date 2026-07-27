package usecase

import (
	"context"
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
	seed, err = newSeed()
	if err != nil {
		return uuid.Nil, 0, err
	}

	token = uuid.New()

	b := &entity.Battle{
		ID:        token,
		UserID:    userID,
		MonsterID: nil, // ボスバトルは NULL
		Seed:      seed,
		Status:    entity.BattleStatusInProgress,
	}
	// 開始時の装備武器を固定する（通常バトルと同じ理由）。
	if user, err := u.userRepo.FindByID(ctx, userID); err == nil && user.Weapon != nil {
		weaponID := user.Weapon.ID
		b.StartWeaponID = &weaponID
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

	// 通常バトルと同じ所持検証。ボスは戦闘中にアイテムを消費していないので、
	// ここでまとめて在庫を確認し、あとで減らす。
	items, weapons, usedItems, err := loadOwnedActionResources(ctx, userID, actions, u.itemRepo, u.weaponRepo)
	if err != nil {
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, entity.BattleStatusExpired)
		return nil, err
	}

	// ボスパラメータをハードコードから使用
	bossMonster := battle.BossParams
	input := battle.SimulatorInput{
		Seed:    b.Seed,
		Actions: actions,
		Monster: bossMonster,
		IsBoss:  true,
		User: battle.UserState{
			Level: user.Level, HitPoint: user.HitPoint, MaxHitPoint: user.HitPoint,
			ExperiencePoint: user.ExperiencePoint,
			// 開始時の装備武器。記録がなければ素手。
			Weapon: entityWeaponToParams(&entity.BareHands),
		},
		Items:   items,
		Weapons: weapons,
	}
	if b.StartWeaponID != nil {
		w, err := u.weaponRepo.FindByID(ctx, *b.StartWeaponID)
		if err != nil {
			return nil, err
		}
		input.User.Weapon = entityWeaponToParams(w)
	}

	_, err = battle.Simulate(input)
	finishTime := time.Now()
	if err != nil {
		status := entity.BattleStatusLost
		if err == battle.ErrBattleInvalid {
			status = entity.BattleStatusExpired
		}
		if err != battle.ErrBattleInvalid {
			consumeUsedItems(ctx, userID, usedItems, u.itemRepo)
		}
		_ = u.battleRepo.UpdateStatus(ctx, b.ID, status)
		return nil, err
	}

	consumeUsedItems(ctx, userID, usedItems, u.itemRepo)

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
		newHP += levelUpHitPointGain(newLevel - user.Level)
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
