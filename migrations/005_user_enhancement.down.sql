-- Drop indexes for new columns
DROP INDEX IF EXISTS idx_users_account_locked_until;
DROP INDEX IF EXISTS idx_users_last_login_at;
DROP INDEX IF EXISTS idx_users_profile_completed;
DROP INDEX IF EXISTS idx_users_email_verified;

-- Remove new columns from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS profile_completion_percentage,
DROP COLUMN IF EXISTS account_locked_until,
DROP COLUMN IF EXISTS failed_login_attempts,
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS password_changed_at,
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS profile_completed,
DROP COLUMN IF EXISTS email_verified_at,
DROP COLUMN IF EXISTS email_verified;

-- Drop tables
DROP TABLE IF EXISTS account_lockouts;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS login_history;
DROP TABLE IF EXISTS password_history;
DROP TABLE IF EXISTS email_verifications; 