-- Migration 025 Down

DROP INDEX IF EXISTS idx_app_assets_category;
DROP INDEX IF EXISTS idx_app_assets_asset_type;
DROP INDEX IF EXISTS idx_app_assets_application_id;
DROP INDEX IF EXISTS idx_app_family_relationship;
DROP INDEX IF EXISTS idx_app_family_application_id;
DROP INDEX IF EXISTS idx_app_edu_history_school_name;
DROP INDEX IF EXISTS idx_app_edu_history_application_id;
DROP INDEX IF EXISTS idx_app_addresses_province;
DROP INDEX IF EXISTS idx_app_addresses_address_type;
DROP INDEX IF EXISTS idx_app_addresses_application_id;
DROP INDEX IF EXISTS idx_app_personal_info_student_id;
DROP INDEX IF EXISTS idx_app_personal_info_citizen_id;
DROP INDEX IF EXISTS idx_app_personal_info_application_id;

DROP TABLE IF EXISTS application_assets;
DROP TABLE IF EXISTS application_family_members;
DROP TABLE IF EXISTS application_education_history;
DROP TABLE IF EXISTS application_addresses;
DROP TABLE IF EXISTS application_personal_info;
