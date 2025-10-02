-- Migration 028 Down

DROP INDEX IF EXISTS idx_payment_schedules_disbursement_id;
DROP INDEX IF EXISTS idx_payment_schedules_status;
DROP INDEX IF EXISTS idx_payment_schedules_due_date;
DROP INDEX IF EXISTS idx_payment_schedules_allocation_id;
DROP INDEX IF EXISTS idx_disbursement_batch_items_disbursement_id;
DROP INDEX IF EXISTS idx_disbursement_batch_items_batch_id;
DROP INDEX IF EXISTS idx_disbursement_batches_scheduled_date;
DROP INDEX IF EXISTS idx_disbursement_batches_status;
DROP INDEX IF EXISTS idx_disbursement_batches_scholarship_id;
DROP INDEX IF EXISTS idx_disbursement_batches_round_id;
DROP INDEX IF EXISTS idx_disbursement_records_transferred_by;
DROP INDEX IF EXISTS idx_disbursement_records_transfer_status;
DROP INDEX IF EXISTS idx_disbursement_records_transfer_date;
DROP INDEX IF EXISTS idx_disbursement_records_account_id;
DROP INDEX IF EXISTS idx_disbursement_records_allocation_id;
DROP INDEX IF EXISTS idx_student_bank_accounts_verified;
DROP INDEX IF EXISTS idx_student_bank_accounts_is_primary;
DROP INDEX IF EXISTS idx_student_bank_accounts_student_id;

DROP TABLE IF EXISTS payment_schedules;
DROP TABLE IF EXISTS disbursement_batch_items;
DROP TABLE IF EXISTS disbursement_batches;
DROP TABLE IF EXISTS disbursement_records;
DROP TABLE IF EXISTS student_bank_accounts;
