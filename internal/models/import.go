package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ImportDetail represents import detail entry
type ImportDetail struct {
	DetailID      uuid.UUID       `json:"detail_id" db:"detail_id"`
	ImportID      int             `json:"import_id" db:"import_id"`
	RowNumber     int             `json:"row_number" db:"row_number"`
	RawData       json.RawMessage `json:"raw_data" db:"raw_data"`
	ProcessedData json.RawMessage `json:"processed_data,omitempty" db:"processed_data"`
	Status        string          `json:"status" db:"status"`
	ErrorMessage  *string         `json:"error_message,omitempty" db:"error_message"`
	Warnings      []string        `json:"warnings,omitempty" db:"warnings"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty" db:"processed_at"`
}

// DataMappingConfig represents data mapping configuration
type DataMappingConfig struct {
	ConfigID          uuid.UUID `json:"config_id" db:"config_id"`
	SourceField       string    `json:"source_field" db:"source_field"`
	TargetField       string    `json:"target_field" db:"target_field"`
	DataType          string    `json:"data_type" db:"data_type"`
	TransformationRule *string  `json:"transformation_rule,omitempty" db:"transformation_rule"`
	ValidationRule    *string   `json:"validation_rule,omitempty" db:"validation_rule"`
	IsRequired        bool      `json:"is_required" db:"is_required"`
	DefaultValue      *string   `json:"default_value,omitempty" db:"default_value"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// ImportRow represents a row being imported
type ImportRow struct {
	RowNumber int                    `json:"row_number"`
	Data      map[string]interface{} `json:"data"`
	Errors    []string               `json:"errors,omitempty"`
	Warnings  []string               `json:"warnings,omitempty"`
	Status    string                 `json:"status"`
}

// Value implements the driver.Valuer interface
func (i ImportDetail) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements the sql.Scanner interface
func (i *ImportDetail) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &i)
}
