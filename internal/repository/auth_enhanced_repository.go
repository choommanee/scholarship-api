package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"scholarship-system/internal/models"
)

type AuthEnhancedRepository struct {
	db *sql.DB
}

func NewAuthEnhancedRepository(db *sql.DB) *AuthEnhancedRepository {
	return &AuthEnhancedRepository{db: db}
}

// Email Verification Methods
func (r *AuthEnhancedRepository) CreateEmailVerification(ctx context.Context, userID, email string) (*models.EmailVerification, error) {
	// Generate verification token
	token := generateSecureToken()
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours expiration

	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	query := `
		INSERT INTO email_verifications (user_id, email, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query, verification.UserID, verification.Email, verification.Token, verification.ExpiresAt).
		Scan(&verification.ID, &verification.CreatedAt)

	return verification, err
}

func (r *AuthEnhancedRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*models.EmailVerification, error) {
	verification := &models.EmailVerification{}
	query := `SELECT id, user_id, email, token, expires_at, verified_at, created_at FROM email_verifications WHERE token = $1 AND verified_at IS NULL AND expires_at > NOW()`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&verification.ID, &verification.UserID, &verification.Email,
		&verification.Token, &verification.ExpiresAt, &verification.VerifiedAt, &verification.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return verification, err
}

func (r *AuthEnhancedRepository) VerifyEmail(ctx context.Context, token string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Mark verification as verified
	_, err = tx.ExecContext(ctx,
		`UPDATE email_verifications SET verified_at = NOW() WHERE token = $1`, token)
	if err != nil {
		return err
	}

	// Update user email_verified status
	_, err = tx.ExecContext(ctx, `
		UPDATE users 
		SET email_verified = true, email_verified_at = NOW() 
		WHERE user_id = (
			SELECT user_id FROM email_verifications WHERE token = $1
		)`, token)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Password History Methods
func (r *AuthEnhancedRepository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	query := `INSERT INTO password_history (user_id, password_hash) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, passwordHash)
	return err
}

func (r *AuthEnhancedRepository) CheckPasswordHistory(ctx context.Context, userID, newPassword string, historyCount int) (bool, error) {
	query := `
		SELECT password_hash 
		FROM password_history 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, historyCount)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Check if new password matches any of the recent passwords
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return false, err
		}
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(newPassword)) == nil {
			return true, nil // Password was used before
		}
	}

	return false, nil // Password is new
}

// Session Management Methods
func (r *AuthEnhancedRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (user_id, session_token, refresh_token, device_info, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		session.UserID, session.SessionToken, session.RefreshToken,
		session.DeviceInfo, session.IPAddress, session.UserAgent, session.ExpiresAt).
		Scan(&session.ID, &session.CreatedAt)

	return err
}

func (r *AuthEnhancedRepository) GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	session := &models.UserSession{}
	query := `SELECT id, user_id, session_token, refresh_token, device_info, ip_address, user_agent, expires_at, is_active, last_accessed, created_at FROM user_sessions WHERE session_token = $1 AND is_active = true AND expires_at > NOW()`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.RefreshToken,
		&session.DeviceInfo, &session.IPAddress, &session.UserAgent, &session.ExpiresAt,
		&session.IsActive, &session.LastAccessed, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return session, err
}

func (r *AuthEnhancedRepository) UpdateSessionAccess(ctx context.Context, sessionID int) error {
	query := `UPDATE user_sessions SET last_accessed = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

func (r *AuthEnhancedRepository) DeactivateSession(ctx context.Context, sessionToken string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE session_token = $1`
	_, err := r.db.ExecContext(ctx, query, sessionToken)
	return err
}

func (r *AuthEnhancedRepository) DeactivateAllUserSessions(ctx context.Context, userID string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// Login History Methods
func (r *AuthEnhancedRepository) RecordLoginAttempt(ctx context.Context, userID, method, ipAddress, userAgent, status, failureReason, sessionID string, deviceInfo *models.DeviceInfo) error {
	// Convert userID to UUID if needed
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO login_history (user_id, login_method, ip_address, user_agent, login_status, failure_reason, provider, session_duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = r.db.ExecContext(ctx, query, userUUID, method, ipAddress, userAgent, status, failureReason, nil, nil)
	return err
}

func (r *AuthEnhancedRepository) GetLoginHistory(ctx context.Context, userID string, limit int) ([]models.LoginHistory, error) {
	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var history []models.LoginHistory
	query := `
		SELECT login_id, user_id, login_method, provider, ip_address, user_agent, login_status, failure_reason, login_time, logout_time, session_duration
		FROM login_history 
		WHERE user_id = $1 
		ORDER BY login_time DESC 
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userUUID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.LoginHistory
		err := rows.Scan(&h.LoginID, &h.UserID, &h.LoginMethod, &h.Provider, &h.IPAddress,
			&h.UserAgent, &h.LoginStatus, &h.FailureReason, &h.LoginTime, &h.LogoutTime, &h.SessionDuration)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}

// Account Lockout Methods
func (r *AuthEnhancedRepository) GetAccountLockout(ctx context.Context, userID string) (*models.AccountLockout, error) {
	lockout := &models.AccountLockout{}
	query := `SELECT id, user_id, failed_attempts, locked_until, locked_at, unlock_token, created_at, updated_at FROM account_lockouts WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&lockout.ID, &lockout.UserID, &lockout.FailedAttempts, &lockout.LockedUntil,
		&lockout.LockedAt, &lockout.UnlockToken, &lockout.CreatedAt, &lockout.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return lockout, err
}

func (r *AuthEnhancedRepository) IncrementFailedAttempts(ctx context.Context, userID string) (*models.AccountLockout, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get or create lockout record
	lockout := &models.AccountLockout{}
	err = tx.QueryRowContext(ctx, `SELECT id, user_id, failed_attempts, locked_until, locked_at, unlock_token, created_at, updated_at FROM account_lockouts WHERE user_id = $1`, userID).
		Scan(&lockout.ID, &lockout.UserID, &lockout.FailedAttempts, &lockout.LockedUntil, &lockout.LockedAt, &lockout.UnlockToken, &lockout.CreatedAt, &lockout.UpdatedAt)
	if err == sql.ErrNoRows {
		// Create new lockout record
		lockout = &models.AccountLockout{
			UserID:         userID,
			FailedAttempts: 1,
		}
		err = tx.QueryRowContext(ctx, `
			INSERT INTO account_lockouts (user_id, failed_attempts) 
			VALUES ($1, $2) 
			RETURNING id, created_at, updated_at`,
			lockout.UserID, lockout.FailedAttempts).
			Scan(&lockout.ID, &lockout.CreatedAt, &lockout.UpdatedAt)
	} else if err == nil {
		// Increment existing record
		lockout.FailedAttempts++

		// Lock account if too many attempts (5 attempts = 30 min lock)
		if lockout.FailedAttempts >= 5 {
			lockedUntil := time.Now().Add(30 * time.Minute)
			lockedAt := time.Now()
			lockout.LockedUntil = &lockedUntil
			lockout.LockedAt = &lockedAt
			lockout.UnlockToken = generateSecureToken()
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE account_lockouts 
			SET failed_attempts = $2, locked_until = $3, locked_at = $4, unlock_token = $5, updated_at = NOW()
			WHERE user_id = $1`,
			lockout.UserID, lockout.FailedAttempts, lockout.LockedUntil, lockout.LockedAt, lockout.UnlockToken)

		// Update user table as well
		_, err = tx.ExecContext(ctx, `
			UPDATE users 
			SET failed_login_attempts = $2, account_locked_until = $3 
			WHERE user_id = $1`,
			lockout.UserID, lockout.FailedAttempts, lockout.LockedUntil)
	}

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return lockout, err
}

