package models

import (
	"time"

	"github.com/google/uuid"
)

// NewsRead represents a record of a user reading a news article
type NewsRead struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	NewsID    uuid.UUID `json:"news_id" db:"news_id"`
	ReadAt    time.Time `json:"read_at" db:"read_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UnreadNewsCount represents the count of unread news for a user
type UnreadNewsCount struct {
	Count int `json:"count"`
}
