-- Migration: Enhance Existing Tables to Match Spec
-- Created: 2025-10-01
-- Description: Adds missing fields to existing tables

-- 1. Enhance users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS line_id VARCHAR(50),
ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS account_locked_until TIMESTAMP,
ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT false;

-- 2. Enhance students table (align with spec)
ALTER TABLE students
ADD COLUMN IF NOT EXISTS current_gpa DECIMAL(3,2),
ADD COLUMN IF NOT EXISTS advisor_id UUID REFERENCES users(user_id),
ADD COLUMN IF NOT EXISTS scholarship_history JSONB,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- 3. Enhance scholarships table
ALTER TABLE scholarships
ADD COLUMN IF NOT EXISTS priority_score INTEGER DEFAULT 5,
ADD COLUMN IF NOT EXISTS auto_approval_threshold DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS max_applications_per_user INTEGER DEFAULT 1;

-- Rename columns to match spec (only if they exist with old names)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name='scholarships' AND column_name='scholarship_name') THEN
        ALTER TABLE scholarships RENAME COLUMN scholarship_name TO name;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name='scholarships' AND column_name='scholarship_type') THEN
        ALTER TABLE scholarships RENAME COLUMN scholarship_type TO type;
    END IF;
END $$;

-- 4. Enhance scholarship_applications table
ALTER TABLE scholarship_applications
ADD COLUMN IF NOT EXISTS gpa_verification_status VARCHAR(30) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS reference_check_status VARCHAR(30) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS automated_score DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS manual_override BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS risk_assessment VARCHAR(20);

-- 5. Enhance application_documents table
ALTER TABLE application_documents
ADD COLUMN IF NOT EXISTS original_filename VARCHAR(255),
ADD COLUMN IF NOT EXISTS verification_status VARCHAR(30) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS verified_by UUID REFERENCES users(user_id);

-- 6. Enhance interview_schedules table
ALTER TABLE interview_schedules
ADD COLUMN IF NOT EXISTS interview_time TIME,
ADD COLUMN IF NOT EXISTS interviewer_ids UUID[],
ADD COLUMN IF NOT EXISTS status VARCHAR(30) DEFAULT 'scheduled',
ADD COLUMN IF NOT EXISTS meeting_type VARCHAR(30) DEFAULT 'in_person',
ADD COLUMN IF NOT EXISTS meeting_link VARCHAR(500),
ADD COLUMN IF NOT EXISTS duration_minutes INTEGER DEFAULT 30,
ADD COLUMN IF NOT EXISTS preparation_notes TEXT;

-- 7. Enhance interview_appointments table
ALTER TABLE interview_appointments
ADD COLUMN IF NOT EXISTS time_slot TIME,
ADD COLUMN IF NOT EXISTS status VARCHAR(30) DEFAULT 'scheduled',
ADD COLUMN IF NOT EXISTS confirmation_sent BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS reminder_sent BOOLEAN DEFAULT false;

-- 8. Enhance interview_results table
ALTER TABLE interview_results
ADD COLUMN IF NOT EXISTS individual_scores JSONB,
ADD COLUMN IF NOT EXISTS interview_duration INTEGER,
ADD COLUMN IF NOT EXISTS technical_issues TEXT,
ADD COLUMN IF NOT EXISTS follow_up_required BOOLEAN DEFAULT false;

-- 9. Enhance scholarship_allocations table
ALTER TABLE scholarship_allocations
ADD COLUMN IF NOT EXISTS academic_year VARCHAR(10),
ADD COLUMN IF NOT EXISTS semester VARCHAR(10),
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(30) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS payment_date DATE,
ADD COLUMN IF NOT EXISTS bank_account_verified BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50) DEFAULT 'bank_transfer',
ADD COLUMN IF NOT EXISTS payment_frequency VARCHAR(30) DEFAULT 'lump_sum',
ADD COLUMN IF NOT EXISTS installments INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS conditions JSONB;

-- 10. Enhance academic_progress_tracking table
ALTER TABLE academic_progress_tracking
ADD COLUMN IF NOT EXISTS academic_year VARCHAR(10),
ADD COLUMN IF NOT EXISTS credits_completed INTEGER,
ADD COLUMN IF NOT EXISTS total_credits INTEGER,
ADD COLUMN IF NOT EXISTS status VARCHAR(30) DEFAULT 'in_progress',
ADD COLUMN IF NOT EXISTS compliance_status VARCHAR(30) DEFAULT 'compliant',
ADD COLUMN IF NOT EXISTS warning_issued BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS probation_status VARCHAR(30);

-- 11. Enhance notifications table
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS related_entity_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS related_entity_id UUID,
ADD COLUMN IF NOT EXISTS action_url VARCHAR(500),
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS delivery_method VARCHAR(30) DEFAULT 'in_app',
ADD COLUMN IF NOT EXISTS sent_via_email BOOLEAN DEFAULT false;

-- 12. Enhance messages table
ALTER TABLE messages
ADD COLUMN IF NOT EXISTS is_read BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS thread_id UUID,
ADD COLUMN IF NOT EXISTS message_type VARCHAR(30) DEFAULT 'user_message',
ADD COLUMN IF NOT EXISTS attachments JSONB;

-- 13. Enhance sso_sessions table
ALTER TABLE sso_sessions
ADD COLUMN IF NOT EXISTS device_info JSONB,
ADD COLUMN IF NOT EXISTS location_info JSONB,
ADD COLUMN IF NOT EXISTS security_level VARCHAR(20) DEFAULT 'standard';

-- 14. Enhance login_history table
ALTER TABLE login_history
ADD COLUMN IF NOT EXISTS device_fingerprint VARCHAR(255),
ADD COLUMN IF NOT EXISTS risk_score INTEGER;

-- 15. Enhance import_logs table
ALTER TABLE import_logs
ADD COLUMN IF NOT EXISTS file_size BIGINT,
ADD COLUMN IF NOT EXISTS processing_time INTEGER;

-- Create missing indexes
CREATE INDEX IF NOT EXISTS idx_users_line_id ON users(line_id);
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified);
CREATE INDEX IF NOT EXISTS idx_students_advisor ON students(advisor_id);
CREATE INDEX IF NOT EXISTS idx_notifications_entity ON notifications(related_entity_type, related_entity_id);
CREATE INDEX IF NOT EXISTS idx_messages_thread ON messages(thread_id);

COMMENT ON COLUMN users.line_id IS 'Line ID สำหรับการแจ้งเตือน';
COMMENT ON COLUMN users.failed_login_attempts IS 'จำนวนครั้งที่เข้าสู่ระบบไม่สำเร็จ';
COMMENT ON COLUMN users.account_locked_until IS 'ล็อคบัญชีจนถึงเวลานี้';
COMMENT ON COLUMN users.email_verified IS 'ยืนยันอีเมลแล้ว';
