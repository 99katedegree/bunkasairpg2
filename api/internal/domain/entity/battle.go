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
	Seed      int64
	Status    BattleStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BossRecord struct {
	UserID    uuid.UUID
	ClearTime int // ミリ秒
	CreatedAt time.Time
	UpdatedAt time.Time
}
