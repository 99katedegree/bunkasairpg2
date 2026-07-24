package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	EquippedWeaponID *int64
	AvatarImageID    *int64
	Name             string
	Level            int
	HitPoint         int
	ExperiencePoint  int
	RememberToken    *string
	Weapon           *Weapon // JOIN 結果（optional）
	AvatarImageURL   *string // JOIN 結果（optional）
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UpdateUser struct {
	ID               uuid.UUID
	Name             *string
	Level            *int
	HitPoint         *int
	ExperiencePoint  *int
	EquippedWeaponID *int64
	AvatarImageID    *int64
}
