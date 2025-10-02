-- Migration 023: Add Advisors and Student-Advisor Mapping

-- Advisors Table (อาจารย์ที่ปรึกษา)
CREATE TABLE IF NOT EXISTS advisors (
    advisor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    -- Personal info
    prefix_th VARCHAR(50),              -- e.g., "ดร.", "ผศ.ดร."
    prefix_en VARCHAR(50),              -- e.g., "Dr.", "Asst.Prof.Dr."
    first_name_th VARCHAR(100) NOT NULL,
    last_name_th VARCHAR(100) NOT NULL,
    first_name_en VARCHAR(100),
    last_name_en VARCHAR(100),

    -- Contact
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    office_phone VARCHAR(20),
    office_location VARCHAR(255),

    -- Academic info
    department VARCHAR(100),
    faculty VARCHAR(100),
    position VARCHAR(100),              -- e.g., "อาจารย์", "ผู้ช่วยศาสตราจารย์"
    specialization TEXT,                -- สาขาที่เชี่ยวชาญ

    -- Advisor capacity
    max_students INTEGER DEFAULT 20,   -- จำนวนนักศึกษาสูงสุดที่รับปรึกษา
    current_students INTEGER DEFAULT 0,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_available BOOLEAN DEFAULT TRUE,  -- พร้อมรับนักศึกษาใหม่

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Student-Advisor Mapping Table
CREATE TABLE IF NOT EXISTS student_advisors (
    mapping_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id VARCHAR(20) NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    advisor_id UUID NOT NULL REFERENCES advisors(advisor_id) ON DELETE CASCADE,

    -- Assignment details
    assigned_date DATE NOT NULL DEFAULT CURRENT_DATE,
    start_academic_year VARCHAR(10),    -- e.g., "2567"
    end_academic_year VARCHAR(10),

    -- Status
    status VARCHAR(50) DEFAULT 'active' CHECK (
        status IN ('active', 'inactive', 'transferred', 'graduated')
    ),
    is_primary BOOLEAN DEFAULT TRUE,    -- อาจารย์ที่ปรึกษาหลัก

    -- Metadata
    assigned_by UUID REFERENCES users(user_id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_active_primary_advisor UNIQUE (student_id, is_primary, status)
        DEFERRABLE INITIALLY DEFERRED
);

-- Advisor Assignments (for scholarship application review)
CREATE TABLE IF NOT EXISTS advisor_assignments (
    assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,
    advisor_id UUID NOT NULL REFERENCES advisors(advisor_id) ON DELETE CASCADE,

    -- Assignment details
    assigned_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(user_id),
    due_date TIMESTAMP,

    -- Status
    status VARCHAR(50) DEFAULT 'pending' CHECK (
        status IN ('pending', 'in_review', 'completed', 'declined', 'expired')
    ),

    -- Review
    reviewed_at TIMESTAMP,
    recommendation VARCHAR(50) CHECK (
        recommendation IN ('strongly_support', 'support', 'neutral', 'not_support')
    ),
    comments TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_advisor_assignment UNIQUE (application_id, advisor_id)
);

-- Indexes for advisors
CREATE INDEX idx_advisors_user_id ON advisors(user_id);
CREATE INDEX idx_advisors_department ON advisors(department);
CREATE INDEX idx_advisors_faculty ON advisors(faculty);
CREATE INDEX idx_advisors_is_active ON advisors(is_active);
CREATE INDEX idx_advisors_is_available ON advisors(is_available);

-- Indexes for student_advisors
CREATE INDEX idx_student_advisors_student_id ON student_advisors(student_id);
CREATE INDEX idx_student_advisors_advisor_id ON student_advisors(advisor_id);
CREATE INDEX idx_student_advisors_status ON student_advisors(status);
CREATE INDEX idx_student_advisors_is_primary ON student_advisors(is_primary);

-- Indexes for advisor_assignments
CREATE INDEX idx_advisor_assignments_application_id ON advisor_assignments(application_id);
CREATE INDEX idx_advisor_assignments_advisor_id ON advisor_assignments(advisor_id);
CREATE INDEX idx_advisor_assignments_status ON advisor_assignments(status);
CREATE INDEX idx_advisor_assignments_due_date ON advisor_assignments(due_date);

-- Comments
COMMENT ON TABLE advisors IS 'Academic advisors (อาจารย์ที่ปรึกษา)';
COMMENT ON TABLE student_advisors IS 'Mapping between students and their advisors';
COMMENT ON TABLE advisor_assignments IS 'Assignment of advisors to review scholarship applications';
