package entity

import "time"

type Item struct {
	ID          int64
	Name        string
	IndexNumber string
	EffectType  string   // heal | buff | debuff
	Amount      *int     // heal のみ
	Rate        *float64 // buff/debuff のみ
	Target      *string  // buff/debuff のみ (slash|blow|shoot|neutral|flame|water|wood|shine|dark)
	ImageURL    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserItem struct {
	Item
	Count int
}
