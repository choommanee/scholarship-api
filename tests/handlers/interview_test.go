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

type InterviewTestSuite struct {
	tests.TestSuite
	interviewHandler  *handlers.InterviewHandler
	testUserID        string
	testStudentID     string
	testInterviewerID string
	testAdminID       string
	testScholarshipID int
	testApplicationID int
}

func (s *InterviewTestSuite) SetupSuite() {
	s.TestSuite.SetupSuite()
	s.interviewHandler = handlers.NewInterviewHandler(s.Config)
}

func (s *InterviewTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	
	// Create test users
	s.testUserID = s.CreateTestUser("student@example.com", "student")
	s.testStudentID = s.CreateTestStudent(s.testUserID)
	s.testInterviewerID = s.CreateTestUser("interviewer@example.com", "interviewer")
	s.testAdminID = s.CreateTestUser("admin@example.com", "admin")
	
	// Create test scholarship and application
	s.testScholarshipID = s.CreateTestScholarship(s.testAdminID)
	s.testApplicationID = s.createTestApplication()
}

func (s *InterviewTestSuite) createTestApplication() int {
	var applicationID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, 
			application_data, family_income, monthly_expenses, siblings_count)
		VALUES ($1, $2, 'submitted', 'Test application for interview', 50000, 15000, 2)
		RETURNING application_id
	`, s.testStudentID, s.testScholarshipID).Scan(&applicationID)
	s.Require().NoError(err)
	return applicationID
}

func (s *InterviewTestSuite) createTestSchedule() int {
	var scheduleID int
	interviewDate := time.Now().AddDate(0, 0, 7) // 7 days from now
	err := s.DB.QueryRow(`
		INSERT INTO interview_schedules (scholarship_id, interview_date, start_time, end_time, 
			location, max_applicants, interviewer_ids, notes, created_by, is_active)
		VALUES ($1, $2, '09:00', '17:00', 'Interview Room A', 10, 
			'["` + s.testInterviewerID + `"]', 'Test interview schedule', $3, true)
		RETURNING schedule_id
	`, s.testScholarshipID, interviewDate.Format("2006-01-02"), s.testAdminID).Scan(&scheduleID)
	s.Require().NoError(err)
	return scheduleID
}

func (s *InterviewTestSuite) TestCreateSchedule() {
	// Setup routes with admin role simulation
	app := fiber.New()
	app.Post("/interviews/schedules", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.interviewHandler.CreateSchedule(c)
	})

	interviewDate := time.Now().AddDate(0, 0, 14) // 14 days from now

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful schedule creation",
			payload: map[string]interface{}{
				"scholarship_id":  s.testScholarshipID,
				"interview_date":  interviewDate.Format("2006-01-02"),
				"start_time":      "09:00",
				"end_time":        "17:00",
				"location":        "Interview Room B",
				"max_applicants":  8,
				"interviewer_ids": `["` + s.testInterviewerID + `"]`,
				"notes":           "Comprehensive interview session",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing required fields",
			payload: map[string]interface{}{
				"scholarship_id": s.testScholarshipID,
				"interview_date": interviewDate.Format("2006-01-02"),
				// missing start_time, end_time
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "minimal valid schedule",
			payload: map[string]interface{}{
				"scholarship_id": s.testScholarshipID,
				"interview_date": interviewDate.Format("2006-01-02"),
				"start_time":     "10:00",
				"end_time":       "16:00",
				"location":       "Main Hall",
				"max_applicants": 5,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/interviews/schedules", bytes.NewBuffer(payload))
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

func (s *InterviewTestSuite) TestGetSchedules() {
	// Create test schedules
	scheduleID1 := s.createTestSchedule()
	
	// Create another schedule for different scholarship
	anotherScholarshipID := s.CreateTestScholarship(s.testAdminID)
	var scheduleID2 int
	interviewDate := time.Now().AddDate(0, 0, 10)
	err := s.DB.QueryRow(`
		INSERT INTO interview_schedules (scholarship_id, interview_date, start_time, end_time, 
			location, max_applicants, interviewer_ids, notes, created_by, is_active)
		VALUES ($1, $2, '08:00', '16:00', 'Interview Room C', 6, 
			'["` + s.testInterviewerID + `"]', 'Another test schedule', $3, true)
		RETURNING schedule_id
	`, anotherScholarshipID, interviewDate.Format("2006-01-02"), s.testAdminID).Scan(&scheduleID2)
	s.NoError(err)

	// Setup routes with admin role simulation
	app := fiber.New()
	app.Get("/interviews/schedules", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testAdminID)
		c.Locals("user_role", "admin")
		return s.interviewHandler.GetSchedules(c)
	})

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		minCount       int
	}{
		{
			name:           "get all schedules",
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
			name:           "filter by non-existent scholarship",
			queryParams:    "?scholarship_id=999999",
			expectedStatus: http.StatusOK,
			minCount:       0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", "/interviews/schedules"+tt.queryParams, nil)
			
			resp, err := app.Test(req)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&response)

			s.Contains(response, "data")
			schedules := response["data"].([]interface{})
			s.GreaterOrEqual(len(schedules), tt.minCount)

			if len(schedules) > 0 {
				schedule := schedules[0].(map[string]interface{})
				s.Contains(schedule, "schedule_id")
				s.Contains(schedule, "scholarship_id")
				s.Contains(schedule, "interview_date")
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id IN ($1, $2)", scheduleID1, scheduleID2)
	s.DB.Exec("DELETE FROM scholarships WHERE scholarship_id = $1", anotherScholarshipID)
}

func (s *InterviewTestSuite) TestBookInterview() {
	// Create test schedule
	scheduleID := s.createTestSchedule()

	// Setup routes with student role simulation
	app := fiber.New()
	app.Post("/interviews/applications/:application_id/schedules/:schedule_id/book", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.interviewHandler.BookInterview(c)
	})

	tests := []struct {
		name           string
		applicationID  string
		scheduleID     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful booking",
			applicationID:  strconv.Itoa(s.testApplicationID),
			scheduleID:     strconv.Itoa(scheduleID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "book with non-existent application",
			applicationID:  "999999",
			scheduleID:     strconv.Itoa(scheduleID),
			expectedStatus: http.StatusForbidden,
			expectedError:  "Application not found or access denied",
		},
		{
			name:           "book with non-existent schedule",
			applicationID:  strconv.Itoa(s.testApplicationID),
			scheduleID:     "999999",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to book interview",
		},
		{
			name:           "duplicate booking",
			applicationID:  strconv.Itoa(s.testApplicationID),
			scheduleID:     strconv.Itoa(scheduleID),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to book interview",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			url := "/interviews/applications/" + tt.applicationID + "/schedules/" + tt.scheduleID + "/book"
			req := httptest.NewRequest("POST", url, nil)
			
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
				s.Contains(response, "appointment_id")
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id = $1", scheduleID)
}

func (s *InterviewTestSuite) TestConfirmInterview() {
	// Create test schedule and appointment
	scheduleID := s.createTestSchedule()
	
	var appointmentID int
	err := s.DB.QueryRow(`
		INSERT INTO interview_appointments (application_id, schedule_id, appointment_status, student_confirmed)
		VALUES ($1, $2, 'scheduled', false)
		RETURNING appointment_id
	`, s.testApplicationID, scheduleID).Scan(&appointmentID)
	s.NoError(err)

	// Setup routes with student role simulation
	app := fiber.New()
	app.Post("/interviews/appointments/:appointment_id/confirm", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.interviewHandler.ConfirmInterview(c)
	})

	tests := []struct {
		name           string
		appointmentID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful confirmation",
			appointmentID:  strconv.Itoa(appointmentID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "confirm non-existent appointment",
			appointmentID:  "999999",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Appointment not found or access denied",
		},
		{
			name:           "invalid appointment ID",
			appointmentID:  "invalid",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Appointment not found or access denied",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("POST", "/interviews/appointments/"+tt.appointmentID+"/confirm", nil)
			
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
				s.Equal("Interview confirmed successfully", response["message"])
			}
		})
	}

	// Cleanup
	s.DB.Exec("DELETE FROM interview_appointments WHERE appointment_id = $1", appointmentID)
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id = $1", scheduleID)
}

func (s *InterviewTestSuite) TestSubmitInterviewResult() {
	// Create test schedule and appointment
	scheduleID := s.createTestSchedule()
	
	var appointmentID int
	err := s.DB.QueryRow(`
		INSERT INTO interview_appointments (application_id, schedule_id, appointment_status, student_confirmed)
		VALUES ($1, $2, 'scheduled', true)
		RETURNING appointment_id
	`, s.testApplicationID, scheduleID).Scan(&appointmentID)
	s.NoError(err)

	// Setup routes with interviewer role simulation
	app := fiber.New()
	app.Post("/interviews/appointments/:appointment_id/result", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testInterviewerID)
		c.Locals("user_role", "interviewer")
		return s.interviewHandler.SubmitInterviewResult(c)
	})

	tests := []struct {
		name           string
		appointmentID  string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful result submission",
			appointmentID: strconv.Itoa(appointmentID),
			payload: map[string]interface{}{
				"scores":          `{"communication": 8, "knowledge": 9, "attitude": 8}`,
				"overall_score":   8.3,
				"comments":        "Excellent candidate with strong communication skills",
				"recommendation":  "recommend",
				"interview_notes": "Well-prepared, answered all questions confidently",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "missing required fields",
			appointmentID: strconv.Itoa(appointmentID),
			payload: map[string]interface{}{
				"overall_score": 7.5,
				// missing other required fields
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "submit for non-existent appointment",
			appointmentID:  "999999",
			payload:        map[string]interface{}{"overall_score": 8.0, "recommendation": "recommend"},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to submit interview result",
		},
		{
			name:           "invalid appointment ID",
			appointmentID:  "invalid",
			payload:        map[string]interface{}{"overall_score": 8.0, "recommendation": "recommend"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid appointment ID format",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/interviews/appointments/"+tt.appointmentID+"/result", bytes.NewBuffer(payload))
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

	// Cleanup
	s.DB.Exec("DELETE FROM interview_results WHERE appointment_id = $1", appointmentID)
	s.DB.Exec("DELETE FROM interview_appointments WHERE appointment_id = $1", appointmentID)
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id = $1", scheduleID)
}

func (s *InterviewTestSuite) TestGetMyInterviewsAsStudent() {
	// Create test schedule and appointment for student
	scheduleID := s.createTestSchedule()
	
	var appointmentID int
	err := s.DB.QueryRow(`
		INSERT INTO interview_appointments (application_id, schedule_id, appointment_status, student_confirmed)
		VALUES ($1, $2, 'scheduled', true)
		RETURNING appointment_id
	`, s.testApplicationID, scheduleID).Scan(&appointmentID)
	s.NoError(err)

	// Setup routes with student role simulation
	app := fiber.New()
	app.Get("/interviews/my", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testUserID)
		c.Locals("user_role", "student")
		return s.interviewHandler.GetMyInterviews(c)
	})

	req := httptest.NewRequest("GET", "/interviews/my", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "data")
	interviews := response["data"].([]interface{})
	s.GreaterOrEqual(len(interviews), 1)

	if len(interviews) > 0 {
		interview := interviews[0].(map[string]interface{})
		s.Contains(interview, "appointment_id")
		s.Contains(interview, "application_id")
		s.Contains(interview, "status")
		s.Contains(interview, "confirmed")
		s.Contains(interview, "scholarship_name")
		s.Contains(interview, "location")
	}

	// Cleanup
	s.DB.Exec("DELETE FROM interview_appointments WHERE appointment_id = $1", appointmentID)
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id = $1", scheduleID)
}

func (s *InterviewTestSuite) TestGetMyInterviewsAsInterviewer() {
	// Create test schedule and appointment for interviewer
	scheduleID := s.createTestSchedule()
	
	var appointmentID int
	err := s.DB.QueryRow(`
		INSERT INTO interview_appointments (application_id, schedule_id, appointment_status, student_confirmed)
		VALUES ($1, $2, 'scheduled', true)
		RETURNING appointment_id
	`, s.testApplicationID, scheduleID).Scan(&appointmentID)
	s.NoError(err)

	// Setup routes with interviewer role simulation
	app := fiber.New()
	app.Get("/interviews/my", func(c *fiber.Ctx) error {
		c.Locals("user_id", s.testInterviewerID)
		c.Locals("user_role", "interviewer")
		return s.interviewHandler.GetMyInterviews(c)
	})

	req := httptest.NewRequest("GET", "/interviews/my", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "data")
	interviews := response["data"].([]interface{})
	s.GreaterOrEqual(len(interviews), 1)

	if len(interviews) > 0 {
		interview := interviews[0].(map[string]interface{})
		s.Contains(interview, "appointment_id")
		s.Contains(interview, "application_id")
		s.Contains(interview, "status")
		s.Contains(interview, "confirmed")
		s.Contains(interview, "scholarship_name")
		s.Contains(interview, "student_name")
	}

	// Cleanup
	s.DB.Exec("DELETE FROM interview_appointments WHERE appointment_id = $1", appointmentID)
	s.DB.Exec("DELETE FROM interview_schedules WHERE schedule_id = $1", scheduleID)
}

func (s *InterviewTestSuite) TestGetMyInterviewsNoAppointments() {
	// Create a new student with no appointments
	newUserID := s.CreateTestUser("newstudent@example.com", "student")
	
	// Setup routes with new student role simulation
	app := fiber.New()
	app.Get("/interviews/my", func(c *fiber.Ctx) error {
		c.Locals("user_id", newUserID)
		c.Locals("user_role", "student")
		return s.interviewHandler.GetMyInterviews(c)
	})

	req := httptest.NewRequest("GET", "/interviews/my", nil)
	
	resp, err := app.Test(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	s.Contains(response, "data")
	interviews := response["data"].([]interface{})
	s.Equal(0, len(interviews))
}

// Run the test suite
func TestInterviewSuite(t *testing.T) {
	suite.Run(t, new(InterviewTestSuite))
}