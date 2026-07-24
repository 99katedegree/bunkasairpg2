package repository

import (
	"context"
	"database/sql"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	domainrepo "github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
	dbgen "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db/sqlc"
)

type imageRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewImageRepository(db *sql.DB) domainrepo.ImageRepository {
	return &imageRepository{db: db, q: dbgen.New(db)}
}

func (r *imageRepository) Create(ctx context.Context, img *entity.Image) (*entity.Image, error) {
	if err := r.q.CreateImage(ctx, dbgen.CreateImageParams{
		Directory: img.Directory,
		Url:       img.URL,
	}); err != nil {
		return nil, err
	}

	row, err := r.q.GetLastInsertImage(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.Image{
		ID:        row.ID,
		Directory: row.Directory,
		URL:       row.Url,
		CreatedAt: row.CreatedAt,
	}, nil
}
