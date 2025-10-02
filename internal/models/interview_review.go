package models

import (
	"encoding/json"
	"time"
)

// Interview Slot Management
type InterviewSlot struct {
	ID              int       `json:"id" db:"id"`
	ScholarshipID   int       `json:"scholarship_id" db:"scholarship_id"`
	InterviewerID   string    `json:"interviewer_id" db:"interviewer_id"`
	InterviewDate   time.Time `json:"interview_date" db:"interview_date"`
	StartTime       string    `json:"start_time" db:"start_time"` // TIME in database
	EndTime         string    `json:"end_time" db:"end_time"`     // TIME in database
	Location        string    `json:"location,omitempty" db:"location"`
	Building        string    `json:"building,omitempty" db:"building"`
	Floor           string    `json:"floor,omitempty" db:"floor"`
	Room            string    `json:"room,omitempty" db:"room"`
	MaxCapacity     int       `json:"max_capacity" db:"max_capacity"`
	CurrentBookings int       `json:"current_bookings" db:"current_bookings"`
	IsAvailable     bool      `json:"is_available" db:"is_available"`
	SlotType        string    `json:"slot_type" db:"slot_type"` // individual, group
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	PreparationTime int       `json:"preparation_time" db:"preparation_time"`
	Notes           string    `json:"notes,omitempty" db:"notes"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Interviewer *User              `json:"interviewer,omitempty"`
	Scholarship *Scholarship       `json:"scholarship,omitempty"`
	Bookings    []InterviewBooking `json:"bookings,omitempty"`
}

// Interview Booking
type InterviewBooking struct {
	ID                    int        `json:"id" db:"id"`
	SlotID                int        `json:"slot_id" db:"slot_id"`
	ApplicationID         int        `json:"application_id" db:"application_id"`
	StudentID             string     `json:"student_id" db:"student_id"`
	BookingStatus         string     `json:"booking_status" db:"booking_status"` // booked, confirmed, cancelled, completed, no_show, rescheduled
	BookedAt              time.Time  `json:"booked_at" db:"booked_at"`
	ConfirmedAt           *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CancelledAt           *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CancellationReason    string     `json:"cancellation_reason,omitempty" db:"cancellation_reason"`
	RescheduledFromSlotID *int       `json:"rescheduled_from_slot_id,omitempty" db:"rescheduled_from_slot_id"`
	RescheduledToSlotID   *int       `json:"rescheduled_to_slot_id,omitempty" db:"rescheduled_to_slot_id"`
	StudentNotes          string     `json:"student_notes,omitempty" db:"student_notes"`
	OfficerNotes          string     `json:"officer_notes,omitempty" db:"officer_notes"`
	ReminderSentAt        *time.Time `json:"reminder_sent_at,omitempty" db:"reminder_sent_at"`
	CheckInTime           *time.Time `json:"check_in_time,omitempty" db:"check_in_time"`
	CheckOutTime          *time.Time `json:"check_out_time,omitempty" db:"check_out_time"`
	ActualDurationMinutes *int       `json:"actual_duration_minutes,omitempty" db:"actual_duration_minutes"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Slot        *InterviewSlot          `json:"slot,omitempty"`
	Application *ScholarshipApplication `json:"application,omitempty"`
	Student     *User                   `json:"student,omitempty"`
	Result      *InterviewResult        `json:"result,omitempty"`
	Scores      []InterviewScore        `json:"scores,omitempty"`
}

