package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JobQueue represents a job in the queue
type JobQueue struct {
	JobID        uuid.UUID       `json:"job_id" db:"job_id"`
	JobType      string          `json:"job_type" db:"job_type"`
	Payload      json.RawMessage `json:"payload" db:"payload"`
	Priority     int             `json:"priority" db:"priority"`
	Status       string          `json:"status" db:"status"`
	Attempts     int             `json:"attempts" db:"attempts"`
	MaxAttempts  int             `json:"max_attempts" db:"max_attempts"`
	ScheduledAt  time.Time       `json:"scheduled_at" db:"scheduled_at"`
	StartedAt    *time.Time      `json:"started_at,omitempty" db:"started_at"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	ErrorMessage *string         `json:"error_message,omitempty" db:"error_message"`
}

// BackgroundTask represents a background task
type BackgroundTask struct {
	TaskID    uuid.UUID  `json:"task_id" db:"task_id"`
	TaskName  string     `json:"task_name" db:"task_name"`
	TaskType  string     `json:"task_type" db:"task_type"`
	Schedule  *string    `json:"schedule,omitempty" db:"schedule"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	LastRun   *time.Time `json:"last_run,omitempty" db:"last_run"`
	NextRun   *time.Time `json:"next_run,omitempty" db:"next_run"`
	RunCount  int        `json:"run_count" db:"run_count"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// JobPayload represents generic job payload
type JobPayload map[string]interface{}

// Value implements the driver.Valuer interface
func (j JobQueue) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JobQueue) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &j)
}
