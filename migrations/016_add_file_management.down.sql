-- Rollback: Drop File Management Tables

DROP TABLE IF EXISTS file_access_logs CASCADE;
DROP TABLE IF EXISTS document_versions CASCADE;
DROP TABLE IF EXISTS file_storage CASCADE;
