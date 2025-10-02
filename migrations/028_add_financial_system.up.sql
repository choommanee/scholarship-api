-- Migration 028: Financial System Tables

-- 1. Student Bank Accounts (บัญชีธนาคารของนักศึกษา)
CREATE TABLE IF NOT EXISTS student_bank_accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id VARCHAR(20) NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,

    -- ข้อมูลธนาคาร
    bank_name VARCHAR(100) NOT NULL,
    bank_code VARCHAR(10),              -- รหัสธนาคาร
    branch_name VARCHAR(255),

    -- ข้อมูลบัญชี
    account_number VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL, -- ชื่อบัญชี (ต้องตรงกับนักศึกษา)
    account_type VARCHAR(50),           -- savings, current

    -- สถานะบัญชี
    is_primary BOOLEAN DEFAULT TRUE,    -- บัญชีหลักสำหรับรับทุน
    is_active BOOLEAN DEFAULT TRUE,
    verified BOOLEAN DEFAULT FALSE,

    -- การยืนยันบัญชี
    verified_by UUID REFERENCES users(user_id),
    verified_at TIMESTAMP,
    verification_method VARCHAR(50),    -- document, test_transfer

    -- เอกสารประกอบ
    bank_book_image TEXT,               -- สำเนาหน้าบัญชี

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_primary_account UNIQUE (student_id, is_primary)
        DEFERRABLE INITIALLY DEFERRED
);

