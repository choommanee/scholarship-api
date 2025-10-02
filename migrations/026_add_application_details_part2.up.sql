-- Migration 026: Application Details Tables (Part 2)

-- 6. Application Guardians (ผู้อุปการะที่ไม่ใช่บิดามารดา)
CREATE TABLE IF NOT EXISTS application_guardians (
    guardian_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ข้อมูลส่วนตัว
    title VARCHAR(50),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    -- ความสัมพันธ์
    relationship VARCHAR(100),          -- ปู่, ย่า, ลุง, ป้า, etc.

    -- ที่อยู่
    address TEXT,

    -- ข้อมูลติดต่อ
    phone VARCHAR(20),

    -- อาชีพและรายได้
    occupation VARCHAR(255),
    position VARCHAR(255),
    workplace VARCHAR(255),
    workplace_phone VARCHAR(20),
    monthly_income DECIMAL(12,2),

    -- หนี้สิน
    debts DECIMAL(15,2),
    debt_details TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 7. Application Siblings (พี่น้อง)
CREATE TABLE IF NOT EXISTS application_siblings (
    sibling_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ลำดับที่
    sibling_order INTEGER NOT NULL,     -- รวมตัวผู้สมัครด้วย (คนที่ 1, 2, 3...)

    -- ข้อมูลพื้นฐาน
    gender VARCHAR(10),                 -- male, female

    -- การศึกษา/การทำงาน
    school_or_workplace VARCHAR(255),
    education_level VARCHAR(100),       -- ระดับการศึกษา

    -- รายได้ (ถ้าทำงานแล้ว)
    is_working BOOLEAN DEFAULT FALSE,
    monthly_income DECIMAL(12,2),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Application Living Situation (สภาพการอยู่อาศัย)
CREATE TABLE IF NOT EXISTS application_living_situation (
    living_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- อยู่กับใคร
    living_with VARCHAR(50) NOT NULL,   -- parents, guardian, dorm, friends, alone, other

    -- รายละเอียด
    living_details TEXT,

    -- ภาพถ่ายบ้าน (3 รูป: ด้านหน้า, ด้านข้าง, ด้านหลัง)
    front_house_image TEXT,
    side_house_image TEXT,
    back_house_image TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_living UNIQUE (application_id)
);

-- 9. Application Financial Info (ข้อมูลการเงิน)
CREATE TABLE IF NOT EXISTS application_financial_info (
    financial_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ค่าใช้จ่าย
    monthly_allowance DECIMAL(12,2),    -- เงินที่ได้จากบ้าน/เดือน
    daily_travel_cost DECIMAL(10,2),    -- ค่าเดินทาง/วัน
    monthly_dorm_cost DECIMAL(12,2),    -- ค่าหอพัก/เดือน
    other_monthly_costs DECIMAL(12,2),  -- ค่าใช้จ่ายอื่นๆ/เดือน

    -- รายได้ของนักศึกษาเอง
    has_income BOOLEAN DEFAULT FALSE,
    income_source VARCHAR(255),         -- แหล่งรายได้
    monthly_income DECIMAL(12,2),       -- รายได้/เดือน

    -- หมายเหตุ
    financial_notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_financial UNIQUE (application_id)
);

-- 10. Application Scholarship History (ประวัติการขอทุน)
CREATE TABLE IF NOT EXISTS application_scholarship_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ทุนที่เคยได้รับ
    scholarship_name VARCHAR(255),
    scholarship_type VARCHAR(100),      -- external, internal
    amount DECIMAL(12,2),
    academic_year VARCHAR(10),          -- ปีการศึกษาที่ได้รับ

    -- ทุนกู้ยืม กยศ.
    has_student_loan BOOLEAN DEFAULT FALSE,
    loan_type VARCHAR(100),             -- กยศ., กรอ.
    loan_year VARCHAR(10),              -- ปีที่กู้
    loan_amount DECIMAL(12,2),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for application_guardians
CREATE INDEX idx_app_guardians_application_id ON application_guardians(application_id);

-- Indexes for application_siblings
CREATE INDEX idx_app_siblings_application_id ON application_siblings(application_id);
CREATE INDEX idx_app_siblings_sibling_order ON application_siblings(sibling_order);

-- Indexes for application_living_situation
CREATE INDEX idx_app_living_application_id ON application_living_situation(application_id);

-- Indexes for application_financial_info
CREATE INDEX idx_app_financial_application_id ON application_financial_info(application_id);

-- Indexes for application_scholarship_history
CREATE INDEX idx_app_scholarship_history_application_id ON application_scholarship_history(application_id);
CREATE INDEX idx_app_scholarship_history_academic_year ON application_scholarship_history(academic_year);

-- Comments
COMMENT ON TABLE application_guardians IS 'ข้อมูลผู้อุปการะที่ไม่ใช่บิดามารดา';
COMMENT ON TABLE application_siblings IS 'ข้อมูลพี่น้อง (รวมตัวผู้สมัคร)';
COMMENT ON TABLE application_living_situation IS 'สภาพการอยู่อาศัยและภาพถ่ายบ้าน';
COMMENT ON TABLE application_financial_info IS 'ข้อมูลรายรับ-รายจ่ายของนักศึกษา';
COMMENT ON TABLE application_scholarship_history IS 'ประวัติการได้รับทุนและทุนกู้ยืม';
