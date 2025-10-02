-- Migration 022 Down: Remove Application Reviews and Rankings

-- Drop indexes
DROP INDEX IF EXISTS idx_review_criteria_display_order;
DROP INDEX IF EXISTS idx_review_criteria_is_active;
DROP INDEX IF EXISTS idx_review_criteria_scholarship_id;
DROP INDEX IF EXISTS idx_application_rankings_is_waitlist;
DROP INDEX IF EXISTS idx_application_rankings_is_awarded;
DROP INDEX IF EXISTS idx_application_rankings_status;
DROP INDEX IF EXISTS idx_application_rankings_total_score;
DROP INDEX IF EXISTS idx_application_rankings_rank_position;
DROP INDEX IF EXISTS idx_application_rankings_application_id;
DROP INDEX IF EXISTS idx_application_rankings_scholarship_id;
DROP INDEX IF EXISTS idx_application_rankings_round_id;
DROP INDEX IF EXISTS idx_application_reviews_reviewed_at;
DROP INDEX IF EXISTS idx_application_reviews_review_status;
DROP INDEX IF EXISTS idx_application_reviews_review_stage;
DROP INDEX IF EXISTS idx_application_reviews_reviewer_id;
DROP INDEX IF EXISTS idx_application_reviews_application_id;

-- Drop tables
DROP TABLE IF EXISTS review_criteria;
DROP TABLE IF EXISTS application_rankings;
DROP TABLE IF EXISTS application_reviews;
