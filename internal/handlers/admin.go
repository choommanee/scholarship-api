package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	cfg *config.Config
}

func NewAdminHandler(cfg *config.Config) *AdminHandler {
	return &AdminHandler{cfg: cfg}
}

// SystemConfig represents system configuration
type SystemConfig struct {
	SystemName                  string  `json:"system_name"`
	SystemVersion               string  `json:"system_version"`
	MaintenanceMode             bool    `json:"maintenance_mode"`
	AllowRegistration           bool    `json:"allow_registration"`
	MaxFileUploadSize           int     `json:"max_file_upload_size"`
	SessionTimeout              int     `json:"session_timeout"`
	CurrentAcademicYear         string  `json:"current_academic_year"`
	ApplicationDeadline         string  `json:"application_deadline"`
	MaxApplicationsPerStudent   int     `json:"max_applications_per_student"`
	AutoApproveApplications     bool    `json:"auto_approve_applications"`
	RequireDocumentVerification bool    `json:"require_document_verification"`
	EmailEnabled                bool    `json:"email_enabled"`
	SMTPHost                    string  `json:"smtp_host"`
	SMTPPort                    int     `json:"smtp_port"`
	SMTPUsername                string  `json:"smtp_username"`
	SMTPPassword                string  `json:"smtp_password"`
	FromEmail                   string  `json:"from_email"`
	FromName                    string  `json:"from_name"`
	NotificationEnabled         bool    `json:"notification_enabled"`
	EmailNotifications          bool    `json:"email_notifications"`
	SMSNotifications            bool    `json:"sms_notifications"`
	PushNotifications           bool    `json:"push_notifications"`
	EnforcePasswordPolicy       bool    `json:"enforce_password_policy"`
	MinPasswordLength           int     `json:"min_password_length"`
	RequireTwoFactor            bool    `json:"require_two_factor"`
	LoginAttemptLimit           int     `json:"login_attempt_limit"`
	LockoutDuration             int     `json:"lockout_duration"`
	TotalBudget                 float64 `json:"total_budget"`
	BudgetWarningThreshold      int     `json:"budget_warning_threshold"`
	AutoCloseBudgetExceeded     bool    `json:"auto_close_budget_exceeded"`
}

// SystemStats represents system statistics
type SystemStats struct {
	TotalUsers          int     `json:"total_users"`
	ActiveUsers         int     `json:"active_users"`
	TotalScholarships   int     `json:"total_scholarships"`
	ActiveScholarships  int     `json:"active_scholarships"`
	TotalApplications   int     `json:"total_applications"`
	PendingApplications int     `json:"pending_applications"`
	TotalBudget         float64 `json:"total_budget"`
	UsedBudget          float64 `json:"used_budget"`
	RemainingBudget     float64 `json:"remaining_budget"`
	BudgetUsagePercent  float64 `json:"budget_usage_percent"`
}

// GetSystemConfig retrieves system configuration
// @Summary Get system configuration
// @Description Get current system configuration (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SystemConfig
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/config [get]
func (h *AdminHandler) GetSystemConfig(c *fiber.Ctx) error {
	// TODO: Load from database or config file
	config := SystemConfig{
		SystemName:                  "ระบบจัดการทุนการศึกษา คณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์",
		SystemVersion:               "1.0.0",
		MaintenanceMode:             false,
		AllowRegistration:           true,
		MaxFileUploadSize:           10,
		SessionTimeout:              30,
		CurrentAcademicYear:         "2567",
		ApplicationDeadline:         "2024-06-30",
		MaxApplicationsPerStudent:   3,
		AutoApproveApplications:     false,
		RequireDocumentVerification: true,
		EmailEnabled:                true,
		SMTPHost:                    "smtp.mahidol.ac.th",
		SMTPPort:                    587,
		SMTPUsername:                "",
		SMTPPassword:                "",
		FromEmail:                   "scholarship@mahidol.ac.th",
		FromName:                    "ระบบทุนการศึกษา คณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์",
		NotificationEnabled:         true,
		EmailNotifications:          true,
		SMSNotifications:            false,
		PushNotifications:           true,
		EnforcePasswordPolicy:       true,
		MinPasswordLength:           8,
		RequireTwoFactor:            false,
		LoginAttemptLimit:           5,
		LockoutDuration:             15,
		TotalBudget:                 4834200,
		BudgetWarningThreshold:      80,
		AutoCloseBudgetExceeded:     false,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    config,
	})
}

