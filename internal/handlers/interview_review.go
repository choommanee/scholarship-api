package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InterviewReviewHandler struct {
	db *sql.DB
}

func NewInterviewReviewHandler(cfg *config.Config) *InterviewReviewHandler {
	return &InterviewReviewHandler{db: database.DB}
}

// ==================== INTERVIEW SLOT MANAGEMENT ====================

// @Summary Create Interview Slot
// @Description Create a new interview slot for a scholarship
// @Tags Interview Management
// @Accept json
// @Produce json
// @Param request body models.CreateInterviewSlotRequest true "Interview slot details"
// @Success 201 {object} models.InterviewSlotResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/slots [post]
func (h *InterviewReviewHandler) CreateInterviewSlot(c *fiber.Ctx) error {
	var req models.CreateInterviewSlotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Validate input
	if req.ScholarshipID <= 0 || req.InterviewerID == "" || req.InterviewDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณากรอกข้อมูลให้ครบถ้วน",
		})
	}

	// Parse and validate date/time
	interviewDate, err := time.Parse("2006-01-02", req.InterviewDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบวันที่ไม่ถูกต้อง",
		})
	}

	// Check if date is in the future
	if interviewDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถสร้างช่วงเวลาในอดีตได้",
		})
	}

	// Validate time format
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบเวลาไม่ถูกต้อง (ใช้ HH:MM)",
		})
	}

	// Check time sequence
	startTime, _ := time.Parse("15:04", req.StartTime)
	endTime, _ := time.Parse("15:04", req.EndTime)
	if !endTime.After(startTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เวลาสิ้นสุดต้องหลังจากเวลาเริ่มต้น",
		})
	}

	// Get user from context
	userIDValue := c.Locals("user_id")
	userUUID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถยืนยันตัวผู้ใช้ได้",
		})
	}

	// Set defaults
	if req.MaxCapacity <= 0 {
		req.MaxCapacity = 1
	}
	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 30
	}
	if req.SlotType == "" {
		req.SlotType = "individual"
	}

	// Check for conflicting slots
	var conflictCount int
	checkConflictQuery := `
		SELECT COUNT(*) FROM interview_slots 
		WHERE interviewer_id = $1 
		AND interview_date = $2 
		AND (
			(start_time <= $3 AND end_time > $3) OR
			(start_time < $4 AND end_time >= $4) OR
			(start_time >= $3 AND end_time <= $4)
		)
	`
	err = h.db.QueryRow(checkConflictQuery,
		req.InterviewerID,
		req.InterviewDate,
		req.StartTime,
		req.EndTime,
	).Scan(&conflictCount)

	if err != nil {
		log.Printf("Error checking conflict: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการตรวจสอบข้อมูล",
		})
	}

	if conflictCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "มีช่วงเวลาที่ซ้อนทับกันอยู่แล้ว",
		})
	}

	// Insert new slot
	insertQuery := `
		INSERT INTO interview_slots (
			scholarship_id, interviewer_id, interview_date, start_time, end_time,
			location, building, floor, room, max_capacity, slot_type,
			duration_minutes, preparation_time, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	var slotID int
	err = h.db.QueryRow(insertQuery,
		req.ScholarshipID,
		req.InterviewerID,
		req.InterviewDate,
		req.StartTime,
		req.EndTime,
		req.Location,
		req.Building,
		req.Floor,
		req.Room,
		req.MaxCapacity,
		req.SlotType,
		req.DurationMinutes,
		req.PreparationTime,
		req.Notes,
		userUUID,
	).Scan(&slotID)

	if err != nil {
		log.Printf("Error creating interview slot: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการสร้างช่วงเวลาสัมภาษณ์",
		})
	}

	// Fetch created slot
	slot, err := h.getInterviewSlotByID(int(slotID))
	if err != nil {
		log.Printf("Error fetching created slot: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "สร้างช่วงเวลาสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.InterviewSlotResponse{
		Success: true,
		Message: "สร้างช่วงเวลาสัมภาษณ์สำเร็จ",
		Data:    *slot,
	})
}

// @Summary Get Interview Slots
// @Description Get list of interview slots with filters
// @Tags Interview Management
// @Accept json
// @Produce json
// @Param scholarship_id query int false "Scholarship ID"
// @Param interviewer_id query string false "Interviewer ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param is_available query bool false "Only available slots"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.InterviewSlotsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/slots [get]
func (h *InterviewReviewHandler) GetInterviewSlots(c *fiber.Ctx) error {
	// Parse query parameters
	filters := models.InterviewSlotFilters{}

	if scholarshipID := c.Query("scholarship_id"); scholarshipID != "" {
		if id, err := strconv.Atoi(scholarshipID); err == nil {
			filters.ScholarshipID = &id
		}
	}

	if interviewerID := c.Query("interviewer_id"); interviewerID != "" {
		filters.InterviewerID = &interviewerID
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters.DateTo = &dateTo
	}

	if isAvailable := c.Query("is_available"); isAvailable != "" {
		if available, err := strconv.ParseBool(isAvailable); err == nil {
			filters.IsAvailable = &available
		}
	}

	if slotType := c.Query("slot_type"); slotType != "" {
		filters.SlotType = &slotType
	}

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Build query
	whereConditions := []string{}
	args := []interface{}{}

	addCondition := func(expression string, value interface{}) {
		placeholder := fmt.Sprintf("$%d", len(args)+1)
		whereConditions = append(whereConditions, fmt.Sprintf(expression, placeholder))
		args = append(args, value)
	}

	if filters.ScholarshipID != nil {
		addCondition("is.scholarship_id = %s", *filters.ScholarshipID)
	}

	if filters.InterviewerID != nil {
		addCondition("is.interviewer_id = %s", *filters.InterviewerID)
	}

	if filters.DateFrom != nil {
		addCondition("is.interview_date >= %s", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		addCondition("is.interview_date <= %s", *filters.DateTo)
	}

	if filters.IsAvailable != nil {
		addCondition("is.is_available = %s", *filters.IsAvailable)
	}

	if filters.SlotType != nil {
		addCondition("is.slot_type = %s", *filters.SlotType)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total items
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM interview_slots is %s
	`, whereClause)

	var totalItems int
	err := h.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting interview slots: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการนับข้อมูล",
		})
	}

	// Fetch slots
	queryArgs := append([]interface{}{}, args...)
	limitPlaceholder := fmt.Sprintf("$%d", len(queryArgs)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(queryArgs)+2)
	query := fmt.Sprintf(`
		SELECT 
			is.id, is.scholarship_id, is.interviewer_id, is.interview_date,
			is.start_time, is.end_time, is.location, is.building, is.floor, is.room,
			is.max_capacity, is.current_bookings, is.is_available, is.slot_type,
			is.duration_minutes, is.preparation_time, is.notes, is.created_by,
			is.created_at, is.updated_at,
			s.scholarship_name,
			u.first_name, u.last_name, u.email
		FROM interview_slots is
		LEFT JOIN scholarships s ON is.scholarship_id = s.scholarship_id
		LEFT JOIN users u ON is.interviewer_id = u.user_id
		%s
		ORDER BY is.interview_date ASC, is.start_time ASC
		LIMIT %s OFFSET %s
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	queryArgs = append(queryArgs, limit, offset)
	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("Error fetching interview slots: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูล",
		})
	}
	defer rows.Close()

	var slots []models.InterviewSlot
	for rows.Next() {
		var slot models.InterviewSlot
		var scholarship models.Scholarship
		var interviewer models.User

		err := rows.Scan(
			&slot.ID, &slot.ScholarshipID, &slot.InterviewerID, &slot.InterviewDate,
			&slot.StartTime, &slot.EndTime, &slot.Location, &slot.Building,
			&slot.Floor, &slot.Room, &slot.MaxCapacity, &slot.CurrentBookings,
			&slot.IsAvailable, &slot.SlotType, &slot.DurationMinutes,
			&slot.PreparationTime, &slot.Notes, &slot.CreatedBy,
			&slot.CreatedAt, &slot.UpdatedAt,
			&scholarship.ScholarshipName,
			&interviewer.FirstName, &interviewer.LastName, &interviewer.Email,
		)
		if err != nil {
			log.Printf("Error scanning interview slot: %v", err)
			continue
		}

		slot.Scholarship = &scholarship
		slot.Interviewer = &interviewer
		slots = append(slots, slot)
	}

	// Calculate pagination meta
	totalPages := (totalItems + limit - 1) / limit
	meta := models.PaginationMeta{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	return c.JSON(models.InterviewSlotsResponse{
		Success: true,
		Data:    slots,
		Meta:    meta,
		Filters: filters,
	})
}

// @Summary Update Interview Slot
// @Description Update interview slot details
// @Tags Interview Management
// @Accept json
// @Produce json
// @Param id path int true "Slot ID"
// @Param request body models.UpdateInterviewSlotRequest true "Update details"
// @Success 200 {object} models.InterviewSlotResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/slots/{id} [put]
func (h *InterviewReviewHandler) UpdateInterviewSlot(c *fiber.Ctx) error {
	slotID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสช่วงเวลาไม่ถูกต้อง",
		})
	}

	var req models.UpdateInterviewSlotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Check if slot exists
	slot, err := h.getInterviewSlotByID(slotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบช่วงเวลาสัมภาษณ์",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Check if slot has bookings and trying to make it unavailable
	if req.IsAvailable != nil && !*req.IsAvailable && slot.CurrentBookings > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถปิดช่วงเวลาที่มีการจองแล้วได้",
		})
	}

	// Build update query
	updateFields := []string{}
	args := []interface{}{}

	addUpdate := func(field string, value interface{}) {
		args = append(args, value)
		placeholder := fmt.Sprintf("$%d", len(args))
		updateFields = append(updateFields, fmt.Sprintf("%s = %s", field, placeholder))
	}

	if req.Location != nil {
		addUpdate("location", *req.Location)
	}
	if req.Building != nil {
		addUpdate("building", *req.Building)
	}
	if req.Floor != nil {
		addUpdate("floor", *req.Floor)
	}
	if req.Room != nil {
		addUpdate("room", *req.Room)
	}
	if req.MaxCapacity != nil {
		// Check if new capacity is not less than current bookings
		if *req.MaxCapacity < slot.CurrentBookings {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("จำนวนที่นั่งต้องไม่น้อยกว่าการจองปัจจุบัน (%d)", slot.CurrentBookings),
			})
		}
		addUpdate("max_capacity", *req.MaxCapacity)
	}
	if req.IsAvailable != nil {
		addUpdate("is_available", *req.IsAvailable)
	}
	if req.DurationMinutes != nil {
		addUpdate("duration_minutes", *req.DurationMinutes)
	}
	if req.PreparationTime != nil {
		addUpdate("preparation_time", *req.PreparationTime)
	}
	if req.Notes != nil {
		addUpdate("notes", *req.Notes)
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่มีข้อมูลที่ต้องอัปเดต",
		})
	}

	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, slotID)

	updateQuery := fmt.Sprintf(`
		UPDATE interview_slots 
		SET %s 
		WHERE id = $%d
	`, strings.Join(updateFields, ", "), len(args))

	_, err = h.db.Exec(updateQuery, args...)
	if err != nil {
		log.Printf("Error updating interview slot: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการอัปเดต",
		})
	}

	// Fetch updated slot
	updatedSlot, err := h.getInterviewSlotByID(slotID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "อัปเดตสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewSlotResponse{
		Success: true,
		Message: "อัปเดตช่วงเวลาสัมภาษณ์สำเร็จ",
		Data:    *updatedSlot,
	})
}

// @Summary Delete Interview Slot
// @Description Delete an interview slot (only if no bookings)
// @Tags Interview Management
// @Accept json
// @Produce json
// @Param id path int true "Slot ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/slots/{id} [delete]
func (h *InterviewReviewHandler) DeleteInterviewSlot(c *fiber.Ctx) error {
	slotID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสช่วงเวลาไม่ถูกต้อง",
		})
	}

	// Check if slot exists and has bookings
	slot, err := h.getInterviewSlotByID(slotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบช่วงเวลาสัมภาษณ์",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if slot.CurrentBookings > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถลบช่วงเวลาที่มีการจองแล้วได้",
		})
	}

	_, err = h.db.Exec("DELETE FROM interview_slots WHERE id = $1", slotID)
	if err != nil {
		log.Printf("Error deleting interview slot: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการลบ",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบช่วงเวลาสัมภาษณ์สำเร็จ",
	})
}

// ==================== INTERVIEW BOOKING ====================

// @Summary Book Interview Slot
// @Description Book an interview slot for student
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param request body models.BookInterviewRequest true "Booking details"
// @Success 201 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/book [post]
func (h *InterviewReviewHandler) BookInterviewSlot(c *fiber.Ctx) error {
	var req models.BookInterviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Get user from context
	userIDValue := c.Locals("user_id")
	userUUID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถยืนยันตัวผู้ใช้ได้",
		})
	}

	// Check if slot exists and is available
	slot, err := h.getInterviewSlotByID(req.SlotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบช่วงเวลาสัมภาษณ์",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if !slot.IsBookingAvailable() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ช่วงเวลานี้ไม่สามารถจองได้",
		})
	}

	// Resolve student's identifier
	var studentID string
	err = h.db.QueryRow(`SELECT student_id FROM students WHERE user_id = $1`, userUUID).Scan(&studentID)
	if err == sql.ErrNoRows {
		studentID = userUUID.String()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถตรวจสอบรหัสนักศึกษาได้",
		})
	}

	// Get student's active application for this scholarship
	var applicationID int
	err = h.db.QueryRow(`
		SELECT application_id FROM scholarship_applications 
		WHERE student_id = $1 AND scholarship_id = $2 
		AND application_status IN ('submitted', 'under_review', 'approved')
		ORDER BY COALESCE(submitted_at, created_at) DESC LIMIT 1
	`, studentID, slot.ScholarshipID).Scan(&applicationID)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบใบสมัครที่เกี่ยวข้อง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Check if student already has an active booking for this application
	var existingBooking int
	err = h.db.QueryRow(`
		SELECT COUNT(*) FROM interview_bookings 
		WHERE application_id = $1 AND booking_status IN ('booked', 'confirmed')
	`, applicationID).Scan(&existingBooking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if existingBooking > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "คุณมีการจองสัมภาษณ์สำหรับทุนนี้แล้ว",
		})
	}

	// Create booking
	insertQuery := `
		INSERT INTO interview_bookings (
			slot_id, application_id, student_id, student_notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var bookingID int
	err = h.db.QueryRow(insertQuery, req.SlotID, applicationID, studentID, req.StudentNotes).Scan(&bookingID)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการจอง",
		})
	}

	// Fetch created booking
	booking, err := h.getInterviewBookingByID(int(bookingID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "จองสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "จองสัมภาษณ์สำเร็จ",
		Data:    *booking,
	})
}

