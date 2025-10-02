package handlers

import (
	"strconv"
	"time"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InterviewHandler struct {
	cfg *config.Config
}

func NewInterviewHandler(cfg *config.Config) *InterviewHandler {
	return &InterviewHandler{cfg: cfg}
}

// CreateSchedule creates new interview schedule
// @Summary Create interview schedule
// @Description Create a new interview schedule for scholarship applications (Admin/Officer only)
// @Tags Interview Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param schedule body models.InterviewSchedule true "Interview schedule data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /interviews/schedules [post]
func (h *InterviewHandler) CreateSchedule(c *fiber.Ctx) error {
	var schedule models.InterviewSchedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get current user ID from context
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	schedule.CreatedBy = userID

	query := `INSERT INTO interview_schedules 
		(scholarship_id, interview_date, start_time, end_time, location, max_applicants, interviewer_ids, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING schedule_id`

	err = database.DB.QueryRow(query,
		schedule.ScholarshipID, schedule.InterviewDate, schedule.StartTime,
		schedule.EndTime, schedule.Location, schedule.MaxApplicants,
		schedule.InterviewerIDs, schedule.Notes, schedule.CreatedBy,
	).Scan(&schedule.ScheduleID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create interview schedule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Interview schedule created successfully",
		"data":    schedule,
	})
}

// GetSchedules retrieves all interview schedules
// @Summary Get interview schedules
// @Description Get list of interview schedules with optional scholarship filter (Admin/Officer only)
// @Tags Interview Management
// @Produce json
// @Security BearerAuth
// @Param scholarship_id query string false "Filter by scholarship ID"
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /interviews/schedules [get]
func (h *InterviewHandler) GetSchedules(c *fiber.Ctx) error {
	scholarshipID := c.Query("scholarship_id")

	query := `SELECT schedule_id, scholarship_id, interview_date, start_time, end_time, 
		location, max_applicants, interviewer_ids, notes, is_active, created_by, created_at
		FROM interview_schedules WHERE is_active = true`

	args := []interface{}{}
	if scholarshipID != "" {
		query += " AND scholarship_id = $1"
		args = append(args, scholarshipID)
	}

	query += " ORDER BY interview_date, start_time"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch interview schedules",
		})
	}
	defer rows.Close()

	var schedules []models.InterviewSchedule
	for rows.Next() {
		var schedule models.InterviewSchedule
		err := rows.Scan(
			&schedule.ScheduleID, &schedule.ScholarshipID, &schedule.InterviewDate,
			&schedule.StartTime, &schedule.EndTime, &schedule.Location,
			&schedule.MaxApplicants, &schedule.InterviewerIDs, &schedule.Notes,
			&schedule.IsActive, &schedule.CreatedBy, &schedule.CreatedAt,
		)
		if err != nil {
			continue
		}
		schedules = append(schedules, schedule)
	}

	return c.JSON(fiber.Map{
		"data": schedules,
	})
}

// BookInterview allows students to book interview slots
// @Summary Book interview slot
// @Description Allow student to book an interview appointment for their application
// @Tags Interview Booking
// @Produce json
// @Security BearerAuth
// @Param application_id path string true "Application ID"
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {object} object{message=string,appointment_id=int}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /interviews/applications/{application_id}/schedules/{schedule_id}/book [post]
func (h *InterviewHandler) BookInterview(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")
	scheduleID := c.Params("schedule_id")

	// Check if application belongs to current user
	userID := c.Locals("user_id").(string)

	var studentID string
	checkQuery := `SELECT s.student_id FROM scholarship_applications sa
		JOIN students s ON sa.student_id = s.student_id
		JOIN users u ON s.user_id = u.user_id
		WHERE sa.application_id = $1 AND u.user_id = $2`

	err := database.DB.QueryRow(checkQuery, applicationID, userID).Scan(&studentID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Application not found or access denied",
		})
	}

	// Insert appointment
	query := `INSERT INTO interview_appointments (application_id, schedule_id)
		VALUES ($1, $2) RETURNING appointment_id`

	var appointmentID int
	err = database.DB.QueryRow(query, applicationID, scheduleID).Scan(&appointmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to book interview",
		})
	}

	return c.JSON(fiber.Map{
		"message":        "Interview booked successfully",
		"appointment_id": appointmentID,
	})
}

// ConfirmInterview allows students to confirm their interview
// @Summary Confirm interview appointment
// @Description Allow student to confirm their scheduled interview appointment
// @Tags Interview Booking
// @Produce json
// @Security BearerAuth
// @Param appointment_id path string true "Appointment ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /interviews/appointments/{appointment_id}/confirm [post]
func (h *InterviewHandler) ConfirmInterview(c *fiber.Ctx) error {
	appointmentID := c.Params("appointment_id")
	userID := c.Locals("user_id").(string)

	// Check if appointment belongs to current user
	checkQuery := `SELECT ia.appointment_id FROM interview_appointments ia
		JOIN scholarship_applications sa ON ia.application_id = sa.application_id
		JOIN students s ON sa.student_id = s.student_id
		JOIN users u ON s.user_id = u.user_id
		WHERE ia.appointment_id = $1 AND u.user_id = $2`

	var id int
	err := database.DB.QueryRow(checkQuery, appointmentID, userID).Scan(&id)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Appointment not found or access denied",
		})
	}

	// Update confirmation
	updateQuery := `UPDATE interview_appointments 
		SET student_confirmed = true, confirmation_date = CURRENT_TIMESTAMP
		WHERE appointment_id = $1`

	_, err = database.DB.Exec(updateQuery, appointmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to confirm interview",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Interview confirmed successfully",
	})
}