// Review Workflow Stage
type ReviewWorkflowStage struct {
	ID                   int       `json:"id" db:"id"`
	StageName            string    `json:"stage_name" db:"stage_name"`
	StageOrder           int       `json:"stage_order" db:"stage_order"`
	Description          string    `json:"description,omitempty" db:"description"`
	RequiredRole         string    `json:"required_role,omitempty" db:"required_role"`
	IsParallel           bool      `json:"is_parallel" db:"is_parallel"`
	AutoAdvance          bool      `json:"auto_advance" db:"auto_advance"`
	MaxDurationDays      *int      `json:"max_duration_days,omitempty" db:"max_duration_days"`
	NotificationTemplate string    `json:"notification_template,omitempty" db:"notification_template"`
	IsActive             bool      `json:"is_active" db:"is_active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// Application Review Workflow
type ApplicationReviewWorkflow struct {
	ID                     int        `json:"id" db:"id"`
	ApplicationID          int        `json:"application_id" db:"application_id"`
	CurrentStageID         int        `json:"current_stage_id" db:"current_stage_id"`
	WorkflowStatus         string     `json:"workflow_status" db:"workflow_status"` // active, completed, cancelled, on_hold
	StartedAt              time.Time  `json:"started_at" db:"started_at"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date,omitempty" db:"expected_completion_date"`
	ActualCompletionDate   *time.Time `json:"actual_completion_date,omitempty" db:"actual_completion_date"`
	TotalStages            int        `json:"total_stages" db:"total_stages"`
	CompletedStages        int        `json:"completed_stages" db:"completed_stages"`
	CurrentStageStartedAt  time.Time  `json:"current_stage_started_at" db:"current_stage_started_at"`
	StageDeadline          *time.Time `json:"stage_deadline,omitempty" db:"stage_deadline"`
	IsOverdue              bool       `json:"is_overdue" db:"is_overdue"`
	PriorityLevel          string     `json:"priority_level" db:"priority_level"` // high, medium, low
	AssignedReviewer       string     `json:"assigned_reviewer,omitempty" db:"assigned_reviewer"`
	WorkflowNotes          string     `json:"workflow_notes,omitempty" db:"workflow_notes"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Application  *ScholarshipApplication `json:"application,omitempty"`
	CurrentStage *ReviewWorkflowStage    `json:"current_stage,omitempty"`
	StageHistory []ReviewStageHistory    `json:"stage_history,omitempty"`
	Reviewer     *User                   `json:"reviewer,omitempty"`
}

// Review Stage History
type ReviewStageHistory struct {
	ID               int             `json:"id" db:"id"`
	WorkflowID       int             `json:"workflow_id" db:"workflow_id"`
	StageID          int             `json:"stage_id" db:"stage_id"`
	EnteredAt        time.Time       `json:"entered_at" db:"entered_at"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	StageStatus      string          `json:"stage_status" db:"stage_status"` // in_progress, completed, skipped, rejected
	ReviewerID       string          `json:"reviewer_id,omitempty" db:"reviewer_id"`
	ReviewDecision   string          `json:"review_decision,omitempty" db:"review_decision"` // approved, rejected, needs_revision, on_hold
	ReviewNotes      string          `json:"review_notes,omitempty" db:"review_notes"`
	TimeSpentMinutes *int            `json:"time_spent_minutes,omitempty" db:"time_spent_minutes"`
	NextStageID      *int            `json:"next_stage_id,omitempty" db:"next_stage_id"`
	StageData        json.RawMessage `json:"stage_data,omitempty" db:"stage_data"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`

	// Loaded relationships
	Stage     *ReviewWorkflowStage `json:"stage,omitempty"`
	Reviewer  *User                `json:"reviewer,omitempty"`
	NextStage *ReviewWorkflowStage `json:"next_stage,omitempty"`
}

// Scoring Criteria
type ScoringCriteria struct {
	ID                        int             `json:"id" db:"id"`
	CriteriaName              string          `json:"criteria_name" db:"criteria_name"`
	Category                  string          `json:"category" db:"category"` // academic, personal, financial, potential
	Description               string          `json:"description,omitempty" db:"description"`
	MaxScore                  int             `json:"max_score" db:"max_score"`
	WeightPercentage          float64         `json:"weight_percentage" db:"weight_percentage"`
	ScoringMethod             string          `json:"scoring_method" db:"scoring_method"` // numeric, scale, boolean, text
	ScaleOptions              json.RawMessage `json:"scale_options,omitempty" db:"scale_options"`
	IsRequired                bool            `json:"is_required" db:"is_required"`
	AppliesToScholarshipTypes json.RawMessage `json:"applies_to_scholarship_types,omitempty" db:"applies_to_scholarship_types"`
	EvaluationGuidelines      string          `json:"evaluation_guidelines,omitempty" db:"evaluation_guidelines"`
	CreatedBy                 string          `json:"created_by" db:"created_by"`
	IsActive                  bool            `json:"is_active" db:"is_active"`
	CreatedAt                 time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time       `json:"updated_at" db:"updated_at"`
}

