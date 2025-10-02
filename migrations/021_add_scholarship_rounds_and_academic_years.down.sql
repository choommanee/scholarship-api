-- Migration 021 Down: Remove Scholarship Rounds and Academic Years

-- Remove column from scholarships
ALTER TABLE scholarships DROP COLUMN IF EXISTS round_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_scholarships_round_id;
DROP INDEX IF EXISTS idx_scholarship_rounds_application_dates;
DROP INDEX IF EXISTS idx_scholarship_rounds_is_active;
DROP INDEX IF EXISTS idx_scholarship_rounds_status;
DROP INDEX IF EXISTS idx_scholarship_rounds_round_number;
DROP INDEX IF EXISTS idx_scholarship_rounds_year_id;
DROP INDEX IF EXISTS idx_academic_years_dates;
DROP INDEX IF EXISTS idx_academic_years_is_active;
DROP INDEX IF EXISTS idx_academic_years_is_current;
DROP INDEX IF EXISTS idx_academic_years_year_code;

-- Drop tables
DROP TABLE IF EXISTS scholarship_rounds;
DROP TABLE IF EXISTS academic_years;
