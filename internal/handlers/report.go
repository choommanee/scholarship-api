package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
)

type ReportHandler struct {
	cfg *config.Config
}

func NewReportHandler(cfg *config.Config) *ReportHandler {
	return &ReportHandler{cfg: cfg}
}

// GetDashboardSummary provides dashboard statistics
// @Summary Get dashboard summary
// @Description Get comprehensive dashboard statistics for admin overview
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/dashboard [get]
func (h *ReportHandler) GetDashboardSummary(c *fiber.Ctx) error {
	var stats struct {
		TotalScholarships    int     `json:"total_scholarships"`
		ActiveScholarships   int     `json:"active_scholarships"`
		TotalApplications    int     `json:"total_applications"`
		PendingApplications  int     `json:"pending_applications"`
		ApprovedApplications int     `json:"approved_applications"`
		TotalBudget          float64 `json:"total_budget"`
		AllocatedBudget      float64 `json:"allocated_budget"`
		RemainingBudget      float64 `json:"remaining_budget"`
		TotalStudents        int     `json:"total_students"`
		InterviewsPending    int     `json:"interviews_pending"`
	}

	// Get scholarship statistics
	database.DB.QueryRow("SELECT COUNT(*) FROM scholarships").Scan(&stats.TotalScholarships)
	database.DB.QueryRow("SELECT COUNT(*) FROM scholarships WHERE is_active = true").Scan(&stats.ActiveScholarships)

	// Get application statistics
	database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications").Scan(&stats.TotalApplications)
	database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'submitted'").Scan(&stats.PendingApplications)
	database.DB.QueryRow("SELECT COUNT(*) FROM scholarship_applications WHERE application_status = 'approved'").Scan(&stats.ApprovedApplications)

	// Get budget statistics
	database.DB.QueryRow("SELECT COALESCE(SUM(total_budget), 0) FROM scholarship_budgets").Scan(&stats.TotalBudget)
	database.DB.QueryRow("SELECT COALESCE(SUM(allocated_budget), 0) FROM scholarship_budgets").Scan(&stats.AllocatedBudget)
	stats.RemainingBudget = stats.TotalBudget - stats.AllocatedBudget

	// Get student statistics
	database.DB.QueryRow("SELECT COUNT(*) FROM students WHERE student_status = 'active'").Scan(&stats.TotalStudents)

	// Get interview statistics
	database.DB.QueryRow("SELECT COUNT(*) FROM interview_appointments WHERE appointment_status = 'scheduled'").Scan(&stats.InterviewsPending)

	return c.JSON(fiber.Map{
		"data": stats,
	})
}

