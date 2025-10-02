package handlers

import (
	"database/sql"
	"math"
	"strconv"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentHandler struct {
	cfg *config.Config
}

func NewStudentHandler(cfg *config.Config) *StudentHandler {
	return &StudentHandler{cfg: cfg}
}

type StudentProfileRequest struct {
	StudentID      string  `json:"student_id" validate:"required"`
	FacultyCode    string  `json:"faculty_code" validate:"required"`
	DepartmentCode string  `json:"department_code" validate:"required"`
	YearLevel      int     `json:"year_level" validate:"required,min=1,max=8"`
	GPA            float64 `json:"gpa" validate:"required,min=0,max=4"`
	AdmissionYear  int     `json:"admission_year" validate:"required"`
	GraduationYear int     `json:"graduation_year"`
}

// PriorityScoreRequest represents the request for calculating priority score
type PriorityScoreRequest struct {
	GPA           float64 `json:"gpa" validate:"required,min=0,max=4"`
	FamilyIncome  float64 `json:"family_income" validate:"required,min=0"`
	ActivityCount int     `json:"activity_count" validate:"min=0"`
}

// PriorityScoreResponse represents the response with calculated score breakdown
type PriorityScoreResponse struct {
	TotalScore      float64  `json:"total_score"`
	GPAScore        float64  `json:"gpa_score"`
	FinancialScore  float64  `json:"financial_score"`
	ActivityScore   float64  `json:"activity_score"`
	ScoreLevel      string   `json:"score_level"`
	Recommendations []string `json:"recommendations"`
}

// CalculatePriorityScore calculates priority score based on PROJECT_RULES.md algorithm
// @Summary Calculate priority score
// @Description Calculate scholarship priority score using GPA (40%), financial need (30%), activities (30%)
// @Tags Student Profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body PriorityScoreRequest true "Score calculation data"
// @Success 200 {object} object{data=PriorityScoreResponse}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /student/calculate-score [post]
func (h *StudentHandler) CalculatePriorityScore(c *fiber.Ctx) error {
	var req PriorityScoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input ranges
	if req.GPA < 0 || req.GPA > 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "GPA must be between 0.00 and 4.00",
		})
	}

	if req.FamilyIncome < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Family income must be positive",
		})
	}

	// Calculate scores using algorithm from PROJECT_RULES.md
	gpaScore := calculateGPAScore(req.GPA)
	financialScore := calculateFinancialScore(req.FamilyIncome)
	activityScore := calculateActivityScore(req.ActivityCount)

	// Calculate weighted total (GPA 40%, Financial 30%, Activities 30%)
	totalScore := (gpaScore * 0.4) + (financialScore * 0.3) + (activityScore * 0.3)
	totalScore = math.Round(totalScore*100) / 100

	// Determine score level
	scoreLevel := getScoreLevel(totalScore)

	// Generate recommendations
	recommendations := generateRecommendations(req.GPA, req.FamilyIncome, req.ActivityCount, totalScore)

	response := PriorityScoreResponse{
		TotalScore:      totalScore,
		GPAScore:        gpaScore,
		FinancialScore:  financialScore,
		ActivityScore:   activityScore,
		ScoreLevel:      scoreLevel,
		Recommendations: recommendations,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Priority score calculated successfully",
	})
}

// calculateGPAScore implements GPA scoring algorithm
func calculateGPAScore(gpa float64) float64 {
	// GPA 4.00 = 100 points, GPA 2.00 = 50 points (linear scale)
	if gpa >= 4.0 {
		return 100.0
	}
	if gpa <= 2.0 {
		return 50.0
	}
	return 50.0 + ((gpa-2.0)/2.0)*50.0
}

// calculateFinancialScore implements financial need scoring
func calculateFinancialScore(income float64) float64 {
	// Lower income = higher score (inverse relationship)
	// Income <= 15,000 = 100 points, Income >= 50,000 = 20 points
	if income <= 15000 {
		return 100.0
	}
	if income >= 50000 {
		return 20.0
	}
	return 100.0 - ((income-15000)/(50000-15000))*80.0
}

