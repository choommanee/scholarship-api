-- Rollback: Drop Payment System Tables

DROP TABLE IF EXISTS payment_confirmations CASCADE;
DROP TABLE IF EXISTS bank_transfer_logs CASCADE;
DROP TABLE IF EXISTS disbursement_schedules CASCADE;
DROP TABLE IF EXISTS payment_transactions CASCADE;
DROP TABLE IF EXISTS payment_methods CASCADE;