// Interview Score
type InterviewScore struct {
	ID               int       `json:"id" db:"id"`
	BookingID        int       `json:"booking_id" db:"booking_id"`
	InterviewerID    string    `json:"interviewer_id" db:"interviewer_id"`
	CriteriaID       int       `json:"criteria_id" db:"criteria_id"`
	Score            float64   `json:"score" db:"score"`
	MaxScore         float64   `json:"max_score" db:"max_score"`
	ScorePercentage  float64   `json:"score_percentage" db:"score_percentage"`
	Comments         string    `json:"comments,omitempty" db:"comments"`
	ScoredAt         time.Time `json:"scored_at" db:"scored_at"`
	IsFinal          bool      `json:"is_final" db:"is_final"`
	ScoringSessionID string    `json:"scoring_session_id,omitempty" db:"scoring_session_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Booking     *InterviewBooking `json:"booking,omitempty"`
	Interviewer *User             `json:"interviewer,omitempty"`
	Criteria    *ScoringCriteria  `json:"criteria,omitempty"`
}

// Interview Result
// Enhanced Interview Result (uses separate table)
type InterviewResultEnhanced struct {
	ID                     int        `json:"id" db:"id"`
	BookingID              int        `json:"booking_id" db:"booking_id"`
	TotalScore             float64    `json:"total_score" db:"total_score"`
	MaxPossibleScore       float64    `json:"max_possible_score" db:"max_possible_score"`
	ScorePercentage        float64    `json:"score_percentage" db:"score_percentage"`
	WeightedScore          *float64   `json:"weighted_score,omitempty" db:"weighted_score"`
	OverallRecommendation  string     `json:"overall_recommendation" db:"overall_recommendation"` // highly_recommended, recommended, conditional, not_recommended
	RecommendationReason   string     `json:"recommendation_reason,omitempty" db:"recommendation_reason"`
	InterviewerFeedback    string     `json:"interviewer_feedback,omitempty" db:"interviewer_feedback"`
	StudentStrengths       string     `json:"student_strengths,omitempty" db:"student_strengths"`
	StudentWeaknesses      string     `json:"student_weaknesses,omitempty" db:"student_weaknesses"`
	ImprovementSuggestions string     `json:"improvement_suggestions,omitempty" db:"improvement_suggestions"`
	FollowUpRequired       bool       `json:"follow_up_required" db:"follow_up_required"`
	FollowUpNotes          string     `json:"follow_up_notes,omitempty" db:"follow_up_notes"`
	ResultStatus           string     `json:"result_status" db:"result_status"` // draft, submitted, approved, rejected
	SubmittedBy            string     `json:"submitted_by,omitempty" db:"submitted_by"`
	SubmittedAt            *time.Time `json:"submitted_at,omitempty" db:"submitted_at"`
	ApprovedBy             string     `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt             *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	RankScore              *float64   `json:"rank_score,omitempty" db:"rank_score"`
	Percentile             *float64   `json:"percentile,omitempty" db:"percentile"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Booking   *InterviewBooking `json:"booking,omitempty"`
	Submitter *User             `json:"submitter,omitempty"`
	Approver  *User             `json:"approver,omitempty"`
}

// Reviewer Assignment
type ReviewerAssignment struct {
	ID                  int        `json:"id" db:"id"`
	ApplicationID       int        `json:"application_id" db:"application_id"`
	ReviewerID          string     `json:"reviewer_id" db:"reviewer_id"`
	ReviewType          string     `json:"review_type" db:"review_type"` // document, interview, final
	AssignmentDate      time.Time  `json:"assignment_date" db:"assignment_date"`
	AssignedBy          string     `json:"assigned_by" db:"assigned_by"`
	ReviewDeadline      *time.Time `json:"review_deadline,omitempty" db:"review_deadline"`
	ReviewStatus        string     `json:"review_status" db:"review_status"` // assigned, in_progress, completed, overdue
	ReviewStartedAt     *time.Time `json:"review_started_at,omitempty" db:"review_started_at"`
	ReviewCompletedAt   *time.Time `json:"review_completed_at,omitempty" db:"review_completed_at"`
	AutoAssigned        bool       `json:"auto_assigned" db:"auto_assigned"`
	WorkloadWeight      float64    `json:"workload_weight" db:"workload_weight"`
	SpecializationMatch *float64   `json:"specialization_match,omitempty" db:"specialization_match"`
	Notes               string     `json:"notes,omitempty" db:"notes"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded relationships
	Application *ScholarshipApplication `json:"application,omitempty"`
	Reviewer    *User                   `json:"reviewer,omitempty"`
	Assigner    *User                   `json:"assigner,omitempty"`
}

