package entity

import (
	"time"

	"github.com/google/uuid"
)

type Monster struct {
	ID              uuid.UUID
	WeaponID        *int64
	ItemID          *int64
	IndexNumber     string
	Name            string
	Attack          int
	HitPoint        int
	ExperiencePoint int
	Slash           float64
	Blow            float64
	Shoot           float64
	Neutral         float64
	Flame           float64
	Water           float64
	Wood            float64
	Shine           float64
	Dark            float64
	Weapon          *Weapon // JOIN 結果（optional）
	Item            *Item   // JOIN 結果（optional）
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// MonsterCatalogEntry はユーザーの図鑑エントリ（閲覧済みのみ）
type MonsterCatalogEntry struct {
	MonsterID *uuid.UUID // nil = 未登録スロット
}
