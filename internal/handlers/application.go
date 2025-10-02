package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ApplicationHandler struct {
	cfg             *config.Config
	applicationRepo *repository.ApplicationRepository
	scholarshipRepo *repository.ScholarshipRepository
	userRepo        *repository.UserRepository
}

func NewApplicationHandler(cfg *config.Config) *ApplicationHandler {
	return &ApplicationHandler{
		cfg:             cfg,
		applicationRepo: repository.NewApplicationRepository(),
		scholarshipRepo: repository.NewScholarshipRepository(),
		userRepo:        repository.NewUserRepository(),
	}
}

type CreateApplicationRequest struct {
	ScholarshipID           uint     `json:"scholarship_id" validate:"required"`
	FamilyIncome            *float64 `json:"family_income"`
	MonthlyExpenses         *float64 `json:"monthly_expenses"`
	SiblingsCount           *int     `json:"siblings_count"`
	SpecialAbilities        string   `json:"special_abilities"`
	ActivitiesParticipation string   `json:"activities_participation"`
	ApplicationData         string   `json:"application_data"`
}

// CreateApplication creates a new scholarship application
// @Summary Create scholarship application
// @Description Create a new scholarship application (Student only)
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param application body CreateApplicationRequest true "Application data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get student info
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// For now, use email as student_id if no student record exists
	// TODO: Implement proper student repository and validation
	studentID := user.Email

	// Check if scholarship exists and is available
	scholarship, err := h.scholarshipRepo.GetByID(req.ScholarshipID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Scholarship not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch scholarship",
		})
	}

	// Check if scholarship is still accepting applications
	now := time.Now()
	if now.Before(scholarship.ApplicationStartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application period has not started yet",
		})
	}

	if now.After(scholarship.ApplicationEndDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application deadline has passed",
		})
	}

	if scholarship.AvailableQuota <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No available quota for this scholarship",
		})
	}

	// Check if student already applied for this scholarship
	if existingApp, err := h.applicationRepo.GetByStudentAndScholarship(studentID, req.ScholarshipID); err == nil && existingApp != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "You have already applied for this scholarship",
		})
	}

	application := &models.ScholarshipApplication{
		StudentID:               studentID,
		ScholarshipID:           req.ScholarshipID,
		ApplicationStatus:       "draft",
		ApplicationData:         &req.ApplicationData,
		FamilyIncome:            req.FamilyIncome,
		MonthlyExpenses:         req.MonthlyExpenses,
		SiblingsCount:           req.SiblingsCount,
		SpecialAbilities:        &req.SpecialAbilities,
		ActivitiesParticipation: &req.ActivitiesParticipation,
	}

	if err := h.applicationRepo.Create(application); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create application",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Application created successfully",
		"application": application,
	})
}

