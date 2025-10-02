-- Schema Alignment Migration
-- Aligns existing database with new comprehensive schema design

-- 1. Add new ENUM types
CREATE TYPE scholarship_round AS ENUM ('round1_main', 'round2_freshman', 'round3_makeup');
CREATE TYPE document_type_enum AS ENUM (
    'personal_history',           -- 1. เอกสารการบรรยายประวัติ
    'house_map',                  -- 2. แผนผังแสดงตำแหน่งบ้าน
    'house_photos',               -- 3. ภาพถ่ายบ้าน 3 ภาพ
    'income_certificate_fixed',   -- 4. เอกสารรับรองรายได้ (รายได้แน่นอน)
    'income_certificate_variable',-- 5. หนังสือรับรองรายได้ (รายได้ไม่แน่นอน)
    'guarantor_id_card',          -- 6. สำเนาบัตรประจำตัวผู้รับรองรายได้
    'national_id_card',           -- 7. สำเนาบัตรประจำตัวประชาชน
    'house_registration',         -- 8. สำเนาทะเบียนบ้าน
    'student_id_card',            -- 9. สำเนาบัตรนักศึกษา
    'transcript',                 -- 10. สำเนาใบรายงานผลการเรียน
    'bank_account',               -- 11. สำเนาบัญชีเงินฝากธนาคาร
    'activity_certificate'       -- 12. หลักฐานการร่วมกิจกรรม
);

CREATE TYPE workflow_step AS ENUM (
    'document_verification',     -- ตรวจสอบเอกสาร
    'eligibility_check',        -- ตรวจสอบคุณสมบัติ
    'initial_screening',        -- คัดกรองเบื้องต้น
    'advisor_review',           -- อาจารย์ให้ความเห็น
    'committee_review',         -- คณะกรรมการพิจารณา
    'interview_scheduling',     -- จัดตารางสัมภาษณ์
    'interview_conducted',      -- ดำเนินการสัมภาษณ์
    'final_evaluation',         -- ประเมินผลสุดท้าย
    'result_approval',          -- อนุมัติผล
    'result_announcement'       -- ประกาศผล
);

CREATE TYPE step_status AS ENUM ('pending', 'in_progress', 'completed', 'rejected', 'on_hold');
CREATE TYPE advisor_recommendation AS ENUM ('strongly_recommend', 'recommend', 'neutral', 'not_recommend', 'strongly_not_recommend');
CREATE TYPE committee_recommendation AS ENUM ('approve', 'conditional_approve', 'reject', 'need_more_info');

-- 2. Add round column to scholarships table
ALTER TABLE scholarships 
ADD COLUMN IF NOT EXISTS round scholarship_round;

-- 3. Enhance scholarship_applications table with 18 data categories
ALTER TABLE scholarship_applications 
ADD COLUMN IF NOT EXISTS personal_info JSONB,
ADD COLUMN IF NOT EXISTS academic_info JSONB,
ADD COLUMN IF NOT EXISTS address_info JSONB,
ADD COLUMN IF NOT EXISTS education_history JSONB,
ADD COLUMN IF NOT EXISTS family_info JSONB,
ADD COLUMN IF NOT EXISTS assets_liabilities JSONB,
ADD COLUMN IF NOT EXISTS guardian_info JSONB,
ADD COLUMN IF NOT EXISTS siblings_info JSONB,
ADD COLUMN IF NOT EXISTS living_condition JSONB,
ADD COLUMN IF NOT EXISTS financial_info JSONB,
ADD COLUMN IF NOT EXISTS scholarship_history JSONB,
ADD COLUMN IF NOT EXISTS activities_skills JSONB,
ADD COLUMN IF NOT EXISTS reference_person JSONB,
ADD COLUMN IF NOT EXISTS special_abilities_detailed JSONB,
ADD COLUMN IF NOT EXISTS health_issues JSONB,
ADD COLUMN IF NOT EXISTS funding_needs JSONB,
ADD COLUMN IF NOT EXISTS advisor_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS terms_agreement BOOLEAN DEFAULT FALSE;

-- 4. Create application_workflow table
CREATE TABLE IF NOT EXISTS application_workflow (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id),
    step workflow_step NOT NULL,
    status step_status NOT NULL,
    assigned_to UUID REFERENCES users(user_id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Create advisor_reviews table
CREATE TABLE IF NOT EXISTS advisor_reviews (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id),
    advisor_id UUID NOT NULL REFERENCES users(user_id),
    recommendation advisor_recommendation NOT NULL,
    comments TEXT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(application_id, advisor_id)
);

-- 6. Create committee_evaluations table
CREATE TABLE IF NOT EXISTS committee_evaluations (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id),
    committee_member_id UUID NOT NULL REFERENCES users(user_id),
    evaluation_criteria JSONB NOT NULL,
    total_score DECIMAL(5,2),
    recommendation committee_recommendation,
    comments TEXT,
    evaluated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(application_id, committee_member_id)
);

