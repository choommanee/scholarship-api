package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Phone        *string    `json:"phone" db:"phone"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	SSOProvider  *string    `json:"sso_provider" db:"sso_provider"`
	SSOUserID    *string    `json:"sso_user_id" db:"sso_user_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`

	// Loaded relationships
	UserRoles []UserRole               `json:"user_roles,omitempty"`
	Student   *Student                 `json:"student,omitempty"`
	Roles     []map[string]interface{} `json:"roles,omitempty"`
}

type Role struct {
	RoleID          uint      `json:"role_id" db:"role_id"`
	RoleName        string    `json:"role_name" db:"role_name"`
	RoleDescription *string   `json:"role_description" db:"role_description"`
	Permissions     *string   `json:"permissions" db:"permissions"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID     uint       `json:"role_id" db:"role_id"`
	AssignedAt time.Time  `json:"assigned_at" db:"assigned_at"`
	AssignedBy *uuid.UUID `json:"assigned_by" db:"assigned_by"`
	IsActive   bool       `json:"is_active" db:"is_active"`

	// Loaded relationships
	Role *Role `json:"role,omitempty"`
}

type Student struct {
	StudentID      string     `json:"student_id" db:"student_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	FacultyCode    *string    `json:"faculty_code" db:"faculty_code"`
	DepartmentCode *string    `json:"department_code" db:"department_code"`
	YearLevel      *int       `json:"year_level" db:"year_level"`
	GPA            *float64   `json:"gpa" db:"gpa"`
	AdmissionYear  *int       `json:"admission_year" db:"admission_year"`
	GraduationYear *int       `json:"graduation_year" db:"graduation_year"`
	StudentStatus  string     `json:"student_status" db:"student_status"`
}

type SSOSession struct {
	SessionID      string     `json:"session_id" db:"session_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	SSOSessionID   *string    `json:"sso_session_id" db:"sso_session_id"`
	Provider       string     `json:"provider" db:"provider"`
	AccessToken    *string    `json:"access_token" db:"access_token"`
	RefreshToken   *string    `json:"refresh_token" db:"refresh_token"`
	TokenExpiresAt *time.Time `json:"token_expires_at" db:"token_expires_at"`
	SessionData    *string    `json:"session_data" db:"session_data"`
	IPAddress      *string    `json:"ip_address" db:"ip_address"`
	UserAgent      *string    `json:"user_agent" db:"user_agent"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	LastAccessed   time.Time  `json:"last_accessed" db:"last_accessed"`
	IsActive       bool       `json:"is_active" db:"is_active"`
}

type LoginHistory struct {
	LoginID         uint       `json:"login_id" db:"login_id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	LoginMethod     string     `json:"login_method" db:"login_method"`
	Provider        *string    `json:"provider" db:"provider"`
	IPAddress       *string    `json:"ip_address" db:"ip_address"`
	UserAgent       *string    `json:"user_agent" db:"user_agent"`
	LoginStatus     string     `json:"login_status" db:"login_status"`
	FailureReason   *string    `json:"failure_reason" db:"failure_reason"`
	LoginTime       time.Time  `json:"login_time" db:"login_time"`
	LogoutTime      *time.Time `json:"logout_time" db:"logout_time"`
	SessionDuration *int       `json:"session_duration" db:"session_duration"`
}