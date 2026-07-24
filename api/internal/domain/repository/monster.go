package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type MonsterRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Monster, error)
	FindCatalogByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.MonsterCatalogEntry, int64, error)
	RegisterEntry(ctx context.Context, userID, monsterID uuid.UUID) error
	IsEntryRegistered(ctx context.Context, userID, monsterID uuid.UUID) (bool, error)
}
