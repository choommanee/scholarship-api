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

type AllocationTestSuite struct {
	tests.TestSuite
	allocationHandler *handlers.AllocationHandler
	testUserID        string
	testStudentID     string
	testAdminID       string
	testScholarshipID int
	testApplicationID int
	testSourceID      int
}

func (s *AllocationTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.allocationHandler = handlers.NewAllocationHandler(s.Config)
}

func (s *AllocationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	
	// Create test users
	s.testUserID = s.CreateTestUser("student@example.com", "student")
	s.testStudentID = s.CreateTestStudent(s.testUserID)
	s.testAdminID = s.CreateTestUser("admin@example.com", "admin")
	
	// Create test scholarship and application
	s.testScholarshipID = s.CreateTestScholarship(s.testAdminID)
	s.testApplicationID = s.createApprovedApplication()
	
	// Create budget for testing
	s.createTestBudget()
}

func (s *AllocationTestSuite) createApprovedApplication() int {
	var applicationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'approved', 'Test approved application', 50000, 15000, 2)
		RETURNING application_id
	`, s.testStudentID, s.testScholarshipID).Scan(&applicationID)
	s.Require().NoError(err)
	return applicationID
}

func (s *AllocationTestSuite) createTestBudget() {
	_, err := s.DB.Exec(`
		INSERT INTO scholarship_budgets (scholarship_id, budget_year, total_budget, allocated_budget, remaining_budget)
		VALUES ($1, '2024', 100000.00, 0.00, 100000.00)
		ON CONFLICT (scholarship_id, budget_year) DO NOTHING
	`, s.testScholarshipID)
	s.Require().NoError(err)
}

func (s *AllocationTestSuite) createTestAllocation() int {
	var allocationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_allocations (application_id, scholarship_id, allocated_amount, 
			allocation_date, disbursement_method, bank_account, bank_name, allocated_by, 
			allocation_status, notes)
		VALUES ($1, $2, 15000.00, $3, 'bank_transfer', '1234567890', 'Test Bank', $4, 
			'pending', 'Test allocation notes')
		RETURNING allocation_id
	`, s.testApplicationID, s.testScholarshipID, time.Now().Format("2006-01-02"), s.testAdminID).Scan(&allocationID)
	s.Require().NoError(err)
	return allocationID
}

func (s *AllocationTestSuite) TestCreateAllocation() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/allocations", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.CreateAllocation(c)
	})

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful allocation creation",
			payload: map[string]interface{}{
				"application_id":       s.testApplicationID,
				"scholarship_id":       s.testScholarshipID,
				"allocated_amount":     12000.0,
				"allocation_date":      time.Now().Format("2006-01-02"),
				"disbursement_method":  "bank_transfer",
				"bank_account":         "9876543210",
				"bank_name":            "Main Bank",
				"notes":                "Merit-based scholarship allocation",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"application_id":  s.testApplicationID,
				"allocated_amount": 10000.0,
				// missing scholarship_id, allocation_date
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "application not approved",
			payload: map[string]interface{}{
				"application_id":      999999, // non-existent application
				"scholarship_id":      s.testScholarshipID,
				"allocated_amount":    10000.0,
				"allocation_date":     time.Now().Format("2006-01-02"),
				"disbursement_method": "bank_transfer",
				"bank_account":        "1111111111",
				"bank_name":           "Test Bank",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Application is not approved for allocation",
		},
		{
			name: "invalid scholarship ID",
			payload: map[string]interface{}{
				"application_id":      s.testApplicationID,
				"scholarship_id":      999999,
				"allocated_amount":    10000.0,
				"allocation_date":     time.Now().Format("2006-01-02"),
				"disbursement_method": "bank_transfer",
				"bank_account":        "2222222222",
				"bank_name":           "Test Bank",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create allocation",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/allocations", bytes.NewBuffer(payload))
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
			}
		})
	}
}

