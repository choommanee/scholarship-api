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
    sso_provider VARCHAR(50), -- Name of SSO provider if using SSO
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

-- User roles junction table - maps users to their assigned roles
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE, -- The user
    role_id INTEGER REFERENCES roles(role_id) ON DELETE CASCADE, -- The role assigned
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the role was assigned
    assigned_by UUID REFERENCES users(user_id), -- Who assigned the role
    is_active BOOLEAN DEFAULT TRUE, -- Whether this role assignment is active
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

-- Scholarships table - defines available scholarship programs
CREATE TABLE scholarships (
    scholarship_id SERIAL PRIMARY KEY, -- Unique identifier for the scholarship
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
    is_active BOOLEAN DEFAULT TRUE, -- Whether scholarship is currently active
    created_by UUID REFERENCES users(user_id), -- User who created the scholarship
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Scholarship budgets table - tracks budget allocation and usage
CREATE TABLE scholarship_budgets (
    budget_id SERIAL PRIMARY KEY, -- Unique identifier for budget entry
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id) ON DELETE CASCADE, -- Link to scholarship
    budget_year VARCHAR(10), -- Budget year
    total_budget DECIMAL(12,2) NOT NULL, -- Total allocated budget
    allocated_budget DECIMAL(12,2) DEFAULT 0, -- Amount already allocated
    remaining_budget DECIMAL(12,2) GENERATED ALWAYS AS (total_budget - allocated_budget) STORED, -- Remaining budget
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When budget was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last update timestamp
    UNIQUE(scholarship_id, budget_year) -- One budget per scholarship per year
);

-- Scholarship applications table - tracks student applications
CREATE TABLE scholarship_applications (
    application_id SERIAL PRIMARY KEY, -- Unique identifier for application
    student_id VARCHAR(20) REFERENCES students(student_id), -- Applicant
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id), -- Scholarship applied for
    application_status VARCHAR(30) DEFAULT 'draft', -- Status of application
    application_data JSONB, -- Application form data
    family_income DECIMAL(10,2), -- Annual family income
    monthly_expenses DECIMAL(10,2), -- Monthly expenses
    siblings_count INTEGER, -- Number of siblings
    special_abilities TEXT, -- Special abilities or skills
    activities_participation TEXT, -- Extracurricular activities
    submitted_at TIMESTAMP, -- When application was submitted
    reviewed_by UUID REFERENCES users(user_id), -- Staff who reviewed
    reviewed_at TIMESTAMP, -- When application was reviewed
    review_notes TEXT, -- Notes from review
    priority_score DECIMAL(5,2), -- Application priority score
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When application was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Application documents table - stores documents submitted with applications
CREATE TABLE application_documents (
    document_id SERIAL PRIMARY KEY, -- Unique identifier for document
    application_id INTEGER REFERENCES scholarship_applications(application_id) ON DELETE CASCADE, -- Link to application
    document_type VARCHAR(50) NOT NULL, -- Type of document
    document_name VARCHAR(255) NOT NULL, -- Original filename
    file_path VARCHAR(500) NOT NULL, -- Path to stored file
    file_size INTEGER, -- Size in bytes
    mime_type VARCHAR(100), -- MIME type of file
    upload_status VARCHAR(20) DEFAULT 'pending', -- Upload status
    verification_notes TEXT, -- Notes from verification
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Upload timestamp
    verified_by UUID REFERENCES users(user_id), -- Staff who verified
    verified_at TIMESTAMP -- When document was verified
);

-- Interview schedules table - manages interview scheduling
CREATE TABLE interview_schedules (
    schedule_id SERIAL PRIMARY KEY, -- Unique identifier for schedule
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id), -- Link to scholarship
    interview_date DATE NOT NULL, -- Date of interviews
    start_time TIME NOT NULL, -- Start time
    end_time TIME NOT NULL, -- End time
    location VARCHAR(255), -- Interview location
    max_applicants INTEGER DEFAULT 1, -- Maximum applicants per slot
    interviewer_ids JSONB, -- Array of interviewer IDs
    notes TEXT, -- Additional notes
    is_active BOOLEAN DEFAULT TRUE, -- Whether schedule is active
    created_by UUID REFERENCES users(user_id), -- Staff who created schedule
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Creation timestamp
);

-- Interview appointments table - tracks individual interview appointments
CREATE TABLE interview_appointments (
    appointment_id SERIAL PRIMARY KEY, -- Unique identifier for appointment
    application_id INTEGER REFERENCES scholarship_applications(application_id), -- Link to application
    schedule_id INTEGER REFERENCES interview_schedules(schedule_id), -- Link to schedule
    appointment_status VARCHAR(20) DEFAULT 'scheduled', -- Status of appointment
    student_confirmed BOOLEAN DEFAULT FALSE, -- Whether student confirmed
    confirmation_date TIMESTAMP, -- When student confirmed
    actual_start_time TIMESTAMP, -- Actual start time
    actual_end_time TIMESTAMP, -- Actual end time
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    UNIQUE(application_id, schedule_id) -- One appointment per application per schedule
);