// ==================== HELPER FUNCTIONS ====================

func (h *InterviewReviewHandler) getInterviewSlotByID(id int) (*models.InterviewSlot, error) {
	query := `
		SELECT 
			is.id, is.scholarship_id, is.interviewer_id, is.interview_date,
			is.start_time, is.end_time, is.location, is.building, is.floor, is.room,
			is.max_capacity, is.current_bookings, is.is_available, is.slot_type,
			is.duration_minutes, is.preparation_time, is.notes, is.created_by,
			is.created_at, is.updated_at
		FROM interview_slots is
		WHERE is.id = $1
	`

	var slot models.InterviewSlot
	err := h.db.QueryRow(query, id).Scan(
		&slot.ID, &slot.ScholarshipID, &slot.InterviewerID, &slot.InterviewDate,
		&slot.StartTime, &slot.EndTime, &slot.Location, &slot.Building,
		&slot.Floor, &slot.Room, &slot.MaxCapacity, &slot.CurrentBookings,
		&slot.IsAvailable, &slot.SlotType, &slot.DurationMinutes,
		&slot.PreparationTime, &slot.Notes, &slot.CreatedBy,
		&slot.CreatedAt, &slot.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &slot, nil
}

func (h *InterviewReviewHandler) getInterviewBookingByID(id int) (*models.InterviewBooking, error) {
	query := `
		SELECT 
			ib.id, ib.slot_id, ib.application_id, ib.student_id, ib.booking_status,
			ib.booked_at, ib.confirmed_at, ib.cancelled_at, ib.cancellation_reason,
			ib.rescheduled_from_slot_id, ib.rescheduled_to_slot_id,
			ib.student_notes, ib.officer_notes, ib.reminder_sent_at,
			ib.check_in_time, ib.check_out_time, ib.actual_duration_minutes,
			ib.created_at, ib.updated_at
		FROM interview_bookings ib
		WHERE ib.id = $1
	`

	var booking models.InterviewBooking
	err := h.db.QueryRow(query, id).Scan(
		&booking.ID, &booking.SlotID, &booking.ApplicationID, &booking.StudentID,
		&booking.BookingStatus, &booking.BookedAt, &booking.ConfirmedAt,
		&booking.CancelledAt, &booking.CancellationReason,
		&booking.RescheduledFromSlotID, &booking.RescheduledToSlotID,
		&booking.StudentNotes, &booking.OfficerNotes, &booking.ReminderSentAt,
		&booking.CheckInTime, &booking.CheckOutTime, &booking.ActualDurationMinutes,
		&booking.CreatedAt, &booking.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// ==================== BOOKING MANAGEMENT ====================

// @Summary Get All Bookings
// @Description Get list of all interview bookings (Admin/Officer only)
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param scholarship_id query int false "Scholarship ID"
// @Param status query string false "Booking status"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.InterviewBookingsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings [get]
func (h *InterviewReviewHandler) GetAllBookings(c *fiber.Ctx) error {
	// Parse query parameters
	scholarshipID := c.Query("scholarship_id")
	status := c.Query("status")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Build query
	whereConditions := []string{}
	args := []interface{}{}

	if scholarshipID != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("is.scholarship_id = $%d", len(args)+1))
		args = append(args, scholarshipID)
	}
	if status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("ib.booking_status = $%d", len(args)+1))
		args = append(args, status)
	}
	if dateFrom != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("is.interview_date >= $%d", len(args)+1))
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("is.interview_date <= $%d", len(args)+1))
		args = append(args, dateTo)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total items
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM interview_bookings ib
		JOIN interview_slots is ON ib.slot_id = is.id
		%s
	`, whereClause)

	var totalItems int
	err := h.db.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting bookings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการนับข้อมูล",
		})
	}

	// Fetch bookings with full details
	queryArgs := append([]interface{}{}, args...)
	queryArgs = append(queryArgs, limit, offset)

	query := fmt.Sprintf(`
		SELECT
			ib.id, ib.slot_id, ib.application_id, ib.student_id, ib.booking_status,
			ib.booked_at, ib.confirmed_at, ib.cancelled_at, ib.cancellation_reason,
			ib.student_notes, ib.officer_notes, ib.check_in_time, ib.check_out_time,
			is.interview_date, is.start_time, is.end_time, is.location,
			s.scholarship_name,
			u.first_name, u.last_name, u.email
		FROM interview_bookings ib
		JOIN interview_slots is ON ib.slot_id = is.id
		JOIN scholarship_applications sa ON ib.application_id = sa.application_id
		JOIN scholarships s ON is.scholarship_id = s.scholarship_id
		JOIN users u ON ib.student_id = u.user_id
		%s
		ORDER BY is.interview_date DESC, is.start_time DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(queryArgs)-1, len(queryArgs))

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูล",
		})
	}
	defer rows.Close()

	var bookings []map[string]interface{}
	for rows.Next() {
		var (
			id, slotID, applicationID                         int
			studentID, bookingStatus, scholarshipName         string
			location, firstName, lastName, email              string
			studentNotes, officerNotes, cancellationReason    sql.NullString
			bookedAt, confirmedAt, cancelledAt                sql.NullTime
			checkInTime, checkOutTime                         sql.NullTime
			interviewDate                                     time.Time
			startTime, endTime                                string
		)

		err := rows.Scan(
			&id, &slotID, &applicationID, &studentID, &bookingStatus,
			&bookedAt, &confirmedAt, &cancelledAt, &cancellationReason,
			&studentNotes, &officerNotes, &checkInTime, &checkOutTime,
			&interviewDate, &startTime, &endTime, &location,
			&scholarshipName, &firstName, &lastName, &email,
		)
		if err != nil {
			log.Printf("Error scanning booking: %v", err)
			continue
		}

		booking := map[string]interface{}{
			"id":                id,
			"slot_id":           slotID,
			"application_id":    applicationID,
			"student_id":        studentID,
			"booking_status":    bookingStatus,
			"interview_date":    interviewDate.Format("2006-01-02"),
			"start_time":        startTime,
			"end_time":          endTime,
			"location":          location,
			"scholarship_name":  scholarshipName,
			"student_name":      fmt.Sprintf("%s %s", firstName, lastName),
			"student_email":     email,
		}

		if bookedAt.Valid {
			booking["booked_at"] = bookedAt.Time
		}
		if confirmedAt.Valid {
			booking["confirmed_at"] = confirmedAt.Time
		}
		if cancelledAt.Valid {
			booking["cancelled_at"] = cancelledAt.Time
		}
		if studentNotes.Valid {
			booking["student_notes"] = studentNotes.String
		}
		if officerNotes.Valid {
			booking["officer_notes"] = officerNotes.String
		}
		if cancellationReason.Valid {
			booking["cancellation_reason"] = cancellationReason.String
		}
		if checkInTime.Valid {
			booking["check_in_time"] = checkInTime.Time
		}
		if checkOutTime.Valid {
			booking["check_out_time"] = checkOutTime.Time
		}

		bookings = append(bookings, booking)
	}

	totalPages := (totalItems + limit - 1) / limit
	meta := models.PaginationMeta{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookings,
		"meta":    meta,
	})
}

