package models

import (
	"time"

	"github.com/google/uuid"
)

type ScholarshipSource struct {
	SourceID      uint      `json:"source_id" db:"source_id"`
	SourceName    string    `json:"source_name" db:"source_name"`
	SourceType    string    `json:"source_type" db:"source_type"`
	ContactPerson *string   `json:"contact_person" db:"contact_person"`
	ContactEmail  *string   `json:"contact_email" db:"contact_email"`
	ContactPhone  *string   `json:"contact_phone" db:"contact_phone"`
	Description   *string   `json:"description" db:"description"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Scholarship struct {
	ScholarshipID        uint      `json:"scholarship_id" db:"scholarship_id"`
	SourceID             uint      `json:"source_id" db:"source_id"`
	ScholarshipName      string    `json:"scholarship_name" db:"scholarship_name"`
	ScholarshipType      string    `json:"scholarship_type" db:"scholarship_type"`
	Amount               float64   `json:"amount" db:"amount"`
	TotalQuota           int       `json:"total_quota" db:"total_quota"`
	AvailableQuota       int       `json:"available_quota" db:"available_quota"`
	AcademicYear         string    `json:"academic_year" db:"academic_year"`
	Semester             *string   `json:"semester" db:"semester"`
	EligibilityCriteria  *string   `json:"eligibility_criteria" db:"eligibility_criteria"`
	RequiredDocuments    *string   `json:"required_documents" db:"required_documents"`
	ApplicationStartDate time.Time `json:"application_start_date" db:"application_start_date"`
	ApplicationEndDate   time.Time `json:"application_end_date" db:"application_end_date"`
	InterviewRequired    bool      `json:"interview_required" db:"interview_required"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedBy            uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Source *ScholarshipSource `json:"source,omitempty"`
}

type ScholarshipBudget struct {
	BudgetID        uint      `json:"budget_id" db:"budget_id"`
	ScholarshipID   uint      `json:"scholarship_id" db:"scholarship_id"`
	BudgetYear      string    `json:"budget_year" db:"budget_year"`
	TotalBudget     float64   `json:"total_budget" db:"total_budget"`
	AllocatedBudget float64   `json:"allocated_budget" db:"allocated_budget"`
	RemainingBudget float64   `json:"remaining_budget" db:"remaining_budget"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ScholarshipApplication struct {
	ApplicationID           uint       `json:"application_id" db:"application_id"`
	StudentID               string     `json:"student_id" db:"student_id"`
	ScholarshipID           uint       `json:"scholarship_id" db:"scholarship_id"`
	ApplicationStatus       string     `json:"application_status" db:"application_status"`
	ApplicationData         *string    `json:"application_data" db:"application_data"`
	FamilyIncome            *float64   `json:"family_income" db:"family_income"`
	MonthlyExpenses         *float64   `json:"monthly_expenses" db:"monthly_expenses"`
	SiblingsCount           *int       `json:"siblings_count" db:"siblings_count"`
	SpecialAbilities        *string    `json:"special_abilities" db:"special_abilities"`
	ActivitiesParticipation *string    `json:"activities_participation" db:"activities_participation"`
	SubmittedAt             *time.Time `json:"submitted_at" db:"submitted_at"`
	ReviewedBy              *uuid.UUID `json:"reviewed_by" db:"reviewed_by"`
	ReviewedAt              *time.Time `json:"reviewed_at" db:"reviewed_at"`
	ReviewNotes             *string    `json:"review_notes" db:"review_notes"`
	PriorityScore           *float64   `json:"priority_score" db:"priority_score"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Student     *Student     `json:"student,omitempty"`
	Scholarship *Scholarship `json:"scholarship,omitempty"`
}