-- Interview results table - stores interview outcomes
CREATE TABLE interview_results (
    result_id SERIAL PRIMARY KEY, -- Unique identifier for result
    appointment_id INTEGER REFERENCES interview_appointments(appointment_id), -- Link to appointment
    interviewer_id UUID REFERENCES users(user_id), -- Interviewer who conducted
    scores JSONB, -- Detailed scoring
    overall_score DECIMAL(5,2), -- Overall interview score
    comments TEXT, -- Interview comments
    recommendation VARCHAR(20), -- Interviewer's recommendation
    interview_notes TEXT, -- Additional notes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Creation timestamp
);

-- Scholarship allocations table - tracks scholarship awards
CREATE TABLE scholarship_allocations (
    allocation_id SERIAL PRIMARY KEY, -- Unique identifier for allocation
    application_id INTEGER REFERENCES scholarship_applications(application_id), -- Link to application
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id), -- Link to scholarship
    allocated_amount DECIMAL(10,2) NOT NULL, -- Amount allocated
    allocation_status VARCHAR(20) DEFAULT 'pending', -- Status of allocation
    allocation_date DATE NOT NULL, -- Date of allocation
    disbursement_method VARCHAR(30), -- How funds are disbursed
    bank_account VARCHAR(50), -- Recipient bank account
    bank_name VARCHAR(100), -- Bank name
    transfer_date DATE, -- Date of transfer
    transfer_reference VARCHAR(100), -- Transfer reference number
    allocated_by UUID REFERENCES users(user_id), -- Staff who allocated
    approved_by UUID REFERENCES users(user_id), -- Staff who approved
    notes TEXT, -- Additional notes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Academic progress tracking table - monitors scholarship recipients' progress
CREATE TABLE academic_progress_tracking (
    tracking_id SERIAL PRIMARY KEY, -- Unique identifier for tracking entry
    student_id VARCHAR(20) REFERENCES students(student_id), -- Student being tracked
    allocation_id INTEGER REFERENCES scholarship_allocations(allocation_id), -- Link to allocation
    semester VARCHAR(10) NOT NULL, -- Semester being tracked
    gpa DECIMAL(3,2), -- GPA for semester
    credits_earned INTEGER, -- Credits earned in semester
    academic_status VARCHAR(20), -- Academic standing
    report_date DATE NOT NULL, -- Date of report
    notes TEXT, -- Additional notes
    created_by UUID REFERENCES users(user_id), -- Staff who created entry
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Creation timestamp
);

-- Notifications table - manages system notifications
CREATE TABLE notifications (
    notification_id SERIAL PRIMARY KEY, -- Unique identifier for notification
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE, -- Recipient user
    notification_type VARCHAR(50) NOT NULL, -- Type of notification
    title VARCHAR(255) NOT NULL, -- Notification title
    message TEXT NOT NULL, -- Notification message
    reference_id VARCHAR(100), -- Related entity ID
    reference_type VARCHAR(50), -- Type of related entity
    is_read BOOLEAN DEFAULT FALSE, -- Whether notification was read
    is_email_sent BOOLEAN DEFAULT FALSE, -- Whether email was sent
    email_sent_at TIMESTAMP, -- When email was sent
    priority VARCHAR(10) DEFAULT 'normal', -- Notification priority
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    read_at TIMESTAMP -- When notification was read
);

-- SSO sessions table - tracks Single Sign-On sessions
CREATE TABLE sso_sessions (
    session_id VARCHAR(255) PRIMARY KEY, -- Unique session identifier
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE, -- User in session
    sso_session_id VARCHAR(255), -- SSO provider's session ID
    provider VARCHAR(50) NOT NULL, -- SSO provider name
    access_token TEXT, -- OAuth access token
    refresh_token TEXT, -- OAuth refresh token
    token_expires_at TIMESTAMP, -- Token expiration time
    session_data JSONB, -- Additional session data
    ip_address INET, -- Client IP address
    user_agent TEXT, -- Client user agent
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Session start time
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last activity time
    is_active BOOLEAN DEFAULT TRUE -- Whether session is active
);

-- Login history table - tracks authentication attempts
CREATE TABLE login_history (
    login_id SERIAL PRIMARY KEY, -- Unique identifier for login attempt
    user_id UUID REFERENCES users(user_id), -- User attempting login
    login_method VARCHAR(20) NOT NULL, -- Authentication method used
    provider VARCHAR(50), -- SSO provider if applicable
    ip_address INET, -- Client IP address
    user_agent TEXT, -- Client user agent
    login_status VARCHAR(20) NOT NULL, -- Success/failure status
    failure_reason VARCHAR(100), -- Reason for failure if applicable
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When login occurred
    logout_time TIMESTAMP, -- When user logged out
    session_duration INTEGER -- Session duration in seconds
);

-- Import logs table - tracks data import activities
CREATE TABLE import_logs (
    import_id SERIAL PRIMARY KEY, -- Unique identifier for import
    import_type VARCHAR(50) NOT NULL, -- Type of data imported
    file_name VARCHAR(255) NOT NULL, -- Name of imported file
    total_records INTEGER, -- Total records in import
    successful_records INTEGER, -- Successfully imported records
    failed_records INTEGER, -- Failed records
    error_details JSONB, -- Details of any errors
    imported_by UUID REFERENCES users(user_id), -- Staff who performed import
    import_status VARCHAR(20) DEFAULT 'processing', -- Status of import
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Start time
    completed_at TIMESTAMP -- Completion time
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
