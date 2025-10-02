-- Drop triggers
DROP TRIGGER IF EXISTS trigger_calculate_result_score_percentage ON interview_results_enhanced;
DROP TRIGGER IF EXISTS trigger_calculate_interview_score_percentage ON interview_scores;
DROP TRIGGER IF EXISTS trigger_update_slot_capacity ON interview_bookings;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_score_percentage();
DROP FUNCTION IF EXISTS update_interview_slot_capacity();

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS notification_queue;
DROP TABLE IF EXISTS notification_rules;
DROP TABLE IF EXISTS reviewer_assignments;
DROP TABLE IF EXISTS interview_results_enhanced;
DROP TABLE IF EXISTS interview_scores;
DROP TABLE IF EXISTS scoring_criteria;
DROP TABLE IF EXISTS review_stage_history;
DROP TABLE IF EXISTS application_review_workflow;
DROP TABLE IF EXISTS review_workflow_stages;
DROP TABLE IF EXISTS interview_bookings;
DROP TABLE IF EXISTS interview_slots; 