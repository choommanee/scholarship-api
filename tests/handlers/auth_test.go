package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	"scholarship-system/internal/handlers"
	"scholarship-system/tests"
)

type AuthTestSuite struct {
	tests.TestSuite
	authHandler *handlers.AuthHandler
}

func (s *AuthTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.authHandler = handlers.NewAuthHandler(s.Config)
}

func (s *AuthTestSuite) TestRegister() {
	// Setup routes
	app := fiber.New()
	app.Post("/register", s.authHandler.Register)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful registration",
			payload: map[string]interface{}{
				"username":   "testuser",
				"email":      "test@example.com",
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				// missing password, first_name, last_name
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
		{
			name: "invalid email format",
			payload: map[string]interface{}{
				"username":   "testuser",
				"email":      "invalid-email",
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
		{
			name: "duplicate email",
			payload: map[string]interface{}{
				"username":   "testuser2",
				"email":      "test@example.com", // same email as first test
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "User already exists",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Create request
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			// Assert status code
			s.Equal(tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Equal("User registered successfully", response["message"])
			}
		})
	}
}

func (s *AuthTestSuite) TestLogin() {
	// Create test user first
	userID := s.CreateTestUser("login@example.com", "student")

	// Setup routes
	app := fiber.New()
	app.Post("/login", s.authHandler.Login)

	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful login",
			email:          "login@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid email",
			email:          "nonexistent@example.com",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name:           "invalid password",
			email:          "login@example.com",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name:           "missing email",
			email:          "",
			password:       "password123",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation failed",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]string{
				"email":    tt.email,
				"password": tt.password,
			}

			// Create request
			payloadBytes, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			// Assert status code
			s.Equal(tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "token")
				s.Contains(response, "user")
				s.NotEmpty(response["token"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM users WHERE user_id = $1", userID)
}

func (s *AuthTestSuite) TestGetProfile() {
	// Create test user
	userID := s.CreateTestUser("profile@example.com", "student")

	// Setup routes with auth middleware simulation
	app := fiber.New()
	app.Get("/profile", func(c *fiber.Ctx) error {
		// Simulate JWT middleware setting user_id
		c.Locals("user_id", userID)
		return s.authHandler.GetProfile(c)
	})

	// Test successful profile retrieval
	req := httptest.NewRequest("GET", "/profile", nil)
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "email")
	s.Equal("profile@example.com", response["email"])
	s.Contains(response, "first_name")
	s.Contains(response, "last_name")
}

func (s *AuthTestSuite) TestUpdateProfile() {
	// Create test user
	userID := s.CreateTestUser("update@example.com", "student")

	// Setup routes with auth middleware simulation
	app := fiber.New()
	app.Put("/profile", func(c *fiber.Ctx) error {
		// Simulate JWT middleware setting user_id
		c.Locals("user_id", userID)
		return s.authHandler.UpdateProfile(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful update",
			payload: map[string]interface{}{
				"first_name": "Updated",
				"last_name":  "Name",
				"phone":      "0123456789",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "partial update",
			payload: map[string]interface{}{
				"first_name": "Partial",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty update",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "No fields to update",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Create request
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			// Assert status code
			s.Equal(tt.expectedStatus, resp.StatusCode)

			// Parse response
			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "message")
				s.Equal("Profile updated successfully", response["message"])
			}
		})
	}
}

// Run the test suite
func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}