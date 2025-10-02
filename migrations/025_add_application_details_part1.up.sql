-- Migration 025: Application Details Tables (Part 1)

-- 1. Application Personal Info (ข้อมูลส่วนตัว)
CREATE TABLE IF NOT EXISTS application_personal_info (
    info_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ชื่อ-นามสกุล
    prefix_th VARCHAR(50),              -- คำนำหน้าไทย
    prefix_en VARCHAR(50),              -- คำนำหน้าอังกฤษ
    first_name_th VARCHAR(100) NOT NULL,
    last_name_th VARCHAR(100) NOT NULL,
    first_name_en VARCHAR(100),
    last_name_en VARCHAR(100),

    -- ข้อมูลติดต่อ
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    line_id VARCHAR(100),

    -- ข้อมูลประจำตัว
    citizen_id VARCHAR(13),             -- เลขบัตรประชาชน
    student_id VARCHAR(20),

    -- ข้อมูลการศึกษา
    faculty VARCHAR(100),
    department VARCHAR(100),
    major VARCHAR(100),
    year_level INTEGER,                 -- ชั้นปี

    -- ประเภทการรับเข้า
    admission_type VARCHAR(50),         -- Portfolio, Quota, Admission
    admission_details TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_application_personal_info UNIQUE (application_id)
);

-- 2. Application Addresses (ที่อยู่)
CREATE TABLE IF NOT EXISTS application_addresses (
    address_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ประเภทที่อยู่
    address_type VARCHAR(50) NOT NULL,  -- registered (ตามภูมิลำเนา), current (ปัจจุบัน)

    -- ที่อยู่
    house_number VARCHAR(50),
    village_number VARCHAR(50),
    alley VARCHAR(100),
    road VARCHAR(100),
    subdistrict VARCHAR(100),
    district VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),

    -- ที่อยู่แบบเต็ม
    address_line1 TEXT,
    address_line2 TEXT,

    -- GPS Location
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),

    -- แผนที่ตั้งบ้าน
    map_image_url TEXT,                 -- Google Map screenshot

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Application Education History (ประวัติการศึกษา)
CREATE TABLE IF NOT EXISTS application_education_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ระดับการศึกษา
    education_level VARCHAR(50) NOT NULL, -- มัธยมต้น, มัธยมปลาย

    -- โรงเรียน
    school_name VARCHAR(255) NOT NULL,
    school_province VARCHAR(100),

    -- ผลการเรียน
    gpa DECIMAL(3,2),

    -- ปีที่จบ
    graduation_year VARCHAR(10),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Application Family Members (ข้อมูลครอบครัว)
CREATE TABLE IF NOT EXISTS application_family_members (
    member_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ความสัมพันธ์
    relationship VARCHAR(50) NOT NULL,  -- father, mother, guardian

    -- ข้อมูลส่วนตัว
    title VARCHAR(50),                  -- นาย, นาง, นางสาว
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age INTEGER,

    -- สถานะ
    living_status VARCHAR(50),          -- alive, deceased

    -- อาชีพและรายได้
    occupation VARCHAR(255),
    position VARCHAR(255),
    workplace VARCHAR(255),
    workplace_province VARCHAR(100),
    monthly_income DECIMAL(12,2),

    -- ข้อมูลติดต่อ
    phone VARCHAR(20),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Application Assets (ทรัพย์สินและหนี้สิน)
CREATE TABLE IF NOT EXISTS application_assets (
    asset_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- ประเภททรัพย์สิน/หนี้สิน
    asset_type VARCHAR(50) NOT NULL,    -- own_house, own_land, rent_house, rent_land, debt
    category VARCHAR(50),                -- asset, liability

    -- รายละเอียด
    description TEXT,

    -- มูลค่า
    value DECIMAL(15,2),

    -- ค่าใช้จ่ายรายเดือน (สำหรับค่าเช่า)
    monthly_cost DECIMAL(12,2),

    -- หมายเหตุ
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for application_personal_info
CREATE INDEX idx_app_personal_info_application_id ON application_personal_info(application_id);
CREATE INDEX idx_app_personal_info_citizen_id ON application_personal_info(citizen_id);
CREATE INDEX idx_app_personal_info_student_id ON application_personal_info(student_id);

-- Indexes for application_addresses
CREATE INDEX idx_app_addresses_application_id ON application_addresses(application_id);
CREATE INDEX idx_app_addresses_address_type ON application_addresses(address_type);
CREATE INDEX idx_app_addresses_province ON application_addresses(province);

-- Indexes for application_education_history
CREATE INDEX idx_app_edu_history_application_id ON application_education_history(application_id);
CREATE INDEX idx_app_edu_history_school_name ON application_education_history(school_name);

-- Indexes for application_family_members
CREATE INDEX idx_app_family_application_id ON application_family_members(application_id);
CREATE INDEX idx_app_family_relationship ON application_family_members(relationship);

-- Indexes for application_assets
CREATE INDEX idx_app_assets_application_id ON application_assets(application_id);
CREATE INDEX idx_app_assets_asset_type ON application_assets(asset_type);
CREATE INDEX idx_app_assets_category ON application_assets(category);

-- Comments
COMMENT ON TABLE application_personal_info IS 'ข้อมูลส่วนตัวของผู้สมัครทุน';
COMMENT ON TABLE application_addresses IS 'ที่อยู่ของผู้สมัคร (ตามทะเบียนบ้านและปัจจุบัน)';
COMMENT ON TABLE application_education_history IS 'ประวัติการศึกษาโดยย่อของผู้สมัคร';
COMMENT ON TABLE application_family_members IS 'ข้อมูลบิดา มารดา และผู้ปกครอง';
COMMENT ON TABLE application_assets IS 'ทรัพย์สินและหนี้สินของครอบครัว';