// @Summary Get Booking Details
// @Description Get detailed information about a specific booking
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id} [get]
func (h *InterviewReviewHandler) GetBookingByID(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Data:    *booking,
	})
}

// @Summary Update Booking
// @Description Update booking information (Admin/Officer only)
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body models.UpdateBookingRequest true "Update details"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id} [put]
func (h *InterviewReviewHandler) UpdateBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	var req struct {
		OfficerNotes *string `json:"officer_notes"`
		BookingStatus *string `json:"booking_status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	_, err = h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Build update query
	updateFields := []string{}
	args := []interface{}{}

	if req.OfficerNotes != nil {
		args = append(args, *req.OfficerNotes)
		updateFields = append(updateFields, fmt.Sprintf("officer_notes = $%d", len(args)))
	}
	if req.BookingStatus != nil {
		args = append(args, *req.BookingStatus)
		updateFields = append(updateFields, fmt.Sprintf("booking_status = $%d", len(args)))
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่มีข้อมูลที่ต้องอัปเดต",
		})
	}

	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, bookingID)

	updateQuery := fmt.Sprintf(`
		UPDATE interview_bookings
		SET %s
		WHERE id = $%d
	`, strings.Join(updateFields, ", "), len(args))

	_, err = h.db.Exec(updateQuery, args...)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการอัปเดต",
		})
	}

	// Fetch updated booking
	updatedBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "อัปเดตสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "อัปเดตข้อมูลการจองสำเร็จ",
		Data:    *updatedBooking,
	})
}

// @Summary Cancel Booking
// @Description Cancel an interview booking
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body models.CancelBookingRequest true "Cancellation details"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id} [delete]
func (h *InterviewReviewHandler) CancelBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	var req struct {
		CancellationReason string `json:"cancellation_reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if booking.BookingStatus == "cancelled" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "การจองนี้ถูกยกเลิกแล้ว",
		})
	}

	// Update booking to cancelled
	updateQuery := `
		UPDATE interview_bookings
		SET booking_status = 'cancelled',
		    cancelled_at = CURRENT_TIMESTAMP,
		    cancellation_reason = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err = h.db.Exec(updateQuery, req.CancellationReason, bookingID)
	if err != nil {
		log.Printf("Error cancelling booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการยกเลิก",
		})
	}

	// Update slot availability
	h.db.Exec("UPDATE interview_slots SET current_bookings = current_bookings - 1 WHERE id = $1", booking.SlotID)

	// Fetch updated booking
	cancelledBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ยกเลิกสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "ยกเลิกการจองสำเร็จ",
		Data:    *cancelledBooking,
	})
}

// @Summary Reschedule Booking
// @Description Reschedule an interview to a different slot
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body models.RescheduleBookingRequest true "New slot details"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id}/reschedule [post]
func (h *InterviewReviewHandler) RescheduleBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	var req struct {
		NewSlotID int `json:"new_slot_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Check if new slot exists and is available
	newSlot, err := h.getInterviewSlotByID(req.NewSlotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบช่วงเวลาสัมภาษณ์ใหม่",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if !newSlot.IsBookingAvailable() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ช่วงเวลาใหม่ไม่สามารถจองได้",
		})
	}

	// Begin transaction
	tx, err := h.db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}
	defer tx.Rollback()

	// Update booking
	updateQuery := `
		UPDATE interview_bookings
		SET slot_id = $1,
		    rescheduled_from_slot_id = $2,
		    rescheduled_to_slot_id = $1,
		    booking_status = 'booked',
		    confirmed_at = NULL,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err = tx.Exec(updateQuery, req.NewSlotID, booking.SlotID, bookingID)
	if err != nil {
		log.Printf("Error rescheduling booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการเปลี่ยนนัด",
		})
	}

	// Update old slot bookings count
	_, err = tx.Exec("UPDATE interview_slots SET current_bookings = current_bookings - 1 WHERE id = $1", booking.SlotID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Update new slot bookings count
	_, err = tx.Exec("UPDATE interview_slots SET current_bookings = current_bookings + 1 WHERE id = $1", req.NewSlotID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	// Fetch updated booking
	rescheduledBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เปลี่ยนนัดสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "เปลี่ยนนัดสัมภาษณ์สำเร็จ",
		Data:    *rescheduledBooking,
	})
}

// @Summary Confirm Booking
// @Description Confirm an interview booking
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id}/confirm [post]
func (h *InterviewReviewHandler) ConfirmBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if booking.BookingStatus == "confirmed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "การจองนี้ได้รับการยืนยันแล้ว",
		})
	}

	if booking.BookingStatus == "cancelled" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถยืนยันการจองที่ถูกยกเลิกได้",
		})
	}

	// Update booking to confirmed
	updateQuery := `
		UPDATE interview_bookings
		SET booking_status = 'confirmed',
		    confirmed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = h.db.Exec(updateQuery, bookingID)
	if err != nil {
		log.Printf("Error confirming booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการยืนยัน",
		})
	}

	// Fetch updated booking
	confirmedBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ยืนยันสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "ยืนยันการจองสำเร็จ",
		Data:    *confirmedBooking,
	})
}