// GetApplicationReport generates application report with filters
// @Summary Get application report
// @Description Generate comprehensive application report with filters (Admin/Officer only)
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Filter from date (YYYY-MM-DD)"
// @Param end_date query string false "Filter to date (YYYY-MM-DD)"
// @Param scholarship_id query string false "Filter by scholarship ID"
// @Param status query string false "Filter by application status"
// @Param faculty_code query string false "Filter by faculty code"
// @Success 200 {object} object{data=[]object,summary=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/applications [get]
func (h *ReportHandler) GetApplicationReport(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	scholarshipID := c.Query("scholarship_id")
	status := c.Query("status")
	facultyCode := c.Query("faculty_code")

	query := `SELECT sa.application_id, sa.student_id, sa.scholarship_id,
		sa.application_status, sa.submitted_at, sa.priority_score,
		s.scholarship_name, s.amount,
		u.first_name, u.last_name, u.email,
		st.faculty_code, st.department_code, st.gpa,
		CASE WHEN ia.appointment_id IS NOT NULL THEN 'Yes' ELSE 'No' END as has_interview,
		CASE WHEN sal.allocation_id IS NOT NULL THEN sal.allocated_amount ELSE 0 END as allocated_amount
		FROM scholarship_applications sa
		JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		JOIN students st ON sa.student_id = st.student_id
		JOIN users u ON st.user_id = u.user_id
		LEFT JOIN interview_appointments ia ON sa.application_id = ia.application_id
		LEFT JOIN scholarship_allocations sal ON sa.application_id = sal.application_id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if startDate != "" {
		argCount++
		query += " AND sa.submitted_at >= $" + strconv.Itoa(argCount)
		args = append(args, startDate)
	}

	if endDate != "" {
		argCount++
		query += " AND sa.submitted_at <= $" + strconv.Itoa(argCount)
		args = append(args, endDate+" 23:59:59")
	}

	if scholarshipID != "" {
		argCount++
		query += " AND sa.scholarship_id = $" + strconv.Itoa(argCount)
		args = append(args, scholarshipID)
	}

	if status != "" {
		argCount++
		query += " AND sa.application_status = $" + strconv.Itoa(argCount)
		args = append(args, status)
	}

	if facultyCode != "" {
		argCount++
		query += " AND st.faculty_code = $" + strconv.Itoa(argCount)
		args = append(args, facultyCode)
	}

	query += " ORDER BY sa.submitted_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate application report",
		})
	}
	defer rows.Close()

	var applications []map[string]interface{}
	for rows.Next() {
		var app map[string]interface{} = make(map[string]interface{})
		var applicationID, scholarshipID int
		var studentID, applicationStatus, scholarshipName, firstName, lastName, email string
		var facultyCode, departmentCode, hasInterview string
		var submittedAt *time.Time
		var priorityScore, amount, gpa, allocatedAmount *float64

		err := rows.Scan(
			&applicationID, &studentID, &scholarshipID,
			&applicationStatus, &submittedAt, &priorityScore,
			&scholarshipName, &amount,
			&firstName, &lastName, &email,
			&facultyCode, &departmentCode, &gpa,
			&hasInterview, &allocatedAmount,
		)
		if err != nil {
			continue
		}

		app["application_id"] = applicationID
		app["student_id"] = studentID
		app["scholarship_id"] = scholarshipID
		app["application_status"] = applicationStatus
		app["submitted_at"] = submittedAt
		app["priority_score"] = priorityScore
		app["scholarship_name"] = scholarshipName
		app["scholarship_amount"] = amount
		app["student_name"] = firstName + " " + lastName
		app["email"] = email
		app["faculty_code"] = facultyCode
		app["department_code"] = departmentCode
		app["gpa"] = gpa
		app["has_interview"] = hasInterview
		app["allocated_amount"] = allocatedAmount

		applications = append(applications, app)
	}

	return c.JSON(fiber.Map{
		"data": applications,
		"summary": fiber.Map{
			"total_applications": len(applications),
			"filters_applied": fiber.Map{
				"start_date":     startDate,
				"end_date":       endDate,
				"scholarship_id": scholarshipID,
				"status":         status,
				"faculty_code":   facultyCode,
			},
		},
	})
}

// GetScholarshipReport generates scholarship utilization report
// @Summary Get scholarship report
// @Description Generate scholarship utilization and statistics report (Admin/Officer only)
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param academic_year query string false "Filter by academic year"
// @Success 200 {object} object{data=[]object,summary=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/scholarships [get]
func (h *ReportHandler) GetScholarshipReport(c *fiber.Ctx) error {
	academicYear := c.Query("academic_year")

	query := `SELECT s.scholarship_id, s.scholarship_name, s.scholarship_type,
		s.amount, s.total_quota, s.available_quota,
		s.academic_year, s.application_start_date, s.application_end_date,
		ss.source_name,
		COUNT(sa.application_id) as total_applications,
		COUNT(CASE WHEN sa.application_status = 'approved' THEN 1 END) as approved_applications,
		COALESCE(sb.total_budget, 0) as total_budget,
		COALESCE(sb.allocated_budget, 0) as allocated_budget,
		COALESCE(sb.remaining_budget, 0) as remaining_budget
		FROM scholarships s
		LEFT JOIN scholarship_sources ss ON s.source_id = ss.source_id
		LEFT JOIN scholarship_applications sa ON s.scholarship_id = sa.scholarship_id
		LEFT JOIN scholarship_budgets sb ON s.scholarship_id = sb.scholarship_id
		WHERE s.is_active = true`

	args := []interface{}{}
	if academicYear != "" {
		query += " AND s.academic_year = $1"
		args = append(args, academicYear)
	}

	query += ` GROUP BY s.scholarship_id, s.scholarship_name, s.scholarship_type,
		s.amount, s.total_quota, s.available_quota, s.academic_year,
		s.application_start_date, s.application_end_date, ss.source_name,
		sb.total_budget, sb.allocated_budget, sb.remaining_budget
		ORDER BY s.scholarship_name`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate scholarship report",
		})
	}
	defer rows.Close()

	var scholarships []map[string]interface{}
	var totalBudget, totalAllocated float64
	var totalQuota, totalApplications, totalApproved int

	for rows.Next() {
		var scholarship map[string]interface{} = make(map[string]interface{})
		var scholarshipID, quota, availableQuota, totalApps, approvedApps int
		var name, scholarshipType, academicYear, sourceName string
		var amount, budget, allocated, remaining float64
		var startDate, endDate time.Time

		err := rows.Scan(
			&scholarshipID, &name, &scholarshipType,
			&amount, &quota, &availableQuota,
			&academicYear, &startDate, &endDate,
			&sourceName, &totalApps, &approvedApps,
			&budget, &allocated, &remaining,
		)
		if err != nil {
			continue
		}

		scholarship["scholarship_id"] = scholarshipID
		scholarship["scholarship_name"] = name
		scholarship["scholarship_type"] = scholarshipType
		scholarship["amount"] = amount
		scholarship["total_quota"] = quota
		scholarship["available_quota"] = availableQuota
		scholarship["academic_year"] = academicYear
		scholarship["application_start_date"] = startDate
		scholarship["application_end_date"] = endDate
		scholarship["source_name"] = sourceName
		scholarship["total_applications"] = totalApps
		scholarship["approved_applications"] = approvedApps
		scholarship["total_budget"] = budget
		scholarship["allocated_budget"] = allocated
		scholarship["remaining_budget"] = remaining
		scholarship["utilization_rate"] = (allocated / budget) * 100
		scholarship["quota_utilization"] = ((quota - availableQuota) / quota) * 100

		scholarships = append(scholarships, scholarship)

		// Aggregate totals
		totalBudget += budget
		totalAllocated += allocated
		totalQuota += quota
		totalApplications += totalApps
		totalApproved += approvedApps
	}

	return c.JSON(fiber.Map{
		"data": scholarships,
		"summary": fiber.Map{
			"total_scholarships":   len(scholarships),
			"total_budget":         totalBudget,
			"total_allocated":      totalAllocated,
			"total_remaining":      totalBudget - totalAllocated,
			"overall_utilization":  (totalAllocated / totalBudget) * 100,
			"total_quota":          totalQuota,
			"total_applications":   totalApplications,
			"total_approved":       totalApproved,
			"academic_year":        academicYear,
		},
	})
}

// GetBudgetReport generates budget allocation report
// @Summary Get budget report
// @Description Generate budget allocation and utilization report (Admin/Officer only)
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param budget_year query string false "Filter by budget year" default(current year)
// @Success 200 {object} object{data=[]object,summary=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/budget [get]
func (h *ReportHandler) GetBudgetReport(c *fiber.Ctx) error {
	budgetYear := c.Query("budget_year", strconv.Itoa(time.Now().Year()))

	query := `SELECT sb.scholarship_id, s.scholarship_name, sb.budget_year,
		sb.total_budget, sb.allocated_budget, sb.remaining_budget,
		ss.source_name, ss.source_type,
		COUNT(sal.allocation_id) as allocation_count,
		COALESCE(SUM(sal.allocated_amount), 0) as total_disbursed
		FROM scholarship_budgets sb
		JOIN scholarships s ON sb.scholarship_id = s.scholarship_id
		LEFT JOIN scholarship_sources ss ON s.source_id = ss.source_id
		LEFT JOIN scholarship_allocations sal ON sb.scholarship_id = sal.scholarship_id
		WHERE sb.budget_year = $1
		GROUP BY sb.scholarship_id, s.scholarship_name, sb.budget_year,
		sb.total_budget, sb.allocated_budget, sb.remaining_budget,
		ss.source_name, ss.source_type
		ORDER BY sb.total_budget DESC`

	rows, err := database.DB.Query(query, budgetYear)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate budget report",
		})
	}
	defer rows.Close()

	var budgets []map[string]interface{}
	var grandTotal, grandAllocated, grandRemaining, grandDisbursed float64

	for rows.Next() {
		var budget map[string]interface{} = make(map[string]interface{})
		var scholarshipID, allocationCount int
		var scholarshipName, year, sourceName, sourceType string
		var totalBudget, allocatedBudget, remainingBudget, totalDisbursed float64

		err := rows.Scan(
			&scholarshipID, &scholarshipName, &year,
			&totalBudget, &allocatedBudget, &remainingBudget,
			&sourceName, &sourceType, &allocationCount, &totalDisbursed,
		)
		if err != nil {
			continue
		}

		budget["scholarship_id"] = scholarshipID
		budget["scholarship_name"] = scholarshipName
		budget["budget_year"] = year
		budget["total_budget"] = totalBudget
		budget["allocated_budget"] = allocatedBudget
		budget["remaining_budget"] = remainingBudget
		budget["source_name"] = sourceName
		budget["source_type"] = sourceType
		budget["allocation_count"] = allocationCount
		budget["total_disbursed"] = totalDisbursed
		budget["utilization_rate"] = (allocatedBudget / totalBudget) * 100
		budget["disbursement_rate"] = (totalDisbursed / totalBudget) * 100

		budgets = append(budgets, budget)

		// Aggregate totals
		grandTotal += totalBudget
		grandAllocated += allocatedBudget
		grandRemaining += remainingBudget
		grandDisbursed += totalDisbursed
	}

	return c.JSON(fiber.Map{
		"data": budgets,
		"summary": fiber.Map{
			"budget_year":               budgetYear,
			"total_scholarships":        len(budgets),
			"grand_total_budget":        grandTotal,
			"grand_allocated_budget":    grandAllocated,
			"grand_remaining_budget":    grandRemaining,
			"grand_total_disbursed":     grandDisbursed,
			"overall_utilization_rate":  (grandAllocated / grandTotal) * 100,
			"overall_disbursement_rate": (grandDisbursed / grandTotal) * 100,
		},
	})
}

// GetStudentReport generates student statistics report
// @Summary Get student report
// @Description Generate student scholarship statistics and success rates report (Admin/Officer only)
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param faculty_code query string false "Filter by faculty code"
// @Param academic_year query string false "Filter by academic year"
// @Success 200 {object} object{data=[]object,summary=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/students [get]
func (h *ReportHandler) GetStudentReport(c *fiber.Ctx) error {
	facultyCode := c.Query("faculty_code")
	academicYear := c.Query("academic_year")

	query := `SELECT st.student_id, u.first_name, u.last_name, u.email,
		st.faculty_code, st.department_code, st.year_level, st.gpa,
		st.admission_year, st.student_status,
		COUNT(sa.application_id) as total_applications,
		COUNT(CASE WHEN sa.application_status = 'approved' THEN 1 END) as approved_applications,
		COALESCE(SUM(sal.allocated_amount), 0) as total_received
		FROM students st
		JOIN users u ON st.user_id = u.user_id
		LEFT JOIN scholarship_applications sa ON st.student_id = sa.student_id
		LEFT JOIN scholarship_allocations sal ON sa.application_id = sal.application_id
		WHERE st.student_status = 'active'`

	args := []interface{}{}
	argCount := 0

	if facultyCode != "" {
		argCount++
		query += " AND st.faculty_code = $" + strconv.Itoa(argCount)
		args = append(args, facultyCode)
	}

	if academicYear != "" {
		argCount++
		query += " AND st.admission_year = $" + strconv.Itoa(argCount)
		args = append(args, academicYear)
	}

	query += ` GROUP BY st.student_id, u.first_name, u.last_name, u.email,
		st.faculty_code, st.department_code, st.year_level, st.gpa,
		st.admission_year, st.student_status
		ORDER BY st.faculty_code, st.student_id`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate student report",
		})
	}
	defer rows.Close()

	var students []map[string]interface{}
	var totalReceived float64
	var totalApplications, totalApproved int

	for rows.Next() {
		var student map[string]interface{} = make(map[string]interface{})
		var studentID, firstName, lastName, email, faculty, department, status string
		var yearLevel, admissionYear, applications, approved int
		var gpa, received float64

		err := rows.Scan(
			&studentID, &firstName, &lastName, &email,
			&faculty, &department, &yearLevel, &gpa,
			&admissionYear, &status, &applications, &approved, &received,
		)
		if err != nil {
			continue
		}

		student["student_id"] = studentID
		student["student_name"] = firstName + " " + lastName
		student["email"] = email
		student["faculty_code"] = faculty
		student["department_code"] = department
		student["year_level"] = yearLevel
		student["gpa"] = gpa
		student["admission_year"] = admissionYear
		student["student_status"] = status
		student["total_applications"] = applications
		student["approved_applications"] = approved
		student["total_received"] = received
		student["success_rate"] = float64(approved) / float64(applications) * 100

		students = append(students, student)

		// Aggregate totals
		totalReceived += received
		totalApplications += applications
		totalApproved += approved
	}

	return c.JSON(fiber.Map{
		"data": students,
		"summary": fiber.Map{
			"total_students":       len(students),
			"total_applications":   totalApplications,
			"total_approved":       totalApproved,
			"total_amount_received": totalReceived,
			"overall_success_rate": float64(totalApproved) / float64(totalApplications) * 100,
			"filters": fiber.Map{
				"faculty_code":   facultyCode,
				"academic_year":  academicYear,
			},
		},
	})
}

// ExportReport exports report data (placeholder for CSV/Excel export)
// @Summary Export report
// @Description Export report data in various formats (Admin/Officer only)
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param type path string true "Report type (application/scholarship/budget/student)"
// @Param format query string false "Export format (csv/excel)" default(csv)
// @Success 200 {object} object{message=string,report_type=string,format=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /reports/export/{type} [get]
func (h *ReportHandler) ExportReport(c *fiber.Ctx) error {
	reportType := c.Params("type") // application, scholarship, budget, student
	format := c.Query("format", "csv") // csv, excel

	// This is a placeholder implementation
	// In a real application, you would generate actual CSV or Excel files
	
	return c.JSON(fiber.Map{
		"message": "Export functionality not yet implemented",
		"report_type": reportType,
		"format": format,
		"note": "This endpoint would generate downloadable files in production",
	})
}