-- Scholarship Management System Database Schema
-- Created for PostgreSQL

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table - stores all user information including students, staff, and administrators
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- Unique identifier for each user
    username VARCHAR(100) UNIQUE NOT NULL, -- Unique username for login
    email VARCHAR(255) UNIQUE NOT NULL, -- User's email address, must be unique
    password_hash VARCHAR(255), -- Hashed password for local authentication
    first_name VARCHAR(100) NOT NULL, -- User's first name
    last_name VARCHAR(100) NOT NULL, -- User's last name
    phone VARCHAR(20), -- Contact phone number
    is_active BOOLEAN DEFAULT TRUE, -- Whether the user account is active
    sso_provider VARCHAR(50), -- Name of SSO provider if using SSO (e.g., 'google', 'microsoft')
    sso_user_id VARCHAR(100), -- User ID from SSO provider
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the user account was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last time user details were updated
    last_login TIMESTAMP -- Last time user logged in
);

-- Roles table - defines system roles and their permissions
CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY, -- Auto-incrementing unique identifier for each role
    role_name VARCHAR(50) UNIQUE NOT NULL, -- Name of the role (e.g., 'admin', 'student')
    role_description TEXT, -- Detailed description of the role's responsibilities
    permissions JSONB, -- JSON array of permission strings for the role
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When the role was created
);

-- User roles junction table
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(role_id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(user_id),
    is_active BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (user_id, role_id)
);

-- Students table - stores additional information specific to students
CREATE TABLE students (
    student_id VARCHAR(20) PRIMARY KEY, -- University student ID number
    user_id UUID UNIQUE REFERENCES users(user_id) ON DELETE CASCADE, -- Link to user account
    faculty_code VARCHAR(20), -- Faculty/College code (e.g., 'ENG' for Engineering)
    department_code VARCHAR(20), -- Department code (e.g., 'CPE' for Computer Engineering)
    year_level INTEGER, -- Current year of study
    gpa DECIMAL(3,2), -- Current Grade Point Average
    admission_year INTEGER, -- Year the student was admitted
    graduation_year INTEGER, -- Expected graduation year
    student_status VARCHAR(20) DEFAULT 'active' -- Current status (active/graduated/suspended)
);

-- Scholarship sources table - stores information about scholarship funding sources
CREATE TABLE scholarship_sources (
    source_id SERIAL PRIMARY KEY, -- Unique identifier for each source
    source_name VARCHAR(255) NOT NULL, -- Name of the funding source
    source_type VARCHAR(50) NOT NULL, -- Type of source (internal/external/government)
    contact_person VARCHAR(100), -- Primary contact person for the source
    contact_email VARCHAR(255), -- Contact email address
    contact_phone VARCHAR(20), -- Contact phone number
    description TEXT, -- Detailed description of the funding source
    is_active BOOLEAN DEFAULT TRUE, -- Whether the source is currently active
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the source was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Scholarships table - stores information about individual scholarship programs
CREATE TABLE scholarships (
    scholarship_id SERIAL PRIMARY KEY, -- Unique identifier for each scholarship
    source_id INTEGER REFERENCES scholarship_sources(source_id), -- Link to funding source
    scholarship_name VARCHAR(255) NOT NULL, -- Name of the scholarship
    scholarship_type VARCHAR(50) NOT NULL, -- Type (merit/need-based/research)
    amount DECIMAL(10,2) NOT NULL, -- Amount per recipient
    total_quota INTEGER NOT NULL, -- Total number of scholarships available
    available_quota INTEGER NOT NULL, -- Remaining number of scholarships
    academic_year VARCHAR(10) NOT NULL, -- Academic year for the scholarship
    semester VARCHAR(10), -- Semester (if applicable)
    eligibility_criteria JSONB, -- JSON of eligibility requirements
    required_documents JSONB, -- JSON of required document types
    application_start_date DATE NOT NULL, -- When applications open
    application_end_date DATE NOT NULL, -- Application deadline
    interview_required BOOLEAN DEFAULT FALSE, -- Whether interviews are required
    is_active BOOLEAN DEFAULT TRUE, -- Whether the scholarship is currently active
    created_by UUID REFERENCES users(user_id), -- User who created the scholarship
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Scholarship budgets table
CREATE TABLE scholarship_budgets (
    budget_id SERIAL PRIMARY KEY,
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id) ON DELETE CASCADE,
    budget_year VARCHAR(10),
    total_budget DECIMAL(12,2) NOT NULL,
    allocated_budget DECIMAL(12,2) DEFAULT 0,
    remaining_budget DECIMAL(12,2) GENERATED ALWAYS AS (total_budget - allocated_budget) STORED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scholarship_id, budget_year)
);

