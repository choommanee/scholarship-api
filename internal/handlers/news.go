package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type NewsHandler struct {
	cfg      *config.Config
	newsRepo *repository.NewsRepository
	readRepo *repository.NewsReadRepository
}

func NewNewsHandler(cfg *config.Config) *NewsHandler {
	return &NewsHandler{
		cfg:      cfg,
		newsRepo: repository.NewNewsRepository(),
		readRepo: repository.NewNewsReadRepository(),
	}
}

// CreateNewsRequest represents the request body for creating a news article
type CreateNewsRequest struct {
	Title       string   `json:"title" validate:"required"`
	Content     string   `json:"content" validate:"required"`
	Summary     string   `json:"summary" validate:"required"`
	ImageURL    string   `json:"image_url"`
	PublishDate string   `json:"publish_date" validate:"required"`
	ExpireDate  string   `json:"expire_date"`
	Category    string   `json:"category" validate:"required"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
}

// CreateNews creates a new news article
// @Summary Create news article
// @Description Create a new news article (Admin/Superadmin only)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param news body CreateNewsRequest true "News data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /news [post]
func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	// User ID and roles are set by the middleware
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Roles are already checked by the RequireRole middleware
	// But we can double-check if needed
	roles, ok := c.Locals("roles").([]string)
	if !ok || len(roles) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	hasPermission := false
	for _, role := range roles {
		if role == "admin" || role == "superadmin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	var req CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse dates
	publishDate, err := time.Parse("2006-01-02", req.PublishDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid publish date format (use YYYY-MM-DD)",
		})
	}

	var expireDate *time.Time
	if req.ExpireDate != "" {
		parsedExpireDate, err := time.Parse("2006-01-02", req.ExpireDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid expire date format (use YYYY-MM-DD)",
			})
		}
		expireDate = &parsedExpireDate
	}

	// Create image URL pointer
	var imageURL *string
	if req.ImageURL != "" {
		imageURL = &req.ImageURL
	}

	// userID is already a uuid.UUID from the middleware
	userIDValue, ok := userID.(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	news := &models.News{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		ImageURL:    imageURL,
		PublishDate: publishDate,
		ExpireDate:  expireDate,
		Category:    req.Category,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
		CreatedBy:   userIDValue.String(),
	}

	if err := h.newsRepo.CreateNews(news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create news article",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "News article created successfully",
		"news":    news,
	})
}

// GetNews retrieves a news article by ID
// @Summary Get news article
// @Description Get a news article by ID
// @Tags News
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} models.News
// @Failure 404 {object} object{error=string}
// @Router /news/{id} [get]
func (h *NewsHandler) GetNews(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	news, err := h.newsRepo.GetNewsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News article not found",
		})
	}

	return c.JSON(news)
}

// GetUnreadNewsCount retrieves the count of unread news for the current user
// @Summary Get unread news count
// @Description Get the count of unread news articles for the current user
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UnreadNewsCount
// @Failure 401 {object} object{error=string}
// @Router /news/unread/count [get]
func (h *NewsHandler) GetUnreadNewsCount(c *fiber.Ctx) error {
	// Get user ID from context
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID := userIDValue.(uuid.UUID).String()

	// Get unread news count
	count, err := h.readRepo.GetUnreadNewsCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread news count",
		})
	}

	return c.JSON(models.UnreadNewsCount{
		Count: count,
	})
}

// MarkNewsAsRead marks a news article as read by the current user
// @Summary Mark news as read
// @Description Mark a news article as read by the current user
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "News ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /news/{id}/read [post]
func (h *NewsHandler) MarkNewsAsRead(c *fiber.Ctx) error {
	// Get user ID from context
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID := userIDValue.(uuid.UUID).String()

	// Get news ID from path
	newsIDStr := c.Params("id")
	if newsIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "News ID is required",
		})
	}

	// Parse news ID
	newsID, err := uuid.Parse(newsIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	// Check if news exists
	_, err = h.newsRepo.GetNewsByID(newsID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	// Mark news as read
	err = h.readRepo.MarkNewsAsRead(userID, newsID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark news as read",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News marked as read",
	})
}

// UpdateNewsRequest represents the request body for updating a news article
type UpdateNewsRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Summary     string   `json:"summary"`
	ImageURL    string   `json:"image_url"`
	PublishDate string   `json:"publish_date"`
	ExpireDate  string   `json:"expire_date"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
}

