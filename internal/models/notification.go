package models

import "time"

type Notification struct {
	NotificationID   int        `json:"notification_id" db:"notification_id"`
	UserID           string     `json:"user_id" db:"user_id"`
	NotificationType string     `json:"notification_type" db:"notification_type"`
	Title            string     `json:"title" db:"title"`
	Message          string     `json:"message" db:"message"`
	ReferenceID      *string    `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType    *string    `json:"reference_type,omitempty" db:"reference_type"`
	IsRead           bool       `json:"is_read" db:"is_read"`
	IsEmailSent      bool       `json:"is_email_sent" db:"is_email_sent"`
	EmailSentAt      *time.Time `json:"email_sent_at,omitempty" db:"email_sent_at"`
	Priority         string     `json:"priority" db:"priority"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	ReadAt           *time.Time `json:"read_at,omitempty" db:"read_at"`
}