// Notification Rule
type NotificationRule struct {
	ID                  int             `json:"id" db:"id"`
	RuleName            string          `json:"rule_name" db:"rule_name"`
	TriggerEvent        string          `json:"trigger_event" db:"trigger_event"`   // interview_booked, review_assigned, deadline_approaching, etc.
	RecipientType       string          `json:"recipient_type" db:"recipient_type"` // student, reviewer, officer, admin
	RecipientCondition  json.RawMessage `json:"recipient_condition,omitempty" db:"recipient_condition"`
	NotificationType    string          `json:"notification_type" db:"notification_type"` // email, sms, push, in_app
	TemplateName        string          `json:"template_name,omitempty" db:"template_name"`
	SendBeforeHours     *int            `json:"send_before_hours,omitempty" db:"send_before_hours"`
	IsRepeating         bool            `json:"is_repeating" db:"is_repeating"`
	RepeatIntervalHours *int            `json:"repeat_interval_hours,omitempty" db:"repeat_interval_hours"`
	MaxRepetitions      int             `json:"max_repetitions" db:"max_repetitions"`
	IsActive            bool            `json:"is_active" db:"is_active"`
	CreatedBy           string          `json:"created_by" db:"created_by"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

// Notification Queue
type NotificationQueue struct {
	ID               int             `json:"id" db:"id"`
	RuleID           int             `json:"rule_id" db:"rule_id"`
	RecipientID      string          `json:"recipient_id" db:"recipient_id"`
	NotificationType string          `json:"notification_type" db:"notification_type"`
	Subject          string          `json:"subject,omitempty" db:"subject"`
	Message          string          `json:"message,omitempty" db:"message"`
	TemplateData     json.RawMessage `json:"template_data,omitempty" db:"template_data"`
	ScheduledFor     time.Time       `json:"scheduled_for" db:"scheduled_for"`
	SentAt           *time.Time      `json:"sent_at,omitempty" db:"sent_at"`
	DeliveryStatus   string          `json:"delivery_status" db:"delivery_status"` // pending, sent, delivered, failed, cancelled
	FailureReason    string          `json:"failure_reason,omitempty" db:"failure_reason"`
	RetryCount       int             `json:"retry_count" db:"retry_count"`
	MaxRetries       int             `json:"max_retries" db:"max_retries"`
	PriorityLevel    string          `json:"priority_level" db:"priority_level"` // high, normal, low
	ReferenceType    string          `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID      string          `json:"reference_id,omitempty" db:"reference_id"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`

	// Loaded relationships
	Rule      *NotificationRule `json:"rule,omitempty"`
	Recipient *User             `json:"recipient,omitempty"`
}

// Request/Response Models for APIs