-- 7. Create messages table for internal communication
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id UUID NOT NULL REFERENCES users(user_id),
    recipient_id UUID NOT NULL REFERENCES users(user_id),
    subject VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    read_at TIMESTAMP,
    replied_at TIMESTAMP,
    parent_message_id INTEGER REFERENCES messages(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Enhance notifications table with new types
ALTER TABLE notifications 
ADD COLUMN IF NOT EXISTS data JSONB,
ADD COLUMN IF NOT EXISTS sent_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS sms_sent BOOLEAN DEFAULT FALSE;

-- Update notification_type to include new types
ALTER TABLE notifications 
ALTER COLUMN notification_type TYPE VARCHAR(50);

-- 9. Create system_settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value JSONB NOT NULL,
    description TEXT,
    category VARCHAR(20) DEFAULT 'general',
    is_active BOOLEAN DEFAULT TRUE,
    updated_by UUID REFERENCES users(user_id),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    action VARCHAR(20) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(100),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    session_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_application_workflow_app_id ON application_workflow(application_id);
CREATE INDEX IF NOT EXISTS idx_application_workflow_step ON application_workflow(step, status);
CREATE INDEX IF NOT EXISTS idx_advisor_reviews_advisor ON advisor_reviews(advisor_id);
CREATE INDEX IF NOT EXISTS idx_committee_evaluations_member ON committee_evaluations(committee_member_id);
CREATE INDEX IF NOT EXISTS idx_messages_recipient ON messages(recipient_id, read_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- 12. Insert default system settings
INSERT INTO system_settings (setting_key, setting_value, description, category) VALUES 
('application_deadline_notification', '{"days_before": 3, "enabled": true}', 'แจ้งเตือนก่อนปิดรับสมัคร', 'notification'),
('max_file_upload_size', '{"size_mb": 10, "allowed_types": ["pdf", "jpg", "png"]}', 'ขนาดไฟล์สูงสุดที่อัปโหลดได้', 'general'),
('workflow_auto_advance', '{"enabled": false, "conditions": []}', 'การเลื่อนขั้นตอนอัตโนมัติ', 'workflow'),
('interview_booking_window', '{"days_before": 7, "days_after": 1}', 'ช่วงเวลาจองสัมภาษณ์', 'interview')
ON CONFLICT (setting_key) DO NOTHING;

-- 13. Update existing roles to include new permissions
UPDATE roles 
SET permissions = permissions || '["manage_workflow", "advisor_review", "committee_evaluation"]'::jsonb
WHERE role_name = 'scholarship_officer';

INSERT INTO roles (role_name, role_description, permissions) VALUES 
('advisor', 'Academic Advisor', '["view_student_applications", "provide_recommendations", "track_student_progress"]'),
('committee_member', 'Committee Member', '["evaluate_applications", "conduct_interviews", "make_decisions"]')
ON CONFLICT (role_name) DO NOTHING;

-- 14. Create triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_application_workflow_updated_at 
    BEFORE UPDATE ON application_workflow 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 15. Migrate existing application data to new structure
-- This would need to be customized based on existing data format
UPDATE scholarship_applications 
SET 
    personal_info = COALESCE(
        jsonb_build_object(
            'student_id', scholarship_applications.student_id,
            'family_income', family_income,
            'monthly_expenses', monthly_expenses,
            'siblings_count', siblings_count
        ),
        '{}'::jsonb
    ),
    activities_skills = COALESCE(
        jsonb_build_object(
            'special_abilities', special_abilities,
            'activities_participation', activities_participation
        ),
        '{}'::jsonb
    )
WHERE personal_info IS NULL;

-- 16. Create initial workflow states for existing applications
INSERT INTO application_workflow (application_id, step, status, started_at)
SELECT 
    application_id,
    CASE 
        WHEN application_status = 'draft' THEN 'document_verification'::workflow_step
        WHEN application_status = 'submitted' THEN 'eligibility_check'::workflow_step
        WHEN application_status = 'under_review' THEN 'advisor_review'::workflow_step
        WHEN application_status = 'approved' THEN 'result_announcement'::workflow_step
        ELSE 'document_verification'::workflow_step
    END,
    CASE 
        WHEN application_status IN ('draft', 'submitted') THEN 'pending'::step_status
        WHEN application_status = 'under_review' THEN 'in_progress'::step_status
        WHEN application_status = 'approved' THEN 'completed'::step_status
        ELSE 'pending'::step_status
    END,
    created_at
FROM scholarship_applications
WHERE NOT EXISTS (
    SELECT 1 FROM application_workflow 
    WHERE application_workflow.application_id = scholarship_applications.application_id
);

COMMENT ON TABLE application_workflow IS 'ติดตามสถานะการดำเนินงานใบสมัครตาม 10 ขั้นตอน';
COMMENT ON TABLE advisor_reviews IS 'ความเห็นของอาจารย์ที่ปรึกษาต่อใบสมัครทุน';
COMMENT ON TABLE committee_evaluations IS 'การประเมินของคณะกรรมการพิจารณาทุน';
COMMENT ON TABLE messages IS 'ระบบข้อความภายในระบบ';
COMMENT ON TABLE system_settings IS 'การตั้งค่าระบบแบบ dynamic';
COMMENT ON TABLE audit_logs IS 'บันทึกการใช้งานระบบเพื่อการตรวจสอบ';
