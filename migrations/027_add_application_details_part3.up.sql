-- Migration 027: Application Details Tables (Part 3)

-- 11. Application Activities (กิจกรรมและความสามารถพิเศษ)
CREATE TABLE IF NOT EXISTS application_activities (
    activity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ประเภท
    activity_type VARCHAR(50) NOT NULL, -- activity (กิจกรรม), special_ability (ความสามารถพิเศษ)

    -- รายละเอียด
    activity_name VARCHAR(255),
    description TEXT,

    -- ผลงาน/รางวัล
    achievement TEXT,
    award_level VARCHAR(50),            -- school, district, province, national, international

    -- ปีที่ได้รับ
    year VARCHAR(10),

    -- หลักฐาน/เอกสารประกอบ
    evidence_url TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. Application References (บุคคลที่ให้ข้อมูลเพิ่มเติม)
CREATE TABLE IF NOT EXISTS application_references (
    reference_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ข้อมูลส่วนตัว
    title VARCHAR(50),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    -- ความสัมพันธ์
    relationship VARCHAR(100),

    -- ที่อยู่
    address TEXT,

    -- ข้อมูลติดต่อ
    phone VARCHAR(20),
    email VARCHAR(255),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 13. Application Health Info (ข้อมูลปัญหาด้านสุขภาพ)
CREATE TABLE IF NOT EXISTS application_health_info (
    health_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- มีปัญหาสุขภาพหรือไม่
    has_health_issues BOOLEAN DEFAULT FALSE,

    -- รายละเอียดปัญหาสุขภาพ
    health_condition VARCHAR(255),      -- โรคประจำตัว
    health_details TEXT,                -- รายละเอียดเพิ่มเติม

    -- ผลกระทบต่อการเรียน
    affects_study BOOLEAN DEFAULT FALSE,
    study_impact_details TEXT,

    -- ค่ารักษาพยาบาล/เดือน
    monthly_medical_cost DECIMAL(12,2),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_health UNIQUE (application_id)
);

-- 14. Application Funding Needs (ความต้องการอุดหนุนทุน)
CREATE TABLE IF NOT EXISTS application_funding_needs (
    need_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ความต้องการทุน
    tuition_support DECIMAL(12,2),      -- ค่าเทอม
    monthly_support DECIMAL(12,2),      -- ค่าใช้จ่ายรายเดือน
    book_support DECIMAL(12,2),         -- ค่าหนังสือ/อุปกรณ์
    dorm_support DECIMAL(12,2),         -- ค่าหอพัก
    other_support DECIMAL(12,2),        -- อื่นๆ

    -- รายละเอียดเพิ่มเติม
    other_details TEXT,

    -- ยอดรวมที่ต้องการ
    total_requested DECIMAL(12,2),

    -- เหตุผลความจำเป็น
    necessity_reason TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_funding_needs UNIQUE (application_id)
);

-- 15. Application House Documents (เอกสารเกี่ยวกับบ้าน)
CREATE TABLE IF NOT EXISTS application_house_documents (
    doc_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ประเภทเอกสาร
    document_type VARCHAR(50) NOT NULL, -- map (แผนผัง), photo_front, photo_side, photo_back, income_cert

    -- URL เอกสาร
    document_url TEXT NOT NULL,
    file_name VARCHAR(255),
    file_size INTEGER,                  -- ขนาดไฟล์ (bytes)
    mime_type VARCHAR(100),

    -- คำอธิบาย
    description TEXT,

    -- สถานะการตรวจสอบ
    verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES users(user_id),
    verified_at TIMESTAMP,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 16. Application Income Certificates (หนังสือรับรองรายได้)
CREATE TABLE IF NOT EXISTS application_income_certificates (
    cert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- เจ้าของรายได้
    owner_name VARCHAR(255) NOT NULL,   -- ชื่อผู้มีรายได้
    relationship VARCHAR(50),           -- father, mother, guardian

    -- ประเภทรายได้
    income_type VARCHAR(50),            -- fixed (แน่นอน), variable (ไม่แน่นอน)

    -- รายได้
    monthly_income DECIMAL(12,2),

    -- ผู้รับรอง
    certified_by VARCHAR(255),          -- ผู้นำชุมชน, ข้าราชการ
    certifier_position VARCHAR(255),
    certifier_id_card VARCHAR(13),      -- เลขบัตรผู้รับรอง

    -- เอกสารแนบ
    certificate_url TEXT,               -- URL หนังสือรับรอง
    id_card_copy_url TEXT,              -- สำเนาบัตรผู้รับรอง

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for application_activities
CREATE INDEX idx_app_activities_application_id ON application_activities(application_id);
CREATE INDEX idx_app_activities_activity_type ON application_activities(activity_type);

-- Indexes for application_references
CREATE INDEX idx_app_references_application_id ON application_references(application_id);

-- Indexes for application_health_info
CREATE INDEX idx_app_health_application_id ON application_health_info(application_id);

-- Indexes for application_funding_needs
CREATE INDEX idx_app_funding_needs_application_id ON application_funding_needs(application_id);

-- Indexes for application_house_documents
CREATE INDEX idx_app_house_docs_application_id ON application_house_documents(application_id);
CREATE INDEX idx_app_house_docs_document_type ON application_house_documents(document_type);
CREATE INDEX idx_app_house_docs_verified ON application_house_documents(verified);

-- Indexes for application_income_certificates
CREATE INDEX idx_app_income_certs_application_id ON application_income_certificates(application_id);
CREATE INDEX idx_app_income_certs_income_type ON application_income_certificates(income_type);

-- Comments
COMMENT ON TABLE application_activities IS 'กิจกรรมที่เข้าร่วมและความสามารถพิเศษ';
COMMENT ON TABLE application_references IS 'บุคคลที่ให้ข้อมูลเพิ่มเติมได้';
COMMENT ON TABLE application_health_info IS 'ข้อมูลปัญหาด้านสุขภาพ';
COMMENT ON TABLE application_funding_needs IS 'ความต้องการอุดหนุนทุนการศึกษา';
COMMENT ON TABLE application_house_documents IS 'เอกสารเกี่ยวกับบ้านและที่อยู่อาศัย';
COMMENT ON TABLE application_income_certificates IS 'หนังสือรับรองรายได้และเอกสารประกอบ';
