package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type WeaponRepository interface {
	// FindAllSummaries は管理画面の一覧用。id / name / 図鑑番号だけを返す。
	FindAllSummaries(ctx context.Context) ([]entity.WeaponSummary, error)
	FindByID(ctx context.Context, id int64) (*entity.Weapon, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Weapon, error)
	FindIndexByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Weapon, int64, error)
	IsOwnedByUser(ctx context.Context, userID uuid.UUID, weaponID int64) (bool, error)
	GrantToUser(ctx context.Context, userID uuid.UUID, weaponID int64) error
}