// UpdateNews updates a news article
// @Summary Update news article
// @Description Update an existing news article (Admin/Superadmin only)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "News ID"
// @Param news body UpdateNewsRequest true "News data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /news/{id} [put]
func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	// User ID and roles are set by the middleware
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Roles are already checked by the RequireRole middleware
	// But we can double-check if needed
	roles, ok := c.Locals("roles").([]string)
	if !ok || len(roles) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	hasPermission := false
	for _, role := range roles {
		if role == "admin" || role == "superadmin" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	// Get existing news
	existingNews, err := h.newsRepo.GetNewsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News article not found",
		})
	}

	var req UpdateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields if provided
	if req.Title != "" {
		existingNews.Title = req.Title
	}
	if req.Content != "" {
		existingNews.Content = req.Content
	}
	if req.Summary != "" {
		existingNews.Summary = req.Summary
	}
	if req.ImageURL != "" {
		existingNews.ImageURL = &req.ImageURL
	}
	if req.Category != "" {
		existingNews.Category = req.Category
	}
	if req.Tags != nil {
		existingNews.Tags = req.Tags
	}

	// Update publish date if provided
	if req.PublishDate != "" {
		publishDate, err := time.Parse("2006-01-02", req.PublishDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid publish date format (use YYYY-MM-DD)",
			})
		}
		existingNews.PublishDate = publishDate
	}

	// Update expire date if provided
	if req.ExpireDate != "" {
		expireDate, err := time.Parse("2006-01-02", req.ExpireDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid expire date format (use YYYY-MM-DD)",
			})
		}
		existingNews.ExpireDate = &expireDate
	}

	// Update published status
	existingNews.IsPublished = req.IsPublished

	if err := h.newsRepo.UpdateNews(existingNews); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update news article",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News article updated successfully",
		"news":    existingNews,
	})
}

// DeleteNews deletes a news article
// @Summary Delete news article
// @Description Delete a news article (Admin/Superadmin only)
// @Tags News
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "News ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /news/{id} [delete]
func (h *NewsHandler) DeleteNews(c *fiber.Ctx) error {
	// User ID and roles are set by the middleware
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Roles are already checked by the RequireRole middleware
	// But we can double-check if needed
	roles, ok := c.Locals("roles").([]string)
	if !ok || len(roles) == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	hasPermission := false
	for _, role := range roles {
		if role == "admin" || role == "superadmin" {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	if err := h.newsRepo.DeleteNews(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete news article",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News article deleted successfully",
	})
}

// ListNews retrieves a paginated list of news articles
// @Summary List news articles
// @Description Get a paginated list of news articles with optional filtering
// @Tags News
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default 10)"
// @Param page query int false "Page number (default 1)"
// @Param category query string false "Filter by category"
// @Param search query string false "Search in title and content"
// @Param published_only query bool false "Show only published articles"
// @Success 200 {object} object{news=[]models.News,pagination=object}
// @Router /news [get]
func (h *NewsHandler) ListNews(c *fiber.Ctx) error {
	// Parse query parameters
	limitStr := c.Query("limit", "10")
	pageStr := c.Query("page", "1")
	category := c.Query("category", "")
	search := c.Query("search", "")
	publishedOnlyStr := c.Query("published_only", "false")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	publishedOnly := publishedOnlyStr == "true"

	// Calculate offset
	offset := (page - 1) * limit

	// Get news articles
	news, total, err := h.newsRepo.ListNews(limit, offset, category, search, publishedOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch news articles",
		})
	}

	// Calculate total pages
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"news": news,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}