func (r *AuthEnhancedRepository) ResetFailedAttempts(ctx context.Context, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Reset lockout record
	_, err = tx.ExecContext(ctx, `
		UPDATE account_lockouts 
		SET failed_attempts = 0, locked_until = NULL, locked_at = NULL, unlock_token = NULL, updated_at = NOW()
		WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Reset user table
	_, err = tx.ExecContext(ctx, `
		UPDATE users 
		SET failed_login_attempts = 0, account_locked_until = NULL, last_login_at = NOW()
		WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// User Registration Methods
func (r *AuthEnhancedRepository) CreateStudentUser(ctx context.Context, req *models.StudentRegistrationRequest) (*models.EnhancedUser, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userID := uuid.New().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, phone, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, req.Email, req.Email, string(passwordHash), req.FirstName, req.LastName, req.Phone, true)
	if err != nil {
		return nil, err
	}

	// Create student record
	_, err = tx.ExecContext(ctx, `
		INSERT INTO students (student_id, user_id, faculty_code, department_code, year_level, student_status)
		VALUES ($1, $2, $3, $4, $5, 'active')`,
		req.StudentID, userID, req.Faculty, req.Department, req.YearLevel)
	if err != nil {
		return nil, err
	}

	// Assign student role
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_roles (user_id, role_id, is_active)
		SELECT $1, role_id, true FROM roles WHERE role_name = 'student'`,
		userID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return created user
	return r.GetEnhancedUserByID(ctx, userID)
}

// Profile Management Methods
func (r *AuthEnhancedRepository) UpdateUserProfile(ctx context.Context, userID string, req *models.ProfileUpdateRequest) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.FirstName != "" {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argIndex))
		args = append(args, req.FirstName)
		argIndex++
	}

	if req.LastName != "" {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", argIndex))
		args = append(args, req.LastName)
		argIndex++
	}

	if req.Phone != "" {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, req.Phone)
		argIndex++
	}

	if req.AvatarURL != "" {
		setParts = append(setParts, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, req.AvatarURL)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	// Add updated_at and user_id
	setParts = append(setParts, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = $%d",
		fmt.Sprintf("%s", setParts), argIndex)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *AuthEnhancedRepository) GetEnhancedUserByID(ctx context.Context, userID string) (*models.EnhancedUser, error) {
	user := &models.EnhancedUser{}
	query := `
		SELECT 
			user_id,
			username,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			is_active,
			created_at,
			updated_at,
			COALESCE(email_verified, false) as email_verified,
			email_verified_at,
			COALESCE(profile_completed, false) as profile_completed,
			COALESCE(avatar_url, '') as avatar_url,
			COALESCE(password_changed_at, created_at) as password_changed_at,
			last_login_at,
			COALESCE(failed_login_attempts, 0) as failed_login_attempts,
			account_locked_until,
			COALESCE(profile_completion_percentage, 0) as profile_completion_percentage
		FROM users WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &user.EmailVerified, &user.EmailVerifiedAt,
		&user.ProfileCompleted, &user.AvatarURL, &user.PasswordChangedAt,
		&user.LastLoginAt, &user.FailedLoginAttempts, &user.AccountLockedUntil,
		&user.ProfileCompletionPercentage)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

// Utility functions
func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
