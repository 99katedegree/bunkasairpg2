package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
	"github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/storage"
)

type ImageUsecase struct {
	imageRepo repository.ImageRepository
	r2        *storage.R2Client
}

func NewImageUsecase(imageRepo repository.ImageRepository, r2 *storage.R2Client) *ImageUsecase {
	return &ImageUsecase{imageRepo: imageRepo, r2: r2}
}

func (u *ImageUsecase) Upload(ctx context.Context, directory, filename string, body io.Reader, contentType string) (*entity.Image, error) {
	key := fmt.Sprintf("%s/%s", directory, filename)
	url, err := u.r2.Upload(ctx, key, body, contentType)
	if err != nil {
		return nil, err
	}
	img := &entity.Image{Directory: directory, URL: url}
	return u.imageRepo.Create(ctx, img)
}