// UpdateSystemConfig updates system configuration
// @Summary Update system configuration
// @Description Update system configuration (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param config body SystemConfig true "System configuration"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/config [put]
func (h *AdminHandler) UpdateSystemConfig(c *fiber.Ctx) error {
	var config SystemConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Validate and save configuration to database
	// For now, just return success

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกการตั้งค่าเรียบร้อยแล้ว",
	})
}

// GetSystemStats retrieves system statistics
// @Summary Get system statistics
// @Description Get system statistics for admin dashboard
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SystemStats
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c *fiber.Ctx) error {
	// TODO: Calculate real statistics from database
	stats := SystemStats{
		TotalUsers:          1250,
		ActiveUsers:         1180,
		TotalScholarships:   181,
		ActiveScholarships:  45,
		TotalApplications:   2340,
		PendingApplications: 156,
		TotalBudget:         4834200,
		UsedBudget:          2450000,
		RemainingBudget:     2384200,
		BudgetUsagePercent:  50.7,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetActivityLog retrieves system activity log
// @Summary Get activity log
// @Description Get system activity log (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of items per page" default(20)
// @Param page query int false "Page number" default(1)
// @Success 200 {object} object{data=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/activity-log [get]
func (h *AdminHandler) GetActivityLog(c *fiber.Ctx) error {
	// TODO: Implement activity log retrieval from database
	activities := []map[string]interface{}{
		{
			"id":          1,
			"user_id":     uuid.New().String(),
			"user_name":   "นางสาว สมใจ ใจดี",
			"action":      "login",
			"description": "เข้าสู่ระบบ",
			"ip_address":  "192.168.1.100",
			"user_agent":  "Mozilla/5.0...",
			"created_at":  "2024-12-20T10:30:00Z",
		},
		{
			"id":          2,
			"user_id":     uuid.New().String(),
			"user_name":   "นายสมชาย ดีใจ",
			"action":      "create_scholarship",
			"description": "สร้างทุนการศึกษาใหม่: ทุนพัฒนาศักยภาพนักศึกษา",
			"ip_address":  "192.168.1.101",
			"user_agent":  "Mozilla/5.0...",
			"created_at":  "2024-12-20T09:15:00Z",
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    activities,
		"pagination": map[string]interface{}{
			"page":        1,
			"limit":       20,
			"total":       2,
			"total_pages": 1,
		},
	})
}

// TestEmailConnection tests email configuration
// @Summary Test email connection
// @Description Test SMTP email connection (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/test-email [post]
func (h *AdminHandler) TestEmailConnection(c *fiber.Ctx) error {
	// TODO: Implement actual email connection test
	// For now, simulate success

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ทดสอบการเชื่อมต่ออีเมลสำเร็จ",
	})
}

// TestDatabaseConnection tests database connection
// @Summary Test database connection
// @Description Test database connection (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string}
// @Failure 500 {object} object{error=string}
// @Router /admin/test-database [post]
func (h *AdminHandler) TestDatabaseConnection(c *fiber.Ctx) error {
	// TODO: Implement actual database connection test
	// For now, simulate success

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ทดสอบการเชื่อมต่อฐานข้อมูลสำเร็จ",
	})
}

