package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents a payment method
type PaymentMethod struct {
	MethodID      uuid.UUID      `json:"method_id" db:"method_id"`
	MethodName    string         `json:"method_name" db:"method_name"`
	MethodCode    string         `json:"method_code" db:"method_code"`
	Description   string         `json:"description" db:"description"`
	IsActive      bool           `json:"is_active" db:"is_active"`
	Configuration json.RawMessage `json:"configuration" db:"configuration"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
}

// PaymentTransaction represents a payment transaction
type PaymentTransaction struct {
	TransactionID   uuid.UUID `json:"transaction_id" db:"transaction_id"`
	AllocationID    int       `json:"allocation_id" db:"allocation_id"`
	Amount          float64   `json:"amount" db:"amount"`
	PaymentMethod   string    `json:"payment_method" db:"payment_method"`
	BankCode        *string   `json:"bank_code,omitempty" db:"bank_code"`
	AccountNumber   *string   `json:"account_number,omitempty" db:"account_number"`
	PaymentDate     time.Time `json:"payment_date" db:"payment_date"`
	PaymentStatus   string    `json:"payment_status" db:"payment_status"`
	ReferenceNumber *string   `json:"reference_number,omitempty" db:"reference_number"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// DisbursementSchedule represents a disbursement schedule
type DisbursementSchedule struct {
	ScheduleID       uuid.UUID  `json:"schedule_id" db:"schedule_id"`
	AllocationID     int        `json:"allocation_id" db:"allocation_id"`
	InstallmentNumber int       `json:"installment_number" db:"installment_number"`
	DueDate          time.Time  `json:"due_date" db:"due_date"`
	Amount           float64    `json:"amount" db:"amount"`
	Status           string     `json:"status" db:"status"`
	PaidDate         *time.Time `json:"paid_date,omitempty" db:"paid_date"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// BankTransferLog represents a bank transfer log
type BankTransferLog struct {
	TransferID       uuid.UUID `json:"transfer_id" db:"transfer_id"`
	TransactionID    uuid.UUID `json:"transaction_id" db:"transaction_id"`
	BankCode         string    `json:"bank_code" db:"bank_code"`
	BankName         string    `json:"bank_name" db:"bank_name"`
	AccountNumber    string    `json:"account_number" db:"account_number"`
	AccountName      string    `json:"account_name" db:"account_name"`
	TransferAmount   float64   `json:"transfer_amount" db:"transfer_amount"`
	TransferFee      float64   `json:"transfer_fee" db:"transfer_fee"`
	TransferDate     time.Time `json:"transfer_date" db:"transfer_date"`
	ConfirmationCode *string   `json:"confirmation_code,omitempty" db:"confirmation_code"`
	TransferStatus   string    `json:"transfer_status" db:"transfer_status"`
	ErrorMessage     *string   `json:"error_message,omitempty" db:"error_message"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PaymentConfirmation represents a payment confirmation
type PaymentConfirmation struct {
	ConfirmationID     uuid.UUID  `json:"confirmation_id" db:"confirmation_id"`
	TransactionID      uuid.UUID  `json:"transaction_id" db:"transaction_id"`
	ConfirmedBy        uuid.UUID  `json:"confirmed_by" db:"confirmed_by"`
	ConfirmationMethod string     `json:"confirmation_method" db:"confirmation_method"`
	ConfirmationDate   time.Time  `json:"confirmation_date" db:"confirmation_date"`
	SlipImagePath      *string    `json:"slip_image_path,omitempty" db:"slip_image_path"`
	Notes              *string    `json:"notes,omitempty" db:"notes"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// Value implements the driver.Valuer interface for database serialization
func (p PaymentMethod) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for database deserialization
func (p *PaymentMethod) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, &p)
}
