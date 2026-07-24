package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type MonsterUsecase struct {
	monsterRepo repository.MonsterRepository
}

func NewMonsterUsecase(monsterRepo repository.MonsterRepository) *MonsterUsecase {
	return &MonsterUsecase{monsterRepo: monsterRepo}
}

func (u *MonsterUsecase) GetDetail(ctx context.Context, monsterID uuid.UUID) (*entity.Monster, error) {
	return u.monsterRepo.FindByID(ctx, monsterID)
}

func (u *MonsterUsecase) GetCatalog(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.MonsterCatalogEntry, int64, error) {
	return u.monsterRepo.FindCatalogByUserID(ctx, userID, offset, limit)
}
