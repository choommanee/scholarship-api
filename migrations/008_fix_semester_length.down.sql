-- Rollback semester field length changes
ALTER TABLE scholarships ALTER COLUMN semester TYPE VARCHAR(10);
ALTER TABLE scholarships ALTER COLUMN academic_year TYPE VARCHAR(10); 