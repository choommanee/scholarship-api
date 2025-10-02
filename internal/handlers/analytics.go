package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

// AnalyticsHandler handles analytics-related requests
type AnalyticsHandler struct {
	repo *repository.AnalyticsRepository
	cfg  *config.Config
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		repo: repository.NewAnalyticsRepository(database.DB),
		cfg:  cfg,
	}
}

// GetScholarshipStatistics retrieves scholarship statistics
// @Summary Get scholarship statistics
// @Description Get statistics for a specific academic year and round
// @Tags analytics
// @Produce json
// @Param year query string true "Academic Year"
// @Param round query string true "Scholarship Round"
// @Success 200 {object} models.ScholarshipStatistics
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/statistics [get]
func (h *AnalyticsHandler) GetScholarshipStatistics(c *fiber.Ctx) error {
	year := c.Query("year")
	round := c.Query("round")

	if year == "" || round == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Year and round are required",
		})
	}

	stats, err := h.repo.GetStatistics(year, round)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Statistics not found",
		})
	}

	return c.JSON(stats)
}

// GetAllStatistics retrieves all statistics
// @Summary Get all statistics
// @Description Get all scholarship statistics
// @Tags analytics
// @Produce json
// @Success 200 {array} models.ScholarshipStatistics
// @Security BearerAuth
// @Router /api/v1/analytics/statistics/all [get]
func (h *AnalyticsHandler) GetAllStatistics(c *fiber.Ctx) error {
	statistics, err := h.repo.GetAllStatistics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve statistics",
		})
	}

	return c.JSON(statistics)
}

// CreateStatistics creates or updates scholarship statistics
// @Summary Create/Update statistics
// @Description Create or update scholarship statistics
// @Tags analytics
// @Accept json
// @Produce json
// @Param statistics body models.ScholarshipStatistics true "Statistics data"
// @Success 201 {object} models.ScholarshipStatistics
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/statistics [post]
func (h *AnalyticsHandler) CreateStatistics(c *fiber.Ctx) error {
	var stats models.ScholarshipStatistics
	if err := c.BodyParser(&stats); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if stats.StatID == uuid.Nil {
		stats.StatID = uuid.New()
	}

	if err := h.repo.CreateStatistics(&stats); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create statistics",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(stats)
}

// GetApplicationAnalytics retrieves application analytics
// @Summary Get application analytics
// @Description Get analytics for a specific application
// @Tags analytics
// @Produce json
// @Param application_id path int true "Application ID"
// @Success 200 {object} models.ApplicationAnalytics
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/applications/{application_id} [get]
func (h *AnalyticsHandler) GetApplicationAnalytics(c *fiber.Ctx) error {
	appID, err := c.ParamsInt("application_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	analytics, err := h.repo.GetApplicationAnalytics(appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Analytics not found",
		})
	}

	return c.JSON(analytics)
}

// CreateApplicationAnalytics creates application analytics
// @Summary Create application analytics
// @Description Create analytics for an application
// @Tags analytics
// @Accept json
// @Produce json
// @Param analytics body models.ApplicationAnalytics true "Analytics data"
// @Success 201 {object} models.ApplicationAnalytics
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/applications [post]
func (h *AnalyticsHandler) CreateApplicationAnalytics(c *fiber.Ctx) error {
	var analytics models.ApplicationAnalytics
	if err := c.BodyParser(&analytics); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if analytics.AnalyticsID == uuid.Nil {
		analytics.AnalyticsID = uuid.New()
	}

	if err := h.repo.CreateApplicationAnalytics(&analytics); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create analytics",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(analytics)
}

// GetAverageProcessingTime gets average processing time
// @Summary Get average processing time
// @Description Get average processing time across all applications
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/processing-time [get]
func (h *AnalyticsHandler) GetAverageProcessingTime(c *fiber.Ctx) error {
	avg, err := h.repo.GetAverageProcessingTime()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate average",
		})
	}

	return c.JSON(fiber.Map{
		"average_processing_time_days": avg,
	})
}

// GetBottleneckSteps gets common bottleneck steps
// @Summary Get bottleneck steps
// @Description Get most common bottleneck steps in the application process
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/bottlenecks [get]
func (h *AnalyticsHandler) GetBottleneckSteps(c *fiber.Ctx) error {
	bottlenecks, err := h.repo.GetBottleneckSteps()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bottlenecks",
		})
	}

	return c.JSON(fiber.Map{
		"bottlenecks": bottlenecks,
	})
}

// GetDashboardSummary gets dashboard summary
// @Summary Get dashboard summary
// @Description Get summary data for dashboard
// @Tags analytics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboardSummary(c *fiber.Ctx) error {
	// Get latest statistics
	stats, err := h.repo.GetAllStatistics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve statistics",
		})
	}

	// Get average processing time
	avgTime, _ := h.repo.GetAverageProcessingTime()

	// Get bottlenecks
	bottlenecks, _ := h.repo.GetBottleneckSteps()

	var latestStats *models.ScholarshipStatistics
	if len(stats) > 0 {
		latestStats = &stats[0]
	}

	return c.JSON(fiber.Map{
		"latest_statistics":       latestStats,
		"average_processing_time": avgTime,
		"bottlenecks":             bottlenecks,
		"total_periods":           len(stats),
	})
}