// @Summary Check-in Booking
// @Description Check-in for an interview
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id}/checkin [post]
func (h *InterviewReviewHandler) CheckInBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if booking.CheckInTime != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เช็คอินแล้ว",
		})
	}

	// Update check-in time
	updateQuery := `
		UPDATE interview_bookings
		SET check_in_time = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = h.db.Exec(updateQuery, bookingID)
	if err != nil {
		log.Printf("Error checking in booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการเช็คอิน",
		})
	}

	// Fetch updated booking
	checkedInBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เช็คอินสำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "เช็คอินสำเร็จ",
		Data:    *checkedInBooking,
	})
}

// @Summary Check-out Booking
// @Description Check-out from an interview
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} models.InterviewBookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/bookings/{id}/checkout [post]
func (h *InterviewReviewHandler) CheckOutBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสการจองไม่ถูกต้อง",
		})
	}

	// Check if booking exists
	booking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบข้อมูลการจอง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาด",
		})
	}

	if booking.CheckInTime == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ต้องเช็คอินก่อนเช็คเอาท์",
		})
	}

	if booking.CheckOutTime != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "เช็คเอาท์แล้ว",
		})
	}

	// Calculate duration
	checkOutTime := time.Now()
	duration := int(checkOutTime.Sub(*booking.CheckInTime).Minutes())

	// Update check-out time
	updateQuery := `
		UPDATE interview_bookings
		SET check_out_time = CURRENT_TIMESTAMP,
		    actual_duration_minutes = $1,
		    booking_status = 'completed',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err = h.db.Exec(updateQuery, duration, bookingID)
	if err != nil {
		log.Printf("Error checking out booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการเช็คเอาท์",
		})
	}

	// Fetch updated booking
	checkedOutBooking, err := h.getInterviewBookingByID(bookingID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เช็คเอาท์สำเร็จ แต่ไม่สามารถดึงข้อมูลได้",
		})
	}

	return c.JSON(models.InterviewBookingResponse{
		Success: true,
		Message: "เช็คเอาท์สำเร็จ",
		Data:    *checkedOutBooking,
	})
}

