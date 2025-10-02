package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FileStorage represents a file storage entry
type FileStorage struct {
	FileID       uuid.UUID  `json:"file_id" db:"file_id"`
	OriginalName string     `json:"original_name" db:"original_name"`
	StoredName   string     `json:"stored_name" db:"stored_name"`
	StoredPath   string     `json:"stored_path" db:"stored_path"`
	FileSize     int64      `json:"file_size" db:"file_size"`
	MimeType     string     `json:"mime_type" db:"mime_type"`
	FileHash     string     `json:"file_hash" db:"file_hash"`
	UploadedBy   uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
	RelatedTable *string    `json:"related_table,omitempty" db:"related_table"`
	RelatedID    *uuid.UUID `json:"related_id,omitempty" db:"related_id"`
	StorageType  string     `json:"storage_type" db:"storage_type"`
	AccessLevel  string     `json:"access_level" db:"access_level"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// FileVersion represents a document version in file storage
type FileVersion struct {
	VersionID         uuid.UUID  `json:"version_id" db:"version_id"`
	FileID            uuid.UUID  `json:"file_id" db:"file_id"`
	VersionNumber     int        `json:"version_number" db:"version_number"`
	ChangeDescription *string    `json:"change_description,omitempty" db:"change_description"`
	UploadedBy        uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
	FileSize          int64      `json:"file_size" db:"file_size"`
	IsCurrent         bool       `json:"is_current" db:"is_current"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ReplacedAt        *time.Time `json:"replaced_at,omitempty" db:"replaced_at"`
}

// FileAccessLog represents a file access log entry
type FileAccessLog struct {
	AccessID   uuid.UUID `json:"access_id" db:"access_id"`
	FileID     uuid.UUID `json:"file_id" db:"file_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	AccessTime time.Time `json:"access_time" db:"access_time"`
	Action     string    `json:"action" db:"action"`
	IPAddress  *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent  *string   `json:"user_agent,omitempty" db:"user_agent"`
	Success    bool      `json:"success" db:"success"`
}

// Value implements the driver.Valuer interface
func (f FileStorage) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements the sql.Scanner interface
func (f *FileStorage) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &f)
}
