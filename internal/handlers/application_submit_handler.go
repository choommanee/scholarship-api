package handlers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ApplicationSubmitHandler struct {
	cfg                    *config.Config
	applicationRepo        *repository.ApplicationRepository
	applicationDetailsRepo *repository.ApplicationDetailsRepository
	userRepo               *repository.UserRepository
}

func NewApplicationSubmitHandler(cfg *config.Config) *ApplicationSubmitHandler {
	return &ApplicationSubmitHandler{
		cfg:                    cfg,
		applicationRepo:        repository.NewApplicationRepository(),
		applicationDetailsRepo: repository.NewApplicationDetailsRepository(),
		userRepo:               repository.NewUserRepository(),
	}
}

// SubmitApplicationRequest represents the submit request
type SubmitApplicationRequest struct {
	TermsAccepted        bool `json:"terms_accepted"`
	DeclarationAccepted  bool `json:"declaration_accepted"`
}

// SubmitApplication submits the draft application for review
// @Summary Submit application
// @Description Submit the draft application for review
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param request body SubmitApplicationRequest true "Submit request"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/submit [post]
func (h *ApplicationSubmitHandler) SubmitApplication(c *fiber.Ctx) error {
	// Get user_id from context
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

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	var req SubmitApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate terms and declaration accepted
	if !req.TermsAccepted || !req.DeclarationAccepted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You must accept the terms and declaration to submit the application",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only submit your own applications",
		})
	}

	// Get application
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

	// Check if application is in draft status
	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Application has already been %s", application.ApplicationStatus),
		})
	}

	// Validate that required information is complete
	form, err := h.applicationDetailsRepo.GetCompleteForm(uint(applicationID))
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate application",
		})
	}

	// Validate required fields
	validationErrors := h.validateApplication(form)
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Please complete all required information before submitting",
			"errors": validationErrors,
		})
	}

	// Validate required documents
	documents, err := h.applicationRepo.GetDocuments(uint(applicationID))
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate documents",
		})
	}

	docValidationErrors := h.validateDocuments(documents)
	if len(docValidationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Please upload all required documents before submitting",
			"errors": docValidationErrors,
		})
	}

	// Generate reference number
	referenceNumber := h.generateReferenceNumber(application)

	// Update status to 'submitted'
	now := time.Now()
	application.ApplicationStatus = "submitted"
	application.SubmittedAt = &now
	application.UpdatedAt = now

	if err := h.applicationRepo.Update(application); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit application",
		})
	}

	// Create workflow record (step 1: submission)
	if err := h.createWorkflowRecord(uint(applicationID), "submitted"); err != nil {
		// Log error but don't fail the submission
		fmt.Printf("Warning: Failed to create workflow record: %v\n", err)
	}

	// TODO: Send email notification to student and officers
	// This would require implementing email service

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Application submitted successfully",
		"data": fiber.Map{
			"application_id":     application.ApplicationID,
			"application_status": application.ApplicationStatus,
			"submitted_at":       application.SubmittedAt,
			"reference_number":   referenceNumber,
		},
	})
}

// validateApplication checks if all required sections are completed
func (h *ApplicationSubmitHandler) validateApplication(form *models.CompleteApplicationForm) []string {
	errors := []string{}

	if form == nil {
		errors = append(errors, "Application data is incomplete")
		return errors
	}

	// Check personal info
	if form.PersonalInfo == nil {
		errors = append(errors, "Personal information is required")
	} else {
		if form.PersonalInfo.FirstNameTH == "" {
			errors = append(errors, "First name (Thai) is required")
		}
		if form.PersonalInfo.LastNameTH == "" {
			errors = append(errors, "Last name (Thai) is required")
		}
		if form.PersonalInfo.Email == "" {
			errors = append(errors, "Email is required")
		}
	}

	// Check addresses
	if len(form.Addresses) == 0 {
		errors = append(errors, "At least one address is required")
	}

	// Check education history
	if len(form.EducationHistory) == 0 {
		errors = append(errors, "Education history is required")
	}

	// Check family members
	if len(form.FamilyMembers) == 0 {
		errors = append(errors, "Family information is required")
	}

	// Check financial info
	if form.FinancialInfo == nil {
		errors = append(errors, "Financial information is required")
	}

	return errors
}

// validateDocuments checks if all required documents are uploaded
func (h *ApplicationSubmitHandler) validateDocuments(documents []models.ApplicationDocument) []string {
	errors := []string{}
	requiredDocs := map[string]bool{
		"id_card":    false,
		"transcript": false,
	}

	for _, doc := range documents {
		if _, required := requiredDocs[doc.DocumentType]; required {
			requiredDocs[doc.DocumentType] = true
		}
	}

	for docType, uploaded := range requiredDocs {
		if !uploaded {
			errors = append(errors, fmt.Sprintf("Required document '%s' is missing", docType))
		}
	}

	return errors
}

// generateReferenceNumber generates a unique reference number for the application
func (h *ApplicationSubmitHandler) generateReferenceNumber(application *models.ScholarshipApplication) string {
	year := time.Now().Year()
	return fmt.Sprintf("SCH-%d-%06d", year, application.ApplicationID)
}

// createWorkflowRecord creates a workflow record for the application
func (h *ApplicationSubmitHandler) createWorkflowRecord(applicationID uint, status string) error {
	query := `
		INSERT INTO application_workflow (application_id, workflow_step, step_status, step_started_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := database.DB.Exec(query, applicationID, 1, status, time.Now())
	return err
}

// verifyApplicationOwnership checks if the user owns the application
func (h *ApplicationSubmitHandler) verifyApplicationOwnership(applicationID uint, userID uuid.UUID) error {
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		return err
	}

	// Get student ID from students table by user_id
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no student record exists, try matching with email
			user, userErr := h.userRepo.GetByID(userID)
			if userErr != nil {
				return userErr
			}

			if application.StudentID != user.Email {
				return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
			}
			return nil
		}
		return err
	}

	if application.StudentID != studentID {
		return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
	}

	return nil
}
