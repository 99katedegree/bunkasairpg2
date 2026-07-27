package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type MonsterRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Monster, error)
	// FindAllIDs はバトルトークン生成用。QR に埋める UUID が要るのでこちらを使う。
	FindAllIDs(ctx context.Context) ([]uuid.UUID, error)
	// FindAllSummaries は管理画面の一覧用。id / name / 図鑑番号だけを返す。
	FindAllSummaries(ctx context.Context) ([]entity.MonsterSummary, error)
	FindCatalogByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.MonsterCatalogEntry, int64, error)
	RegisterEntry(ctx context.Context, userID, monsterID uuid.UUID) error
	IsEntryRegistered(ctx context.Context, userID, monsterID uuid.UUID) (bool, error)
}
