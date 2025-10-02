-- Migration: Add Payment System Tables
-- Created: 2025-10-01
-- Description: Adds payment transactions, disbursement schedules, bank transfers, and payment confirmations

-- 1. Payment Methods Table
CREATE TABLE IF NOT EXISTS payment_methods (
    method_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    method_name VARCHAR(100) NOT NULL,
    method_code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    configuration JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Payment Transactions Table
CREATE TABLE IF NOT EXISTS payment_transactions (
    transaction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    allocation_id INTEGER REFERENCES scholarship_allocations(allocation_id),
    amount DECIMAL(12,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    bank_code VARCHAR(10),
    account_number VARCHAR(20),
    payment_date TIMESTAMP NOT NULL,
    payment_status VARCHAR(30) DEFAULT 'pending',
    reference_number VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Disbursement Schedules Table
CREATE TABLE IF NOT EXISTS disbursement_schedules (
    schedule_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    allocation_id INTEGER REFERENCES scholarship_allocations(allocation_id),
    installment_number INTEGER NOT NULL,
    due_date DATE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    status VARCHAR(30) DEFAULT 'scheduled',
    paid_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Bank Transfer Logs Table
CREATE TABLE IF NOT EXISTS bank_transfer_logs (
    transfer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES payment_transactions(transaction_id),
    bank_code VARCHAR(10) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    transfer_amount DECIMAL(12,2) NOT NULL,
    transfer_fee DECIMAL(8,2) DEFAULT 0,
    transfer_date TIMESTAMP NOT NULL,
    confirmation_code VARCHAR(50),
    transfer_status VARCHAR(30) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Payment Confirmations Table
CREATE TABLE IF NOT EXISTS payment_confirmations (
    confirmation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID REFERENCES payment_transactions(transaction_id),
    confirmed_by UUID REFERENCES users(user_id),
    confirmation_method VARCHAR(50) NOT NULL,
    confirmation_date TIMESTAMP NOT NULL,
    slip_image_path VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_payment_transactions_allocation ON payment_transactions(allocation_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_status ON payment_transactions(payment_status);
CREATE INDEX IF NOT EXISTS idx_disbursement_schedules_allocation ON disbursement_schedules(allocation_id);
CREATE INDEX IF NOT EXISTS idx_disbursement_schedules_status ON disbursement_schedules(status);
CREATE INDEX IF NOT EXISTS idx_bank_transfer_logs_transaction ON bank_transfer_logs(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payment_confirmations_transaction ON payment_confirmations(transaction_id);

-- Insert default payment methods
INSERT INTO payment_methods (method_name, method_code, description, is_active) VALUES
('Bank Transfer', 'bank_transfer', 'โอนเงินผ่านธนาคาร', true),
('Cheque', 'cheque', 'เช็คธนาคาร', true),
('Cash', 'cash', 'เงินสด', false),
('Mobile Banking', 'mobile_banking', 'Mobile Banking', true)
ON CONFLICT (method_code) DO NOTHING;

-- Add trigger for updated_at
CREATE TRIGGER update_payment_transactions_updated_at
    BEFORE UPDATE ON payment_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_disbursement_schedules_updated_at
    BEFORE UPDATE ON disbursement_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payment_confirmations_updated_at
    BEFORE UPDATE ON payment_confirmations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE payment_methods IS 'วิธีการจ่ายเงินทุน';
COMMENT ON TABLE payment_transactions IS 'บันทึกการจ่ายเงินทุน';
COMMENT ON TABLE disbursement_schedules IS 'ตารางการจ่ายเงินแบบแบ่งงวด';
COMMENT ON TABLE bank_transfer_logs IS 'บันทึกการโอนเงินผ่านธนาคาร';
COMMENT ON TABLE payment_confirmations IS 'การยืนยันการจ่ายเงิน';
