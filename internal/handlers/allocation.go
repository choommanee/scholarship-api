package handlers

import (
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AllocationHandler struct {
	cfg *config.Config
}

func NewAllocationHandler(cfg *config.Config) *AllocationHandler {
	return &AllocationHandler{cfg: cfg}
}

// CreateAllocation creates new scholarship allocation
// @Summary Create scholarship allocation
// @Description Create new scholarship allocation for approved application (Admin/Officer only)
// @Tags Fund Allocation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param allocation body models.ScholarshipAllocation true "Allocation data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /allocations [post]
func (h *AllocationHandler) CreateAllocation(c *fiber.Ctx) error {
	var allocation models.ScholarshipAllocation
	if err := c.BodyParser(&allocation); err != nil {
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
	allocation.AllocatedBy = userID

	// Check if application is approved for allocation
	var appStatus string
	checkQuery := `SELECT application_status FROM scholarship_applications WHERE application_id = $1`
	err = database.DB.QueryRow(checkQuery, allocation.ApplicationID).Scan(&appStatus)
	if err != nil || appStatus != "approved" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application is not approved for allocation",
		})
	}

	query := `INSERT INTO scholarship_allocations 
		(application_id, scholarship_id, allocated_amount, allocation_date, 
		disbursement_method, bank_account, bank_name, allocated_by, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING allocation_id`

	err = database.DB.QueryRow(query,
		allocation.ApplicationID, allocation.ScholarshipID, allocation.AllocatedAmount,
		allocation.AllocationDate, allocation.DisbursementMethod, allocation.BankAccount,
		allocation.BankName, allocation.AllocatedBy, allocation.Notes,
	).Scan(&allocation.AllocationID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create allocation",
		})
	}

	// Update scholarship budget
	updateBudgetQuery := `UPDATE scholarship_budgets 
		SET allocated_budget = allocated_budget + $1 
		WHERE scholarship_id = $2`
	database.DB.Exec(updateBudgetQuery, allocation.AllocatedAmount, allocation.ScholarshipID)

	// Update scholarship available quota
	updateQuotaQuery := `UPDATE scholarships 
		SET available_quota = available_quota - 1 
		WHERE scholarship_id = $1`
	database.DB.Exec(updateQuotaQuery, allocation.ScholarshipID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Allocation created successfully",
		"data":    allocation,
	})
}

