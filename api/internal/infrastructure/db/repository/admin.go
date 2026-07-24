package repository

import (
	"context"
	"database/sql"

	dbgen "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db/sqlc"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
)

type adminRepository struct {
	q *dbgen.Queries
}

func NewAdminRepository(db *sql.DB) *adminRepository {
	return &adminRepository{q: dbgen.New(db)}
}

func (r *adminRepository) FindByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	row, err := r.q.FindAdminByEmail(ctx, email)
	if err != nil {
		return nil, entity.ErrNotFound
	}
	return &entity.Admin{
		ID:        row.ID,
		Email:     row.Email,
		Password:  row.Password,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
