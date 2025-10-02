package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ScholarshipStatistics represents scholarship statistics
type ScholarshipStatistics struct {
	StatID                  uuid.UUID `json:"stat_id" db:"stat_id"`
	AcademicYear            string    `json:"academic_year" db:"academic_year"`
	ScholarshipRound        string    `json:"scholarship_round" db:"scholarship_round"`
	TotalApplications       int       `json:"total_applications" db:"total_applications"`
	ApprovedApplications    int       `json:"approved_applications" db:"approved_applications"`
	RejectedApplications    int       `json:"rejected_applications" db:"rejected_applications"`
	TotalBudget             float64   `json:"total_budget" db:"total_budget"`
	AllocatedBudget         float64   `json:"allocated_budget" db:"allocated_budget"`
	RemainingBudget         float64   `json:"remaining_budget" db:"remaining_budget"`
	AverageAmount           float64   `json:"average_amount" db:"average_amount"`
	SuccessRate             float64   `json:"success_rate" db:"success_rate"`
	ProcessingTimeAvg       int       `json:"processing_time_avg" db:"processing_time_avg"`
	TotalFaculties          int       `json:"total_faculties" db:"total_faculties"`
	MostPopularScholarship  *string   `json:"most_popular_scholarship,omitempty" db:"most_popular_scholarship"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
}

// ApplicationAnalytics represents application analytics
type ApplicationAnalytics struct {
	AnalyticsID        uuid.UUID       `json:"analytics_id" db:"analytics_id"`
	ApplicationID      int             `json:"application_id" db:"application_id"`
	ProcessingTime     int             `json:"processing_time" db:"processing_time"`
	TotalSteps         int             `json:"total_steps" db:"total_steps"`
	CompletedSteps     int             `json:"completed_steps" db:"completed_steps"`
	BottleneckStep     *string         `json:"bottleneck_step,omitempty" db:"bottleneck_step"`
	TimeInEachStep     json.RawMessage `json:"time_in_each_step,omitempty" db:"time_in_each_step"`
	DocumentUploadTime *int            `json:"document_upload_time,omitempty" db:"document_upload_time"`
	ReviewTime         *int            `json:"review_time,omitempty" db:"review_time"`
	InterviewScore     *float64        `json:"interview_score,omitempty" db:"interview_score"`
	FinalScore         *float64        `json:"final_score,omitempty" db:"final_score"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
}

// TimeInStep represents time spent in each step
type TimeInStep struct {
	Step     string `json:"step"`
	Duration int    `json:"duration"` // in minutes
}

// Value implements the driver.Valuer interface
func (s ScholarshipStatistics) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *ScholarshipStatistics) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &s)
}
