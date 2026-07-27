package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type ItemRepository interface {
	FindAllIDs(ctx context.Context) ([]int64, error)
	FindByID(ctx context.Context, id int64) (*entity.Item, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserItem, error)
	FindIndexByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Item, int64, error)
	DecrementUserItem(ctx context.Context, userID uuid.UUID, itemID int64) error
	GrantToUser(ctx context.Context, userID uuid.UUID, itemID int64) error
}