func (s *AllocationTestSuite) TestGetAllocations() {
	// Create test allocations
	allocationID1 := s.createTestAllocation()
	
	// Create another allocation with different status
	var allocationID2 int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_allocations (application_id, scholarship_id, allocated_amount, 
			allocation_date, disbursement_method, bank_account, bank_name, allocated_by, 
			allocation_status, notes)
		VALUES ($1, $2, 18000.00, $3, 'check', '5555555555', 'Another Bank', $4, 
			'approved', 'Another test allocation')
		RETURNING allocation_id
	`, s.testApplicationID, s.testScholarshipID, time.Now().Format("2006-01-02"), s.testAdminID).Scan(&allocationID2)
	s.NoError(err)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Get("/allocations", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.GetAllocations(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minCount       int
	}{
		{
			name:           "get all allocations",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "get with pagination",
			queryParams:    "?page=1&limit=1",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "filter by status",
			queryParams:    "?status=pending",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "filter by scholarship",
			queryParams:    "?scholarship_id=" + strconv.Itoa(s.testScholarshipID),
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "filter by non-existent status",
			queryParams:    "?status=cancelled",
			expectedStatus: http.StatusOK,
			minCount:       0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/allocations"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "data")
			s.Contains(response, "pagination")
			
			allocations := response["data"].([]interface{})
			s.GreaterOrEqual(len(allocations), tt.minCount)

			if len(allocations) > 0 {
				allocation := allocations[0].(map[string]interface{})
				s.Contains(allocation, "allocation_id")
				s.Contains(allocation, "scholarship_name")
				s.Contains(allocation, "student_name")
				s.Contains(allocation, "allocated_amount")
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarship_allocations WHERE allocation_id IN ($1, $2)", allocationID1, allocationID2)
}

func (s *AllocationTestSuite) TestGetAllocationDetails() {
	// Create test allocation
	allocationID := s.createTestAllocation()

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Get("/allocations/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.GetAllocationDetails(c)
	})

	tests := []struct {
		name           string
		allocationID   string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful get allocation details",
			allocationID:   strconv.Itoa(allocationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent allocation",
			allocationID:   "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found",
		},
		{
			name:           "invalid allocation ID",
			allocationID:   "invalid",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/allocations/"+tt.allocationID, nil)
			
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
				s.Contains(data, "allocation_id")
				s.Contains(data, "scholarship_name")
				s.Contains(data, "student_name")
				s.Contains(data, "allocated_amount")
				s.Contains(data, "allocation_status")
				s.Contains(data, "email")
				s.Contains(data, "phone")
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarship_allocations WHERE allocation_id = $1", allocationID)
}

func (s *AllocationTestSuite) TestApproveAllocation() {
	// Create test allocation with pending status
	allocationID := s.createTestAllocation()

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/allocations/:id/approve", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.ApproveAllocation(c)
	})

	tests := []struct {
		name           string
		allocationID   string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful approval",
			allocationID:   strconv.Itoa(allocationID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "approve non-existent allocation",
			allocationID:   "999999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found or already processed",
		},
		{
			name:           "approve already processed allocation",
			allocationID:   strconv.Itoa(allocationID),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found or already processed",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("POST", "/allocations/"+tt.allocationID+"/approve", nil)
			
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
				s.Equal("Allocation approved successfully", response["message"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarship_allocations WHERE allocation_id = $1", allocationID)
}

func (s *AllocationTestSuite) TestDisburseAllocation() {
	// Create approved allocation
	var allocationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_allocations (application_id, scholarship_id, allocated_amount, 
			allocation_date, disbursement_method, bank_account, bank_name, allocated_by, 
			allocation_status, notes)
		VALUES ($1, $2, 20000.00, $3, 'bank_transfer', '6666666666', 'Disburse Bank', $4, 
			'approved', 'Approved allocation for disbursement')
		RETURNING allocation_id
	`, s.testApplicationID, s.testScholarshipID, time.Now().Format("2006-01-02"), s.testAdminID).Scan(&allocationID)
	s.NoError(err)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/allocations/:id/disburse", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.DisburseAllocation(c)
	})

	tests := []struct {
		name           string
		allocationID   string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:         "successful disbursement",
			allocationID: strconv.Itoa(allocationID),
			payload: map[string]interface{}{
				"transfer_date":      time.Now().Format("2006-01-02"),
				"transfer_reference": "TXN123456789",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "missing transfer details",
			allocationID: strconv.Itoa(allocationID),
			payload: map[string]interface{}{
				"transfer_date": time.Now().Format("2006-01-02"),
				// missing transfer_reference
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:         "disburse non-existent allocation",
			allocationID: "999999",
			payload: map[string]interface{}{
				"transfer_date":      time.Now().Format("2006-01-02"),
				"transfer_reference": "TXN987654321",
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found or not approved",
		},
		{
			name:         "disburse already disbursed allocation",
			allocationID: strconv.Itoa(allocationID),
			payload: map[string]interface{}{
				"transfer_date":      time.Now().Format("2006-01-02"),
				"transfer_reference": "TXN111222333",
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Allocation not found or not approved",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/allocations/"+tt.allocationID+"/disburse", bytes.NewBuffer(payload))
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
				s.Equal("Allocation disbursed successfully", response["message"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarship_allocations WHERE allocation_id = $1", allocationID)
}

func (s *AllocationTestSuite) TestGetBudgetSummary() {
	// Create additional test data for budget summary
	anotherScholarshipID := s.CreateTestScholarship(s.testAdminID)
	
	// Create budget for the new scholarship
	_, err := s.DB.Exec(`
		INSERT INTO scholarship_budgets (scholarship_id, budget_year, total_budget, allocated_budget, remaining_budget)
		VALUES ($1, '2024', 75000.00, 25000.00, 50000.00)
	`, anotherScholarshipID)
	s.NoError(err)

	// Create some allocations to test counts
	allocationID1 := s.createTestAllocation()
	
	var allocationID2 int
	err = s.DB.QueryRow(`
		INSERT INTO scholarship_allocations (application_id, scholarship_id, allocated_amount, 
			allocation_date, disbursement_method, bank_account, bank_name, allocated_by, 
			allocation_status, notes)
		VALUES ($1, $2, 10000.00, $3, 'bank_transfer', '7777777777', 'Budget Bank', $4, 
			'pending', 'Budget test allocation')
		RETURNING allocation_id
	`, s.testApplicationID, anotherScholarshipID, time.Now().Format("2006-01-02"), s.testAdminID).Scan(&allocationID2)
	s.NoError(err)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Get("/allocations/budget/summary", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.allocationHandler.GetBudgetSummary(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minCount       int
	}{
		{
			name:           "get all budget summaries",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "filter by scholarship",
			queryParams:    "?scholarship_id=" + strconv.Itoa(s.testScholarshipID),
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
		{
			name:           "filter by year",
			queryParams:    "?year=2024",
			expectedStatus: http.StatusOK,
			minCount:       2,
		},
		{
			name:           "filter by non-existent year",
			queryParams:    "?year=2023",
			expectedStatus: http.StatusOK,
			minCount:       0,
		},
		{
			name:           "combined filters",
			queryParams:    "?scholarship_id=" + strconv.Itoa(anotherScholarshipID) + "&year=2024",
			expectedStatus: http.StatusOK,
			minCount:       1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/allocations/budget/summary"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "data")
			budgets := response["data"].([]interface{})
			s.GreaterOrEqual(len(budgets), tt.minCount)

			if len(budgets) > 0 {
				budget := budgets[0].(map[string]interface{})
				s.Contains(budget, "scholarship_id")
				s.Contains(budget, "scholarship_name")
				s.Contains(budget, "budget_year")
				s.Contains(budget, "total_budget")
				s.Contains(budget, "allocated_budget")
				s.Contains(budget, "remaining_budget")
				s.Contains(budget, "allocation_count")
				s.Contains(budget, "utilization_rate")
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM scholarship_allocations WHERE allocation_id IN ($1, $2)", allocationID1, allocationID2)
	s.DB.Exec("DELETE FROM scholarship_budgets WHERE scholarship_id = $1", anotherScholarshipID)
	s.DB.Exec("DELETE FROM scholarships WHERE scholarship_id = $1", anotherScholarshipID)
}

// Run the test suite
func TestAllocationSuite(t *testing.T) {
	suite.Run(t, new(AllocationTestSuite))
}