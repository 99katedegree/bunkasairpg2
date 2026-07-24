package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type WeaponUsecase struct {
	weaponRepo repository.WeaponRepository
	userRepo   repository.UserRepository
}

func NewWeaponUsecase(weaponRepo repository.WeaponRepository, userRepo repository.UserRepository) *WeaponUsecase {
	return &WeaponUsecase{weaponRepo: weaponRepo, userRepo: userRepo}
}

func (u *WeaponUsecase) GetUserWeapons(ctx context.Context, userID uuid.UUID) ([]*entity.Weapon, error) {
	return u.weaponRepo.FindByUserID(ctx, userID)
}

func (u *WeaponUsecase) GetIndex(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Weapon, int64, error) {
	return u.weaponRepo.FindIndexByUserID(ctx, userID, offset, limit)
}

func (u *WeaponUsecase) ChangeWeapon(ctx context.Context, userID uuid.UUID, weaponID int64) error {
	owned, err := u.weaponRepo.IsOwnedByUser(ctx, userID, weaponID)
	if err != nil {
		return err
	}
	if !owned {
		return entity.ErrWeaponNotOwned
	}
	wid := weaponID
	return u.userRepo.Update(ctx, &entity.UpdateUser{
		ID:               userID,
		EquippedWeaponID: &wid,
	})
}