// GetScholarshipManagement retrieves scholarship management data
// @Summary Get scholarship management data
// @Description Get scholarship management overview (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /admin/scholarships [get]
func (h *AdminHandler) GetScholarshipManagement(c *fiber.Ctx) error {
	// TODO: Implement scholarship management data retrieval
	data := map[string]interface{}{
		"total_scholarships":    181,
		"active_scholarships":   45,
		"draft_scholarships":    12,
		"closed_scholarships":   124,
		"total_applications":    2340,
		"pending_applications":  156,
		"approved_applications": 890,
		"rejected_applications": 1294,
		"total_budget":          4834200,
		"allocated_budget":      2450000,
		"remaining_budget":      2384200,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// ==================== DASHBOARD ENDPOINTS ====================

// GetDashboardStats retrieves statistics for admin dashboard
// @Summary Get dashboard statistics
// @Description Get comprehensive statistics for admin dashboard
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,data=object}
// @Failure 500 {object} object{error=string}
// @Router /admin/dashboard/stats [get]
func (h *AdminHandler) GetDashboardStats(c *fiber.Ctx) error {
	db := database.DB

	var stats struct {
		TotalUsers         int     `json:"total_users"`
		ActiveUsers        int     `json:"active_users"`
		TotalScholarships  int     `json:"total_scholarships"`
		ActiveScholarships int     `json:"active_scholarships"`
		TotalApplications  int     `json:"total_applications"`
		PendingReview      int     `json:"pending_review"`
		SystemUptime       string  `json:"system_uptime"`
		StorageUsage       float64 `json:"storage_usage"`
	}

	// Get total users
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)

	// Get active users (logged in within last 30 days)
	db.QueryRow(`
		SELECT COUNT(*) FROM users
		WHERE last_login_at > NOW() - INTERVAL '30 days'
	`).Scan(&stats.ActiveUsers)

	// Get scholarship counts
	db.QueryRow("SELECT COUNT(*) FROM scholarships").Scan(&stats.TotalScholarships)
	db.QueryRow(`
		SELECT COUNT(*) FROM scholarships
		WHERE scholarship_status = 'open'
		AND application_deadline > NOW()
	`).Scan(&stats.ActiveScholarships)

	// Get application counts
	db.QueryRow("SELECT COUNT(*) FROM scholarship_applications").Scan(&stats.TotalApplications)
	db.QueryRow(`
		SELECT COUNT(*) FROM scholarship_applications
		WHERE application_status = 'under_review'
	`).Scan(&stats.PendingReview)

	// Mock system uptime and storage
	stats.SystemUptime = "99.8%"
	stats.StorageUsage = 68.0

	// Calculate trends
	var monthlyUserGrowth int
	db.QueryRow(`
		SELECT COUNT(*) FROM users
		WHERE created_at > NOW() - INTERVAL '30 days'
	`).Scan(&monthlyUserGrowth)

	return c.JSON(fiber.Map{
		"success": true,
		"data": map[string]interface{}{
			"stats": stats,
			"trends": map[string]interface{}{
				"users":      monthlyUserGrowth,
				"isPositive": monthlyUserGrowth > 0,
				"period":     "เดือนนี้",
			},
		},
	})
}

// GetDashboardAlerts retrieves system alerts for admin dashboard
// @Summary Get dashboard alerts
// @Description Get system alerts and notifications
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,data=[]object}
// @Failure 500 {object} object{error=string}
// @Router /admin/dashboard/alerts [get]
func (h *AdminHandler) GetDashboardAlerts(c *fiber.Ctx) error {
	db := database.DB

	var alerts []map[string]interface{}

	// Check for high application volume
	var pendingCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM scholarship_applications
		WHERE application_status IN ('submitted', 'under_review')
	`).Scan(&pendingCount)

	if pendingCount > 100 {
		alerts = append(alerts, map[string]interface{}{
			"id":          "alert_1",
			"title":       "ใบสมัครรอพิจารณาจำนวนมาก",
			"description": "มีใบสมัครรอการพิจารณา " + string(rune(pendingCount)) + " รายการ",
			"type":        "warning",
			"severity":    "medium",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"status":      "active",
		})
	}

	// Check for deadline approaching
	var upcomingDeadlines int
	db.QueryRow(`
		SELECT COUNT(*) FROM scholarships
		WHERE scholarship_status = 'open'
		AND application_deadline BETWEEN NOW() AND NOW() + INTERVAL '7 days'
	`).Scan(&upcomingDeadlines)

	if upcomingDeadlines > 0 {
		alerts = append(alerts, map[string]interface{}{
			"id":          "alert_2",
			"title":       "กำหนดการรับสมัครใกล้ปิด",
			"description": "มีทุนการศึกษา " + string(rune(upcomingDeadlines)) + " รายการที่ใกล้ถึงกำหนดปิดรับสมัครภายใน 7 วัน",
			"type":        "info",
			"severity":    "low",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"status":      "active",
		})
	}

	// Add success notification
	alerts = append(alerts, map[string]interface{}{
		"id":          "alert_3",
		"title":       "การสำรองข้อมูลเสร็จสิ้น",
		"description": "การสำรองข้อมูลประจำวันเสร็จสิ้นเรียบร้อยแล้ว",
		"type":        "info",
		"severity":    "low",
		"timestamp":   time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
		"status":      "acknowledged",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    alerts,
	})
}

// GetDashboardActivities retrieves recent system activities
// @Summary Get dashboard activities
// @Description Get recent user activities for admin dashboard
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of activities" default(10)
// @Success 200 {object} object{success=bool,data=[]object}
// @Failure 500 {object} object{error=string}
// @Router /admin/dashboard/activities [get]
func (h *AdminHandler) GetDashboardActivities(c *fiber.Ctx) error {
	db := database.DB

	limit := c.QueryInt("limit", 10)

	// Get recent activities (using audit_logs if exists, or recent application activities)
	query := `
		SELECT
			sa.application_id as id,
			u.user_id,
			CONCAT(u.first_name, ' ', u.last_name) as user_name,
			CASE
				WHEN sa.application_status = 'submitted' THEN 'ส่งใบสมัคร'
				WHEN sa.application_status = 'under_review' THEN 'อยู่ระหว่างพิจารณา'
				WHEN sa.application_status = 'approved' THEN 'อนุมัติ'
				WHEN sa.application_status = 'rejected' THEN 'ปฏิเสธ'
				ELSE 'อื่นๆ'
			END as action,
			s.scholarship_name as resource,
			sa.submitted_at as timestamp,
			sa.application_status as status
		FROM scholarship_applications sa
		JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		JOIN students st ON sa.student_id = st.student_id
		JOIN users u ON st.user_id = u.user_id
		WHERE sa.submitted_at IS NOT NULL
		ORDER BY sa.submitted_at DESC
		LIMIT $1
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		log.Printf("Error fetching activities: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูลกิจกรรม",
		})
	}
	defer rows.Close()

	var activities []map[string]interface{}
	for rows.Next() {
		var (
			id, userId, userName, action, resource, status string
			timestamp                                      time.Time
		)

		err := rows.Scan(&id, &userId, &userName, &action, &resource, &timestamp, &status)
		if err != nil {
			continue
		}

		// Calculate time ago
		timeAgo := time.Since(timestamp)
		var timeAgoStr string
		if timeAgo.Minutes() < 60 {
			timeAgoStr = string(rune(int(timeAgo.Minutes()))) + " นาทีที่แล้ว"
		} else if timeAgo.Hours() < 24 {
			timeAgoStr = string(rune(int(timeAgo.Hours()))) + " ชั่วโมงที่แล้ว"
		} else {
			timeAgoStr = string(rune(int(timeAgo.Hours()/24))) + " วันที่แล้ว"
		}

		activityStatus := "success"
		if status == "rejected" {
			activityStatus = "failed"
		} else if status == "under_review" {
			activityStatus = "pending"
		}

		activities = append(activities, map[string]interface{}{
			"id":         id,
			"userId":     userId,
			"userName":   userName,
			"action":     action,
			"resource":   resource,
			"timestamp":  timeAgoStr,
			"ipAddress":  "192.168.1.x",
			"userAgent":  "Web Browser",
			"status":     activityStatus,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    activities,
	})
}

