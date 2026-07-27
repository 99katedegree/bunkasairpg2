package entity

import (
	"time"

	"github.com/google/uuid"
)

type BattleStatus string

const (
	BattleStatusInProgress BattleStatus = "in_progress"
	BattleStatusCompleted  BattleStatus = "completed"
	BattleStatusLost       BattleStatus = "lost"
	BattleStatusExpired    BattleStatus = "expired"
)

type Battle struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	MonsterID *uuid.UUID // nil = ボスバトル
	// StartWeaponID は開始時に装備していた武器。戦闘中の持ち替えで
	// users.equipped_weapon_id が書き換わっても、再計算はここから始める。
	// nil = 素手（entity.BareHands）
	StartWeaponID *int64
	Seed          int64
	Status        BattleStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BossRecord struct {
	UserID    uuid.UUID
	ClearTime int // ミリ秒
	CreatedAt time.Time
	UpdatedAt time.Time
}
