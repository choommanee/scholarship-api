package handlers

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ScholarshipHandler struct {
	cfg             *config.Config
	scholarshipRepo *repository.ScholarshipRepository
}

func NewScholarshipHandler(cfg *config.Config) *ScholarshipHandler {
	return &ScholarshipHandler{
		cfg:             cfg,
		scholarshipRepo: repository.NewScholarshipRepository(),
	}
}

// Scholarship Source Handlers
type CreateSourceRequest struct {
	SourceName    string `json:"source_name" validate:"required"`
	SourceType    string `json:"source_type" validate:"required"`
	ContactPerson string `json:"contact_person"`
	ContactEmail  string `json:"contact_email"`
	ContactPhone  string `json:"contact_phone"`
	Description   string `json:"description"`
}

// CreateSource creates a new scholarship source
// @Summary Create scholarship source
// @Description Create a new scholarship funding source (Admin/Officer only)
// @Tags Scholarship Sources
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param source body CreateSourceRequest true "Source data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /scholarship-sources [post]
func (h *ScholarshipHandler) CreateSource(c *fiber.Ctx) error {
	var req CreateSourceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	source := &models.ScholarshipSource{
		SourceName:    req.SourceName,
		SourceType:    req.SourceType,
		ContactPerson: &req.ContactPerson,
		ContactEmail:  &req.ContactEmail,
		ContactPhone:  &req.ContactPhone,
		Description:   &req.Description,
		IsActive:      true,
	}

	if err := h.scholarshipRepo.CreateSource(source); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create scholarship source",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Scholarship source created successfully",
		"source":  source,
	})
}

func (h *ScholarshipHandler) GetSources(c *fiber.Ctx) error {
	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	search := c.Query("search", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	sources, total, err := h.scholarshipRepo.ListSources(limit, offset, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch scholarship sources",
		})
	}

	return c.JSON(fiber.Map{
		"sources": sources,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// Scholarship Handlers
type CreateScholarshipRequest struct {
	SourceID             uint      `json:"source_id" validate:"required"`
	ScholarshipName      string    `json:"scholarship_name" validate:"required"`
	ScholarshipType      string    `json:"scholarship_type" validate:"required"`
	Amount               float64   `json:"amount" validate:"required,min=0"`
	TotalQuota           int       `json:"total_quota" validate:"required,min=1"`
	AcademicYear         string    `json:"academic_year" validate:"required"`
	Semester             string    `json:"semester"`
	EligibilityCriteria  string    `json:"eligibility_criteria"`
	RequiredDocuments    string    `json:"required_documents"`
	ApplicationStartDate time.Time `json:"application_start_date" validate:"required"`
	ApplicationEndDate   time.Time `json:"application_end_date" validate:"required"`
	InterviewRequired    bool      `json:"interview_required"`
}

func (h *ScholarshipHandler) CreateScholarship(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req CreateScholarshipRequest
	if err := c.BodyParser(&req); err != nil {
		// Log the raw body for debugging
		bodyBytes := c.Body()
		log.Printf("Failed to parse request body: %v, Raw body: %s", err, string(bodyBytes))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Log the parsed request for debugging
	log.Printf("Creating scholarship with data: %+v", req)

	// Validate dates
	if req.ApplicationEndDate.Before(req.ApplicationStartDate) {
		log.Printf("Date validation failed: start=%v, end=%v", req.ApplicationStartDate, req.ApplicationEndDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application end date must be after start date",
		})
	}

	scholarship := &models.Scholarship{
		SourceID:             req.SourceID,
		ScholarshipName:      req.ScholarshipName,
		ScholarshipType:      req.ScholarshipType,
		Amount:               req.Amount,
		TotalQuota:           req.TotalQuota,
		AvailableQuota:       req.TotalQuota, // Initially same as total quota
		AcademicYear:         req.AcademicYear,
		Semester:             &req.Semester,
		EligibilityCriteria:  &req.EligibilityCriteria,
		RequiredDocuments:    &req.RequiredDocuments,
		ApplicationStartDate: req.ApplicationStartDate,
		ApplicationEndDate:   req.ApplicationEndDate,
		InterviewRequired:    req.InterviewRequired,
		IsActive:             true,
		CreatedBy:            userID,
	}

	if err := h.scholarshipRepo.Create(scholarship); err != nil {
		log.Printf("Failed to create scholarship in database: %v", err)
		log.Printf("Scholarship data: %+v", scholarship)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create scholarship",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Scholarship created successfully",
		"scholarship": scholarship,
	})
}

// GetScholarships retrieves scholarships with pagination and filters
// @Summary Get scholarships
// @Description Get list of scholarships with optional filters and pagination
// @Tags Scholarships
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Param search query string false "Search term"
// @Param type query string false "Scholarship type filter"
// @Param academic_year query string false "Academic year filter"
// @Success 200 {object} object{data=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Router /scholarships [get]
func (h *ScholarshipHandler) GetScholarships(c *fiber.Ctx) error {
	// Parse query parameters
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")
	search := c.Query("search", "")
	scholarshipType := c.Query("type", "")
	academicYear := c.Query("academic_year", "")
	activeOnlyStr := c.Query("active_only", "true")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	activeOnly := activeOnlyStr == "true"

	scholarships, total, err := h.scholarshipRepo.List(limit, offset, search, scholarshipType, academicYear, activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch scholarships",
		})
	}

	return c.JSON(fiber.Map{
		"scholarships": scholarships,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *ScholarshipHandler) GetScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

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

	return c.JSON(scholarship)
}

