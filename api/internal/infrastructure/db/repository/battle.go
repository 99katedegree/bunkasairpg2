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

type battleRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewBattleRepository(db *sql.DB) domainrepo.BattleRepository {
	return &battleRepository{db: db, q: dbgen.New(db)}
}

func (r *battleRepository) Create(ctx context.Context, b *entity.Battle) error {
	params := dbgen.CreateBattleParams{
		ID:     b.ID.String(),
		UserID: b.UserID.String(),
		Seed:   b.Seed,
		Status: string(b.Status),
	}
	if b.MonsterID != nil {
		params.MonsterID = sql.NullString{String: b.MonsterID.String(), Valid: true}
	}
	if b.StartWeaponID != nil {
		params.StartWeaponID = sql.NullInt64{Int64: *b.StartWeaponID, Valid: true}
	}
	return r.q.CreateBattle(ctx, params)
}

func (r *battleRepository) FindByToken(ctx context.Context, token uuid.UUID) (*entity.Battle, error) {
	row, err := r.q.GetBattleByID(ctx, token.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(row.UserID)
	if err != nil {
		return nil, err
	}

	b := &entity.Battle{
		ID:        id,
		UserID:    userID,
		Seed:      row.Seed,
		Status:    entity.BattleStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if row.MonsterID.Valid {
		monsterID, err := uuid.Parse(row.MonsterID.String)
		if err != nil {
			return nil, err
		}
		b.MonsterID = &monsterID
	}
	if row.StartWeaponID.Valid {
		weaponID := row.StartWeaponID.Int64
		b.StartWeaponID = &weaponID
	}

	return b, nil
}

func (r *battleRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BattleStatus) error {
	return r.q.UpdateBattleStatus(ctx, dbgen.UpdateBattleStatusParams{
		ID:     id.String(),
		Status: string(status),
	})
}

func (r *battleRepository) UpsertBossRecord(ctx context.Context, rec *entity.BossRecord) error {
	return r.q.UpsertBossRecord(ctx, dbgen.UpsertBossRecordParams{
		UserID:    rec.UserID.String(),
		ClearTime: int32(rec.ClearTime),
	})
}
