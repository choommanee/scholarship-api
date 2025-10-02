package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EmailTemplate represents an email template
type EmailTemplate struct {
	TemplateID   uuid.UUID       `json:"template_id" db:"template_id"`
	TemplateName string          `json:"template_name" db:"template_name"`
	Subject      string          `json:"subject" db:"subject"`
	Body         string          `json:"body" db:"body"`
	Variables    json.RawMessage `json:"variables,omitempty" db:"variables"`
	TemplateType string          `json:"template_type" db:"template_type"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// EmailQueue represents an email in the queue
type EmailQueue struct {
	QueueID        uuid.UUID  `json:"queue_id" db:"queue_id"`
	RecipientEmail string     `json:"recipient_email" db:"recipient_email"`
	RecipientName  *string    `json:"recipient_name,omitempty" db:"recipient_name"`
	SenderEmail    string     `json:"sender_email" db:"sender_email"`
	Subject        string     `json:"subject" db:"subject"`
	Body           string     `json:"body" db:"body"`
	TemplateID     *uuid.UUID `json:"template_id,omitempty" db:"template_id"`
	Priority       int        `json:"priority" db:"priority"`
	Status         string     `json:"status" db:"status"`
	SentAt         *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	ErrorMessage   *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// EmailVariables represents variables for email template
type EmailVariables map[string]interface{}

// Value implements the driver.Valuer interface
func (e EmailTemplate) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Scan implements the sql.Scanner interface
func (e *EmailTemplate) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &e)
}
