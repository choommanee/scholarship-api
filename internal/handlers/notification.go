package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type NotificationHandler struct {
	cfg *config.Config
}

func NewNotificationHandler(cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{cfg: cfg}
}

// GetNotifications retrieves notifications for current user
// @Summary Get notifications
// @Description Get paginated list of notifications for current user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param unread_only query bool false "Show only unread notifications" default(false)
// @Success 200 {object} object{data=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	unreadOnly := c.Query("unread_only", "false") == "true"
	
	offset := (page - 1) * limit

	query := `SELECT notification_id, notification_type, title, message, 
		reference_id, reference_type, is_read, priority, created_at, read_at
		FROM notifications 
		WHERE user_id = $1`

	args := []interface{}{userID}

	if unreadOnly {
		query += " AND is_read = false"
	}

	query += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.NotificationID, &notification.NotificationType,
			&notification.Title, &notification.Message, &notification.ReferenceID,
			&notification.ReferenceType, &notification.IsRead, &notification.Priority,
			&notification.CreatedAt, &notification.ReadAt,
		)
		if err != nil {
			continue
		}
		notifications = append(notifications, notification)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	countArgs := []interface{}{userID}
	
	if unreadOnly {
		countQuery += " AND is_read = false"
	}
	
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return c.JSON(fiber.Map{
		"data": notifications,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// MarkAsRead marks notification as read
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	userID := c.Locals("user_id").(string)

	query := `UPDATE notifications 
		SET is_read = true, read_at = CURRENT_TIMESTAMP 
		WHERE notification_id = $1 AND user_id = $2`

	result, err := database.DB.Exec(query, notificationID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all notifications as read for current user
// @Summary Mark all notifications as read
// @Description Mark all unread notifications as read for current user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{message=string,updated_count=int}
// @Failure 401 {object} object{error=string}
// @Router /notifications/mark-all-read [post]
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	query := `UPDATE notifications 
		SET is_read = true, read_at = CURRENT_TIMESTAMP 
		WHERE user_id = $1 AND is_read = false`

	result, err := database.DB.Exec(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notifications as read",
		})
	}

	rowsAffected, _ := result.RowsAffected()

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
		"updated_count": rowsAffected,
	})
}

// GetUnreadCount gets count of unread notifications
// @Summary Get unread notification count
// @Description Get count of unread notifications for current user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{unread_count=int}
// @Failure 401 {object} object{error=string}
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()

	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	
	err := database.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}

// SendNotification creates and sends a notification (admin/system use)
// @Summary Send notification
// @Description Send a notification to a user (Admin/Officer only)
// @Tags Notification Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body models.Notification true "Notification data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /notifications/send [post]
func (h *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	var notification models.Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	query := `INSERT INTO notifications 
		(user_id, notification_type, title, message, reference_id, reference_type, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING notification_id`

	err := database.DB.QueryRow(query,
		notification.UserID, notification.NotificationType, notification.Title,
		notification.Message, notification.ReferenceID, notification.ReferenceType,
		notification.Priority,
	).Scan(&notification.NotificationID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send notification",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Notification sent successfully",
		"data":    notification,
	})
}

// SendBulkNotification sends notification to multiple users
// @Summary Send bulk notification
// @Description Send notification to multiple users at once (Admin/Officer only)
// @Tags Notification Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body object{user_ids=[]string,notification_type=string,title=string,message=string,reference_id=string,reference_type=string,priority=string} true "Bulk notification data"
// @Success 200 {object} object{message=string,total_users=int,success_count=int}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /notifications/bulk-send [post]
func (h *NotificationHandler) SendBulkNotification(c *fiber.Ctx) error {
	var request struct {
		UserIDs          []string `json:"user_ids"`
		NotificationType string   `json:"notification_type"`
		Title            string   `json:"title"`
		Message          string   `json:"message"`
		ReferenceID      string   `json:"reference_id,omitempty"`
		ReferenceType    string   `json:"reference_type,omitempty"`
		Priority         string   `json:"priority"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(request.UserIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User IDs are required",
		})
	}

	// Insert notifications in batch
	query := `INSERT INTO notifications 
		(user_id, notification_type, title, message, reference_id, reference_type, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	successCount := 0
	for _, userID := range request.UserIDs {
		_, err := database.DB.Exec(query,
			userID, request.NotificationType, request.Title,
			request.Message, request.ReferenceID, request.ReferenceType,
			request.Priority,
		)
		if err == nil {
			successCount++
		}
	}

	return c.JSON(fiber.Map{
		"message": "Bulk notification sent",
		"total_users": len(request.UserIDs),
		"success_count": successCount,
	})
}

// DeleteNotification deletes a notification
// @Summary Delete notification
// @Description Delete a notification for current user
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	userID := c.Locals("user_id").(string)

	query := `DELETE FROM notifications WHERE notification_id = $1 AND user_id = $2`

	result, err := database.DB.Exec(query, notificationID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notification",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification deleted successfully",
	})
}

// GetNotificationTypes returns available notification types
// @Summary Get notification types
// @Description Get list of available notification types
// @Tags Notification Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Router /notifications/types [get]
func (h *NotificationHandler) GetNotificationTypes(c *fiber.Ctx) error {
	types := []map[string]string{
		{"type": "application_received", "description": "Application received"},
		{"type": "document_required", "description": "Document required"},
		{"type": "interview_scheduled", "description": "Interview scheduled"},
		{"type": "interview_reminder", "description": "Interview reminder"},
		{"type": "result_announced", "description": "Result announced"},
		{"type": "fund_allocated", "description": "Fund allocated"},
		{"type": "fund_transferred", "description": "Fund transferred"},
		{"type": "document_approved", "description": "Document approved"},
		{"type": "document_rejected", "description": "Document rejected"},
		{"type": "application_approved", "description": "Application approved"},
		{"type": "application_rejected", "description": "Application rejected"},
		{"type": "system_maintenance", "description": "System maintenance"},
		{"type": "deadline_reminder", "description": "Deadline reminder"},
	}

	return c.JSON(fiber.Map{
		"data": types,
	})
}

// Helper function to create notification (can be called from other handlers)
func CreateNotification(userID, notificationType, title, message, referenceID, referenceType, priority string) error {
	query := `INSERT INTO notifications 
		(user_id, notification_type, title, message, reference_id, reference_type, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := database.DB.Exec(query, userID, notificationType, title, message, referenceID, referenceType, priority)
	return err
}