// Interview Slot Management
type CreateInterviewSlotRequest struct {
	ScholarshipID   int    `json:"scholarship_id" validate:"required"`
	InterviewerID   string `json:"interviewer_id" validate:"required"`
	InterviewDate   string `json:"interview_date" validate:"required"` // YYYY-MM-DD
	StartTime       string `json:"start_time" validate:"required"`     // HH:MM
	EndTime         string `json:"end_time" validate:"required"`       // HH:MM
	Location        string `json:"location"`
	Building        string `json:"building"`
	Floor           string `json:"floor"`
	Room            string `json:"room"`
	MaxCapacity     int    `json:"max_capacity" validate:"min=1"`
	SlotType        string `json:"slot_type" validate:"oneof=individual group"`
	DurationMinutes int    `json:"duration_minutes" validate:"min=15,max=180"`
	PreparationTime int    `json:"preparation_time" validate:"min=0,max=60"`
	Notes           string `json:"notes"`
}

type UpdateInterviewSlotRequest struct {
	Location        *string `json:"location"`
	Building        *string `json:"building"`
	Floor           *string `json:"floor"`
	Room            *string `json:"room"`
	MaxCapacity     *int    `json:"max_capacity"`
	IsAvailable     *bool   `json:"is_available"`
	DurationMinutes *int    `json:"duration_minutes"`
	PreparationTime *int    `json:"preparation_time"`
	Notes           *string `json:"notes"`
}

type InterviewSlotResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    InterviewSlot `json:"data"`
}

type InterviewSlotsResponse struct {
	Success bool                 `json:"success"`
	Data    []InterviewSlot      `json:"data"`
	Meta    PaginationMeta       `json:"meta"`
	Filters InterviewSlotFilters `json:"filters"`
}

type InterviewSlotFilters struct {
	ScholarshipID *int    `json:"scholarship_id"`
	InterviewerID *string `json:"interviewer_id"`
	DateFrom      *string `json:"date_from"`
	DateTo        *string `json:"date_to"`
	IsAvailable   *bool   `json:"is_available"`
	SlotType      *string `json:"slot_type"`
}

// Interview Booking
type BookInterviewRequest struct {
	SlotID       int    `json:"slot_id" validate:"required"`
	StudentNotes string `json:"student_notes"`
}

type UpdateBookingRequest struct {
	BookingStatus      *string `json:"booking_status"`
	CancellationReason *string `json:"cancellation_reason"`
	StudentNotes       *string `json:"student_notes"`
	OfficerNotes       *string `json:"officer_notes"`
}

type RescheduleInterviewRequest struct {
	NewSlotID string `json:"new_slot_id" validate:"required"`
	Reason    string `json:"reason"`
}

type InterviewBookingResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    InterviewBooking `json:"data"`
}

// Interview Scoring
type SubmitScoreRequest struct {
	BookingID        int                    `json:"booking_id" validate:"required"`
	ScoringSessionID string                 `json:"scoring_session_id"`
	Scores           []InterviewScoreSubmit `json:"scores" validate:"required,dive"`
}

type InterviewScoreSubmit struct {
	CriteriaID int     `json:"criteria_id" validate:"required"`
	Score      float64 `json:"score" validate:"required,min=0"`
	Comments   string  `json:"comments"`
}

type SubmitInterviewResultRequest struct {
	BookingID              int    `json:"booking_id" validate:"required"`
	OverallRecommendation  string `json:"overall_recommendation" validate:"required,oneof=highly_recommended recommended conditional not_recommended"`
	RecommendationReason   string `json:"recommendation_reason"`
	InterviewerFeedback    string `json:"interviewer_feedback"`
	StudentStrengths       string `json:"student_strengths"`
	StudentWeaknesses      string `json:"student_weaknesses"`
	ImprovementSuggestions string `json:"improvement_suggestions"`
	FollowUpRequired       bool   `json:"follow_up_required"`
	FollowUpNotes          string `json:"follow_up_notes"`
}

// Review Assignment
type AssignReviewerRequest struct {
	ApplicationID       int      `json:"application_id" validate:"required"`
	ReviewerID          string   `json:"reviewer_id" validate:"required"`
	ReviewType          string   `json:"review_type" validate:"required,oneof=document interview final"`
	ReviewDeadline      *string  `json:"review_deadline"` // YYYY-MM-DD
	WorkloadWeight      *float64 `json:"workload_weight"`
	SpecializationMatch *float64 `json:"specialization_match"`
	Notes               string   `json:"notes"`
}

