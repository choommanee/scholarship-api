package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	"scholarship-system/internal/handlers"
	"scholarship-system/tests"
)

type NotificationTestSuite struct {
	tests.TestSuite
	notificationHandler *handlers.NotificationHandler
	testUserID          string
	testAdminID         string
	testUserID2         string
}

func (s *NotificationTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.notificationHandler = handlers.NewNotificationHandler(s.Config)
}

func (s *NotificationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	
	// Create test users
	s.testUserID = s.CreateTestUser("user@example.com", "student")
	s.testUserID2 = s.CreateTestUser("user2@example.com", "student") 
	s.testAdminID = s.CreateTestUser("admin@example.com", "admin")
}

func (s *NotificationTestSuite) createTestNotification(userID string, isRead bool) int {
	var notificationID int
	readValue := "false"
	if isRead {
		readValue = "true"
	}
	
	err := s.DB.QueryRow(`
		INSERT INTO notifications (user_id, notification_type, title, message, 
			reference_id, reference_type, priority, is_read)
		VALUES ($1, 'application_received', 'Application Received', 
			'Your scholarship application has been received', '123', 'application', 'normal', $2)
		RETURNING notification_id
	`, userID, readValue).Scan(&notificationID)
	s.Require().NoError(err)
	return notificationID
}

func (s *NotificationTestSuite) TestGetNotifications() {
	// Create test notifications
	notificationID1 := s.createTestNotification(s.testUserID, false) // unread
	notificationID2 := s.createTestNotification(s.testUserID, true)  // read
	notificationID3 := s.createTestNotification(s.testUserID2, false) // different user

	// Setup routes
	app := fiber.New()
	app.Get("/notifications", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.GetNotifications(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "get all notifications",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  2, // Only user's notifications
		},
		{
			name:           "get with pagination",
			queryParams:    "?page=1&limit=1",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get unread only",
			queryParams:    "?unread_only=true",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get second page",
			queryParams:    "?page=2&limit=1",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get with high page number",
			queryParams:    "?page=10&limit=5",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/notifications"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "data")
			s.Contains(response, "pagination")
			
			notifications := response["data"].([]interface{})
			s.Equal(tt.expectedCount, len(notifications))

			if len(notifications) > 0 {
				notification := notifications[0].(map[string]interface{})
				s.Contains(notification, "notification_id")
				s.Contains(notification, "notification_type")
				s.Contains(notification, "title")
				s.Contains(notification, "message")
				s.Contains(notification, "is_read")
			}

			pagination := response["pagination"].(map[string]interface{})
			s.Contains(pagination, "total")
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id IN ($1, $2, $3)", notificationID1, notificationID2, notificationID3)
}

func (s *NotificationTestSuite) TestMarkAsRead() {
	// Create test notification (unread)
	notificationID := s.createTestNotification(s.testUserID, false)

	// Setup routes
	app := fiber.New()
	app.Post("/notifications/:id/read", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.MarkAsRead(c)
	})

	tests := []struct {
		name           string
		notificationID string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful mark as read",
			notificationID: strconv.Itoa(notificationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "mark non-existent notification",
			notificationID: "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Notification not found",
		},
		{
			name:           "mark already read notification",
			notificationID: strconv.Itoa(notificationID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Notification not found",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("POST", "/notifications/"+tt.notificationID+"/read", nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Equal("Notification marked as read", response["message"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id = $1", notificationID)
}

func (s *NotificationTestSuite) TestMarkAllAsRead() {
	// Create multiple test notifications (unread)
	notificationID1 := s.createTestNotification(s.testUserID, false)
	notificationID2 := s.createTestNotification(s.testUserID, false)
	notificationID3 := s.createTestNotification(s.testUserID, true) // already read
	notificationID4 := s.createTestNotification(s.testUserID2, false) // different user

	// Setup routes
	app := fiber.New()
	app.Post("/notifications/mark-all-read", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.MarkAllAsRead(c)
	})

	req := httptest.NewRequest("POST", "/notifications/mark-all-read", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "message")
	s.Contains(response, "updated_count")
	s.Equal("All notifications marked as read", response["message"])
	s.Equal(float64(2), response["updated_count"]) // Only 2 unread notifications for testUserID

	// Verify notifications are marked as read
	var readCount int
	err = s.DB.QueryRow(`
		SELECT COUNT(*) FROM notifications 
		WHERE user_id = $1 AND is_read = true
	`, s.testUserID).Scan(&readCount)
	s.NoError(err)
	s.Equal(3, readCount) // All 3 notifications for testUserID should now be read

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id IN ($1, $2, $3, $4)", notificationID1, notificationID2, notificationID3, notificationID4)
}

func (s *NotificationTestSuite) TestGetUnreadCount() {
	// Create test notifications
	notificationID1 := s.createTestNotification(s.testUserID, false) // unread
	notificationID2 := s.createTestNotification(s.testUserID, false) // unread
	notificationID3 := s.createTestNotification(s.testUserID, true)  // read
	notificationID4 := s.createTestNotification(s.testUserID2, false) // different user

	// Setup routes
	app := fiber.New()
	app.Get("/notifications/unread-count", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.GetUnreadCount(c)
	})

	req := httptest.NewRequest("GET", "/notifications/unread-count", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "unread_count")
	s.Equal(float64(2), response["unread_count"]) // 2 unread notifications for testUserID

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id IN ($1, $2, $3, $4)", notificationID1, notificationID2, notificationID3, notificationID4)
}

func (s *NotificationTestSuite) TestSendNotification() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/notifications/send", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.notificationHandler.SendNotification(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful notification sending",
			payload: map[string]interface{}{
				"user_id":           s.testUserID,
				"notification_type": "document_required",
				"title":             "Document Required",
				"message":           "Please upload your transcript",
				"reference_id":      "APP123",
				"reference_type":    "application",
				"priority":          "high",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"notification_type": "interview_scheduled",
				"title":             "Interview Scheduled",
				// missing user_id, message
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "minimal valid notification",
			payload: map[string]interface{}{
				"user_id":           s.testUserID2,
				"notification_type": "system_maintenance",
				"title":             "System Maintenance",
				"message":           "System will be down for maintenance",
				"priority":          "normal",
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/notifications/send", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Contains(response, "data")
				s.Equal("Notification sent successfully", response["message"])
			}
		})
	}
}