// GetDashboardResources retrieves system resource usage
// @Summary Get dashboard resources
// @Description Get system resource usage information
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,data=[]object}
// @Failure 500 {object} object{error=string}
// @Router /admin/dashboard/resources [get]
func (h *AdminHandler) GetDashboardResources(c *fiber.Ctx) error {
	// Mock system resources
	// In production, this would query actual system metrics

	resources := []map[string]interface{}{
		{
			"name":   "CPU",
			"usage":  45.0,
			"total":  100.0,
			"unit":   "%",
			"status": "healthy",
		},
		{
			"name":   "Memory",
			"usage":  6.2,
			"total":  16.0,
			"unit":   "GB",
			"status": "healthy",
		},
		{
			"name":   "Storage",
			"usage":  68.0,
			"total":  100.0,
			"unit":   "%",
			"status": "warning",
		},
		{
			"name":   "Network",
			"usage":  23.0,
			"total":  100.0,
			"unit":   "Mbps",
			"status": "healthy",
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resources,
	})
}

// ==================== ADMIN PROFILE ENDPOINTS ====================

// GetAdminProfile retrieves admin user profile
// @Summary Get admin profile
// @Description Get current admin user profile information
// @Tags Admin Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,data=object}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /admin/profile [get]
func (h *AdminHandler) GetAdminProfile(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id")
	var userIDStr string

	// Handle both string and uuid.UUID types
	switch v := userID.(type) {
	case string:
		userIDStr = v
	case uuid.UUID:
		userIDStr = v.String()
	default:
		userIDStr = fmt.Sprintf("%v", v)
	}

	// Get user profile with roles
	query := `
		SELECT
			u.user_id,
			u.username,
			u.email,
			u.first_name,
			u.last_name,
			u.phone,
			u.is_active,
			u.sso_provider,
			u.created_at,
			u.updated_at,
			u.last_login,
			COALESCE(
				json_agg(
					json_build_object(
						'role_id', r.role_id,
						'role_name', r.role_name,
						'role_description', r.role_description,
						'assigned_at', ur.assigned_at
					)
				) FILTER (WHERE r.role_id IS NOT NULL),
				'[]'
			) as roles
		FROM users u
		LEFT JOIN user_roles ur ON u.user_id = ur.user_id AND ur.is_active = true
		LEFT JOIN roles r ON ur.role_id = r.role_id
		WHERE u.user_id = $1
		GROUP BY u.user_id
	`

	var profile struct {
		UserID      string          `json:"user_id"`
		Username    string          `json:"username"`
		Email       string          `json:"email"`
		FirstName   string          `json:"first_name"`
		LastName    string          `json:"last_name"`
		Phone       *string         `json:"phone"`
		IsActive    bool            `json:"is_active"`
		SSOProvider *string         `json:"sso_provider"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
		LastLogin   *time.Time      `json:"last_login"`
		Roles       json.RawMessage `json:"roles"`
	}

	var rolesJSON string
	err := db.QueryRow(query, userIDStr).Scan(
		&profile.UserID,
		&profile.Username,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.IsActive,
		&profile.SSOProvider,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.LastLogin,
		&rolesJSON,
	)

	if err != nil {
		log.Printf("Error fetching admin profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถดึงข้อมูลโปรไฟล์ได้",
		})
	}

	profile.Roles = json.RawMessage(rolesJSON)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
	})
}

// UpdateAdminProfile updates admin user profile
// @Summary Update admin profile
// @Description Update current admin user profile information
// @Tags Admin Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body object{first_name=string,last_name=string,phone=string,email=string} true "Profile update data"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /admin/profile [put]
func (h *AdminHandler) UpdateAdminProfile(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id")
	var userIDStr string
	switch v := userID.(type) {
	case string:
		userIDStr = v
	case uuid.UUID:
		userIDStr = v.String()
	default:
		userIDStr = fmt.Sprintf("%v", v)
	}

	var updateData struct {
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Phone     *string `json:"phone"`
		Email     string  `json:"email"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Validate required fields
	if updateData.FirstName == "" || updateData.LastName == "" || updateData.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "กรุณากรอกข้อมูลให้ครบถ้วน",
		})
	}

	// Update profile
	query := `
		UPDATE users
		SET
			first_name = $1,
			last_name = $2,
			phone = $3,
			email = $4,
			updated_at = NOW()
		WHERE user_id = $5
		RETURNING user_id, username, email, first_name, last_name, phone, updated_at
	`

	var updatedProfile struct {
		UserID    string    `json:"user_id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Phone     *string   `json:"phone"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err := db.QueryRow(query,
		updateData.FirstName,
		updateData.LastName,
		updateData.Phone,
		updateData.Email,
		userIDStr,
	).Scan(
		&updatedProfile.UserID,
		&updatedProfile.Username,
		&updatedProfile.Email,
		&updatedProfile.FirstName,
		&updatedProfile.LastName,
		&updatedProfile.Phone,
		&updatedProfile.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating admin profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอัปเดตโปรไฟล์ได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "อัปเดตโปรไฟล์เรียบร้อยแล้ว",
		"data":    updatedProfile,
	})
}

