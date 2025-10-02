package models

import (
	"time"

	"github.com/google/uuid"
)

type ScholarshipAllocation struct {
	AllocationID       uint       `json:"allocation_id" db:"allocation_id"`
	ApplicationID      uint       `json:"application_id" db:"application_id"`
	ScholarshipID      uint       `json:"scholarship_id" db:"scholarship_id"`
	AllocatedAmount    float64    `json:"allocated_amount" db:"allocated_amount"`
	AllocationStatus   string     `json:"allocation_status" db:"allocation_status"`
	AllocationDate     time.Time  `json:"allocation_date" db:"allocation_date"`
	DisbursementMethod *string    `json:"disbursement_method" db:"disbursement_method"`
	BankAccount        *string    `json:"bank_account" db:"bank_account"`
	BankName           *string    `json:"bank_name" db:"bank_name"`
	TransferDate       *time.Time `json:"transfer_date" db:"transfer_date"`
	TransferReference  *string    `json:"transfer_reference" db:"transfer_reference"`
	AllocatedBy        uuid.UUID  `json:"allocated_by" db:"allocated_by"`
	ApprovedBy         *uuid.UUID `json:"approved_by" db:"approved_by"`
	Notes              *string    `json:"notes" db:"notes"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

type AcademicProgressTracking struct {
	TrackingID     uint      `json:"tracking_id" db:"tracking_id"`
	StudentID      string    `json:"student_id" db:"student_id"`
	AllocationID   uint      `json:"allocation_id" db:"allocation_id"`
	Semester       string    `json:"semester" db:"semester"`
	GPA            *float64  `json:"gpa" db:"gpa"`
	CreditsEarned  *int      `json:"credits_earned" db:"credits_earned"`
	AcademicStatus *string   `json:"academic_status" db:"academic_status"`
	ReportDate     time.Time `json:"report_date" db:"report_date"`
	Notes          *string   `json:"notes" db:"notes"`
	CreatedBy      uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ImportLog struct {
	ImportID          uint       `json:"import_id" db:"import_id"`
	ImportType        string     `json:"import_type" db:"import_type"`
	FileName          string     `json:"file_name" db:"file_name"`
	TotalRecords      *int       `json:"total_records" db:"total_records"`
	SuccessfulRecords *int       `json:"successful_records" db:"successful_records"`
	FailedRecords     *int       `json:"failed_records" db:"failed_records"`
	ErrorDetails      *string    `json:"error_details" db:"error_details"`
	ImportedBy        uuid.UUID  `json:"imported_by" db:"imported_by"`
	ImportStatus      string     `json:"import_status" db:"import_status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	CompletedAt       *time.Time `json:"completed_at" db:"completed_at"`
}
