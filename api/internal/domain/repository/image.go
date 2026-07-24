package repository

import (
	"context"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
)

type ImageRepository interface {
	Create(ctx context.Context, img *entity.Image) (*entity.Image, error)
}
