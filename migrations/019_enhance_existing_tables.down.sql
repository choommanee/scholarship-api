-- Rollback: Remove enhancements from existing tables

-- Note: This is a complex rollback, use with caution
-- Some columns may have data that would be lost

ALTER TABLE users
DROP COLUMN IF EXISTS line_id,
DROP COLUMN IF EXISTS failed_login_attempts,
DROP COLUMN IF EXISTS account_locked_until,
DROP COLUMN IF EXISTS email_verified;

ALTER TABLE students
DROP COLUMN IF EXISTS current_gpa,
DROP COLUMN IF EXISTS advisor_id,
DROP COLUMN IF EXISTS scholarship_history;

ALTER TABLE scholarships
DROP COLUMN IF EXISTS priority_score,
DROP COLUMN IF EXISTS auto_approval_threshold,
DROP COLUMN IF EXISTS max_applications_per_user;

ALTER TABLE scholarship_applications
DROP COLUMN IF EXISTS gpa_verification_status,
DROP COLUMN IF EXISTS reference_check_status,
DROP COLUMN IF EXISTS automated_score,
DROP COLUMN IF EXISTS manual_override,
DROP COLUMN IF EXISTS risk_assessment;

ALTER TABLE application_documents
DROP COLUMN IF EXISTS original_filename,
DROP COLUMN IF EXISTS verification_status,
DROP COLUMN IF EXISTS verified_by;

ALTER TABLE interview_schedules
DROP COLUMN IF EXISTS interview_time,
DROP COLUMN IF EXISTS interviewer_ids,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS meeting_type,
DROP COLUMN IF EXISTS meeting_link,
DROP COLUMN IF EXISTS duration_minutes,
DROP COLUMN IF EXISTS preparation_notes;

ALTER TABLE interview_appointments
DROP COLUMN IF EXISTS time_slot,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS confirmation_sent,
DROP COLUMN IF EXISTS reminder_sent;

ALTER TABLE interview_results
DROP COLUMN IF EXISTS individual_scores,
DROP COLUMN IF EXISTS interview_duration,
DROP COLUMN IF EXISTS technical_issues,
DROP COLUMN IF EXISTS follow_up_required;

ALTER TABLE scholarship_allocations
DROP COLUMN IF EXISTS academic_year,
DROP COLUMN IF EXISTS semester,
DROP COLUMN IF EXISTS payment_status,
DROP COLUMN IF EXISTS payment_date,
DROP COLUMN IF EXISTS bank_account_verified,
DROP COLUMN IF EXISTS payment_method,
DROP COLUMN IF EXISTS payment_frequency,
DROP COLUMN IF EXISTS installments,
DROP COLUMN IF EXISTS conditions;

ALTER TABLE academic_progress_tracking
DROP COLUMN IF EXISTS academic_year,
DROP COLUMN IF EXISTS credits_completed,
DROP COLUMN IF EXISTS total_credits,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS compliance_status,
DROP COLUMN IF EXISTS warning_issued,
DROP COLUMN IF EXISTS probation_status;

ALTER TABLE notifications
DROP COLUMN IF EXISTS related_entity_type,
DROP COLUMN IF EXISTS related_entity_id,
DROP COLUMN IF EXISTS action_url,
DROP COLUMN IF EXISTS expires_at,
DROP COLUMN IF EXISTS delivery_method,
DROP COLUMN IF EXISTS sent_via_email;

ALTER TABLE messages
DROP COLUMN IF EXISTS is_read,
DROP COLUMN IF EXISTS thread_id,
DROP COLUMN IF EXISTS message_type,
DROP COLUMN IF EXISTS attachments;

ALTER TABLE sso_sessions
DROP COLUMN IF EXISTS device_info,
DROP COLUMN IF EXISTS location_info,
DROP COLUMN IF EXISTS security_level;

ALTER TABLE login_history
DROP COLUMN IF EXISTS device_fingerprint,
DROP COLUMN IF EXISTS risk_score;

ALTER TABLE import_logs
DROP COLUMN IF EXISTS file_size,
DROP COLUMN IF EXISTS processing_time;
