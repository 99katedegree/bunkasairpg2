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

type monsterRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewMonsterRepository(db *sql.DB) domainrepo.MonsterRepository {
	return &monsterRepository{db: db, q: dbgen.New(db)}
}

func parseDecimal(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func (r *monsterRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Monster, error) {
	row, err := r.q.GetMonsterByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	monsterID, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, err
	}

	m := &entity.Monster{
		ID:              monsterID,
		IndexNumber:     row.IndexNumber,
		Name:            row.Name,
		Attack:          int(row.Attack),
		HitPoint:        int(row.HitPoint),
		ExperiencePoint: int(row.ExperiencePoint),
		Slash:           parseDecimal(row.Slash),
		Blow:            parseDecimal(row.Blow),
		Shoot:           parseDecimal(row.Shoot),
		Neutral:         parseDecimal(row.Neutral),
		Flame:           parseDecimal(row.Flame),
		Water:           parseDecimal(row.Water),
		Wood:            parseDecimal(row.Wood),
		Shine:           parseDecimal(row.Shine),
		Dark:            parseDecimal(row.Dark),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	if row.WeaponID.Valid {
		v := row.WeaponID.Int64
		m.WeaponID = &v
	}
	if row.ItemID.Valid {
		v := row.ItemID.Int64
		m.ItemID = &v
	}

	// Build joined Weapon if present
	if row.WeaponIDJ.Valid {
		physicsAttack := float64(row.PhysicsAttack.Int32)
		w := &entity.Weapon{
			ID:            row.WeaponIDJ.Int64,
			Name:          row.WeaponName.String,
			IndexNumber:   row.WeaponIndexNumber.String,
			PhysicsAttack: physicsAttack,
			PhysicsType:   row.PhysicsType.String,
			ElementType:   row.ElementType.String,
		}
		if row.ElementAttack.Valid {
			ea := float64(row.ElementAttack.Int32)
			w.ElementAttack = &ea
		}
		if row.WeaponCreatedAt.Valid {
			w.CreatedAt = row.WeaponCreatedAt.Time
		}
		if row.WeaponUpdatedAt.Valid {
			w.UpdatedAt = row.WeaponUpdatedAt.Time
		}
		m.Weapon = w
	}

	// Build joined Item if present
	if row.ItemIDJ.Valid {
		item := &entity.Item{
			ID:          row.ItemIDJ.Int64,
			Name:        row.ItemName.String,
			IndexNumber: row.ItemIndexNumber.String,
			EffectType:  row.EffectType.String,
		}
		if row.ItemCreatedAt.Valid {
			item.CreatedAt = row.ItemCreatedAt.Time
		}
		if row.ItemUpdatedAt.Valid {
			item.UpdatedAt = row.ItemUpdatedAt.Time
		}
		if row.Amount.Valid {
			v := int(row.Amount.Int32)
			item.Amount = &v
		}
		if row.BuffRate.Valid {
			item.Rate = parseDecimalString(row.BuffRate)
			if row.BuffTarget.Valid {
				item.Target = &row.BuffTarget.String
			}
		} else if row.DebuffRate.Valid {
			item.Rate = parseDecimalString(row.DebuffRate)
			if row.DebuffTarget.Valid {
				item.Target = &row.DebuffTarget.String
			}
		}
		m.Item = item
	}

	return m, nil
}

func (r *monsterRepository) FindAllIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.q.GetAllMonsterIDs(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, idStr := range rows {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *monsterRepository) FindCatalogByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.MonsterCatalogEntry, int64, error) {
	total, err := r.q.CountMonsterCatalogByUserID(ctx, userID.String())
	if err != nil {
		return nil, 0, err
	}

	monsterIDs, err := r.q.GetMonsterCatalogByUserID(ctx, dbgen.GetMonsterCatalogByUserIDParams{
		UserID: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.MonsterCatalogEntry, 0, len(monsterIDs))
	for _, idStr := range monsterIDs {
		parsed, err := uuid.Parse(idStr)
		if err != nil {
			return nil, 0, err
		}
		p := parsed
		result = append(result, &entity.MonsterCatalogEntry{MonsterID: &p})
	}
	return result, total, nil
}

func (r *monsterRepository) RegisterEntry(ctx context.Context, userID, monsterID uuid.UUID) error {
	return r.q.RegisterMonsterEntry(ctx, dbgen.RegisterMonsterEntryParams{
		UserID:    userID.String(),
		MonsterID: monsterID.String(),
	})
}

func (r *monsterRepository) IsEntryRegistered(ctx context.Context, userID, monsterID uuid.UUID) (bool, error) {
	count, err := r.q.IsMonsterEntryRegistered(ctx, dbgen.IsMonsterEntryRegisteredParams{
		UserID:    userID.String(),
		MonsterID: monsterID.String(),
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