// GetAllocations retrieves scholarship allocations with filters
// @Summary Get scholarship allocations
// @Description Get paginated list of scholarship allocations with filters (Admin/Officer only)
// @Tags Fund Allocation
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by allocation status"
// @Param scholarship_id query string false "Filter by scholarship ID"
// @Success 200 {object} object{data=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /allocations [get]
func (h *AllocationHandler) GetAllocations(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	scholarshipID := c.Query("scholarship_id")

	offset := (page - 1) * limit

	query := `SELECT sa.allocation_id, sa.application_id, sa.scholarship_id, 
		sa.allocated_amount, sa.allocation_status, sa.allocation_date,
		sa.disbursement_method, sa.bank_account, sa.bank_name,
		sa.transfer_date, sa.transfer_reference, sa.notes,
		s.scholarship_name, u.first_name, u.last_name, st.student_id
		FROM scholarship_allocations sa
		JOIN scholarship_applications app ON sa.application_id = app.application_id
		JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		JOIN students st ON app.student_id = st.student_id
		JOIN users u ON st.user_id = u.user_id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if status != "" {
		argCount++
		query += " AND sa.allocation_status = $" + strconv.Itoa(argCount)
		args = append(args, status)
	}

	if scholarshipID != "" {
		argCount++
		query += " AND sa.scholarship_id = $" + strconv.Itoa(argCount)
		args = append(args, scholarshipID)
	}

	query += " ORDER BY sa.allocation_date DESC"

	argCount++
	query += " LIMIT $" + strconv.Itoa(argCount)
	args = append(args, limit)

	argCount++
	query += " OFFSET $" + strconv.Itoa(argCount)
	args = append(args, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch allocations",
		})
	}
	defer rows.Close()

	var allocations []map[string]interface{}
	for rows.Next() {
		var allocation map[string]interface{} = make(map[string]interface{})
		var allocationID, applicationID, scholarshipID int
		var allocatedAmount float64
		var allocationStatus, disbursementMethod, bankAccount, bankName string
		var allocationDate, transferDate *string
		var transferReference, notes, scholarshipName, firstName, lastName, studentID string

		err := rows.Scan(
			&allocationID, &applicationID, &scholarshipID,
			&allocatedAmount, &allocationStatus, &allocationDate,
			&disbursementMethod, &bankAccount, &bankName,
			&transferDate, &transferReference, &notes,
			&scholarshipName, &firstName, &lastName, &studentID,
		)
		if err != nil {
			continue
		}

		allocation["allocation_id"] = allocationID
		allocation["application_id"] = applicationID
		allocation["scholarship_id"] = scholarshipID
		allocation["allocated_amount"] = allocatedAmount
		allocation["allocation_status"] = allocationStatus
		allocation["allocation_date"] = allocationDate
		allocation["disbursement_method"] = disbursementMethod
		allocation["bank_account"] = bankAccount
		allocation["bank_name"] = bankName
		allocation["transfer_date"] = transferDate
		allocation["transfer_reference"] = transferReference
		allocation["notes"] = notes
		allocation["scholarship_name"] = scholarshipName
		allocation["student_name"] = firstName + " " + lastName
		allocation["student_id"] = studentID

		allocations = append(allocations, allocation)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM scholarship_allocations sa WHERE 1=1`
	countArgs := []interface{}{}

	if status != "" {
		countQuery += " AND allocation_status = $1"
		countArgs = append(countArgs, status)
	}
	if scholarshipID != "" {
		if len(countArgs) > 0 {
			countQuery += " AND scholarship_id = $2"
		} else {
			countQuery += " AND scholarship_id = $1"
		}
		countArgs = append(countArgs, scholarshipID)
	}

	var total int
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return c.JSON(fiber.Map{
		"data": allocations,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ApproveAllocation approves a pending allocation
// @Summary Approve allocation
// @Description Approve a pending scholarship allocation (Admin/Officer only)
// @Tags Fund Allocation
// @Produce json
// @Security BearerAuth
// @Param id path string true "Allocation ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /allocations/{id}/approve [post]
func (h *AllocationHandler) ApproveAllocation(c *fiber.Ctx) error {
	allocationID := c.Params("id")
	userID := c.Locals("user_id").(string)

	query := `UPDATE scholarship_allocations 
		SET allocation_status = 'approved', approved_by = $1, updated_at = CURRENT_TIMESTAMP
		WHERE allocation_id = $2 AND allocation_status = 'pending'`

	result, err := database.DB.Exec(query, userID, allocationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve allocation",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Allocation not found or already processed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Allocation approved successfully",
	})
}

// DisburseAllocation marks allocation as disbursed with transfer details
// @Summary Disburse allocation
// @Description Mark allocation as disbursed with transfer reference details (Admin/Officer only)
// @Tags Fund Allocation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Allocation ID"
// @Param disbursement body object{transfer_date=string,transfer_reference=string} true "Disbursement details"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /allocations/{id}/disburse [post]
func (h *AllocationHandler) DisburseAllocation(c *fiber.Ctx) error {
	allocationID := c.Params("id")

	var disbursement struct {
		TransferDate      string `json:"transfer_date"`
		TransferReference string `json:"transfer_reference"`
	}

	if err := c.BodyParser(&disbursement); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	query := `UPDATE scholarship_allocations 
		SET allocation_status = 'disbursed', 
		    transfer_date = $1, 
		    transfer_reference = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE allocation_id = $3 AND allocation_status = 'approved'`

	result, err := database.DB.Exec(query, disbursement.TransferDate, disbursement.TransferReference, allocationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update disbursement",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Allocation not found or not approved",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Allocation disbursed successfully",
	})
}

// GetAllocationDetails retrieves detailed allocation information
// @Summary Get allocation details
// @Description Get detailed information about a specific allocation (Admin/Officer only)
// @Tags Fund Allocation
// @Produce json
// @Security BearerAuth
// @Param id path string true "Allocation ID"
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /allocations/{id} [get]
func (h *AllocationHandler) GetAllocationDetails(c *fiber.Ctx) error {
	allocationID := c.Params("id")

	query := `SELECT sa.allocation_id, sa.application_id, sa.scholarship_id,
		sa.allocated_amount, sa.allocation_status, sa.allocation_date,
		sa.disbursement_method, sa.bank_account, sa.bank_name,
		sa.transfer_date, sa.transfer_reference, sa.notes,
		s.scholarship_name, u.first_name, u.last_name, st.student_id,
		u.email, u.phone, st.faculty_code, st.department_code
		FROM scholarship_allocations sa
		JOIN scholarship_applications app ON sa.application_id = app.application_id
		JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		JOIN students st ON app.student_id = st.student_id
		JOIN users u ON st.user_id = u.user_id
		WHERE sa.allocation_id = $1`

	var allocation map[string]interface{} = make(map[string]interface{})
	var allocationID_int, applicationID, scholarshipID int
	var allocatedAmount float64
	var allocationStatus, disbursementMethod, bankAccount, bankName string
	var allocationDate, transferDate *string
	var transferReference, notes, scholarshipName, firstName, lastName, studentID string
	var email, phone, facultyCode, departmentCode string

	err := database.DB.QueryRow(query, allocationID).Scan(
		&allocationID_int, &applicationID, &scholarshipID,
		&allocatedAmount, &allocationStatus, &allocationDate,
		&disbursementMethod, &bankAccount, &bankName,
		&transferDate, &transferReference, &notes,
		&scholarshipName, &firstName, &lastName, &studentID,
		&email, &phone, &facultyCode, &departmentCode,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Allocation not found",
		})
	}

	allocation["allocation_id"] = allocationID_int
	allocation["application_id"] = applicationID
	allocation["scholarship_id"] = scholarshipID
	allocation["allocated_amount"] = allocatedAmount
	allocation["allocation_status"] = allocationStatus
	allocation["allocation_date"] = allocationDate
	allocation["disbursement_method"] = disbursementMethod
	allocation["bank_account"] = bankAccount
	allocation["bank_name"] = bankName
	allocation["transfer_date"] = transferDate
	allocation["transfer_reference"] = transferReference
	allocation["notes"] = notes
	allocation["scholarship_name"] = scholarshipName
	allocation["student_id"] = studentID
	allocation["student_name"] = firstName + " " + lastName
	allocation["email"] = email
	allocation["phone"] = phone
	allocation["faculty_code"] = facultyCode
	allocation["department_code"] = departmentCode

	return c.JSON(fiber.Map{
		"data": allocation,
	})
}

