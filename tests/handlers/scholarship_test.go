package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"

	"scholarship-system/internal/handlers"
	"scholarship-system/tests"
)

type ScholarshipTestSuite struct {
	tests.TestSuite
	scholarshipHandler *handlers.ScholarshipHandler
	testUserID         string
	testAdminID        string
	testSourceID       int
}

func (s *ScholarshipTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.scholarshipHandler = handlers.NewScholarshipHandler(s.Config)
}

func (s *ScholarshipTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	
	// Create test users
	s.testUserID = s.CreateTestUser("user@example.com", "student")
	s.testAdminID = s.CreateTestUser("admin@example.com", "admin")
	
	// Create test source
	s.testSourceID = s.createTestSource()
}

func (s *ScholarshipTestSuite) createTestSource() int {
	var sourceID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_sources (source_name, source_type, contact_person, contact_email)
		VALUES ('Test Foundation', 'foundation', 'Test Contact', 'test@foundation.com')
		RETURNING source_id
	`).Scan(&sourceID)
	s.Require().NoError(err)
	return sourceID
}

func (s *ScholarshipTestSuite) TestCreateSource() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/scholarship-sources", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.scholarshipHandler.CreateSource(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful source creation",
			payload: map[string]interface{}{
				"source_name":    "New Foundation",
				"source_type":    "foundation",
				"contact_person": "John Doe",
				"contact_email":  "john@newfoundation.com",
				"contact_phone":  "0123456789",
				"description":    "Test foundation for scholarship",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"source_name": "Incomplete Foundation",
				// missing source_type
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "minimal valid source",
			payload: map[string]interface{}{
				"source_name": "Minimal Foundation",
				"source_type": "government",
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/scholarship-sources", bytes.NewBuffer(payload))
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
				s.Contains(response, "source")
			}
		})
	}
}

func (s *ScholarshipTestSuite) TestGetSources() {
	// Create additional test sources
	s.DB.Exec(`
		INSERT INTO scholarship_sources (source_name, source_type, contact_person, contact_email)
		VALUES 
			('Government Agency', 'government', 'Jane Smith', 'jane@gov.com'),
			('Private Company', 'private', 'Bob Johnson', 'bob@company.com')
	`)

	// Setup routes
	app := fiber.New()
	app.Get("/scholarship-sources", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.scholarshipHandler.GetSources(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minCount       int
	}{
		{
			name:           "get all sources",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			minCount:       3, // testSourceID + 2 additional
		},
		{
			name:           "get with pagination",
			queryParams:    "?limit=2&offset=0",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "search sources",
			queryParams:    "?search=Government",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "high offset",
			queryParams:    "?limit=5&offset=10",
			expectedStatus: http.StatusOK,
			minCount:       0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/scholarship-sources"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "sources")
			sources := response["sources"].([]interface{})
			s.GreaterOrEqual(len(sources), tt.minCount)
			s.Contains(response, "total")
		})
	}
}

func (s *ScholarshipTestSuite) TestCreateScholarship() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/scholarships", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.scholarshipHandler.CreateScholarship(c)
	})

	startDate := time.Now().AddDate(0, 1, 0)
	endDate := time.Now().AddDate(0, 6, 0)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful scholarship creation",
			payload: map[string]interface{}{
				"source_id":              s.testSourceID,
				"scholarship_name":       "Excellence Scholarship",
				"scholarship_type":       "merit",
				"amount":                 25000.0,
				"total_quota":            5,
				"academic_year":          "2024",
				"semester":               "1",
				"eligibility_criteria":   "GPA >= 3.5",
				"required_documents":     "Transcript, Essay",
				"application_start_date": startDate.Format(time.RFC3339),
				"application_end_date":   endDate.Format(time.RFC3339),
				"interview_required":     true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"scholarship_name": "Incomplete Scholarship",
				"amount":           10000.0,
				// missing source_id, scholarship_type, etc.
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "invalid date range",
			payload: map[string]interface{}{
				"source_id":              s.testSourceID,
				"scholarship_name":       "Invalid Date Scholarship",
				"scholarship_type":       "need",
				"amount":                 15000.0,
				"total_quota":            3,
				"academic_year":          "2024",
				"application_start_date": endDate.Format(time.RFC3339),
				"application_end_date":   startDate.Format(time.RFC3339), // End before start
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Application end date must be after start date",
		},
		{
			name: "invalid amount",
			payload: map[string]interface{}{
				"source_id":        s.testSourceID,
				"scholarship_name": "Zero Amount Scholarship",
				"scholarship_type": "merit",
				"amount":           -1000.0,
				"total_quota":      1,
				"academic_year":    "2024",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/scholarships", bytes.NewBuffer(payload))
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
				s.Contains(response, "scholarship")
			}
		})
	}
}

func (s *ScholarshipTestSuite) TestGetScholarships() {
	// Create test scholarships
	scholarshipID1 := s.CreateTestScholarship(s.testAdminID)
	
	// Create another scholarship with different type
	var scholarshipID2 int
	err := s.DB.QueryRow(`
		INSERT INTO scholarships (source_id, scholarship_name, scholarship_type, amount, 
			total_quota, available_quota, academic_year, application_start_date, 
			application_end_date, is_active, created_by)
		VALUES ($1, 'Merit Scholarship', 'merit', 20000.00, 3, 3, '2024', 
			'2024-01-01', '2024-12-31', true, $2)
		RETURNING scholarship_id
	`, s.testSourceID, s.testAdminID).Scan(&scholarshipID2)
	s.NoError(err)

	// Setup routes
	app := fiber.New()
	app.Get("/scholarships", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.scholarshipHandler.GetScholarships(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minCount       int
	}{
		{
			name:           "get all scholarships",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "get with pagination",
			queryParams:    "?limit=1&offset=0",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "filter by type",
			queryParams:    "?type=merit",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "filter by academic year",
			queryParams:    "?academic_year=2024",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "search scholarships",
			queryParams:    "?search=Merit",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/scholarships"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "scholarships")
			scholarships := response["scholarships"].([]interface{})
			s.GreaterOrEqual(len(scholarships), tt.minCount)
			s.Contains(response, "total")
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarships WHERE scholarship_id IN ($1, $2)", scholarshipID1, scholarshipID2)
}

func (s *ScholarshipTestSuite) TestGetScholarship() {
	// Create test scholarship
	scholarshipID := s.CreateTestScholarship(s.testAdminID)

	// Setup routes
	app := fiber.New()
	app.Get("/scholarships/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.scholarshipHandler.GetScholarship(c)
	})

	tests := []struct {
		name           string
		scholarshipID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful get scholarship",
			scholarshipID:  strconv.Itoa(scholarshipID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent scholarship",
			scholarshipID:  "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Scholarship not found",
		},
		{
			name:           "invalid scholarship ID",
			scholarshipID:  "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid scholarship ID",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/scholarships/"+tt.scholarshipID, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			if tt.expectedError != "" {
				s.Contains(response["error"], tt.expectedError)
			} else {
				s.Contains(response, "scholarship_id")
				s.Contains(response, "scholarship_name")
			}
		})
	}
}

func (s *ScholarshipTestSuite) TestUpdateScholarship() {
	// Create test scholarship
	scholarshipID := s.CreateTestScholarship(s.testAdminID)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Put("/scholarships/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.scholarshipHandler.UpdateScholarship(c)
	})

	startDate := time.Now().AddDate(0, 2, 0)
	endDate := time.Now().AddDate(0, 7, 0)

	tests := []struct {
		name           string
		scholarshipID  string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful update",
			scholarshipID: strconv.Itoa(scholarshipID),
			payload: map[string]interface{}{
				"source_id":              s.testSourceID,
				"scholarship_name":       "Updated Scholarship",
				"scholarship_type":       "merit",
				"amount":                 30000.0,
				"total_quota":            8,
				"academic_year":          "2024",
				"semester":               "2",
				"eligibility_criteria":   "GPA >= 3.8",
				"required_documents":     "Updated documents",
				"application_start_date": startDate.Format(time.RFC3339),
				"application_end_date":   endDate.Format(time.RFC3339),
				"interview_required":     false,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "update non-existent scholarship",
			scholarshipID:  "999999",
			payload:        map[string]interface{}{"scholarship_name": "Updated"},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Scholarship not found",
		},
		{
			name:          "invalid date range",
			scholarshipID: strconv.Itoa(scholarshipID),
			payload: map[string]interface{}{
				"source_id":              s.testSourceID,
				"scholarship_name":       "Invalid Update",
				"scholarship_type":       "merit",
				"amount":                 25000.0,
				"total_quota":            5,
				"academic_year":          "2024",
				"application_start_date": endDate.Format(time.RFC3339),
				"application_end_date":   startDate.Format(time.RFC3339),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Application end date must be after start date",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("PUT", "/scholarships/"+tt.scholarshipID, bytes.NewBuffer(payload))
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
				s.Contains(response, "scholarship")
			}
		})
	}
}

func (s *ScholarshipTestSuite) TestGetAvailableScholarships() {
	// Create test scholarships - one active, one inactive
	activeScholarshipID := s.CreateTestScholarship(s.testAdminID)
	
	var inactiveScholarshipID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarships (source_id, scholarship_name, scholarship_type, amount, 
			total_quota, available_quota, academic_year, application_start_date, 
			application_end_date, is_active, created_by)
		VALUES ($1, 'Inactive Scholarship', 'need', 15000.00, 2, 2, '2024', 
			'2024-01-01', '2024-12-31', false, $2)
		RETURNING scholarship_id
	`, s.testSourceID, s.testAdminID).Scan(&inactiveScholarshipID)
	s.NoError(err)

	// Setup routes
	app := fiber.New()
	app.Get("/scholarships/available", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.scholarshipHandler.GetAvailableScholarships(c)
	})

	req := httptest.NewRequest("GET", "/scholarships/available", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "scholarships")
	s.Contains(response, "count")
	
	scholarships := response["scholarships"].([]interface{})
	s.GreaterOrEqual(len(scholarships), 1) // At least the active scholarship

	// Cleanup
	s.DB.Exec("DELETE FROM scholarships WHERE scholarship_id IN ($1, $2)", activeScholarshipID, inactiveScholarshipID)
}

func (s *ScholarshipTestSuite) TestToggleScholarshipStatus() {
	// Create test scholarship
	scholarshipID := s.CreateTestScholarship(s.testAdminID)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/scholarships/:id/toggle", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.scholarshipHandler.ToggleScholarshipStatus(c)
	})

	tests := []struct {
		name           string
		scholarshipID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful toggle",
			scholarshipID:  strconv.Itoa(scholarshipID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "toggle non-existent scholarship",
			scholarshipID:  "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Scholarship not found",
		},
		{
			name:           "invalid scholarship ID",
			scholarshipID:  "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid scholarship ID",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("POST", "/scholarships/"+tt.scholarshipID+"/toggle", nil)
			
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
				s.Contains(response, "scholarship")
			}
		})
	}
}

// Run the test suite
func TestScholarshipSuite(t *testing.T) {
	suite.Run(t, new(ScholarshipTestSuite))
}