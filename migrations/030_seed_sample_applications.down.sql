-- Migration 030 Down: Remove sample application data

-- Delete in reverse order due to foreign key constraints
DELETE FROM application_references WHERE application_id IN (1, 2);
DELETE FROM application_funding_needs WHERE application_id IN (1, 2);
DELETE FROM application_health_info WHERE application_id IN (1, 2);
DELETE FROM application_living_situation WHERE application_id IN (1, 2);
DELETE FROM application_assets WHERE application_id IN (1, 2);
DELETE FROM application_documents WHERE application_id IN (1, 2);
DELETE FROM application_activities WHERE application_id IN (1, 2);
DELETE FROM application_scholarship_history WHERE application_id IN (1, 2);
DELETE FROM application_financial_info WHERE application_id IN (1, 2);
DELETE FROM application_family_members WHERE application_id IN (1, 2);
DELETE FROM application_education_history WHERE application_id IN (1, 2);
DELETE FROM application_addresses WHERE application_id IN (1, 2);
DELETE FROM application_personal_info WHERE application_id IN (1, 2);

-- Reset application status
UPDATE scholarship_applications
SET
    application_status = 'draft',
    priority_score = NULL,
    automated_score = NULL,
    submitted_at = NULL
WHERE application_id IN (1, 2);