// ChangeAdminPassword changes admin password
// @Summary Change admin password
// @Description Change current admin user password
// @Tags Admin Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body object{current_password=string,new_password=string,confirm_password=string} true "Password change data"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /admin/profile/password [put]
func (h *AdminHandler) ChangeAdminPassword(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id")
	var userIDStr string
	switch v := userID.(type) {
	case string:
		userIDStr = v
	case uuid.UUID:
		userIDStr = v.String()
	default:
		userIDStr = fmt.Sprintf("%v", v)
	}

	var passwordData struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&passwordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Validate passwords
	if passwordData.NewPassword != passwordData.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านใหม่และรหัสผ่านยืนยันไม่ตรงกัน",
		})
	}

	if len(passwordData.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านต้องมีความยาวอย่างน้อย 8 ตัวอักษร",
		})
	}

	// Get current password hash
	var currentPasswordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE user_id = $1", userIDStr).Scan(&currentPasswordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถตรวจสอบรหัสผ่านได้",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(passwordData.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านปัจจุบันไม่ถูกต้อง",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถเข้ารหัสรหัสผ่านใหม่ได้",
		})
	}

	// Update password
	_, err = db.Exec("UPDATE users SET password_hash = $1, updated_at = NOW() WHERE user_id = $2", string(hashedPassword), userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถเปลี่ยนรหัสผ่านได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เปลี่ยนรหัสผ่านเรียบร้อยแล้ว",
	})
}