// GetMyApplications retrieves current user's applications
// @Summary Get my applications
// @Description Get current user's scholarship applications (Student only)
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} object{applications=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Router /applications/my [get]
func (h *ApplicationHandler) GetMyApplications(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get user info to get student ID
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		fmt.Printf("Debug - GetMyApplications: Failed to get user info for userID=%s, error=%v\n", userID.String(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	studentID := user.Email // TODO: Fix this with proper student repository
	fmt.Printf("Debug - GetMyApplications: userID=%s, studentID=%s\n", userID.String(), studentID)

	applications, total, err := h.applicationRepo.ListByStudent(studentID, limit, offset)
	if err != nil {
		fmt.Printf("Debug - GetMyApplications: Failed to list applications for studentID=%s, error=%v\n", studentID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch applications",
		})
	}

	fmt.Printf("Debug - GetMyApplications: Found %d applications (total=%d)\n", len(applications), total)

	// Transform applications to match frontend interface
	var transformedApplications []map[string]interface{}
	for _, app := range applications {
		// Get scholarship details
		scholarship, scholarshipErr := h.scholarshipRepo.GetByID(app.ScholarshipID)
		scholarshipName := "Unknown Scholarship"
		scholarshipAmount := 0.0
		if scholarshipErr == nil && scholarship != nil {
			scholarshipName = scholarship.ScholarshipName
			scholarshipAmount = scholarship.Amount
		}

		// Calculate progress based on status
		progress := calculateProgress(app.ApplicationStatus)

		// Get document status (default to pending for now)
		documentsStatus := map[string]string{
			"transcript":         "pending",
			"income_certificate": "pending",
			"id_card":            "pending",
			"photo":              "pending",
			"recommendation":     "pending",
		}

		// Format dates
		applicationDate := ""
		lastUpdate := ""
		if app.SubmittedAt != nil {
			applicationDate = app.SubmittedAt.Format("2006-01-02")
			lastUpdate = app.SubmittedAt.Format("2006-01-02")
		} else {
			applicationDate = app.CreatedAt.Format("2006-01-02")
			lastUpdate = app.UpdatedAt.Format("2006-01-02")
		}

		transformedApp := map[string]interface{}{
			"id":                strconv.Itoa(int(app.ApplicationID)),
			"scholarshipId":     app.ScholarshipID,
			"scholarshipName":   scholarshipName,
			"scholarshipAmount": scholarshipAmount,
			"applicationDate":   applicationDate,
			"lastUpdate":        lastUpdate,
			"status":            app.ApplicationStatus,
			"progress":          progress,
			"documentsStatus":   documentsStatus,
			"notes":             getApplicationNotes(app.ApplicationStatus),
		}

		// Add optional fields from review notes if available
		if app.ReviewNotes != nil && *app.ReviewNotes != "" {
			transformedApp["rejectionReason"] = *app.ReviewNotes
		}

		transformedApplications = append(transformedApplications, transformedApp)
	}

	return c.JSON(fiber.Map{
		"applications": transformedApplications,
		"pagination": fiber.Map{
			"page":       (offset / limit) + 1,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// Helper function to calculate progress based on status
func calculateProgress(status string) int {
	switch status {
	case "draft":
		return 20
	case "submitted":
		return 40
	case "under_review":
		return 60
	case "document_pending":
		return 50
	case "interview_scheduled":
		return 80
	case "approved":
		return 100
	case "rejected":
		return 100
	default:
		return 0
	}
}

// Helper function to get application notes based on status
func getApplicationNotes(status string) string {
	switch status {
	case "draft":
		return "ร่างใบสมัคร ยังไม่ส่ง"
	case "submitted":
		return "ส่งใบสมัครแล้ว รอการตรวจสอบเอกสาร"
	case "under_review":
		return "อยู่ระหว่างการพิจารณา"
	case "document_pending":
		return "รอเอกสารเพิ่มเติม"
	case "interview_scheduled":
		return "ผ่านการตรวจสอบเอกสารแล้ว มีนัดสัมภาษณ์"
	case "approved":
		return "ได้รับการอนุมัติแล้ว จะได้รับเงินทุนภายใน 30 วัน"
	case "rejected":
		return "ไม่ได้รับการอนุมัติ"
	default:
		return ""
	}
}

func (h *ApplicationHandler) GetApplications(c *fiber.Ctx) error {
	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	status := c.Query("status", "")
	scholarshipType := c.Query("scholarship_type", "")
	scholarshipIDStr := c.Query("scholarship_id", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var scholarshipID *uint
	if scholarshipIDStr != "" {
		id, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			scholarshipID = &uid
		}
	}

	applications, total, err := h.applicationRepo.List(limit, offset, status, scholarshipType, scholarshipID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch applications",
		})
	}

	return c.JSON(fiber.Map{
		"applications": applications,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Get documents
	documents, err := h.applicationRepo.GetDocuments(uint(applicationID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application documents",
		})
	}

	return c.JSON(fiber.Map{
		"application": application,
		"documents":   documents,
	})
}

func (h *ApplicationHandler) UpdateApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Get existing application
	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Check if user owns this application (for students)
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	roles := c.Locals("roles").([]string)
	isStudent := false
	for _, role := range roles {
		if role == "student" {
			isStudent = true
			break
		}
	}

	if isStudent && application.StudentID != user.Email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only update your own applications",
		})
	}

	// Check if application can be updated
	if application.ApplicationStatus == "submitted" && isStudent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot update submitted application",
		})
	}

	var req CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update application fields
	application.ApplicationData = &req.ApplicationData
	application.FamilyIncome = req.FamilyIncome
	application.MonthlyExpenses = req.MonthlyExpenses
	application.SiblingsCount = req.SiblingsCount
	application.SpecialAbilities = &req.SpecialAbilities
	application.ActivitiesParticipation = &req.ActivitiesParticipation

	if err := h.applicationRepo.Update(application); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update application",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Application updated successfully",
		"application": application,
	})
}

