package entity

import "time"

type Image struct {
	ID        int64
	Directory string
	URL       string
	CreatedAt time.Time
}
