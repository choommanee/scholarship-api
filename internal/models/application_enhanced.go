package models

import (
	"encoding/json"
	"time"
)

// ApplicationDraft represents draft application data
type ApplicationDraft struct {
	ID            int             `json:"id" db:"id"`
	UserID        string          `json:"user_id" db:"user_id"`
	ScholarshipID int             `json:"scholarship_id" db:"scholarship_id"`
	DraftData     json.RawMessage `json:"draft_data" db:"draft_data"`
	CurrentStep   int             `json:"current_step" db:"current_step"`
	TotalSteps    int             `json:"total_steps" db:"total_steps"`
	LastSavedAt   time.Time       `json:"last_saved_at" db:"last_saved_at"`
	ExpiresAt     time.Time       `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

// DocumentVersion represents document version control
type DocumentVersion struct {
	ID               int             `json:"id" db:"id"`
	DocumentID       int             `json:"document_id" db:"document_id"`
	VersionNumber    int             `json:"version_number" db:"version_number"`
	FileName         string          `json:"file_name" db:"file_name"`
	FilePath         string          `json:"file_path" db:"file_path"`
	FileSize         int64           `json:"file_size" db:"file_size"`
	FileType         string          `json:"file_type" db:"file_type"`
	Checksum         string          `json:"checksum" db:"checksum"`
	UploadStatus     string          `json:"upload_status" db:"upload_status"`
	ValidationStatus string          `json:"validation_status" db:"validation_status"`
	ValidationErrors json.RawMessage `json:"validation_errors,omitempty" db:"validation_errors"`
	UploadedBy       string          `json:"uploaded_by" db:"uploaded_by"`
	UploadedAt       time.Time       `json:"uploaded_at" db:"uploaded_at"`
	IsCurrent        bool            `json:"is_current" db:"is_current"`
}

// ValidationRule represents dynamic validation rules
type ValidationRule struct {
	ID               int             `json:"id" db:"id"`
	RuleName         string          `json:"rule_name" db:"rule_name"`
	RuleType         string          `json:"rule_type" db:"rule_type"`
	TargetField      string          `json:"target_field" db:"target_field"`
	RuleConfig       json.RawMessage `json:"rule_config" db:"rule_config"`
	ErrorMessage     string          `json:"error_message" db:"error_message"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	AppliesTo        string          `json:"applies_to" db:"applies_to"`
	ScholarshipTypes json.RawMessage `json:"scholarship_types,omitempty" db:"scholarship_types"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// ApplicationValidationResult represents validation results
type ApplicationValidationResult struct {
	ID               int             `json:"id" db:"id"`
	ApplicationID    int             `json:"application_id" db:"application_id"`
	RuleID           int             `json:"rule_id" db:"rule_id"`
	ValidationStatus string          `json:"validation_status" db:"validation_status"`
	ErrorDetails     json.RawMessage `json:"error_details,omitempty" db:"error_details"`
	ValidatedAt      time.Time       `json:"validated_at" db:"validated_at"`
	ValidatedBy      string          `json:"validated_by" db:"validated_by"`
}

// ApplicationWorkflowState represents application state transitions
type ApplicationWorkflowState struct {
	ID            int             `json:"id" db:"id"`
	ApplicationID int             `json:"application_id" db:"application_id"`
	StateName     string          `json:"state_name" db:"state_name"`
	StateData     json.RawMessage `json:"state_data,omitempty" db:"state_data"`
	EnteredAt     time.Time       `json:"entered_at" db:"entered_at"`
	EnteredBy     string          `json:"entered_by,omitempty" db:"entered_by"`
	Notes         string          `json:"notes,omitempty" db:"notes"`
	IsCurrent     bool            `json:"is_current" db:"is_current"`
}

// DocumentValidationRule represents document-specific validation
type DocumentValidationRule struct {
	ID               int             `json:"id" db:"id"`
	DocumentType     string          `json:"document_type" db:"document_type"`
	RuleName         string          `json:"rule_name" db:"rule_name"`
	ValidationConfig json.RawMessage `json:"validation_config" db:"validation_config"`
	ErrorMessage     string          `json:"error_message" db:"error_message"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// BulkUploadSession represents bulk file upload sessions
type BulkUploadSession struct {
	ID             int             `json:"id" db:"id"`
	SessionToken   string          `json:"session_token" db:"session_token"`
	UserID         string          `json:"user_id" db:"user_id"`
	ApplicationID  *int            `json:"application_id,omitempty" db:"application_id"`
	TotalFiles     int             `json:"total_files" db:"total_files"`
	UploadedFiles  int             `json:"uploaded_files" db:"uploaded_files"`
	FailedFiles    int             `json:"failed_files" db:"failed_files"`
	SessionStatus  string          `json:"session_status" db:"session_status"`
	UploadProgress float64         `json:"upload_progress" db:"upload_progress"`
	ErrorSummary   json.RawMessage `json:"error_summary,omitempty" db:"error_summary"`
	StartedAt      time.Time       `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	ExpiresAt      time.Time       `json:"expires_at" db:"expires_at"`
}

// Enhanced Application model
type EnhancedApplication struct {
	ScholarshipApplication
	DraftID                *int       `json:"draft_id,omitempty" db:"draft_id"`
	CompletionPercentage   float64    `json:"completion_percentage" db:"completion_percentage"`
	ValidationScore        float64    `json:"validation_score" db:"validation_score"`
	AutoDisqualified       bool       `json:"auto_disqualified" db:"auto_disqualified"`
	DisqualificationReason string     `json:"disqualification_reason,omitempty" db:"disqualification_reason"`
	LastValidationAt       *time.Time `json:"last_validation_at,omitempty" db:"last_validation_at"`
	SubmissionDeadline     *time.Time `json:"submission_deadline,omitempty" db:"submission_deadline"`
	CanEditAfterSubmit     bool       `json:"can_edit_after_submit" db:"can_edit_after_submit"`
	TotalDocumentsRequired int        `json:"total_documents_required" db:"total_documents_required"`
	DocumentsUploaded      int        `json:"documents_uploaded" db:"documents_uploaded"`
	DocumentsValidated     int        `json:"documents_validated" db:"documents_validated"`
}

// Enhanced Document model
type EnhancedDocument struct {
	ApplicationDocument
	DocumentCategory      string          `json:"document_category,omitempty" db:"document_category"`
	IsRequired            bool            `json:"is_required" db:"is_required"`
	ValidationRules       json.RawMessage `json:"validation_rules,omitempty" db:"validation_rules"`
	AutoValidationEnabled bool            `json:"auto_validation_enabled" db:"auto_validation_enabled"`
	MaxFileSize           int64           `json:"max_file_size" db:"max_file_size"`
	AllowedFileTypes      json.RawMessage `json:"allowed_file_types" db:"allowed_file_types"`
	CurrentVersionID      *int            `json:"current_version_id,omitempty" db:"current_version_id"`
	TotalVersions         int             `json:"total_versions" db:"total_versions"`
}

// Request/Response models for API

// MultiStepApplicationRequest represents step-by-step application data
type MultiStepApplicationRequest struct {
	ScholarshipID int             `json:"scholarship_id" validate:"required"`
	Step          int             `json:"step" validate:"required,min=1,max=5"`
	StepData      json.RawMessage `json:"step_data" validate:"required"`
	IsComplete    bool            `json:"is_complete"`
	SaveAsDraft   bool            `json:"save_as_draft"`
}

// ApplicationStepData represents different step data structures
type PersonalInfoStep struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	StudentID   string `json:"student_id" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
	Nationality string `json:"nationality"`
}

type AcademicInfoStep struct {
	Faculty       string  `json:"faculty" validate:"required"`
	Department    string  `json:"department" validate:"required"`
	YearLevel     int     `json:"year_level" validate:"required,min=1,max=8"`
	GPA           float64 `json:"gpa" validate:"required,min=0,max=4"`
	AdmissionYear int     `json:"admission_year" validate:"required"`
	Transcript    string  `json:"transcript"` // File upload reference
}

type FinancialInfoStep struct {
	FamilyIncome          float64  `json:"family_income" validate:"required"`
	NumberOfSiblings      int      `json:"number_of_siblings"`
	ParentOccupation      string   `json:"parent_occupation"`
	MonthlyExpenses       float64  `json:"monthly_expenses"`
	HasOtherScholarships  bool     `json:"has_other_scholarships"`
	OtherScholarshipsList string   `json:"other_scholarships_list"`
	IncomeDocuments       []string `json:"income_documents"` // File upload references
}

type ActivityInfoStep struct {
	Extracurricular []ActivityRecord `json:"extracurricular"`
	Volunteer       []ActivityRecord `json:"volunteer"`
	Awards          []AwardRecord    `json:"awards"`
	Skills          []string         `json:"skills"`
	Languages       []LanguageSkill  `json:"languages"`
}

type DocumentsStep struct {
	RequiredDocuments []DocumentUpload `json:"required_documents"`
	OptionalDocuments []DocumentUpload `json:"optional_documents"`
	PersonalStatement string           `json:"personal_statement"`
}

type ActivityRecord struct {
	Name        string `json:"name" validate:"required"`
	Role        string `json:"role"`
	Duration    string `json:"duration"`
	Description string `json:"description"`
	Certificate string `json:"certificate"` // File upload reference
}

type AwardRecord struct {
	Name        string `json:"name" validate:"required"`
	Issuer      string `json:"issuer"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Certificate string `json:"certificate"` // File upload reference
}

type LanguageSkill struct {
	Language string `json:"language" validate:"required"`
	Level    string `json:"level" validate:"required"` // beginner, intermediate, advanced, native
}

type DocumentUpload struct {
	DocumentType string `json:"document_type" validate:"required"`
	FileName     string `json:"file_name" validate:"required"`
	FileID       string `json:"file_id" validate:"required"`
	IsRequired   bool   `json:"is_required"`
	Status       string `json:"status"` // uploaded, validated, rejected
}

// Application Draft Request/Response
type SaveDraftRequest struct {
	ScholarshipID int             `json:"scholarship_id" validate:"required"`
	CurrentStep   int             `json:"current_step" validate:"required"`
	DraftData     json.RawMessage `json:"draft_data" validate:"required"`
	AutoSave      bool            `json:"auto_save"`
}

type DraftResponse struct {
	Success     bool            `json:"success"`
	Message     string          `json:"message"`
	DraftID     int             `json:"draft_id"`
	CurrentStep int             `json:"current_step"`
	TotalSteps  int             `json:"total_steps"`
	LastSavedAt time.Time       `json:"last_saved_at"`
	ExpiresAt   time.Time       `json:"expires_at"`
	DraftData   json.RawMessage `json:"draft_data"`
}

// Document Upload Request/Response
type BulkUploadRequest struct {
	ApplicationID *int     `json:"application_id"`
	DocumentTypes []string `json:"document_types" validate:"required"`
}

type BulkUploadResponse struct {
	Success      bool      `json:"success"`
	SessionToken string    `json:"session_token"`
	UploadURL    string    `json:"upload_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UploadProgressResponse struct {
	SessionToken  string   `json:"session_token"`
	TotalFiles    int      `json:"total_files"`
	UploadedFiles int      `json:"uploaded_files"`
	FailedFiles   int      `json:"failed_files"`
	Progress      float64  `json:"progress"`
	Status        string   `json:"status"`
	Errors        []string `json:"errors,omitempty"`
}

// Validation Request/Response
type ValidateApplicationRequest struct {
	ApplicationID int      `json:"application_id" validate:"required"`
	ValidateAll   bool     `json:"validate_all"`
	RuleTypes     []string `json:"rule_types"`
}

type ValidationResponse struct {
	Success         bool                `json:"success"`
	ValidationScore float64             `json:"validation_score"`
	IsValid         bool                `json:"is_valid"`
	Errors          []ValidationError   `json:"errors,omitempty"`
	Warnings        []ValidationWarning `json:"warnings,omitempty"`
	ValidatedAt     time.Time           `json:"validated_at"`
}

type ValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	RuleType string `json:"rule_type"`
	Severity string `json:"severity"` // error, warning
}