// GetBudgetSummary retrieves budget summary for scholarships
// @Summary Get budget summary
// @Description Get budget allocation summary with utilization rates (Admin/Officer only)
// @Tags Fund Allocation
// @Produce json
// @Security BearerAuth
// @Param scholarship_id query string false "Filter by scholarship ID"
// @Param year query string false "Filter by budget year"
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /allocations/budget/summary [get]
func (h *AllocationHandler) GetBudgetSummary(c *fiber.Ctx) error {
	scholarshipID := c.Query("scholarship_id")
	year := c.Query("year")

	query := `SELECT sb.scholarship_id, s.scholarship_name, sb.budget_year,
		sb.total_budget, sb.allocated_budget, sb.remaining_budget,
		COUNT(sa.allocation_id) as allocation_count
		FROM scholarship_budgets sb
		JOIN scholarships s ON sb.scholarship_id = s.scholarship_id
		LEFT JOIN scholarship_allocations sa ON sb.scholarship_id = sa.scholarship_id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if scholarshipID != "" {
		argCount++
		query += " AND sb.scholarship_id = $" + strconv.Itoa(argCount)
		args = append(args, scholarshipID)
	}

	if year != "" {
		argCount++
		query += " AND sb.budget_year = $" + strconv.Itoa(argCount)
		args = append(args, year)
	}

	query += ` GROUP BY sb.scholarship_id, s.scholarship_name, sb.budget_year,
		sb.total_budget, sb.allocated_budget, sb.remaining_budget
		ORDER BY sb.budget_year DESC, s.scholarship_name`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch budget summary",
		})
	}
	defer rows.Close()

	var budgets []map[string]interface{}
	for rows.Next() {
		var budget map[string]interface{} = make(map[string]interface{})
		var scholarshipID int
		var scholarshipName, budgetYear string
		var totalBudget, allocatedBudget, remainingBudget float64
		var allocationCount int

		err := rows.Scan(
			&scholarshipID, &scholarshipName, &budgetYear,
			&totalBudget, &allocatedBudget, &remainingBudget,
			&allocationCount,
		)
		if err != nil {
			continue
		}

		budget["scholarship_id"] = scholarshipID
		budget["scholarship_name"] = scholarshipName
		budget["budget_year"] = budgetYear
		budget["total_budget"] = totalBudget
		budget["allocated_budget"] = allocatedBudget
		budget["remaining_budget"] = remainingBudget
		budget["allocation_count"] = allocationCount
		budget["utilization_rate"] = (allocatedBudget / totalBudget) * 100

		budgets = append(budgets, budget)
	}

	return c.JSON(fiber.Map{
		"data": budgets,
	})
}
