-- Drop triggers
DROP TRIGGER IF EXISTS update_validation_rules_updated_at ON validation_rules;
DROP TRIGGER IF EXISTS update_application_drafts_updated_at ON application_drafts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Remove new columns from existing tables
ALTER TABLE documents 
DROP COLUMN IF EXISTS total_versions,
DROP COLUMN IF EXISTS current_version_id,
DROP COLUMN IF EXISTS allowed_file_types,
DROP COLUMN IF EXISTS max_file_size,
DROP COLUMN IF EXISTS auto_validation_enabled,
DROP COLUMN IF EXISTS validation_rules,
DROP COLUMN IF EXISTS is_required,
DROP COLUMN IF EXISTS document_category;

-- Remove foreign key constraint first
ALTER TABLE scholarship_applications 
DROP CONSTRAINT IF EXISTS fk_application_draft;

-- Remove new columns from scholarship_applications
ALTER TABLE scholarship_applications 
DROP COLUMN IF EXISTS documents_validated,
DROP COLUMN IF EXISTS documents_uploaded,
DROP COLUMN IF EXISTS total_documents_required,
DROP COLUMN IF EXISTS can_edit_after_submit,
DROP COLUMN IF EXISTS submission_deadline,
DROP COLUMN IF EXISTS last_validation_at,
DROP COLUMN IF EXISTS disqualification_reason,
DROP COLUMN IF EXISTS auto_disqualified,
DROP COLUMN IF EXISTS validation_score,
DROP COLUMN IF EXISTS completion_percentage,
DROP COLUMN IF EXISTS draft_id;

-- Drop new tables
DROP TABLE IF EXISTS bulk_upload_sessions;
DROP TABLE IF EXISTS document_validation_rules;
DROP TABLE IF EXISTS application_workflow_states;
DROP TABLE IF EXISTS application_validation_results;
DROP TABLE IF EXISTS validation_rules;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS application_drafts; 