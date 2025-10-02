-- Migration 024 Down

DROP INDEX IF EXISTS idx_interviewer_availability_is_available;
DROP INDEX IF EXISTS idx_interviewer_availability_date;
DROP INDEX IF EXISTS idx_interviewer_availability_round_id;
DROP INDEX IF EXISTS idx_interviewer_availability_interviewer_id;
DROP INDEX IF EXISTS idx_committee_assignments_status;
DROP INDEX IF EXISTS idx_committee_assignments_interview_id;
DROP INDEX IF EXISTS idx_committee_assignments_committee_id;
DROP INDEX IF EXISTS idx_committee_members_status;
DROP INDEX IF EXISTS idx_committee_members_role;
DROP INDEX IF EXISTS idx_committee_members_interviewer_id;
DROP INDEX IF EXISTS idx_committee_members_committee_id;
DROP INDEX IF EXISTS idx_interview_committees_status;
DROP INDEX IF EXISTS idx_interview_committees_chair_id;
DROP INDEX IF EXISTS idx_interview_committees_scholarship_id;
DROP INDEX IF EXISTS idx_interview_committees_round_id;
DROP INDEX IF EXISTS idx_interviewers_is_available;
DROP INDEX IF EXISTS idx_interviewers_is_active;
DROP INDEX IF EXISTS idx_interviewers_faculty;
DROP INDEX IF EXISTS idx_interviewers_department;
DROP INDEX IF EXISTS idx_interviewers_user_id;

DROP TABLE IF EXISTS interviewer_availability;
DROP TABLE IF EXISTS committee_interview_assignments;
DROP TABLE IF EXISTS committee_members;
DROP TABLE IF EXISTS interview_committees;
DROP TABLE IF EXISTS interviewers;
