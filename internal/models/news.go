package models

import (
	"time"

	"github.com/google/uuid"
)

// News represents a news article in the system
type News struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Summary     string     `json:"summary" db:"summary"`
	ImageURL    *string    `json:"image_url" db:"image_url"`
	PublishDate time.Time  `json:"publish_date" db:"publish_date"`
	ExpireDate  *time.Time `json:"expire_date" db:"expire_date"`
	Category    string     `json:"category" db:"category"`
	Tags        []string   `json:"tags" db:"tags"`
	IsPublished bool       `json:"is_published" db:"is_published"`
	CreatedBy   string     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
