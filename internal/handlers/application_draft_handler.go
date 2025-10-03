package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ApplicationDraftHandler struct {
	cfg             *config.Config
	applicationRepo *repository.ApplicationRepository
	scholarshipRepo *repository.ScholarshipRepository
	userRepo        *repository.UserRepository
}

func NewApplicationDraftHandler(cfg *config.Config) *ApplicationDraftHandler {
	return &ApplicationDraftHandler{
		cfg:             cfg,
		applicationRepo: repository.NewApplicationRepository(),
		scholarshipRepo: repository.NewScholarshipRepository(),
		userRepo:        repository.NewUserRepository(),
	}
}

// CreateDraftRequest represents the request to create a draft application
type CreateDraftRequest struct {
	ScholarshipID uint `json:"scholarship_id" validate:"required"`
}

// CreateDraft creates or returns existing draft application
// @Summary Create draft application
// @Description Create a new draft application or return existing one for a scholarship
// @Tags Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDraftRequest true "Draft request"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/draft [post]
func (h *ApplicationDraftHandler) CreateDraft(c *fiber.Ctx) error {
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

	var req CreateDraftRequest
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

	// Get student_id from students table
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Fallback to email if no student record
			studentID = user.Email
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get student information",
			})
		}
	}

	// Check if scholarship exists
	_, err = h.scholarshipRepo.GetByID(req.ScholarshipID)
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

	// Check if student already has application for this scholarship
	existingApp, err := h.applicationRepo.GetByStudentAndScholarship(studentID, req.ScholarshipID)
	if err == nil && existingApp != nil {
		// If status is 'draft', return existing
		if existingApp.ApplicationStatus == "draft" {
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Draft application already exists",
				"data": fiber.Map{
					"application_id":     existingApp.ApplicationID,
					"scholarship_id":     existingApp.ScholarshipID,
					"student_id":        existingApp.StudentID,
					"application_status": existingApp.ApplicationStatus,
					"created_at":        existingApp.CreatedAt,
					"updated_at":        existingApp.UpdatedAt,
				},
			})
		}

		// If status is 'submitted' or later, return error
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fmt.Sprintf("You already have a %s application for this scholarship", existingApp.ApplicationStatus),
		})
	}

	// Create new draft application
	application := &models.ScholarshipApplication{
		StudentID:         studentID,
		ScholarshipID:     req.ScholarshipID,
		ApplicationStatus: "draft",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := h.applicationRepo.Create(application); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create draft application",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Draft created successfully",
		"data": fiber.Map{
			"application_id":     application.ApplicationID,
			"scholarship_id":     application.ScholarshipID,
			"student_id":        application.StudentID,
			"application_status": application.ApplicationStatus,
			"created_at":        application.CreatedAt,
		},
	})
}

// GetDraft retrieves draft application for a scholarship
// @Summary Get draft application
// @Description Get draft application for a specific scholarship
// @Tags Applications
// @Produce json
// @Security BearerAuth
// @Param scholarship_id query int true "Scholarship ID"
// @Success 200 {object} object{success=bool,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/draft [get]
func (h *ApplicationDraftHandler) GetDraft(c *fiber.Ctx) error {
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

	scholarshipIDStr := c.Query("scholarship_id")
	if scholarshipIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "scholarship_id query parameter is required",
		})
	}

	var scholarshipID uint
	if _, err := fmt.Sscanf(scholarshipIDStr, "%d", &scholarshipID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship_id",
		})
	}

	// Get student info
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// Get student_id from students table
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Fallback to email if no student record
			studentID = user.Email
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get student information",
			})
		}
	}

	// Get draft application
	application, err := h.applicationRepo.GetByStudentAndScholarship(studentID, scholarshipID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Draft application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch draft application",
		})
	}

	// Check if it's actually a draft
	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No draft application found",
		})
	}

	// Get all application details
	detailsRepo := repository.NewApplicationDetailsRepository()
	completeForm, err := detailsRepo.GetCompleteForm(uint(application.ApplicationID))
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application details",
		})
	}

	// Calculate current step based on completed sections
	currentStep := 1
	if completeForm != nil {
		if completeForm.PersonalInfo != nil {
			currentStep = 2
		}
		if len(completeForm.Addresses) > 0 {
			currentStep = 3
		}
		if len(completeForm.EducationHistory) > 0 {
			currentStep = 4
		}
		if len(completeForm.FamilyMembers) > 0 {
			currentStep = 5
		}
		if completeForm.FinancialInfo != nil {
			currentStep = 6
		}
		if len(completeForm.Activities) > 0 {
			currentStep = 7
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"application_id": application.ApplicationID,
			"scholarship_id": application.ScholarshipID,
			"current_step":   currentStep,
			"draft_data":     completeForm,
			"updated_at":     application.UpdatedAt,
		},
	})
}

