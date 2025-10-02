package models

import (
	"time"

	"github.com/google/uuid"
)

type InterviewSchedule struct {
	ScheduleID     uint      `json:"schedule_id" db:"schedule_id"`
	ScholarshipID  uint      `json:"scholarship_id" db:"scholarship_id"`
	InterviewDate  time.Time `json:"interview_date" db:"interview_date"`
	StartTime      string    `json:"start_time" db:"start_time"`
	EndTime        string    `json:"end_time" db:"end_time"`
	Location       *string   `json:"location" db:"location"`
	MaxApplicants  int       `json:"max_applicants" db:"max_applicants"`
	InterviewerIDs *string   `json:"interviewer_ids" db:"interviewer_ids"`
	Notes          *string   `json:"notes" db:"notes"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedBy      uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type InterviewAppointment struct {
	AppointmentID     uint       `json:"appointment_id" db:"appointment_id"`
	ApplicationID     uint       `json:"application_id" db:"application_id"`
	ScheduleID        uint       `json:"schedule_id" db:"schedule_id"`
	AppointmentStatus string     `json:"appointment_status" db:"appointment_status"`
	StudentConfirmed  bool       `json:"student_confirmed" db:"student_confirmed"`
	ConfirmationDate  *time.Time `json:"confirmation_date" db:"confirmation_date"`
	ActualStartTime   *time.Time `json:"actual_start_time" db:"actual_start_time"`
	ActualEndTime     *time.Time `json:"actual_end_time" db:"actual_end_time"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

type InterviewResult struct {
	ResultID       uint      `json:"result_id" db:"result_id"`
	AppointmentID  uint      `json:"appointment_id" db:"appointment_id"`
	InterviewerID  uuid.UUID `json:"interviewer_id" db:"interviewer_id"`
	Scores         *string   `json:"scores" db:"scores"`
	OverallScore   *float64  `json:"overall_score" db:"overall_score"`
	Comments       *string   `json:"comments" db:"comments"`
	Recommendation *string   `json:"recommendation" db:"recommendation"`
	InterviewNotes *string   `json:"interview_notes" db:"interview_notes"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}