-- Migration 026 Down

DROP INDEX IF EXISTS idx_app_scholarship_history_academic_year;
DROP INDEX IF EXISTS idx_app_scholarship_history_application_id;
DROP INDEX IF EXISTS idx_app_financial_application_id;
DROP INDEX IF EXISTS idx_app_living_application_id;
DROP INDEX IF EXISTS idx_app_siblings_sibling_order;
DROP INDEX IF EXISTS idx_app_siblings_application_id;
DROP INDEX IF EXISTS idx_app_guardians_application_id;

DROP TABLE IF EXISTS application_scholarship_history;
DROP TABLE IF EXISTS application_financial_info;
DROP TABLE IF EXISTS application_living_situation;
DROP TABLE IF EXISTS application_siblings;
DROP TABLE IF EXISTS application_guardians;