func (s *NotificationTestSuite) TestSendBulkNotification() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/notifications/bulk-send", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.notificationHandler.SendBulkNotification(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful bulk notification",
			payload: map[string]interface{}{
				"user_ids":          []string{s.testUserID, s.testUserID2},
				"notification_type": "deadline_reminder",
				"title":             "Application Deadline Reminder",
				"message":           "Application deadline is approaching in 3 days",
				"reference_type":    "deadline",
				"priority":          "high",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty user_ids array",
			payload: map[string]interface{}{
				"user_ids":          []string{},
				"notification_type": "system_maintenance",
				"title":             "System Maintenance",
				"message":           "System will be down",
				"priority":          "normal",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User IDs are required",
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"user_ids": []string{s.testUserID},
				"title":    "Test Title",
				// missing notification_type, message
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "single user bulk notification",
			payload: map[string]interface{}{
				"user_ids":          []string{s.testUserID},
				"notification_type": "result_announced",
				"title":             "Results Announced",
				"message":           "Scholarship results have been announced",
				"priority":          "high",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/notifications/bulk-send", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Contains(response, "total_users")
				s.Contains(response, "success_count")
				s.Equal("Bulk notification sent", response["message"])
			}
		})
	}
}

func (s *NotificationTestSuite) TestDeleteNotification() {
	// Create test notification
	notificationID := s.createTestNotification(s.testUserID, false)
	// Create notification for different user
	otherNotificationID := s.createTestNotification(s.testUserID2, false)

	// Setup routes
	app := fiber.New()
	app.Delete("/notifications/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.DeleteNotification(c)
	})

	tests := []struct {
		name           string
		notificationID string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful deletion",
			notificationID: strconv.Itoa(notificationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "delete non-existent notification",
			notificationID: "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Notification not found",
		},
		{
			name:           "delete other user's notification",
			notificationID: strconv.Itoa(otherNotificationID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Notification not found",
		},
		{
			name:           "delete already deleted notification",
			notificationID: strconv.Itoa(notificationID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Notification not found",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("DELETE", "/notifications/"+tt.notificationID, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Equal("Notification deleted successfully", response["message"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id = $1", otherNotificationID)
}

func (s *NotificationTestSuite) TestGetNotificationTypes() {
	// Setup routes
	app := fiber.New()
	app.Get("/notifications/types", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.GetNotificationTypes(c)
	})

	req := httptest.NewRequest("GET", "/notifications/types", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "data")
	types := response["data"].([]interface{})
	s.Greater(len(types), 0)

	// Check that each type has required fields
	for _, typeInterface := range types {
		typeObj := typeInterface.(map[string]interface{})
		s.Contains(typeObj, "type")
		s.Contains(typeObj, "description")
	}

	// Check for specific notification types
	typeNames := make([]string, len(types))
	for i, typeInterface := range types {
		typeObj := typeInterface.(map[string]interface{})
		typeNames[i] = typeObj["type"].(string)
	}

	expectedTypes := []string{
		"application_received",
		"document_required", 
		"interview_scheduled",
		"result_announced",
		"fund_allocated",
		"system_maintenance",
	}

	for _, expectedType := range expectedTypes {
		s.Contains(typeNames, expectedType)
	}
}

func (s *NotificationTestSuite) TestGetNotificationsWithDifferentUsers() {
	// Create notifications for different users
	notificationID1 := s.createTestNotification(s.testUserID, false)
	notificationID2 := s.createTestNotification(s.testUserID2, false)

	// Setup routes for testUserID
	app1 := fiber.New()
	app1.Get("/notifications", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.notificationHandler.GetNotifications(c)
	})

	// Setup routes for testUserID2
	app2 := fiber.New()
	app2.Get("/notifications", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID2)
		c.Locals("user_role", "student")
		return s.notificationHandler.GetNotifications(c)
	})

	// Test testUserID notifications
	req1 := httptest.NewRequest("GET", "/notifications", nil)
	resp1, err := app1.Test(req1)
	s.NoError(err)
	defer resp1.Body.Close()

	s.Equal(http.StatusOK, resp1.StatusCode)

	var response1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&response1)
	notifications1 := response1["data"].([]interface{})
	s.Equal(1, len(notifications1)) // Only one notification for testUserID

	// Test testUserID2 notifications
	req2 := httptest.NewRequest("GET", "/notifications", nil)
	resp2, err := app2.Test(req2)
	s.NoError(err)
	defer resp2.Body.Close()

	s.Equal(http.StatusOK, resp2.StatusCode)

	var response2 map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&response2)
	notifications2 := response2["data"].([]interface{})
	s.Equal(1, len(notifications2)) // Only one notification for testUserID2

	// Cleanup
	s.DB.Exec("DELETE FROM notifications WHERE notification_id IN ($1, $2)", notificationID1, notificationID2)
}

// Run the test suite
func TestNotificationSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}