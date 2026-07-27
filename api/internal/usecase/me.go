package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type MeUsecase struct {
	userRepo   repository.UserRepository
	weaponRepo repository.WeaponRepository
}

func NewMeUsecase(userRepo repository.UserRepository, weaponRepo repository.WeaponRepository) *MeUsecase {
	return &MeUsecase{userRepo: userRepo, weaponRepo: weaponRepo}
}

func (u *MeUsecase) Get(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}

// Update はプロフィール項目だけを更新する。
// レベル・HP・経験値はバトル終了時にサーバーが決めるため、ここでは受け付けない。
func (u *MeUsecase) Update(ctx context.Context, req *entity.UpdateUser) error {
	// 装備は所持している武器に限る。WeaponUsecase.ChangeWeapon と同じ制約。
	if req.EquippedWeaponID != nil {
		owned, err := u.weaponRepo.IsOwnedByUser(ctx, req.ID, *req.EquippedWeaponID)
		if err != nil {
			return err
		}
		if !owned {
			return entity.ErrWeaponNotOwned
		}
	}
	return u.userRepo.Update(ctx, req)
}
