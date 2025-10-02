-- Migration 027 Down

DROP INDEX IF EXISTS idx_app_income_certs_income_type;
DROP INDEX IF EXISTS idx_app_income_certs_application_id;
DROP INDEX IF EXISTS idx_app_house_docs_verified;
DROP INDEX IF EXISTS idx_app_house_docs_document_type;
DROP INDEX IF EXISTS idx_app_house_docs_application_id;
DROP INDEX IF EXISTS idx_app_funding_needs_application_id;
DROP INDEX IF EXISTS idx_app_health_application_id;
DROP INDEX IF EXISTS idx_app_references_application_id;
DROP INDEX IF EXISTS idx_app_activities_activity_type;
DROP INDEX IF EXISTS idx_app_activities_application_id;

DROP TABLE IF EXISTS application_income_certificates;
DROP TABLE IF EXISTS application_house_documents;
DROP TABLE IF EXISTS application_funding_needs;
DROP TABLE IF EXISTS application_health_info;
DROP TABLE IF EXISTS application_references;
DROP TABLE IF EXISTS application_activities;
