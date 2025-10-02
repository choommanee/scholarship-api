package models

import (
	"encoding/json"
	"time"
)

// EmailVerification represents email verification records
type EmailVerification struct {
	ID         int        `json:"id" db:"id"`
	UserID     string     `json:"user_id" db:"user_id"`
	Email      string     `json:"email" db:"email"`
	Token      string     `json:"token" db:"token"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// PasswordHistory represents password change history
type PasswordHistory struct {
	ID           int       `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	Platform   string `json:"platform"`
	Browser    string `json:"browser"`
	Version    string `json:"version"`
	DeviceType string `json:"device_type"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
}

// UserSession represents active user sessions
type UserSession struct {
	ID           int         `json:"id" db:"id"`
	UserID       string      `json:"user_id" db:"user_id"`
	SessionToken string      `json:"-" db:"session_token"` // Hidden from JSON
	RefreshToken string      `json:"-" db:"refresh_token"` // Hidden from JSON
	DeviceInfo   *DeviceInfo `json:"device_info,omitempty" db:"device_info"`
	IPAddress    string      `json:"ip_address" db:"ip_address"`
	UserAgent    string      `json:"user_agent" db:"user_agent"`
	ExpiresAt    time.Time   `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	LastAccessed time.Time   `json:"last_accessed" db:"last_accessed"`
	IsActive     bool        `json:"is_active" db:"is_active"`
}

// AccountLockout represents account lockout information
type AccountLockout struct {
	ID             int        `json:"id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	FailedAttempts int        `json:"failed_attempts" db:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	LockedAt       *time.Time `json:"locked_at,omitempty" db:"locked_at"`
	UnlockToken    string     `json:"-" db:"unlock_token"` // Hidden from JSON
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Enhanced User model with new fields
type EnhancedUser struct {
	User
	EmailVerified               bool       `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt             *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	ProfileCompleted            bool       `json:"profile_completed" db:"profile_completed"`
	AvatarURL                   string     `json:"avatar_url,omitempty" db:"avatar_url"`
	PasswordChangedAt           time.Time  `json:"password_changed_at" db:"password_changed_at"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	FailedLoginAttempts         int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	AccountLockedUntil          *time.Time `json:"account_locked_until,omitempty" db:"account_locked_until"`
	ProfileCompletionPercentage int        `json:"profile_completion_percentage" db:"profile_completion_percentage"`
}

// Registration request models
type StudentRegistrationRequest struct {
	StudentID  string `json:"student_id" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Phone      string `json:"phone,omitempty"`
	Faculty    string `json:"faculty" validate:"required"`
	Department string `json:"department,omitempty"`
	YearLevel  int    `json:"year_level" validate:"required,min=1,max=8"`
}

type StaffRegistrationRequest struct {
	EmployeeID string `json:"employee_id" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department" validate:"required"`
	Role       string `json:"role" validate:"required,oneof=scholarship_officer interviewer admin"`
}

// Password management requests
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirm struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// Profile management requests
type ProfileUpdateRequest struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	Faculty    string `json:"faculty,omitempty"`
	Department string `json:"department,omitempty"`
	YearLevel  int    `json:"year_level,omitempty"`
}

// Response models
type RegistrationResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	UserID           string `json:"user_id,omitempty"`
	VerificationSent bool   `json:"verification_sent"`
}

type EmailVerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ProfileResponse struct {
	*EnhancedUser
	CompletionSteps []ProfileCompletionStep `json:"completion_steps"`
}

type ProfileCompletionStep struct {
	Step        string `json:"step"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Required    bool   `json:"required"`
}

// Helper methods for DeviceInfo
func (d *DeviceInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, d)
	case string:
		return json.Unmarshal([]byte(v), d)
	}
	return nil
}

func (d DeviceInfo) Value() (interface{}, error) {
	return json.Marshal(d)
}

// Helper methods for Enhanced User
func (u *EnhancedUser) IsAccountLocked() bool {
	if u.AccountLockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.AccountLockedUntil)
}

func (u *EnhancedUser) CalculateProfileCompletion() int {
	completion := 0
	total := 10

	if u.FirstName != "" {
		completion++
	}
	if u.LastName != "" {
		completion++
	}
	if u.Email != "" {
		completion++
	}
	if u.EmailVerified {
		completion++
	}
	if u.Phone != nil && *u.Phone != "" {
		completion++
	}
	if u.AvatarURL != "" {
		completion++
	}

	// Additional checks for students
	if u.Student != nil {
		if u.Student.FacultyCode != nil && *u.Student.FacultyCode != "" {
			completion++
		}
		if u.Student.DepartmentCode != nil && *u.Student.DepartmentCode != "" {
			completion++
		}
		if u.Student.YearLevel != nil && *u.Student.YearLevel > 0 {
			completion++
		}
		if u.Student.GPA != nil && *u.Student.GPA > 0 {
			completion++
		}
	} else {
		completion += 4 // Skip student-specific fields for non-students
	}

	return (completion * 100) / total
}

func (u *EnhancedUser) GetCompletionSteps() []ProfileCompletionStep {
	steps := []ProfileCompletionStep{
		{
			Step:        "basic_info",
			Description: "ข้อมูลพื้นฐาน (ชื่อ-นามสกุล)",
			Completed:   u.FirstName != "" && u.LastName != "",
			Required:    true,
		},
		{
			Step:        "email_verification",
			Description: "ยืนยันอีเมล",
			Completed:   u.EmailVerified,
			Required:    true,
		},
		{
			Step:        "contact_info",
			Description: "ข้อมูลติดต่อ (เบอร์โทร)",
			Completed:   u.Phone != nil && *u.Phone != "",
			Required:    false,
		},
		{
			Step:        "profile_picture",
			Description: "รูปโปรไฟล์",
			Completed:   u.AvatarURL != "",
			Required:    false,
		},
	}

	return steps
}
