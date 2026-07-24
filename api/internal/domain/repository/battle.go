package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type BattleRepository interface {
	Create(ctx context.Context, b *entity.Battle) error
	FindByToken(ctx context.Context, token uuid.UUID) (*entity.Battle, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BattleStatus) error
	UpsertBossRecord(ctx context.Context, rec *entity.BossRecord) error
}
