-- Migration 024: Add Interviewers and Interview Committees

-- Interviewers Table (กรรมการสัมภาษณ์)
CREATE TABLE IF NOT EXISTS interviewers (
    interviewer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    -- Personal info
    prefix_th VARCHAR(50),
    prefix_en VARCHAR(50),
    first_name_th VARCHAR(100) NOT NULL,
    last_name_th VARCHAR(100) NOT NULL,
    first_name_en VARCHAR(100),
    last_name_en VARCHAR(100),

    -- Contact
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),

    -- Professional info
    department VARCHAR(100),
    faculty VARCHAR(100),
    position VARCHAR(100),
    expertise TEXT[],                   -- สาขาความเชี่ยวชาญ (array)

    -- Interviewer capacity
    max_interviews_per_day INTEGER DEFAULT 10,
    max_interviews_per_round INTEGER DEFAULT 50,
    current_interviews INTEGER DEFAULT 0,

    -- Availability
    is_active BOOLEAN DEFAULT TRUE,
    is_available BOOLEAN DEFAULT TRUE,
    available_days TEXT[],              -- e.g., ['monday', 'wednesday', 'friday']
    available_times JSONB,              -- เวลาว่างแต่ละวัน

    -- Rating
    average_rating DECIMAL(3,2),        -- คะแนนเฉลี่ยจากผู้สมัคร (0-5)
    total_interviews INTEGER DEFAULT 0,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Interview Committees Table
CREATE TABLE IF NOT EXISTS interview_committees (
    committee_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    round_id INTEGER REFERENCES scholarship_rounds(round_id) ON DELETE CASCADE,
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id) ON DELETE CASCADE,

    -- Committee details
    committee_name VARCHAR(255) NOT NULL,
    committee_code VARCHAR(50) UNIQUE,
    description TEXT,

    -- Chair
    chair_id UUID REFERENCES interviewers(interviewer_id),

    -- Committee size
    min_members INTEGER DEFAULT 3,
    max_members INTEGER DEFAULT 5,
    current_members INTEGER DEFAULT 0,

    -- Interview slots
    total_slots INTEGER DEFAULT 0,
    used_slots INTEGER DEFAULT 0,

    -- Status
    status VARCHAR(50) DEFAULT 'forming' CHECK (
        status IN ('forming', 'active', 'completed', 'disbanded')
    ),

    -- Metadata
    created_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Committee Members Table
CREATE TABLE IF NOT EXISTS committee_members (
    member_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    committee_id UUID NOT NULL REFERENCES interview_committees(committee_id) ON DELETE CASCADE,
    interviewer_id UUID NOT NULL REFERENCES interviewers(interviewer_id) ON DELETE CASCADE,

    -- Member role
    role VARCHAR(50) DEFAULT 'member' CHECK (
        role IN ('chair', 'co_chair', 'member', 'secretary')
    ),

    -- Assignment
    assigned_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(user_id),

    -- Status
    status VARCHAR(50) DEFAULT 'active' CHECK (
        status IN ('active', 'inactive', 'resigned')
    ),

    -- Performance
    interviews_conducted INTEGER DEFAULT 0,
    average_duration_minutes INTEGER,

    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_committee_member UNIQUE (committee_id, interviewer_id)
);

-- Committee Assignments (assign committee to interview sessions)
CREATE TABLE IF NOT EXISTS committee_interview_assignments (
    assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    committee_id UUID NOT NULL REFERENCES interview_committees(committee_id) ON DELETE CASCADE,
    interview_id INTEGER NOT NULL REFERENCES interview_schedules(schedule_id) ON DELETE CASCADE,

    -- Assignment details
    assigned_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(user_id),

    -- Status
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (
        status IN ('scheduled', 'in_progress', 'completed', 'cancelled')
    ),

    -- Results
    completed_at TIMESTAMP,
    duration_minutes INTEGER,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_committee_interview UNIQUE (committee_id, interview_id)
);

-- Interviewer Availability Slots
CREATE TABLE IF NOT EXISTS interviewer_availability (
    availability_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interviewer_id UUID NOT NULL REFERENCES interviewers(interviewer_id) ON DELETE CASCADE,
    round_id INTEGER REFERENCES scholarship_rounds(round_id) ON DELETE CASCADE,

    -- Time slot
    available_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    -- Capacity
    max_interviews INTEGER DEFAULT 1,
    booked_interviews INTEGER DEFAULT 0,

    -- Status
    is_available BOOLEAN DEFAULT TRUE,
    is_blocked BOOLEAN DEFAULT FALSE,   -- blocked by interviewer
    block_reason TEXT,

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_time_slot CHECK (end_time > start_time)
);

-- Indexes for interviewers
CREATE INDEX idx_interviewers_user_id ON interviewers(user_id);
CREATE INDEX idx_interviewers_department ON interviewers(department);
CREATE INDEX idx_interviewers_faculty ON interviewers(faculty);
CREATE INDEX idx_interviewers_is_active ON interviewers(is_active);
CREATE INDEX idx_interviewers_is_available ON interviewers(is_available);

-- Indexes for interview_committees
CREATE INDEX idx_interview_committees_round_id ON interview_committees(round_id);
CREATE INDEX idx_interview_committees_scholarship_id ON interview_committees(scholarship_id);
CREATE INDEX idx_interview_committees_chair_id ON interview_committees(chair_id);
CREATE INDEX idx_interview_committees_status ON interview_committees(status);

-- Indexes for committee_members
CREATE INDEX idx_committee_members_committee_id ON committee_members(committee_id);
CREATE INDEX idx_committee_members_interviewer_id ON committee_members(interviewer_id);
CREATE INDEX idx_committee_members_role ON committee_members(role);
CREATE INDEX idx_committee_members_status ON committee_members(status);

-- Indexes for committee_interview_assignments
CREATE INDEX idx_committee_assignments_committee_id ON committee_interview_assignments(committee_id);
CREATE INDEX idx_committee_assignments_interview_id ON committee_interview_assignments(interview_id);
CREATE INDEX idx_committee_assignments_status ON committee_interview_assignments(status);

-- Indexes for interviewer_availability
CREATE INDEX idx_interviewer_availability_interviewer_id ON interviewer_availability(interviewer_id);
CREATE INDEX idx_interviewer_availability_round_id ON interviewer_availability(round_id);
CREATE INDEX idx_interviewer_availability_date ON interviewer_availability(available_date);
CREATE INDEX idx_interviewer_availability_is_available ON interviewer_availability(is_available);

-- Comments
COMMENT ON TABLE interviewers IS 'Interview committee members (กรรมการสัมภาษณ์)';
COMMENT ON TABLE interview_committees IS 'Interview committees for scholarship rounds';
COMMENT ON TABLE committee_members IS 'Members of interview committees';
COMMENT ON TABLE committee_interview_assignments IS 'Assignment of committees to interview sessions';
COMMENT ON TABLE interviewer_availability IS 'Availability slots for interviewers';
