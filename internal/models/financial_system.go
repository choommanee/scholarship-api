package models

import (
	"time"

	"github.com/google/uuid"
)

// StudentBankAccount represents a student's bank account for receiving scholarships
type StudentBankAccount struct {
	AccountID     uuid.UUID  `json:"account_id" db:"account_id"`
	StudentID     string     `json:"student_id" db:"student_id"`
	BankName      string     `json:"bank_name" db:"bank_name"`
	BankCode      *string    `json:"bank_code" db:"bank_code"`
	BranchName    *string    `json:"branch_name" db:"branch_name"`
	AccountNumber string     `json:"account_number" db:"account_number"`
	AccountName   string     `json:"account_name" db:"account_name"`
	AccountType   *string    `json:"account_type" db:"account_type"`
	IsPrimary     bool       `json:"is_primary" db:"is_primary"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	Verified      bool       `json:"verified" db:"verified"`
	VerifiedBy    *uuid.UUID `json:"verified_by" db:"verified_by"`
	VerifiedAt    *time.Time `json:"verified_at" db:"verified_at"`
	VerificationMethod *string `json:"verification_method" db:"verification_method"`
	BankBookImage *string    `json:"bank_book_image" db:"bank_book_image"`
	Notes         *string    `json:"notes" db:"notes"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// DisbursementRecord represents a scholarship disbursement transaction
type DisbursementRecord struct {
	DisbursementID   uuid.UUID  `json:"disbursement_id" db:"disbursement_id"`
	AllocationID     uint       `json:"allocation_id" db:"allocation_id"`
	AccountID        *uuid.UUID `json:"account_id" db:"account_id"`
	TransferDate     time.Time  `json:"transfer_date" db:"transfer_date"`
	TransferTime     *time.Time `json:"transfer_time" db:"transfer_time"`
	Amount           float64    `json:"amount" db:"amount"`
	BankName         *string    `json:"bank_name" db:"bank_name"`
	AccountNumber    *string    `json:"account_number" db:"account_number"`
	AccountName      *string    `json:"account_name" db:"account_name"`
	TransferRef      *string    `json:"transfer_ref" db:"transfer_ref"`
	TransactionID    *string    `json:"transaction_id" db:"transaction_id"`
	TransferStatus   string     `json:"transfer_status" db:"transfer_status"`
	TransferProofURL *string    `json:"transfer_proof_url" db:"transfer_proof_url"`
	ReceiptURL       *string    `json:"receipt_url" db:"receipt_url"`
	TransferredBy    *uuid.UUID `json:"transferred_by" db:"transferred_by"`
	ApprovedBy       *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at" db:"approved_at"`
	FailureReason    *string    `json:"failure_reason" db:"failure_reason"`
	RetryCount       int        `json:"retry_count" db:"retry_count"`
	Notes            *string    `json:"notes" db:"notes"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// DisbursementBatch represents a batch of disbursements for multiple students
type DisbursementBatch struct {
	BatchID        uuid.UUID  `json:"batch_id" db:"batch_id"`
	BatchName      string     `json:"batch_name" db:"batch_name"`
	BatchCode      *string    `json:"batch_code" db:"batch_code"`
	RoundID        *uint      `json:"round_id" db:"round_id"`
	ScholarshipID  *uint      `json:"scholarship_id" db:"scholarship_id"`
	TotalStudents  int        `json:"total_students" db:"total_students"`
	TotalAmount    float64    `json:"total_amount" db:"total_amount"`
	CompletedCount int        `json:"completed_count" db:"completed_count"`
	FailedCount    int        `json:"failed_count" db:"failed_count"`
	Status         string     `json:"status" db:"status"`
	ScheduledDate  *time.Time `json:"scheduled_date" db:"scheduled_date"`
	ProcessedDate  *time.Time `json:"processed_date" db:"processed_date"`
	CreatedBy      *uuid.UUID `json:"created_by" db:"created_by"`
	ApprovedBy     *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at" db:"approved_at"`
	Notes          *string    `json:"notes" db:"notes"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// DisbursementBatchItem represents an individual disbursement within a batch
type DisbursementBatchItem struct {
	ItemID           uuid.UUID  `json:"item_id" db:"item_id"`
	BatchID          uuid.UUID  `json:"batch_id" db:"batch_id"`
	DisbursementID   uuid.UUID  `json:"disbursement_id" db:"disbursement_id"`
	SequenceNumber   *int       `json:"sequence_number" db:"sequence_number"`
	ProcessingStatus string     `json:"processing_status" db:"processing_status"`
	ProcessedAt      *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// PaymentSchedule represents an installment payment schedule for scholarships
type PaymentSchedule struct {
	ScheduleID        uuid.UUID  `json:"schedule_id" db:"schedule_id"`
	AllocationID      uint       `json:"allocation_id" db:"allocation_id"`
	InstallmentNumber int        `json:"installment_number" db:"installment_number"`
	TotalInstallments int        `json:"total_installments" db:"total_installments"`
	Amount            float64    `json:"amount" db:"amount"`
	DueDate           time.Time  `json:"due_date" db:"due_date"`
	Status            string     `json:"status" db:"status"`
	DisbursementID    *uuid.UUID `json:"disbursement_id" db:"disbursement_id"`
	PaidDate          *time.Time `json:"paid_date" db:"paid_date"`
	PaidAmount        *float64   `json:"paid_amount" db:"paid_amount"`
	Notes             *string    `json:"notes" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}