// SubmitInterviewResult allows interviewers to submit results
// @Summary Submit interview result
// @Description Submit interview evaluation and scoring (Interviewer/Admin/Officer only)
// @Tags Interview Results
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param appointment_id path string true "Appointment ID"
// @Param result body models.InterviewResult true "Interview result data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /interviews/appointments/{appointment_id}/result [post]
func (h *InterviewHandler) SubmitInterviewResult(c *fiber.Ctx) error {
	appointmentID := c.Params("appointment_id")

	var result models.InterviewResult
	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get interviewer ID from context and convert to UUID
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	result.InterviewerID = userID
	appointmentIDInt, err := strconv.ParseUint(appointmentID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid appointment ID format",
		})
	}
	result.AppointmentID = uint(appointmentIDInt)

	query := `INSERT INTO interview_results 
		(appointment_id, interviewer_id, scores, overall_score, comments, recommendation, interview_notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING result_id`

	err = database.DB.QueryRow(query,
		result.AppointmentID, result.InterviewerID, result.Scores,
		result.OverallScore, result.Comments, result.Recommendation, result.InterviewNotes,
	).Scan(&result.ResultID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit interview result",
		})
	}

	// Update appointment status
	updateQuery := `UPDATE interview_appointments 
		SET appointment_status = 'completed', actual_end_time = CURRENT_TIMESTAMP
		WHERE appointment_id = $1`

	database.DB.Exec(updateQuery, appointmentID)

	return c.JSON(fiber.Map{
		"message": "Interview result submitted successfully",
		"data":    result,
	})
}

// GetMyInterviews retrieves interviews for current user (student or interviewer)
// @Summary Get my interviews
// @Description Get interview appointments for current user (different view for students vs interviewers)
// @Tags Interview Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Router /interviews/my [get]
func (h *InterviewHandler) GetMyInterviews(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("user_role").(string)

	var query string
	var args []interface{}

	if role == "student" {
		query = `SELECT ia.appointment_id, ia.application_id, ia.schedule_id, 
			ia.appointment_status, ia.student_confirmed, ia.confirmation_date,
			isch.interview_date, isch.start_time, isch.end_time, isch.location,
			s.scholarship_name
			FROM interview_appointments ia
			JOIN interview_schedules isch ON ia.schedule_id = isch.schedule_id
			JOIN scholarship_applications sa ON ia.application_id = sa.application_id
			JOIN scholarships s ON isch.scholarship_id = s.scholarship_id
			JOIN students st ON sa.student_id = st.student_id
			WHERE st.user_id = $1
			ORDER BY isch.interview_date, isch.start_time`
		args = []interface{}{userID}
	} else {
		// For interviewers
		query = `SELECT ia.appointment_id, ia.application_id, ia.schedule_id,
			ia.appointment_status, ia.student_confirmed,
			isch.interview_date, isch.start_time, isch.end_time, isch.location,
			s.scholarship_name, u.first_name, u.last_name
			FROM interview_appointments ia
			JOIN interview_schedules isch ON ia.schedule_id = isch.schedule_id
			JOIN scholarships s ON isch.scholarship_id = s.scholarship_id
			JOIN scholarship_applications sa ON ia.application_id = sa.application_id
			JOIN students st ON sa.student_id = st.student_id
			JOIN users u ON st.user_id = u.user_id
			WHERE JSON_EXTRACT(isch.interviewer_ids, '$[*]') LIKE CONCAT('%', $1, '%')
			ORDER BY isch.interview_date, isch.start_time`
		args = []interface{}{userID}
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch interviews",
		})
	}
	defer rows.Close()

	var interviews []map[string]interface{}
	for rows.Next() {
		interview := make(map[string]interface{})
		var cols []interface{}

		if role == "student" {
			var appointmentID, applicationID, scheduleID int
			var status string
			var confirmed bool
			var confirmationDate *time.Time
			var interviewDate time.Time
			var startTime, endTime, location, scholarshipName string

			cols = []interface{}{
				&appointmentID, &applicationID, &scheduleID,
				&status, &confirmed, &confirmationDate,
				&interviewDate, &startTime, &endTime, &location,
				&scholarshipName,
			}

			if err := rows.Scan(cols...); err != nil {
				continue
			}

			interview["appointment_id"] = appointmentID
			interview["application_id"] = applicationID
			interview["schedule_id"] = scheduleID
			interview["status"] = status
			interview["confirmed"] = confirmed
			interview["confirmation_date"] = confirmationDate
			interview["interview_date"] = interviewDate
			interview["start_time"] = startTime
			interview["end_time"] = endTime
			interview["location"] = location
			interview["scholarship_name"] = scholarshipName
		} else {
			var appointmentID, applicationID, scheduleID int
			var status string
			var confirmed bool
			var interviewDate time.Time
			var startTime, endTime, location, scholarshipName, firstName, lastName string

			cols = []interface{}{
				&appointmentID, &applicationID, &scheduleID,
				&status, &confirmed,
				&interviewDate, &startTime, &endTime, &location,
				&scholarshipName, &firstName, &lastName,
			}

			if err := rows.Scan(cols...); err != nil {
				continue
			}

			interview["appointment_id"] = appointmentID
			interview["application_id"] = applicationID
			interview["schedule_id"] = scheduleID
			interview["status"] = status
			interview["confirmed"] = confirmed
			interview["interview_date"] = interviewDate
			interview["start_time"] = startTime
			interview["end_time"] = endTime
			interview["location"] = location
			interview["scholarship_name"] = scholarshipName
			interview["student_name"] = firstName + " " + lastName
		}

		interviews = append(interviews, interview)
	}

	return c.JSON(fiber.Map{
		"data": interviews,
	})
}