// calculateActivityScore implements activity scoring
func calculateActivityScore(activityCount int) float64 {
	// Each activity = 20 points, maximum 100 points
	score := float64(activityCount * 20)
	if score > 100 {
		return 100.0
	}
	return score
}

// getScoreLevel determines score level description
func getScoreLevel(score float64) string {
	if score >= 80 {
		return "สูงมาก"
	}
	if score >= 60 {
		return "สูง"
	}
	if score >= 40 {
		return "ปานกลาง"
	}
	return "ต่ำ"
}

// generateRecommendations provides improvement suggestions
func generateRecommendations(gpa, income float64, activities int, totalScore float64) []string {
	var recommendations []string

	if gpa < 3.0 {
		recommendations = append(recommendations, "เพิ่มเกรดเฉลี่ยให้สูงขึ้นเพื่อเพิ่มโอกาสได้รับทุน")
	}

	if activities < 3 {
		recommendations = append(recommendations, "เข้าร่วมกิจกรรมเสริมหลักสูตรเพิ่มเติมเพื่อเพิ่มคะแนน")
	}

	if income > 30000 {
		recommendations = append(recommendations, "พิจารณาสมัครทุนประเภทเรียนดีมากกว่าทุนขาดแคลน")
	}

	if totalScore < 60 {
		recommendations = append(recommendations, "ควรปรับปรุงผลการเรียนและเพิ่มกิจกรรมเพื่อเพิ่มโอกาสได้รับทุน")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "คะแนนของคุณอยู่ในระดับดี สามารถสมัครทุนได้อย่างมั่นใจ")
	}

	return recommendations
}

// CreateStudentProfile creates or updates student profile
// @Summary Create/Update student profile
// @Description Create or update student academic profile
// @Tags Student Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body StudentProfileRequest true "Student profile data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /student/profile [post]
func (h *StudentHandler) CreateStudentProfile(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()

	var req StudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if student profile already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM students WHERE user_id = $1)`
	database.DB.QueryRow(checkQuery, userID).Scan(&exists)

	if exists {
		// Update existing profile
		updateQuery := `UPDATE students SET 
			faculty_code = $1, department_code = $2, year_level = $3, 
			gpa = $4, admission_year = $5, graduation_year = $6
			WHERE user_id = $7`

		_, err := database.DB.Exec(updateQuery,
			req.FacultyCode, req.DepartmentCode, req.YearLevel,
			req.GPA, req.AdmissionYear, req.GraduationYear, userID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update student profile",
			})
		}
	} else {
		// Create new profile
		insertQuery := `INSERT INTO students 
			(student_id, user_id, faculty_code, department_code, year_level, gpa, admission_year, graduation_year)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := database.DB.Exec(insertQuery,
			req.StudentID, userID, req.FacultyCode, req.DepartmentCode,
			req.YearLevel, req.GPA, req.AdmissionYear, req.GraduationYear)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create student profile",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Student profile saved successfully",
		"data":    req,
	})
}

// GetStudentProfile retrieves current user's student profile
// @Summary Get student profile
// @Description Get current user's student academic profile
// @Tags Student Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /student/profile [get]
func (h *StudentHandler) GetStudentProfile(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()

	query := `SELECT student_id, faculty_code, department_code, year_level, 
		gpa, admission_year, graduation_year, student_status
		FROM students WHERE user_id = $1`

	var profile struct {
		StudentID      string  `json:"student_id"`
		FacultyCode    string  `json:"faculty_code"`
		DepartmentCode string  `json:"department_code"`
		YearLevel      int     `json:"year_level"`
		GPA            float64 `json:"gpa"`
		AdmissionYear  int     `json:"admission_year"`
		GraduationYear *int    `json:"graduation_year"`
		StudentStatus  string  `json:"student_status"`
	}

	err := database.DB.QueryRow(query, userID).Scan(
		&profile.StudentID, &profile.FacultyCode, &profile.DepartmentCode,
		&profile.YearLevel, &profile.GPA, &profile.AdmissionYear,
		&profile.GraduationYear, &profile.StudentStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Student profile not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch student profile",
		})
	}

	return c.JSON(fiber.Map{
		"data": profile,
	})
}