// CheckEligibility checks if student is eligible for a scholarship
// @Summary Check scholarship eligibility
// @Description Check if student meets scholarship eligibility criteria
// @Tags Scholarships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Scholarship ID"
// @Param data body object{student_data=object} true "Student data"
// @Success 200 {object} object{success=bool,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/scholarships/{id}/check-eligibility [post]
func (h *ApplicationDraftHandler) CheckEligibility(c *fiber.Ctx) error {
	scholarshipID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

	var req struct {
		StudentData map[string]interface{} `json:"student_data"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get scholarship with eligibility criteria
	scholarship, err := h.scholarshipRepo.GetByID(uint(scholarshipID))
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

	// Parse eligibility criteria
	var criteria map[string]interface{}
	if scholarship.EligibilityCriteria != nil && *scholarship.EligibilityCriteria != "" {
		if err := json.Unmarshal([]byte(*scholarship.EligibilityCriteria), &criteria); err != nil {
			// If not JSON, treat as simple criteria
			criteria = make(map[string]interface{})
		}
	} else {
		criteria = make(map[string]interface{})
	}

	// Check eligibility
	isEligible := true
	criteriaResults := []map[string]interface{}{}
	missingRequirements := []string{}
	eligibilityScore := 100.0

	// Check min GPA
	if minGPA, ok := criteria["min_gpa"].(float64); ok {
		if gpa, exists := req.StudentData["gpa"].(float64); exists {
			passed := gpa >= minGPA
			if !passed {
				isEligible = false
				eligibilityScore -= 20
			}
			criteriaResults = append(criteriaResults, map[string]interface{}{
				"criteria": "min_gpa",
				"required": minGPA,
				"actual":   gpa,
				"passed":   passed,
			})
		} else {
			missingRequirements = append(missingRequirements, "GPA information")
		}
	}

	// Check max family income
	if maxIncome, ok := criteria["max_family_income"].(float64); ok {
		if income, exists := req.StudentData["family_income"].(float64); exists {
			passed := income <= maxIncome
			if !passed {
				isEligible = false
				eligibilityScore -= 30
			}
			criteriaResults = append(criteriaResults, map[string]interface{}{
				"criteria": "max_family_income",
				"required": maxIncome,
				"actual":   income,
				"passed":   passed,
			})
		} else {
			missingRequirements = append(missingRequirements, "Family income information")
		}
	}

	// Check allowed faculties
	if allowedFaculties, ok := criteria["allowed_faculties"].([]interface{}); ok {
		if faculty, exists := req.StudentData["faculty"].(string); exists {
			passed := false
			for _, f := range allowedFaculties {
				if fStr, ok := f.(string); ok && fStr == faculty {
					passed = true
					break
				}
			}
			if !passed {
				isEligible = false
				eligibilityScore -= 25
			}
			criteriaResults = append(criteriaResults, map[string]interface{}{
				"criteria": "allowed_faculties",
				"required": allowedFaculties,
				"actual":   faculty,
				"passed":   passed,
			})
		} else {
			missingRequirements = append(missingRequirements, "Faculty information")
		}
	}

	// Check year level
	if minYear, ok := criteria["min_year_level"].(float64); ok {
		if year, exists := req.StudentData["year_level"].(float64); exists {
			passed := year >= minYear
			if !passed {
				isEligible = false
				eligibilityScore -= 15
			}
			criteriaResults = append(criteriaResults, map[string]interface{}{
				"criteria": "min_year_level",
				"required": minYear,
				"actual":   year,
				"passed":   passed,
			})
		} else {
			missingRequirements = append(missingRequirements, "Year level information")
		}
	}

	// Ensure score doesn't go below 0
	if eligibilityScore < 0 {
		eligibilityScore = 0
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_eligible":          isEligible,
			"eligibility_score":    eligibilityScore,
			"criteria_results":     criteriaResults,
			"missing_requirements": missingRequirements,
		},
	})
}