type ValidationWarning struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// Application Preview Request/Response
type ApplicationPreviewRequest struct {
	ApplicationID    int    `json:"application_id" validate:"required"`
	Format           string `json:"format"` // html, pdf
	IncludeDocuments bool   `json:"include_documents"`
}

type ApplicationPreviewResponse struct {
	Success     bool   `json:"success"`
	PreviewHTML string `json:"preview_html,omitempty"`
	PreviewURL  string `json:"preview_url,omitempty"`
	PDFUrl      string `json:"pdf_url,omitempty"`
}

// Helper methods

// CalculateCompletionPercentage calculates application completion percentage
func (app *EnhancedApplication) CalculateCompletionPercentage() float64 {
	totalFields := 20.0 // Define based on required fields
	completedFields := 0.0

	// Count completed fields based on application data
	// This would be implemented based on specific requirements

	return (completedFields / totalFields) * 100
}

// IsEligible checks if application meets basic eligibility criteria
func (app *EnhancedApplication) IsEligible() bool {
	return !app.AutoDisqualified && app.ValidationScore >= 70.0
}

// GetCurrentWorkflowState returns the current workflow state
func (app *EnhancedApplication) GetCurrentWorkflowState() string {
	// This would query the workflow states table
	return app.ApplicationStatus // use ApplicationStatus from ScholarshipApplication
}

// CanBeEdited checks if application can still be edited
func (app *EnhancedApplication) CanBeEdited() bool {
	if app.CanEditAfterSubmit {
		return true
	}

	// Check if before submission deadline
	if app.SubmissionDeadline != nil && time.Now().After(*app.SubmissionDeadline) {
		return false
	}

	// Check workflow state
	editableStates := []string{"draft", "documents_required"}
	currentState := app.GetCurrentWorkflowState()

	for _, state := range editableStates {
		if currentState == state {
			return true
		}
	}

	return false
}

// GetDocumentUploadProgress returns document upload progress
func (app *EnhancedApplication) GetDocumentUploadProgress() float64 {
	if app.TotalDocumentsRequired == 0 {
		return 100.0
	}
	return (float64(app.DocumentsUploaded) / float64(app.TotalDocumentsRequired)) * 100
}
