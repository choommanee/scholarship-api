-- Migration 023 Down

DROP INDEX IF EXISTS idx_advisor_assignments_due_date;
DROP INDEX IF EXISTS idx_advisor_assignments_status;
DROP INDEX IF EXISTS idx_advisor_assignments_advisor_id;
DROP INDEX IF EXISTS idx_advisor_assignments_application_id;
DROP INDEX IF EXISTS idx_student_advisors_is_primary;
DROP INDEX IF EXISTS idx_student_advisors_status;
DROP INDEX IF EXISTS idx_student_advisors_advisor_id;
DROP INDEX IF EXISTS idx_student_advisors_student_id;
DROP INDEX IF EXISTS idx_advisors_is_available;
DROP INDEX IF EXISTS idx_advisors_is_active;
DROP INDEX IF EXISTS idx_advisors_faculty;
DROP INDEX IF EXISTS idx_advisors_department;
DROP INDEX IF EXISTS idx_advisors_user_id;

DROP TABLE IF EXISTS advisor_assignments;
DROP TABLE IF EXISTS student_advisors;
DROP TABLE IF EXISTS advisors;