// @Summary Get Interview Statistics
// @Description Get statistics about interviews
// @Tags Interview Management
// @Accept json
// @Produce json
// @Param scholarship_id query int false "Scholarship ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.InterviewStatisticsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/statistics [get]
func (h *InterviewReviewHandler) GetStatistics(c *fiber.Ctx) error {
	scholarshipID := c.Query("scholarship_id")
	dateFrom := c.Query("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().Format("2006-01-02"))

	whereConditions := []string{"is.interview_date BETWEEN $1 AND $2"}
	args := []interface{}{dateFrom, dateTo}

	if scholarshipID != "" {
		args = append(args, scholarshipID)
		whereConditions = append(whereConditions, fmt.Sprintf("is.scholarship_id = $%d", len(args)))
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get overall statistics
	statsQuery := fmt.Sprintf(`
		SELECT
			COUNT(DISTINCT is.id) as total_slots,
			COUNT(DISTINCT CASE WHEN is.is_available = true THEN is.id END) as available_slots,
			COUNT(DISTINCT ib.id) as total_bookings,
			COUNT(DISTINCT CASE WHEN ib.booking_status = 'confirmed' THEN ib.id END) as confirmed_bookings,
			COUNT(DISTINCT CASE WHEN ib.booking_status = 'cancelled' THEN ib.id END) as cancelled_bookings,
			COUNT(DISTINCT CASE WHEN ib.booking_status = 'completed' THEN ib.id END) as completed_bookings,
			COUNT(DISTINCT CASE WHEN ib.check_in_time IS NOT NULL THEN ib.id END) as checked_in,
			AVG(CASE WHEN ib.actual_duration_minutes IS NOT NULL THEN ib.actual_duration_minutes END) as avg_duration
		FROM interview_slots is
		LEFT JOIN interview_bookings ib ON is.id = ib.slot_id
		WHERE %s
	`, whereClause)

	var stats struct {
		TotalSlots        int
		AvailableSlots    int
		TotalBookings     int
		ConfirmedBookings int
		CancelledBookings int
		CompletedBookings int
		CheckedIn         int
		AvgDuration       sql.NullFloat64
	}

	err := h.db.QueryRow(statsQuery, args...).Scan(
		&stats.TotalSlots, &stats.AvailableSlots, &stats.TotalBookings,
		&stats.ConfirmedBookings, &stats.CancelledBookings, &stats.CompletedBookings,
		&stats.CheckedIn, &stats.AvgDuration,
	)

	if err != nil {
		log.Printf("Error fetching statistics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูลสถิติ",
		})
	}

	avgDuration := 0.0
	if stats.AvgDuration.Valid {
		avgDuration = stats.AvgDuration.Float64
	}

	result := map[string]interface{}{
		"total_slots":        stats.TotalSlots,
		"available_slots":    stats.AvailableSlots,
		"total_bookings":     stats.TotalBookings,
		"confirmed_bookings": stats.ConfirmedBookings,
		"cancelled_bookings": stats.CancelledBookings,
		"completed_bookings": stats.CompletedBookings,
		"checked_in":         stats.CheckedIn,
		"avg_duration":       avgDuration,
		"utilization_rate":   0.0,
	}

	if stats.TotalSlots > 0 {
		result["utilization_rate"] = float64(stats.TotalBookings) / float64(stats.TotalSlots) * 100
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// @Summary Get Available Interview Slots
// @Description Get available interview slots for a scholarship
// @Tags Interview Booking
// @Accept json
// @Produce json
// @Param scholarship_id query int true "Scholarship ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.AvailabilityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/interview/availability [get]
func (h *InterviewReviewHandler) GetAvailableSlots(c *fiber.Ctx) error {
	scholarshipIDStr := c.Query("scholarship_id")
	if scholarshipIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณาระบุรหัสทุนการศึกษา",
		})
	}

	scholarshipID, err := strconv.Atoi(scholarshipIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสทุนการศึกษาไม่ถูกต้อง",
		})
	}

	dateFrom := c.Query("date_from", time.Now().Format("2006-01-02"))
	dateTo := c.Query("date_to", time.Now().AddDate(0, 1, 0).Format("2006-01-02"))

	query := `
		SELECT 
			is.id, is.interview_date, is.start_time, is.end_time,
			is.location, is.building, is.room, is.max_capacity, is.current_bookings,
			is.duration_minutes, u.first_name, u.last_name
		FROM interview_slots is
		LEFT JOIN users u ON is.interviewer_id = u.user_id
		WHERE is.scholarship_id = $1 
		AND is.interview_date BETWEEN $2 AND $3
		AND is.is_available = TRUE
		AND is.current_bookings < is.max_capacity
		ORDER BY is.interview_date ASC, is.start_time ASC
	`

	rows, err := h.db.Query(query, scholarshipID, dateFrom, dateTo)
	if err != nil {
		log.Printf("Error fetching available slots: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูล",
		})
	}
	defer rows.Close()

	dayMap := make(map[string]*models.DayAvailability)

	for rows.Next() {
		var slotID int
		var interviewDate time.Time
		var startTime, endTime, location, building, room string
		var maxCapacity, currentBookings, duration int
		var firstName, lastName string

		err := rows.Scan(
			&slotID, &interviewDate, &startTime, &endTime,
			&location, &building, &room, &maxCapacity, &currentBookings,
			&duration, &firstName, &lastName,
		)
		if err != nil {
			continue
		}

		dateStr := interviewDate.Format("2006-01-02")
		dayOfWeek := interviewDate.Weekday().String()

		if dayMap[dateStr] == nil {
			dayMap[dateStr] = &models.DayAvailability{
				Date:           dateStr,
				DayOfWeek:      dayOfWeek,
				TotalSlots:     0,
				AvailableSlots: 0,
				BookedSlots:    0,
				TimeSlots:      []models.TimeSlotInfo{},
			}
		}

		dayAvail := dayMap[dateStr]
		dayAvail.TotalSlots++

		isBooked := currentBookings >= maxCapacity
		if !isBooked {
			dayAvail.AvailableSlots++
		} else {
			dayAvail.BookedSlots++
		}

		timeSlot := models.TimeSlotInfo{
			SlotID:      slotID,
			StartTime:   startTime,
			EndTime:     endTime,
			IsAvailable: !isBooked,
			IsBooked:    isBooked,
			Location:    fmt.Sprintf("%s %s %s", building, room, location),
			Interviewer: fmt.Sprintf("%s %s", firstName, lastName),
			Duration:    duration,
		}

		dayAvail.TimeSlots = append(dayAvail.TimeSlots, timeSlot)
	}

	// Convert map to slice
	var availability []models.DayAvailability
	for _, day := range dayMap {
		availability = append(availability, *day)
	}

	return c.JSON(models.AvailabilityResponse{
		Success: true,
		Data:    availability,
	})
}