-- Scholarship applications table
CREATE TABLE scholarship_applications (
    application_id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) REFERENCES students(student_id),
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id),
    application_status VARCHAR(30) DEFAULT 'draft',
    application_data JSONB,
    family_income DECIMAL(10,2),
    monthly_expenses DECIMAL(10,2),
    siblings_count INTEGER,
    special_abilities TEXT,
    activities_participation TEXT,
    submitted_at TIMESTAMP,
    reviewed_by UUID REFERENCES users(user_id),
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    priority_score DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Application documents table
CREATE TABLE application_documents (
    document_id SERIAL PRIMARY KEY,
    application_id INTEGER REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    document_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size INTEGER,
    mime_type VARCHAR(100),
    upload_status VARCHAR(20) DEFAULT 'pending',
    verification_notes TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_by UUID REFERENCES users(user_id),
    verified_at TIMESTAMP
);

-- Interview schedules table
CREATE TABLE interview_schedules (
    schedule_id SERIAL PRIMARY KEY,
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id),
    interview_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    location VARCHAR(255),
    max_applicants INTEGER DEFAULT 1,
    interviewer_ids JSONB,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Interview appointments table
CREATE TABLE interview_appointments (
    appointment_id SERIAL PRIMARY KEY,
    application_id INTEGER REFERENCES scholarship_applications(application_id),
    schedule_id INTEGER REFERENCES interview_schedules(schedule_id),
    appointment_status VARCHAR(20) DEFAULT 'scheduled',
    student_confirmed BOOLEAN DEFAULT FALSE,
    confirmation_date TIMESTAMP,
    actual_start_time TIMESTAMP,
    actual_end_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(application_id, schedule_id)
);

-- Interview results table
CREATE TABLE interview_results (
    result_id SERIAL PRIMARY KEY,
    appointment_id INTEGER REFERENCES interview_appointments(appointment_id),
    interviewer_id UUID REFERENCES users(user_id),
    scores JSONB,
    overall_score DECIMAL(5,2),
    comments TEXT,
    recommendation VARCHAR(20),
    interview_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scholarship allocations table
CREATE TABLE scholarship_allocations (
    allocation_id SERIAL PRIMARY KEY,
    application_id INTEGER REFERENCES scholarship_applications(application_id),
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id),
    allocated_amount DECIMAL(10,2) NOT NULL,
    allocation_status VARCHAR(20) DEFAULT 'pending',
    allocation_date DATE NOT NULL,
    disbursement_method VARCHAR(30),
    bank_account VARCHAR(50),
    bank_name VARCHAR(100),
    transfer_date DATE,
    transfer_reference VARCHAR(100),
    allocated_by UUID REFERENCES users(user_id),
    approved_by UUID REFERENCES users(user_id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Academic progress tracking table
CREATE TABLE academic_progress_tracking (
    tracking_id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) REFERENCES students(student_id),
    allocation_id INTEGER REFERENCES scholarship_allocations(allocation_id),
    semester VARCHAR(10) NOT NULL,
    gpa DECIMAL(3,2),
    credits_earned INTEGER,
    academic_status VARCHAR(20),
    report_date DATE NOT NULL,
    notes TEXT,
    created_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notifications table
CREATE TABLE notifications (
    notification_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    reference_id VARCHAR(100),
    reference_type VARCHAR(50),
    is_read BOOLEAN DEFAULT FALSE,
    is_email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP,
    priority VARCHAR(10) DEFAULT 'normal',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- SSO sessions table
CREATE TABLE sso_sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    sso_session_id VARCHAR(255),
    provider VARCHAR(50) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    session_data JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Login history table
CREATE TABLE login_history (
    login_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    login_method VARCHAR(20) NOT NULL,
    provider VARCHAR(50),
    ip_address INET,
    user_agent TEXT,
    login_status VARCHAR(20) NOT NULL,
    failure_reason VARCHAR(100),
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    logout_time TIMESTAMP,
    session_duration INTEGER
);

-- Import logs table
CREATE TABLE import_logs (
    import_id SERIAL PRIMARY KEY,
    import_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    total_records INTEGER,
    successful_records INTEGER,
    failed_records INTEGER,
    error_details JSONB,
    imported_by UUID REFERENCES users(user_id),
    import_status VARCHAR(20) DEFAULT 'processing',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_students_faculty ON students(faculty_code);
CREATE INDEX idx_applications_status ON scholarship_applications(application_status);
CREATE INDEX idx_applications_student ON scholarship_applications(student_id);
CREATE INDEX idx_scholarships_active ON scholarships(is_active, application_start_date, application_end_date);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_created ON notifications(created_at DESC);

-- Insert default roles
INSERT INTO roles (role_name, role_description, permissions) VALUES 
('admin', 'System Administrator', '["manage_users", "manage_scholarships", "manage_budget", "view_all_reports", "system_config"]'),
('scholarship_officer', 'Scholarship Officer', '["manage_applications", "review_documents", "schedule_interviews", "allocate_funds", "generate_reports"]'),
('interviewer', 'Interviewer', '["view_applications", "conduct_interviews", "submit_scores"]'),
('student', 'Student', '["apply_scholarship", "view_own_applications", "upload_documents", "schedule_interview"]');