-- 2. Disbursement Records (บันทึกการโอนเงินทุน)
CREATE TABLE IF NOT EXISTS disbursement_records (
    disbursement_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    allocation_id INTEGER NOT NULL REFERENCES scholarship_allocations(allocation_id) ON DELETE CASCADE,
    account_id UUID REFERENCES student_bank_accounts(account_id),

    -- รายละเอียดการโอน
    transfer_date DATE NOT NULL,
    transfer_time TIME,
    amount DECIMAL(12,2) NOT NULL,

    -- ข้อมูลธนาคาร
    bank_name VARCHAR(100),
    account_number VARCHAR(50),
    account_name VARCHAR(255),

    -- หมายเลขอ้างอิง
    transfer_ref VARCHAR(100),          -- เลขที่อ้างอิงการโอน
    transaction_id VARCHAR(100),        -- Transaction ID จากธนาคาร

    -- สถานะการโอน
    transfer_status VARCHAR(50) DEFAULT 'pending' CHECK (
        transfer_status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')
    ),

    -- หลักฐานการโอน
    transfer_proof_url TEXT,            -- สลิปโอนเงิน
    receipt_url TEXT,                   -- ใบเสร็จ

    -- ผู้ดำเนินการ
    transferred_by UUID REFERENCES users(user_id),
    approved_by UUID REFERENCES users(user_id),
    approved_at TIMESTAMP,

    -- กรณีล้มเหลว
    failure_reason TEXT,
    retry_count INTEGER DEFAULT 0,

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Disbursement Batches (ชุดการโอนเงินทุน - สำหรับโอนครั้งละหลายคน)
CREATE TABLE IF NOT EXISTS disbursement_batches (
    batch_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ข้อมูลชุดการโอน
    batch_name VARCHAR(255) NOT NULL,
    batch_code VARCHAR(50) UNIQUE,

    -- รอบทุน
    round_id INTEGER REFERENCES scholarship_rounds(round_id),
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id),

    -- สถิติ
    total_students INTEGER DEFAULT 0,
    total_amount DECIMAL(15,2) DEFAULT 0.00,
    completed_count INTEGER DEFAULT 0,
    failed_count INTEGER DEFAULT 0,

    -- สถานะ
    status VARCHAR(50) DEFAULT 'draft' CHECK (
        status IN ('draft', 'pending', 'processing', 'completed', 'partial', 'failed')
    ),

    -- กำหนดการ
    scheduled_date DATE,
    processed_date DATE,

    -- ผู้ดำเนินการ
    created_by UUID REFERENCES users(user_id),
    approved_by UUID REFERENCES users(user_id),
    approved_at TIMESTAMP,

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Disbursement Batch Items (รายการในชุดการโอน)
CREATE TABLE IF NOT EXISTS disbursement_batch_items (
    item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES disbursement_batches(batch_id) ON DELETE CASCADE,
    disbursement_id UUID NOT NULL REFERENCES disbursement_records(disbursement_id) ON DELETE CASCADE,

    -- ลำดับในชุด
    sequence_number INTEGER,

    -- สถานะการประมวลผล
    processing_status VARCHAR(50) DEFAULT 'pending',
    processed_at TIMESTAMP,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Payment Schedule (ตารางการจ่ายทุน - กรณีทุนแบ่งจ่ายหลายงวด)
CREATE TABLE IF NOT EXISTS payment_schedules (
    schedule_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    allocation_id INTEGER NOT NULL REFERENCES scholarship_allocations(allocation_id) ON DELETE CASCADE,

    -- งวดที่
    installment_number INTEGER NOT NULL,
    total_installments INTEGER NOT NULL,

    -- จำนวนเงิน
    amount DECIMAL(12,2) NOT NULL,

    -- กำหนดจ่าย
    due_date DATE NOT NULL,

    -- สถานะ
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (
        status IN ('scheduled', 'pending', 'paid', 'overdue', 'cancelled')
    ),

    -- การจ่ายจริง
    disbursement_id UUID REFERENCES disbursement_records(disbursement_id),
    paid_date DATE,
    paid_amount DECIMAL(12,2),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_installment CHECK (installment_number <= total_installments)
);

-- Indexes for student_bank_accounts
CREATE INDEX idx_student_bank_accounts_student_id ON student_bank_accounts(student_id);
CREATE INDEX idx_student_bank_accounts_is_primary ON student_bank_accounts(is_primary);
CREATE INDEX idx_student_bank_accounts_verified ON student_bank_accounts(verified);

-- Indexes for disbursement_records
CREATE INDEX idx_disbursement_records_allocation_id ON disbursement_records(allocation_id);
CREATE INDEX idx_disbursement_records_account_id ON disbursement_records(account_id);
CREATE INDEX idx_disbursement_records_transfer_date ON disbursement_records(transfer_date);
CREATE INDEX idx_disbursement_records_transfer_status ON disbursement_records(transfer_status);
CREATE INDEX idx_disbursement_records_transferred_by ON disbursement_records(transferred_by);

-- Indexes for disbursement_batches
CREATE INDEX idx_disbursement_batches_round_id ON disbursement_batches(round_id);
CREATE INDEX idx_disbursement_batches_scholarship_id ON disbursement_batches(scholarship_id);
CREATE INDEX idx_disbursement_batches_status ON disbursement_batches(status);
CREATE INDEX idx_disbursement_batches_scheduled_date ON disbursement_batches(scheduled_date);

-- Indexes for disbursement_batch_items
CREATE INDEX idx_disbursement_batch_items_batch_id ON disbursement_batch_items(batch_id);
CREATE INDEX idx_disbursement_batch_items_disbursement_id ON disbursement_batch_items(disbursement_id);

-- Indexes for payment_schedules
CREATE INDEX idx_payment_schedules_allocation_id ON payment_schedules(allocation_id);
CREATE INDEX idx_payment_schedules_due_date ON payment_schedules(due_date);
CREATE INDEX idx_payment_schedules_status ON payment_schedules(status);
CREATE INDEX idx_payment_schedules_disbursement_id ON payment_schedules(disbursement_id);

-- Comments
COMMENT ON TABLE student_bank_accounts IS 'บัญชีธนาคารของนักศึกษาสำหรับรับทุน';
COMMENT ON TABLE disbursement_records IS 'บันทึกการโอนเงินทุนให้นักศึกษา';
COMMENT ON TABLE disbursement_batches IS 'ชุดการโอนเงินทุนแบบกลุ่ม';
COMMENT ON TABLE disbursement_batch_items IS 'รายการโอนเงินในแต่ละชุด';
COMMENT ON TABLE payment_schedules IS 'ตารางการจ่ายทุนแบบแบ่งงวด';