type ReviewProgressRequest struct {
	ApplicationID int    `json:"application_id" validate:"required"`
	StageID       int    `json:"stage_id" validate:"required"`
	Decision      string `json:"decision" validate:"required,oneof=approved rejected needs_revision on_hold"`
	ReviewNotes   string `json:"review_notes"`
	NextStageID   *int   `json:"next_stage_id"`
}

// Calendar and Availability
type AvailabilityResponse struct {
	Success bool              `json:"success"`
	Data    []DayAvailability `json:"data"`
}

type DayAvailability struct {
	Date           string         `json:"date"`
	DayOfWeek      string         `json:"day_of_week"`
	TotalSlots     int            `json:"total_slots"`
	AvailableSlots int            `json:"available_slots"`
	BookedSlots    int            `json:"booked_slots"`
	TimeSlots      []TimeSlotInfo `json:"time_slots"`
}

type TimeSlotInfo struct {
	SlotID      int    `json:"slot_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsAvailable bool   `json:"is_available"`
	IsBooked    bool   `json:"is_booked"`
	Location    string `json:"location"`
	Interviewer string `json:"interviewer"`
	Duration    int    `json:"duration_minutes"`
}

// Statistics and Reports
type InterviewStatistics struct {
	TotalSlots     int            `json:"total_slots"`
	AvailableSlots int            `json:"available_slots"`
	BookedSlots    int            `json:"booked_slots"`
	CompletedSlots int            `json:"completed_slots"`
	CancelledSlots int            `json:"cancelled_slots"`
	NoShowSlots    int            `json:"no_show_slots"`
	ByStatus       map[string]int `json:"by_status"`
	ByInterviewer  map[string]int `json:"by_interviewer"`
	ByScholarship  map[string]int `json:"by_scholarship"`
	UpcomingToday  int            `json:"upcoming_today"`
	UpcomingWeek   int            `json:"upcoming_week"`
}

type ReviewWorkflowStatistics struct {
	TotalApplications  int            `json:"total_applications"`
	InProgress         int            `json:"in_progress"`
	Completed          int            `json:"completed"`
	Overdue            int            `json:"overdue"`
	ByStage            map[string]int `json:"by_stage"`
	ByPriority         map[string]int `json:"by_priority"`
	ByReviewer         map[string]int `json:"by_reviewer"`
	AverageProcessTime float64        `json:"average_process_time_days"`
}

// Utility Types
type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// Helper Methods
func (is *InterviewSlot) IsBookingAvailable() bool {
	return is.IsAvailable && is.CurrentBookings < is.MaxCapacity
}

func (ib *InterviewBooking) GetDuration() int {
	if ib.CheckInTime != nil && ib.CheckOutTime != nil {
		duration := ib.CheckOutTime.Sub(*ib.CheckInTime)
		return int(duration.Minutes())
	}
	return 0
}

func (ir *InterviewResult) GetRecommendationText() string {
	if ir.Recommendation == nil {
		return "ไม่ระบุ"
	}
	switch *ir.Recommendation {
	case "highly_recommended":
		return "แนะนำอย่างยิ่ง"
	case "recommended":
		return "แนะนำ"
	case "conditional":
		return "แนะนำแบบมีเงื่อนไข"
	case "not_recommended":
		return "ไม่แนะนำ"
	default:
		return "ไม่ระบุ"
	}
}

func (arw *ApplicationReviewWorkflow) GetCompletionPercentage() float64 {
	if arw.TotalStages == 0 {
		return 0
	}
	return (float64(arw.CompletedStages) / float64(arw.TotalStages)) * 100
}

func (arw *ApplicationReviewWorkflow) IsDeadlineApproaching() bool {
	if arw.StageDeadline == nil {
		return false
	}
	hoursUntilDeadline := time.Until(*arw.StageDeadline).Hours()
	return hoursUntilDeadline > 0 && hoursUntilDeadline <= 24
}
