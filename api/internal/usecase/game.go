package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
)

type GameUsecase struct {
	userRepo repository.UserRepository
}

func NewGameUsecase(userRepo repository.UserRepository) *GameUsecase {
	return &GameUsecase{userRepo: userRepo}
}

// Archive は全ユーザーをアーカイブし、未ログインユーザーを削除する。
func (u *GameUsecase) Archive(ctx context.Context) error {
	if err := u.userRepo.DeleteInactive(ctx); err != nil {
		return err
	}
	return u.userRepo.ArchiveAll(ctx)
}

// Start はアーカイブ処理後に count 人の新規ユーザーを作成し、そのIDを返す。
func (u *GameUsecase) Start(ctx context.Context, count int) ([]uuid.UUID, error) {
	if err := u.Archive(ctx); err != nil {
		return nil, err
	}
	return u.userRepo.BulkCreate(ctx, count)
}
