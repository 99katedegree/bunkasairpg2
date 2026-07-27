package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type ItemUsecase struct {
	itemRepo repository.ItemRepository
}

func NewItemUsecase(itemRepo repository.ItemRepository) *ItemUsecase {
	return &ItemUsecase{itemRepo: itemRepo}
}

func (u *ItemUsecase) GetAllIDs(ctx context.Context) ([]int64, error) {
	return u.itemRepo.FindAllIDs(ctx)
}

func (u *ItemUsecase) GetUserItems(ctx context.Context, userID uuid.UUID) ([]*entity.UserItem, error) {
	return u.itemRepo.FindByUserID(ctx, userID)
}

func (u *ItemUsecase) GetIndex(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Item, int64, error) {
	return u.itemRepo.FindIndexByUserID(ctx, userID, offset, limit)
}

func (u *ItemUsecase) UseItem(ctx context.Context, userID uuid.UUID, itemID int64) error {
	// アイテムが所持しているか確認
	items, err := u.itemRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	var found *entity.UserItem
	for _, it := range items {
		if it.ID == itemID {
			found = it
			break
		}
	}
	if found == nil {
		return entity.ErrItemNotFound
	}
	if found.Count <= 0 {
		return entity.ErrItemStockEmpty
	}
	return u.itemRepo.DecrementUserItem(ctx, userID, itemID)
}
