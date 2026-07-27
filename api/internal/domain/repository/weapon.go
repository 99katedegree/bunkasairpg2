package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type WeaponRepository interface {
	FindAllIDs(ctx context.Context) ([]int64, error)
	FindByID(ctx context.Context, id int64) (*entity.Weapon, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Weapon, error)
	FindIndexByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Weapon, int64, error)
	IsOwnedByUser(ctx context.Context, userID uuid.UUID, weaponID int64) (bool, error)
	GrantToUser(ctx context.Context, userID uuid.UUID, weaponID int64) error
}
