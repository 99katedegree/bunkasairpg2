package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, u *entity.UpdateUser) error
	Activate(ctx context.Context, id uuid.UUID) error
	ArchiveAll(ctx context.Context) error
	DeleteInactive(ctx context.Context) error
	BulkCreate(ctx context.Context, count int) ([]uuid.UUID, error)
	GetAllIDs(ctx context.Context) ([]uuid.UUID, error)
}