func (h *ScholarshipHandler) UpdateScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

	// Get existing scholarship
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

	var req CreateScholarshipRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate dates
	if req.ApplicationEndDate.Before(req.ApplicationStartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application end date must be after start date",
		})
	}

	// Update scholarship fields
	scholarship.SourceID = req.SourceID
	scholarship.ScholarshipName = req.ScholarshipName
	scholarship.ScholarshipType = req.ScholarshipType
	scholarship.Amount = req.Amount
	scholarship.TotalQuota = req.TotalQuota
	// Don't automatically update available quota - this should be handled separately
	scholarship.AcademicYear = req.AcademicYear
	scholarship.Semester = &req.Semester
	scholarship.EligibilityCriteria = &req.EligibilityCriteria
	scholarship.RequiredDocuments = &req.RequiredDocuments
	scholarship.ApplicationStartDate = req.ApplicationStartDate
	scholarship.ApplicationEndDate = req.ApplicationEndDate
	scholarship.InterviewRequired = req.InterviewRequired

	if err := h.scholarshipRepo.Update(scholarship); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update scholarship",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Scholarship updated successfully",
		"scholarship": scholarship,
	})
}

func (h *ScholarshipHandler) GetAvailableScholarships(c *fiber.Ctx) error {
	scholarships, err := h.scholarshipRepo.GetAvailableScholarships()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch available scholarships",
		})
	}

	return c.JSON(fiber.Map{
		"scholarships": scholarships,
		"count":        len(scholarships),
	})
}

func (h *ScholarshipHandler) ToggleScholarshipStatus(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

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

	// Toggle status
	scholarship.IsActive = !scholarship.IsActive

	if err := h.scholarshipRepo.Update(scholarship); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update scholarship status",
		})
	}

	status := "deactivated"
	if scholarship.IsActive {
		status = "activated"
	}

	return c.JSON(fiber.Map{
		"message":     "Scholarship " + status + " successfully",
		"scholarship": scholarship,
	})
}

// Delete scholarship
func (h *ScholarshipHandler) DeleteScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

	if err := h.scholarshipRepo.Delete(uint(scholarshipID)); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Scholarship not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete scholarship",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Scholarship deleted successfully",
	})
}

// Publish scholarship (set status to open)
func (h *ScholarshipHandler) PublishScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

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

	scholarship.IsActive = true

	if err := h.scholarshipRepo.Update(scholarship); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish scholarship",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Scholarship published successfully",
		"scholarship": scholarship,
	})
}

// Close scholarship applications
func (h *ScholarshipHandler) CloseScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

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

	// Set end date to today to close applications
	now := time.Now()
	scholarship.ApplicationEndDate = now

	if err := h.scholarshipRepo.Update(scholarship); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to close scholarship",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Scholarship closed successfully",
		"scholarship": scholarship,
	})
}

// Suspend scholarship
func (h *ScholarshipHandler) SuspendScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

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

	scholarship.IsActive = false

	if err := h.scholarshipRepo.Update(scholarship); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to suspend scholarship",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Scholarship suspended successfully",
		"scholarship": scholarship,
	})
}

// Duplicate scholarship
func (h *ScholarshipHandler) DuplicateScholarship(c *fiber.Ctx) error {
	scholarshipIDStr := c.Params("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid scholarship ID",
		})
	}

	userID := c.Locals("user_id").(uuid.UUID)

	original, err := h.scholarshipRepo.GetByID(uint(scholarshipID))
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

	// Create duplicate with modified name and reset ID
	duplicate := &models.Scholarship{
		SourceID:             original.SourceID,
		ScholarshipName:      original.ScholarshipName + " (คัดลอก)",
		ScholarshipType:      original.ScholarshipType,
		Amount:               original.Amount,
		TotalQuota:           original.TotalQuota,
		AvailableQuota:       original.TotalQuota,
		AcademicYear:         original.AcademicYear,
		Semester:             original.Semester,
		EligibilityCriteria:  original.EligibilityCriteria,
		RequiredDocuments:    original.RequiredDocuments,
		ApplicationStartDate: original.ApplicationStartDate,
		ApplicationEndDate:   original.ApplicationEndDate,
		InterviewRequired:    original.InterviewRequired,
		IsActive:             false, // Start as inactive
		CreatedBy:            userID,
	}

	if err := h.scholarshipRepo.Create(duplicate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to duplicate scholarship",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Scholarship duplicated successfully",
		"scholarship": duplicate,
	})
}

// GetScholarshipStats retrieves scholarship statistics for admin dashboard
// @Summary Get scholarship statistics
// @Description Get comprehensive scholarship statistics (Admin/Officer only)
// @Tags Scholarships
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{total=int,open=int,closed=int,draft=int,totalBudget=float64}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/scholarships/stats [get]
func (h *ScholarshipHandler) GetScholarshipStats(c *fiber.Ctx) error {
	var stats struct {
		Total       int     `json:"total"`
		Open        int     `json:"open"`
		Closed      int     `json:"closed"`
		Draft       int     `json:"draft"`
		TotalBudget float64 `json:"totalBudget"`
	}

	// Get all scholarships to calculate stats
	scholarships, _, err := h.scholarshipRepo.List(1000, 0, "", "", "", false) // Get all scholarships
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch scholarships for stats",
		})
	}

	stats.Total = len(scholarships)
	now := time.Now()

	for _, scholarship := range scholarships {
		// Calculate total budget
		stats.TotalBudget += scholarship.Amount * float64(scholarship.TotalQuota)

		// Determine status based on dates and activity
		if !scholarship.IsActive {
			stats.Draft++
		} else if scholarship.ApplicationEndDate.Before(now) {
			stats.Closed++
		} else if scholarship.ApplicationStartDate.Before(now) || scholarship.ApplicationStartDate.Equal(now) {
			stats.Open++
		} else {
			stats.Draft++
		}
	}

	return c.JSON(stats)
}
