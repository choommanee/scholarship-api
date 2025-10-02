-- Rollback Schema Alignment Migration

-- 1. Drop triggers
DROP TRIGGER IF EXISTS update_application_workflow_updated_at ON application_workflow;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 2. Drop new tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS committee_evaluations;
DROP TABLE IF EXISTS advisor_reviews;
DROP TABLE IF EXISTS application_workflow;

-- 3. Remove new columns from scholarship_applications
ALTER TABLE scholarship_applications 
DROP COLUMN IF EXISTS terms_agreement,
DROP COLUMN IF EXISTS advisor_name,
DROP COLUMN IF EXISTS funding_needs,
DROP COLUMN IF EXISTS health_issues,
DROP COLUMN IF EXISTS special_abilities_detailed,
DROP COLUMN IF EXISTS reference_person,
DROP COLUMN IF EXISTS activities_skills,
DROP COLUMN IF EXISTS scholarship_history,
DROP COLUMN IF EXISTS financial_info,
DROP COLUMN IF EXISTS living_condition,
DROP COLUMN IF EXISTS siblings_info,
DROP COLUMN IF EXISTS guardian_info,
DROP COLUMN IF EXISTS assets_liabilities,
DROP COLUMN IF EXISTS family_info,
DROP COLUMN IF EXISTS education_history,
DROP COLUMN IF EXISTS address_info,
DROP COLUMN IF EXISTS academic_info,
DROP COLUMN IF EXISTS personal_info;

-- 4. Remove new columns from scholarships
ALTER TABLE scholarships 
DROP COLUMN IF EXISTS round;

-- 5. Remove new columns from notifications
ALTER TABLE notifications 
DROP COLUMN IF EXISTS sms_sent,
DROP COLUMN IF EXISTS sent_at,
DROP COLUMN IF EXISTS data;

-- 6. Drop new roles
DELETE FROM user_roles WHERE role_id IN (
    SELECT role_id FROM roles WHERE role_name IN ('advisor', 'committee_member')
);
DELETE FROM roles WHERE role_name IN ('advisor', 'committee_member');

-- 7. Revert role permissions
UPDATE roles 
SET permissions = jsonb_strip_nulls(permissions - 'manage_workflow' - 'advisor_review' - 'committee_evaluation')
WHERE role_name = 'scholarship_officer';

-- 8. Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_user;
DROP INDEX IF EXISTS idx_messages_recipient;
DROP INDEX IF EXISTS idx_committee_evaluations_member;
DROP INDEX IF EXISTS idx_advisor_reviews_advisor;
DROP INDEX IF EXISTS idx_application_workflow_step;
DROP INDEX IF EXISTS idx_application_workflow_app_id;

-- 9. Drop ENUM types
DROP TYPE IF EXISTS committee_recommendation;
DROP TYPE IF EXISTS advisor_recommendation;
DROP TYPE IF EXISTS step_status;
DROP TYPE IF EXISTS workflow_step;
DROP TYPE IF EXISTS document_type_enum;
DROP TYPE IF EXISTS scholarship_round;
