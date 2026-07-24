package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
)

type AdminRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.Admin, error)
}
