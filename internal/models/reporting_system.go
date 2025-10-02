package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ReportTemplate represents a template for generating reports
type ReportTemplate struct {
	TemplateID      uuid.UUID      `json:"template_id" db:"template_id"`
	TemplateName    string         `json:"template_name" db:"template_name"`
	TemplateCode    string         `json:"template_code" db:"template_code"`
	Description     *string        `json:"description" db:"description"`
	ReportType      string         `json:"report_type" db:"report_type"`
	TemplateConfig  *string        `json:"template_config" db:"template_config"` // JSONB stored as string
	QueryTemplate   *string        `json:"query_template" db:"query_template"`
	OutputFormat    pq.StringArray `json:"output_format" db:"output_format"`
	DefaultFormat   string         `json:"default_format" db:"default_format"`
	Columns         *string        `json:"columns" db:"columns"`         // JSONB stored as string
	Filters         *string        `json:"filters" db:"filters"`         // JSONB stored as string
	Sorting         *string        `json:"sorting" db:"sorting"`         // JSONB stored as string
	Grouping        *string        `json:"grouping" db:"grouping"`       // JSONB stored as string
	AccessibleRoles pq.StringArray `json:"accessible_roles" db:"accessible_roles"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	IsSystem        bool           `json:"is_system" db:"is_system"`
	CreatedBy       *uuid.UUID     `json:"created_by" db:"created_by"`
	UpdatedBy       *uuid.UUID     `json:"updated_by" db:"updated_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// GeneratedReport represents a report that has been generated
type GeneratedReport struct {
	ReportID          uuid.UUID  `json:"report_id" db:"report_id"`
	TemplateID        *uuid.UUID `json:"template_id" db:"template_id"`
	ReportName        string     `json:"report_name" db:"report_name"`
	ReportDescription *string    `json:"report_description" db:"report_description"`
	ReportPeriod      *string    `json:"report_period" db:"report_period"`
	StartDate         *time.Time `json:"start_date" db:"start_date"`
	EndDate           *time.Time `json:"end_date" db:"end_date"`
	FilterParams      *string    `json:"filter_params" db:"filter_params"` // JSONB stored as string
	FilePath          *string    `json:"file_path" db:"file_path"`
	FileName          *string    `json:"file_name" db:"file_name"`
	FileSize          *int       `json:"file_size" db:"file_size"`
	FileFormat        *string    `json:"file_format" db:"file_format"`
	MimeType          *string    `json:"mime_type" db:"mime_type"`
	TotalRecords      *int       `json:"total_records" db:"total_records"`
	TotalPages        *int       `json:"total_pages" db:"total_pages"`
	Status            string     `json:"status" db:"status"`
	ExpiresAt         *time.Time `json:"expires_at" db:"expires_at"`
	IsExpired         bool       `json:"is_expired" db:"is_expired"`
	ErrorMessage      *string    `json:"error_message" db:"error_message"`
	GeneratedBy       uuid.UUID  `json:"generated_by" db:"generated_by"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// ReportSchedule represents a scheduled automatic report generation
type ReportSchedule struct {
	ScheduleID     uuid.UUID  `json:"schedule_id" db:"schedule_id"`
	TemplateID     uuid.UUID  `json:"template_id" db:"template_id"`
	ScheduleName   string     `json:"schedule_name" db:"schedule_name"`
	Description    *string    `json:"description" db:"description"`
	Frequency      string     `json:"frequency" db:"frequency"`
	CronExpression *string    `json:"cron_expression" db:"cron_expression"`
	NextRunDate    *time.Time `json:"next_run_date" db:"next_run_date"`
	LastRunDate    *time.Time `json:"last_run_date" db:"last_run_date"`
	DefaultFilters *string    `json:"default_filters" db:"default_filters"` // JSONB stored as string
	Recipients     *string    `json:"recipients" db:"recipients"`           // JSONB stored as string
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedBy      *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ReportAccessLog represents a log of report access
type ReportAccessLog struct {
	LogID      uuid.UUID  `json:"log_id" db:"log_id"`
	ReportID   *uuid.UUID `json:"report_id" db:"report_id"`
	TemplateID *uuid.UUID `json:"template_id" db:"template_id"`
	UserID     *uuid.UUID `json:"user_id" db:"user_id"`
	Action     string     `json:"action" db:"action"`
	IPAddress  *string    `json:"ip_address" db:"ip_address"`
	UserAgent  *string    `json:"user_agent" db:"user_agent"`
	AccessedAt time.Time  `json:"accessed_at" db:"accessed_at"`
}

// DashboardWidget represents a widget for dashboard display
type DashboardWidget struct {
	WidgetID        uuid.UUID      `json:"widget_id" db:"widget_id"`
	WidgetName      string         `json:"widget_name" db:"widget_name"`
	WidgetType      string         `json:"widget_type" db:"widget_type"`
	Description     *string        `json:"description" db:"description"`
	DataSource      *string        `json:"data_source" db:"data_source"`
	QueryTemplate   *string        `json:"query_template" db:"query_template"`
	APIEndpoint     *string        `json:"api_endpoint" db:"api_endpoint"`
	Config          *string        `json:"config" db:"config"` // JSONB stored as string
	DisplayOrder    *int           `json:"display_order" db:"display_order"`
	Width           int            `json:"width" db:"width"`
	Height          int            `json:"height" db:"height"`
	AccessibleRoles pq.StringArray `json:"accessible_roles" db:"accessible_roles"`
	RefreshInterval *int           `json:"refresh_interval" db:"refresh_interval"`
	CacheDuration   *int           `json:"cache_duration" db:"cache_duration"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	CreatedBy       *uuid.UUID     `json:"created_by" db:"created_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}
