package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"scholarship-system/internal/models"
)

// PaymentRepository handles payment-related database operations
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreateTransaction creates a new payment transaction
func (r *PaymentRepository) CreateTransaction(tx *models.PaymentTransaction) error {
	query := `
		INSERT INTO payment_transactions (
			transaction_id, allocation_id, amount, payment_method,
			bank_code, account_number, payment_date, payment_status,
			reference_number, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(
		query,
		tx.TransactionID, tx.AllocationID, tx.Amount, tx.PaymentMethod,
		tx.BankCode, tx.AccountNumber, tx.PaymentDate, tx.PaymentStatus,
		tx.ReferenceNumber, tx.Notes,
	).Scan(&tx.CreatedAt, &tx.UpdatedAt)
}

// GetTransactionByID retrieves a transaction by ID
func (r *PaymentRepository) GetTransactionByID(id uuid.UUID) (*models.PaymentTransaction, error) {
	query := `
		SELECT transaction_id, allocation_id, amount, payment_method,
			   bank_code, account_number, payment_date, payment_status,
			   reference_number, notes, created_at, updated_at
		FROM payment_transactions
		WHERE transaction_id = $1`

	tx := &models.PaymentTransaction{}
	err := r.db.QueryRow(query, id).Scan(
		&tx.TransactionID, &tx.AllocationID, &tx.Amount, &tx.PaymentMethod,
		&tx.BankCode, &tx.AccountNumber, &tx.PaymentDate, &tx.PaymentStatus,
		&tx.ReferenceNumber, &tx.Notes, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// UpdateTransactionStatus updates transaction status
func (r *PaymentRepository) UpdateTransactionStatus(id uuid.UUID, status string) error {
	query := `
		UPDATE payment_transactions
		SET payment_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE transaction_id = $2`

	_, err := r.db.Exec(query, status, id)
	return err
}

// GetTransactionsByAllocation gets all transactions for an allocation
func (r *PaymentRepository) GetTransactionsByAllocation(allocationID int) ([]models.PaymentTransaction, error) {
	query := `
		SELECT transaction_id, allocation_id, amount, payment_method,
			   bank_code, account_number, payment_date, payment_status,
			   reference_number, notes, created_at, updated_at
		FROM payment_transactions
		WHERE allocation_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, allocationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.PaymentTransaction
	for rows.Next() {
		var tx models.PaymentTransaction
		err := rows.Scan(
			&tx.TransactionID, &tx.AllocationID, &tx.Amount, &tx.PaymentMethod,
			&tx.BankCode, &tx.AccountNumber, &tx.PaymentDate, &tx.PaymentStatus,
			&tx.ReferenceNumber, &tx.Notes, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// CreateDisbursementSchedule creates a disbursement schedule
func (r *PaymentRepository) CreateDisbursementSchedule(schedule *models.DisbursementSchedule) error {
	query := `
		INSERT INTO disbursement_schedules (
			schedule_id, allocation_id, installment_number, due_date,
			amount, status
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(
		query,
		schedule.ScheduleID, schedule.AllocationID, schedule.InstallmentNumber,
		schedule.DueDate, schedule.Amount, schedule.Status,
	).Scan(&schedule.CreatedAt, &schedule.UpdatedAt)
}

// GetDisbursementSchedules gets disbursement schedules for an allocation
func (r *PaymentRepository) GetDisbursementSchedules(allocationID int) ([]models.DisbursementSchedule, error) {
	query := `
		SELECT schedule_id, allocation_id, installment_number, due_date,
			   amount, status, paid_date, created_at, updated_at
		FROM disbursement_schedules
		WHERE allocation_id = $1
		ORDER BY installment_number`

	rows, err := r.db.Query(query, allocationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.DisbursementSchedule
	for rows.Next() {
		var s models.DisbursementSchedule
		err := rows.Scan(
			&s.ScheduleID, &s.AllocationID, &s.InstallmentNumber, &s.DueDate,
			&s.Amount, &s.Status, &s.PaidDate, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// MarkDisbursementPaid marks a disbursement as paid
func (r *PaymentRepository) MarkDisbursementPaid(scheduleID uuid.UUID) error {
	query := `
		UPDATE disbursement_schedules
		SET status = 'paid', paid_date = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE schedule_id = $1`

	_, err := r.db.Exec(query, scheduleID)
	return err
}

// CreateBankTransferLog creates a bank transfer log
func (r *PaymentRepository) CreateBankTransferLog(log *models.BankTransferLog) error {
	query := `
		INSERT INTO bank_transfer_logs (
			transfer_id, transaction_id, bank_code, bank_name,
			account_number, account_name, transfer_amount, transfer_fee,
			transfer_date, confirmation_code, transfer_status, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		log.TransferID, log.TransactionID, log.BankCode, log.BankName,
		log.AccountNumber, log.AccountName, log.TransferAmount, log.TransferFee,
		log.TransferDate, log.ConfirmationCode, log.TransferStatus, log.ErrorMessage,
	).Scan(&log.CreatedAt)
}

// GetPaymentMethods retrieves all active payment methods
func (r *PaymentRepository) GetPaymentMethods() ([]models.PaymentMethod, error) {
	query := `
		SELECT method_id, method_name, method_code, description, is_active,
		       COALESCE(configuration, '{}'::jsonb) as configuration, created_at
		FROM payment_methods
		WHERE is_active = true
		ORDER BY method_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []models.PaymentMethod
	for rows.Next() {
		var m models.PaymentMethod
		err := rows.Scan(
			&m.MethodID, &m.MethodName, &m.MethodCode, &m.Description,
			&m.IsActive, &m.Configuration, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return methods, nil
}

// GetPendingDisbursements gets disbursements due for payment
func (r *PaymentRepository) GetPendingDisbursements(beforeDate time.Time) ([]models.DisbursementSchedule, error) {
	query := `
		SELECT schedule_id, allocation_id, installment_number, due_date,
			   amount, status, paid_date, created_at, updated_at
		FROM disbursement_schedules
		WHERE status = 'scheduled' AND due_date <= $1
		ORDER BY due_date`

	rows, err := r.db.Query(query, beforeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.DisbursementSchedule
	for rows.Next() {
		var s models.DisbursementSchedule
		err := rows.Scan(
			&s.ScheduleID, &s.AllocationID, &s.InstallmentNumber, &s.DueDate,
			&s.Amount, &s.Status, &s.PaidDate, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}