// GetStudentApplicationHistory retrieves student's application history
// @Summary Get application history
// @Description Get student's scholarship application history
// @Tags Student Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Router /student/application-history [get]
func (h *StudentHandler) GetStudentApplicationHistory(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()

	query := `SELECT sa.application_id, sa.scholarship_id, sa.application_status,
		sa.submitted_at, sa.priority_score, s.scholarship_name, s.amount,
		COALESCE(sal.allocated_amount, 0) as allocated_amount,
		sal.allocation_status
		FROM scholarship_applications sa
		JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		JOIN students st ON sa.student_id = st.student_id
		LEFT JOIN scholarship_allocations sal ON sa.application_id = sal.application_id
		WHERE st.user_id = $1
		ORDER BY sa.submitted_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application history",
		})
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var app map[string]interface{} = make(map[string]interface{})
		var applicationID, scholarshipID int
		var applicationStatus, scholarshipName string
		var submittedAt *string
		var priorityScore, amount, allocatedAmount *float64
		var allocationStatus *string

		err := rows.Scan(
			&applicationID, &scholarshipID, &applicationStatus,
			&submittedAt, &priorityScore, &scholarshipName, &amount,
			&allocatedAmount, &allocationStatus,
		)
		if err != nil {
			continue
		}

		app["application_id"] = applicationID
		app["scholarship_id"] = scholarshipID
		app["application_status"] = applicationStatus
		app["submitted_at"] = submittedAt
		app["priority_score"] = priorityScore
		app["scholarship_name"] = scholarshipName
		app["scholarship_amount"] = amount
		app["allocated_amount"] = allocatedAmount
		app["allocation_status"] = allocationStatus

		history = append(history, app)
	}

	return c.JSON(fiber.Map{
		"data": history,
	})
}

// GetEligibleScholarships returns scholarships the student is eligible for
// @Summary Get eligible scholarships
// @Description Get scholarships that the student meets criteria for
// @Tags Student Profile
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param gpa query number false "Student's current GPA"
// @Param income query number false "Student's family income"
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Router /student/eligible-scholarships [get]
func (h *StudentHandler) GetEligibleScholarships(c *fiber.Ctx) error {
	// Get query parameters
	gpaStr := c.Query("gpa", "0")
	incomeStr := c.Query("income", "0")

	gpa, _ := strconv.ParseFloat(gpaStr, 64)
	income, _ := strconv.ParseFloat(incomeStr, 64)

	// Mock eligible scholarships based on criteria
	scholarships := []map[string]interface{}{
		{
			"scholarship_id":   1,
			"scholarship_name": "ทุนวิจัยระดับปริญญาตรี",
			"amount":           25000,
			"deadline":         "2024-12-31",
			"days_left":        16,
			"type":             "research",
			"is_eligible":      gpa >= 3.25,
			"eligibility_reason": func() string {
				if gpa >= 3.25 {
					return "เข้าเกณฑ์คุณสมบัติ"
				}
				return "เกรดเฉลี่ยต้องมากกว่า 3.25"
			}(),
		},
		{
			"scholarship_id":   2,
			"scholarship_name": "ทุนความเป็นเลิศทางวิชาการ",
			"amount":           30000,
			"deadline":         "2025-01-15",
			"days_left":        31,
			"type":             "excellence",
			"is_eligible":      gpa >= 3.75,
			"eligibility_reason": func() string {
				if gpa >= 3.75 {
					return "เข้าเกณฑ์คุณสมบัติ"
				}
				return "เกรดเฉลี่ยต้องมากกว่า 3.75"
			}(),
		},
		{
			"scholarship_id":   3,
			"scholarship_name": "ทุนช่วยเหลือการศึกษา",
			"amount":           12000,
			"deadline":         "2024-12-25",
			"days_left":        10,
			"type":             "financial_aid",
			"is_eligible":      income < 30000,
			"eligibility_reason": func() string {
				if income < 30000 {
					return "เข้าเกณฑ์คุณสมบัติ"
				}
				return "รายได้ครอบครัวต้องน้อยกว่า 30,000 บาท/เดือน"
			}(),
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    scholarships,
		"message": "Eligible scholarships retrieved successfully",
	})
}
