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

var _ domainrepo.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *sql.DB
	q  *dbgen.Queries
}

func NewUserRepository(db *sql.DB) domainrepo.UserRepository {
	return &userRepository{db: db, q: dbgen.New(db)}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row, err := r.q.GetUserByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	userID, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		ID:              userID,
		Name:            row.Name,
		Level:           int(row.Level),
		HitPoint:        int(row.HitPoint),
		ExperiencePoint: int(row.ExperiencePoint),
		IsArchived:      row.IsArchived,
		IsActivated:     row.IsActivated,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	if row.EquippedWeaponID.Valid {
		v := row.EquippedWeaponID.Int64
		u.EquippedWeaponID = &v
	}
	if row.AvatarImageID.Valid {
		v := row.AvatarImageID.Int64
		u.AvatarImageID = &v
	}
	if row.RememberToken.Valid {
		u.RememberToken = &row.RememberToken.String
	}
	if row.AvatarImageUrl.Valid {
		u.AvatarImageURL = &row.AvatarImageUrl.String
	}

	if row.WeaponIDJoin.Valid {
		physicsAttack := float64(row.PhysicsAttack.Int32)
		w := &entity.Weapon{
			ID:            row.WeaponIDJoin.Int64,
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
		u.Weapon = w
	}

	return u, nil
}

func (r *userRepository) Activate(ctx context.Context, id uuid.UUID) error {
	return r.q.ActivateUser(ctx, id.String())
}

func (r *userRepository) ArchiveAll(ctx context.Context) error {
	return r.q.ArchiveAllUsers(ctx)
}

func (r *userRepository) DeleteInactive(ctx context.Context) error {
	return r.q.DeleteInactiveUsers(ctx)
}

func (r *userRepository) BulkCreate(ctx context.Context, count int) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, count)
	for i := range count {
		id := uuid.New()
		if err := r.q.CreateUser(ctx, id.String()); err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (r *userRepository) GetAllIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.q.GetAllUserIDs(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, s := range rows {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *userRepository) Update(ctx context.Context, u *entity.UpdateUser) error {
	// Build params — for COALESCE fields send current-or-zero; we pass pointer values.
	// The SQL uses COALESCE(?, name) so passing NULL keeps the old value.
	// We reconstruct as sql.NullXxx.
	params := dbgen.UpdateUserParams{
		ID: u.ID.String(),
		EquippedWeaponID: sql.NullInt64{
			Int64: 0,
			Valid: false,
		},
		AvatarImageID: sql.NullInt64{
			Int64: 0,
			Valid: false,
		},
	}

	// For COALESCE columns, sqlc generated non-null types.
	// We need the current values — but here we only have what was changed.
	// The COALESCE(?=NULL, col) approach: send nil for unchanged fields.
	// However, sqlc generated non-nullable types for Name/Level/HitPoint/ExperiencePoint.
	// We work around by using a raw exec with nullable params.
	if u.EquippedWeaponID != nil {
		params.EquippedWeaponID = sql.NullInt64{Int64: *u.EquippedWeaponID, Valid: true}
	}
	if u.AvatarImageID != nil {
		params.AvatarImageID = sql.NullInt64{Int64: *u.AvatarImageID, Valid: true}
	}

	// sqlc generated UpdateUserParams with non-null Name/Level etc., but the SQL uses COALESCE.
	// We need to pass nil to keep old values for non-set fields.
	// Use raw query to pass nullable args.
	const q = `UPDATE users SET
    name = COALESCE(?, name),
    level = COALESCE(?, level),
    hit_point = COALESCE(?, hit_point),
    experience_point = COALESCE(?, experience_point),
    equipped_weapon_id = ?,
    avatar_image_id = ?,
    updated_at = NOW()
WHERE id = ?`

	var name interface{}
	if u.Name != nil {
		name = *u.Name
	}
	var level interface{}
	if u.Level != nil {
		level = int32(*u.Level)
	}
	var hitPoint interface{}
	if u.HitPoint != nil {
		hitPoint = int32(*u.HitPoint)
	}
	var expPoint interface{}
	if u.ExperiencePoint != nil {
		expPoint = int32(*u.ExperiencePoint)
	}

	var equippedWeaponID interface{}
	if u.EquippedWeaponID != nil {
		equippedWeaponID = *u.EquippedWeaponID
	}
	var avatarImageID interface{}
	if u.AvatarImageID != nil {
		avatarImageID = *u.AvatarImageID
	}

	_, err := r.db.ExecContext(ctx, q,
		name,
		level,
		hitPoint,
		expPoint,
		equippedWeaponID,
		avatarImageID,
		u.ID.String(),
	)
	return err
}
