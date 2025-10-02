-- Rollback: Drop Background Jobs System

DROP TABLE IF EXISTS background_tasks CASCADE;
DROP TABLE IF EXISTS job_queue CASCADE;