func (h *ApplicationHandler) SubmitApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Get existing application
	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Check if user owns this application
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	if application.StudentID != user.Email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only submit your own applications",
		})
	}

	// Check if application is in draft status
	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application has already been submitted",
		})
	}

	// Submit application
	if err := h.applicationRepo.Submit(uint(applicationID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit application",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Application submitted successfully",
	})
}

func (h *ApplicationHandler) ReviewApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	var req struct {
		Status      string  `json:"status" validate:"required"`
		ReviewNotes string  `json:"review_notes"`
		Score       float64 `json:"score"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	validStatuses := []string{"under_review", "approved", "rejected", "interview_scheduled"}
	isValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application status",
		})
	}

	if err := h.applicationRepo.UpdateStatus(uint(applicationID), req.Status, &userID, &req.ReviewNotes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update application status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Application reviewed successfully",
	})
}

func (h *ApplicationHandler) DeleteApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Get existing application
	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Check if user owns this application
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	roles := c.Locals("roles").([]string)
	isStudent := false
	for _, role := range roles {
		if role == "student" {
			isStudent = true
			break
		}
	}

	if isStudent && application.StudentID != user.Email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own applications",
		})
	}

	// Check if application can be deleted
	if application.ApplicationStatus == "submitted" && isStudent {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete submitted application",
		})
	}

	if err := h.applicationRepo.Delete(uint(applicationID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete application",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Application deleted successfully",
	})
}

// GetApplicationStats retrieves application statistics for admin dashboard
// @Summary Get application statistics
// @Description Get comprehensive application statistics (Admin/Officer only)
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{total=int,pending=int,approved=int,rejected=int,interview=int,overdue=int}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/applications/stats [get]
func (h *ApplicationHandler) GetApplicationStats(c *fiber.Ctx) error {
	var stats struct {
		Total     int `json:"total"`
		Pending   int `json:"pending"`
		Approved  int `json:"approved"`
		Rejected  int `json:"rejected"`
		Interview int `json:"interview"`
		Overdue   int `json:"overdue"`
	}

	// Get total applications
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications").Scan(&stats.Total); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch total applications",
		})
	}

	// Get pending applications
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'submitted'").Scan(&stats.Pending); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pending applications",
		})
	}

	// Get approved applications
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'approved'").Scan(&stats.Approved); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch approved applications",
		})
	}

	// Get rejected applications
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'rejected'").Scan(&stats.Rejected); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rejected applications",
		})
	}

	// Get interview scheduled applications
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'interview_scheduled'").Scan(&stats.Interview); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch interview applications",
		})
	}

	// Get overdue applications (submitted more than 30 days ago without review)
	if err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM scholarship_applications 
		WHERE application_status = 'submitted' 
		AND submitted_at < NOW() - INTERVAL '30 days'
	`).Scan(&stats.Overdue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch overdue applications",
		})
	}

	return c.JSON(stats)
}

// UpdateApplicationStatus updates application status directly
// @Summary Update application status
// @Description Update application status directly (Admin/Officer only)
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Param status body object{status=string,notes=string} true "Status update data"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /admin/applications/{id}/status [put]
func (h *ApplicationHandler) UpdateApplicationStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	var req struct {
		Status string `json:"status" validate:"required"`
		Notes  string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	validStatuses := []string{"submitted", "under_review", "approved", "rejected", "interview_scheduled", "completed"}
	isValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application status",
		})
	}

	// Check if application exists
	_, err = h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Update status
	if err := h.applicationRepo.UpdateStatus(uint(applicationID), req.Status, &userID, &req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update application status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Application status updated successfully",
	})
}
