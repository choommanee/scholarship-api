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

type ApplicationTestSuite struct {
	tests.TestSuite
	applicationHandler *handlers.ApplicationHandler
	testUserID         string
	testStudentID      string
	testScholarshipID  int
}

func (s *ApplicationTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.applicationHandler = handlers.NewApplicationHandler(s.Config)
}

func (s *ApplicationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	
	// Create test user and student
	s.testUserID = s.CreateTestUser("student@example.com", "student")
	s.testStudentID = s.CreateTestStudent(s.testUserID)
	s.testScholarshipID = s.CreateTestScholarship(s.testUserID)
}

func (s *ApplicationTestSuite) TestCreateApplication() {
	// Setup routes with auth middleware simulation
	app := fiber.New()
	app.Post("/applications", func(c *fiber.Ctx) error {
		// Simulate JWT middleware
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.applicationHandler.CreateApplication(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful application",
			payload: map[string]interface{}{
				"scholarship_id":           s.testScholarshipID,
				"application_data":         "Test application data",
				"family_income":            50000.0,
				"monthly_expenses":         15000.0,
				"siblings_count":           2,
				"activities_participation": "Student council, sports",
				"special_abilities":        "Leadership, programming",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required scholarship_id",
			payload: map[string]interface{}{
				"application_data": "Test application data",
				"family_income":    50000.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Scholarship ID is required",
		},
		{
			name: "invalid scholarship_id",
			payload: map[string]interface{}{
				"scholarship_id":   999999,
				"application_data": "Test application data",
				"family_income":    50000.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid scholarship",
		},
		{
			name: "duplicate application",
			payload: map[string]interface{}{
				"scholarship_id":           s.testScholarshipID,
				"application_data":         "Duplicate application",
				"family_income":            50000.0,
				"monthly_expenses":         15000.0,
				"siblings_count":           1,
				"activities_participation": "Activities",
				"special_abilities":        "Abilities",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "already applied",
		},
	}

	for i, tt := range tests {
		s.Run(tt.name, func() {
			// Skip duplicate test on first iteration
			if i == 3 && tt.name == "duplicate application" {
				// First create a successful application
				firstPayload := tests[0].payload
				payloadBytes, _ := json.Marshal(firstPayload)
				req := httptest.NewRequest("POST", "/applications", bytes.NewBuffer(payloadBytes))
				req.Header.Set("Content-Type", "application/json")
				resp, _ := app.Test(req)
				resp.Body.Close()
			}

			// Create request
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/applications", bytes.NewBuffer(payload))
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
				s.Contains(response, "data")
			}
		})
	}
}

func (s *ApplicationTestSuite) TestGetMyApplications() {
	// Create test application first
	_, err := s.DB.Exec(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'draft', 'Test application', 50000, 15000, 2)
	`, s.testStudentID, s.testScholarshipID)
	s.NoError(err)

	// Setup routes with auth middleware simulation
	app := fiber.New()
	app.Get("/applications/my", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.applicationHandler.GetMyApplications(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "get all applications",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get with pagination",
			queryParams:    "?limit=5&offset=0",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get with high offset",
			queryParams:    "?limit=5&offset=10",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/applications/my"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "data")
			data := response["data"].([]interface{})
			s.Equal(tt.expectedCount, len(data))

			if len(data) > 0 {
				application := data[0].(map[string]interface{})
				s.Contains(application, "application_id")
				s.Contains(application, "scholarship_name")
				s.Contains(application, "application_status")
			}
		})
	}
}

func (s *ApplicationTestSuite) TestGetApplication() {
	// Create test application
	var applicationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'submitted', 'Test application detail', 60000, 18000, 1)
		RETURNING application_id
	`, s.testStudentID, s.testScholarshipID).Scan(&applicationID)
	s.NoError(err)

	// Setup routes
	app := fiber.New()
	app.Get("/applications/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.applicationHandler.GetApplication(c)
	})

	tests := []struct {
		name           string
		applicationID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful get application",
			applicationID:  strconv.Itoa(applicationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent application",
			applicationID:  "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Application not found",
		},
		{
			name:           "invalid application ID",
			applicationID:  "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid application ID",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/applications/"+tt.applicationID, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "data")
				data := response["data"].(map[string]interface{})
				s.Contains(data, "application_id")
				s.Contains(data, "application_status")
				s.Equal("submitted", data["application_status"])
			}
		})
	}
}

func (s *ApplicationTestSuite) TestUpdateApplication() {
	// Create test application in draft status
	var applicationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'draft', 'Original data', 50000, 15000, 2)
		RETURNING application_id
	`, s.testStudentID, s.testScholarshipID).Scan(&applicationID)
	s.NoError(err)

	// Setup routes
	app := fiber.New()
	app.Put("/applications/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.applicationHandler.UpdateApplication(c)
	})

	tests := []struct {
		name           string
		applicationID  string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful update",
			applicationID: strconv.Itoa(applicationID),
			payload: map[string]interface{}{
				"application_data":         "Updated application data",
				"family_income":            55000.0,
				"monthly_expenses":         16000.0,
				"siblings_count":           3,
				"activities_participation": "Updated activities",
				"special_abilities":        "Updated abilities",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "partial update",
			applicationID: strconv.Itoa(applicationID),
			payload: map[string]interface{}{
				"family_income": 60000.0,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "update non-existent application",
			applicationID:  "999999",
			payload:        map[string]interface{}{"family_income": 60000.0},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Application not found",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("PUT", "/applications/"+tt.applicationID, bytes.NewBuffer(payload))
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
			}
		})
	}
}

func (s *ApplicationTestSuite) TestSubmitApplication() {
	// Create test application in draft status
	var applicationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'draft', 'Complete application', 50000, 15000, 2)
		RETURNING application_id
	`, s.testStudentID, s.testScholarshipID).Scan(&applicationID)
	s.NoError(err)

	// Setup routes
	app := fiber.New()
	app.Post("/applications/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.applicationHandler.SubmitApplication(c)
	})

	tests := []struct {
		name           string
		applicationID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful submission",
			applicationID:  strconv.Itoa(applicationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "submit non-existent application",
			applicationID:  "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Application not found",
		},
		{
			name:           "submit already submitted application",
			applicationID:  strconv.Itoa(applicationID),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "already submitted",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("POST", "/applications/"+tt.applicationID+"/submit", nil)

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
				s.Equal("Application submitted successfully", response["message"])
			}
		})
	}
}

// Run the test suite
func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationTestSuite))
}