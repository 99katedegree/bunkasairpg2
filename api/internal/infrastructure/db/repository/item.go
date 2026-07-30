package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	domainrepo "github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
	dbgen "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db/sqlc"
)

type itemRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewItemRepository(db *sql.DB) domainrepo.ItemRepository {
	return &itemRepository{db: db, q: dbgen.New(db)}
}

// parseDecimalString converts a nullable DECIMAL string from sqlc to *float64.
func parseDecimalString(s sql.NullString) *float64 {
	if !s.Valid {
		return nil
	}
	v, err := strconv.ParseFloat(s.String, 64)
	if err != nil {
		return nil
	}
	return &v
}

func itemFromRow(id int64, name, indexNumber, effectType string, imageURL sql.NullString, createdAt, updatedAt interface{},
	amount sql.NullInt32, buffRate, buffTarget, debuffRate, debuffTarget sql.NullString,
) *entity.Item {
	item := &entity.Item{
		ID:          id,
		Name:        name,
		IndexNumber: indexNumber,
		EffectType:  effectType,
		ImageURL:    nullStringPtr(imageURL),
	}

	if amount.Valid {
		v := int(amount.Int32)
		item.Amount = &v
	}
	if buffRate.Valid {
		item.Rate = parseDecimalString(buffRate)
		if buffTarget.Valid {
			item.Target = &buffTarget.String
		}
	} else if debuffRate.Valid {
		item.Rate = parseDecimalString(debuffRate)
		if debuffTarget.Valid {
			item.Target = &debuffTarget.String
		}
	}

	return item
}

func (r *itemRepository) FindAllSummaries(ctx context.Context) ([]entity.ItemSummary, error) {
	rows, err := r.q.GetAllItemSummaries(ctx)
	if err != nil {
		return nil, err
	}
	// sqlc は 0 件のとき nil を返す。JSON で null にならないよう空スライスで初期化する。
	out := make([]entity.ItemSummary, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.ItemSummary{
			ID: row.ID, Name: row.Name, IndexNumber: row.IndexNumber,
			ImageURL: nullStringPtr(row.ImageUrl),
		})
	}
	return out, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id int64) (*entity.Item, error) {
	row, err := r.q.GetItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	item := itemFromRow(
		row.ID, row.Name, row.IndexNumber, row.EffectType, row.ImageUrl, row.CreatedAt, row.UpdatedAt,
		row.Amount, row.BuffRate, row.BuffTarget, row.DebuffRate, row.DebuffTarget,
	)
	item.CreatedAt = row.CreatedAt
	item.UpdatedAt = row.UpdatedAt
	return item, nil
}

func (r *itemRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserItem, error) {
	rows, err := r.q.GetItemsByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entity.UserItem, 0, len(rows))
	for _, row := range rows {
		item := itemFromRow(
			row.ID, row.Name, row.IndexNumber, row.EffectType, row.ImageUrl, row.CreatedAt, row.UpdatedAt,
			row.Amount, row.BuffRate, row.BuffTarget, row.DebuffRate, row.DebuffTarget,
		)
		item.CreatedAt = row.CreatedAt
		item.UpdatedAt = row.UpdatedAt
		result = append(result, &entity.UserItem{
			Item:  *item,
			Count: int(row.Count),
		})
	}
	return result, nil
}

func (r *itemRepository) FindIndexByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Item, int64, error) {
	total, err := r.q.CountItemIndexByUserID(ctx, userID.String())
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.GetItemIndexByUserID(ctx, dbgen.GetItemIndexByUserIDParams{
		UserID: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Item, 0, len(rows))
	for _, row := range rows {
		item := itemFromRow(
			row.ID, row.Name, row.IndexNumber, row.EffectType, row.ImageUrl, row.CreatedAt, row.UpdatedAt,
			row.Amount, row.BuffRate, row.BuffTarget, row.DebuffRate, row.DebuffTarget,
		)
		item.CreatedAt = row.CreatedAt
		item.UpdatedAt = row.UpdatedAt
		result = append(result, item)
	}
	return result, total, nil
}

func (r *itemRepository) DecrementUserItem(ctx context.Context, userID uuid.UUID, itemID int64) error {
	if err := r.q.DecrementUserItem(ctx, dbgen.DecrementUserItemParams{
		UserID: userID.String(),
		ItemID: itemID,
	}); err != nil {
		return err
	}
	return r.q.DeleteUserItemIfZero(ctx, dbgen.DeleteUserItemIfZeroParams{
		UserID: userID.String(),
		ItemID: itemID,
	})
}

func (r *itemRepository) GrantToUser(ctx context.Context, userID uuid.UUID, itemID int64) error {
	if err := r.q.GrantItemToUser(ctx, dbgen.GrantItemToUserParams{
		UserID: userID.String(),
		ItemID: itemID,
	}); err != nil {
		return err
	}
	return r.q.GrantItemToEntry(ctx, dbgen.GrantItemToEntryParams{
		UserID: userID.String(),
		ItemID: itemID,
	})
}
