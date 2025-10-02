-- Migration 021: Add Scholarship Rounds and Academic Years
-- For managing 3 rounds of scholarship applications per academic year

-- Academic Years Table
CREATE TABLE IF NOT EXISTS academic_years (
    year_id SERIAL PRIMARY KEY,
    year_code VARCHAR(10) NOT NULL UNIQUE, -- e.g., "2567", "2568"
    year_name_th VARCHAR(100) NOT NULL,    -- e.g., "ปีการศึกษา 2567"
    year_name_en VARCHAR(100),              -- e.g., "Academic Year 2024"
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_year_dates CHECK (end_date > start_date)
);

-- Scholarship Rounds Table
CREATE TABLE IF NOT EXISTS scholarship_rounds (
    round_id SERIAL PRIMARY KEY,
    year_id INTEGER NOT NULL REFERENCES academic_years(year_id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL CHECK (round_number IN (1, 2, 3)),
    round_name_th VARCHAR(100) NOT NULL,    -- e.g., "ทุนการศึกษา รอบที่ 1"
    round_name_en VARCHAR(100),              -- e.g., "Scholarship Round 1"
    semester VARCHAR(20),                    -- e.g., "1", "2", "ภาคฤดูร้อน"

    -- Application period
    application_start_date DATE NOT NULL,
    application_end_date DATE NOT NULL,

    -- Review period
    review_start_date DATE,
    review_end_date DATE,

    -- Interview period
    interview_start_date DATE,
    interview_end_date DATE,

    -- Announcement date
    announcement_date DATE,

    -- Disbursement period
    disbursement_start_date DATE,
    disbursement_end_date DATE,

    -- Budget allocation
    total_budget DECIMAL(15,2),
    allocated_budget DECIMAL(15,2) DEFAULT 0,
    remaining_budget DECIMAL(15,2),

    -- Statistics
    total_quota INTEGER,
    applications_count INTEGER DEFAULT 0,
    approved_count INTEGER DEFAULT 0,

    -- Status
    status VARCHAR(50) DEFAULT 'planning' CHECK (
        status IN ('planning', 'open', 'reviewing', 'interviewing', 'completed', 'cancelled')
    ),
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    description TEXT,
    notes TEXT,
    created_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT valid_application_dates CHECK (application_end_date >= application_start_date),
    CONSTRAINT valid_review_dates CHECK (review_end_date IS NULL OR review_end_date >= review_start_date),
    CONSTRAINT valid_interview_dates CHECK (interview_end_date IS NULL OR interview_end_date >= interview_start_date),
    CONSTRAINT unique_round_per_year UNIQUE (year_id, round_number)
);

-- Link scholarships to rounds
ALTER TABLE scholarships
ADD COLUMN IF NOT EXISTS round_id INTEGER REFERENCES scholarship_rounds(round_id) ON DELETE SET NULL;

-- Indexes for academic_years
CREATE INDEX idx_academic_years_year_code ON academic_years(year_code);
CREATE INDEX idx_academic_years_is_current ON academic_years(is_current);
CREATE INDEX idx_academic_years_is_active ON academic_years(is_active);
CREATE INDEX idx_academic_years_dates ON academic_years(start_date, end_date);

-- Indexes for scholarship_rounds
CREATE INDEX idx_scholarship_rounds_year_id ON scholarship_rounds(year_id);
CREATE INDEX idx_scholarship_rounds_round_number ON scholarship_rounds(round_number);
CREATE INDEX idx_scholarship_rounds_status ON scholarship_rounds(status);
CREATE INDEX idx_scholarship_rounds_is_active ON scholarship_rounds(is_active);
CREATE INDEX idx_scholarship_rounds_application_dates ON scholarship_rounds(application_start_date, application_end_date);

-- Indexes for scholarships
CREATE INDEX idx_scholarships_round_id ON scholarships(round_id);

-- Comments
COMMENT ON TABLE academic_years IS 'Academic years for scholarship management';
COMMENT ON TABLE scholarship_rounds IS 'Scholarship rounds (3 rounds per academic year)';
COMMENT ON COLUMN scholarship_rounds.round_number IS 'Round number: 1, 2, or 3';
COMMENT ON COLUMN scholarship_rounds.status IS 'Round status: planning, open, reviewing, interviewing, completed, cancelled';

-- Insert sample data for current academic year
INSERT INTO academic_years (year_code, year_name_th, year_name_en, start_date, end_date, is_current, is_active)
VALUES
    ('2567', 'ปีการศึกษา 2567', 'Academic Year 2024', '2024-06-01', '2025-05-31', TRUE, TRUE)
ON CONFLICT (year_code) DO NOTHING;

-- Insert 3 rounds for academic year 2567
INSERT INTO scholarship_rounds (
    year_id,
    round_number,
    round_name_th,
    round_name_en,
    semester,
    application_start_date,
    application_end_date,
    status
)
SELECT
    y.year_id,
    r.round_num,
    'ทุนการศึกษา รอบที่ ' || r.round_num,
    'Scholarship Round ' || r.round_num,
    CASE
        WHEN r.round_num = 1 THEN '1'
        WHEN r.round_num = 2 THEN '2'
        ELSE 'ภาคฤดูร้อน'
    END,
    r.start_date,
    r.end_date,
    'planning'
FROM academic_years y
CROSS JOIN (
    VALUES
        (1, '2024-06-01'::DATE, '2024-07-31'::DATE),
        (2, '2024-11-01'::DATE, '2024-12-31'::DATE),
        (3, '2025-03-01'::DATE, '2025-04-30'::DATE)
) AS r(round_num, start_date, end_date)
WHERE y.year_code = '2567'
ON CONFLICT (year_id, round_number) DO NOTHING;
