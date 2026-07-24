package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battletoken"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type MonsterUsecase struct {
	monsterRepo  repository.MonsterRepository
	battleToken  *battletoken.BattleToken
}

func NewMonsterUsecase(monsterRepo repository.MonsterRepository, bt *battletoken.BattleToken) *MonsterUsecase {
	return &MonsterUsecase{monsterRepo: monsterRepo, battleToken: bt}
}

func (u *MonsterUsecase) GetDetail(ctx context.Context, monsterID uuid.UUID) (*entity.Monster, error) {
	return u.monsterRepo.FindByID(ctx, monsterID)
}

func (u *MonsterUsecase) GetCatalog(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.MonsterCatalogEntry, int64, error) {
	return u.monsterRepo.FindCatalogByUserID(ctx, userID, offset, limit)
}

func (u *MonsterUsecase) GetBattleTokens(ctx context.Context) ([]string, error) {
	ids, err := u.monsterRepo.FindAllIDs(ctx)
	if err != nil {
		return nil, err
	}
	tokens := make([]string, 0, len(ids))
	for _, id := range ids {
		token, err := u.battleToken.Encrypt(id)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
