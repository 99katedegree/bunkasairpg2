package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type MeUsecase struct {
	userRepo repository.UserRepository
}

func NewMeUsecase(userRepo repository.UserRepository) *MeUsecase {
	return &MeUsecase{userRepo: userRepo}
}

func (u *MeUsecase) Get(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}

func (u *MeUsecase) Update(ctx context.Context, req *entity.UpdateUser) error {
	return u.userRepo.Update(ctx, req)
}
