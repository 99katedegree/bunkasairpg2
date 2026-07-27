package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
	domainrepo "github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
	dbgen "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db/sqlc"
)

type weaponRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewWeaponRepository(db *sql.DB) domainrepo.WeaponRepository {
	return &weaponRepository{db: db, q: dbgen.New(db)}
}

func weaponFromDB(w dbgen.Weapon) *entity.Weapon {
	e := &entity.Weapon{
		ID:            w.ID,
		Name:          w.Name,
		IndexNumber:   w.IndexNumber,
		PhysicsAttack: float64(w.PhysicsAttack),
		PhysicsType:   w.PhysicsType,
		ElementType:   w.ElementType,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
	if w.ElementAttack.Valid {
		ea := float64(w.ElementAttack.Int32)
		e.ElementAttack = &ea
	}
	return e
}

func (r *weaponRepository) FindAllSummaries(ctx context.Context) ([]entity.WeaponSummary, error) {
	rows, err := r.q.GetAllWeaponSummaries(ctx)
	if err != nil {
		return nil, err
	}
	// sqlc は 0 件のとき nil を返す。JSON で null にならないよう空スライスで初期化する。
	out := make([]entity.WeaponSummary, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.WeaponSummary{ID: row.ID, Name: row.Name, IndexNumber: row.IndexNumber})
	}
	return out, nil
}

func (r *weaponRepository) FindByID(ctx context.Context, id int64) (*entity.Weapon, error) {
	row, err := r.q.GetWeaponByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return weaponFromDB(row), nil
}

func (r *weaponRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Weapon, error) {
	rows, err := r.q.GetWeaponsByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	result := make([]*entity.Weapon, 0, len(rows))
	for _, w := range rows {
		result = append(result, weaponFromDB(w))
	}
	return result, nil
}

func (r *weaponRepository) FindIndexByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Weapon, int64, error) {
	total, err := r.q.CountWeaponIndexByUserID(ctx, userID.String())
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.GetWeaponIndexByUserID(ctx, dbgen.GetWeaponIndexByUserIDParams{
		UserID: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.Weapon, 0, len(rows))
	for _, w := range rows {
		result = append(result, weaponFromDB(w))
	}
	return result, total, nil
}

func (r *weaponRepository) IsOwnedByUser(ctx context.Context, userID uuid.UUID, weaponID int64) (bool, error) {
	count, err := r.q.IsWeaponOwnedByUser(ctx, dbgen.IsWeaponOwnedByUserParams{
		UserID:   userID.String(),
		WeaponID: weaponID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *weaponRepository) GrantToUser(ctx context.Context, userID uuid.UUID, weaponID int64) error {
	if err := r.q.GrantWeaponToUser(ctx, dbgen.GrantWeaponToUserParams{
		UserID:   userID.String(),
		WeaponID: weaponID,
	}); err != nil {
		return err
	}
	return r.q.GrantWeaponToEntry(ctx, dbgen.GrantWeaponToEntryParams{
		UserID:   userID.String(),
		WeaponID: weaponID,
	})
}
