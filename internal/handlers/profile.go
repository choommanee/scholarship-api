package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"

	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	cfg *config.Config
}

func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{cfg: cfg}
}

// ProfileResponse represents user profile with role-specific data
type ProfileResponse struct {
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
	// Student specific fields
	StudentData *StudentProfileData `json:"student_data,omitempty"`
	// Staff specific fields (Officer/Interviewer/Admin)
	StaffData *StaffProfileData `json:"staff_data,omitempty"`
}

type StudentProfileData struct {
	StudentID      string   `json:"student_id"`
	FacultyCode    *string  `json:"faculty_code"`
	DepartmentCode *string  `json:"department_code"`
	YearLevel      *int     `json:"year_level"`
	CurrentGPA     *float64 `json:"current_gpa"`
	AdmissionYear  *int     `json:"admission_year"`
	StudentStatus  *string  `json:"student_status"`
	AdvisorName    *string  `json:"advisor_name"`
}

type StaffProfileData struct {
	Department *string `json:"department"`
	Position   *string `json:"position"`
}

// GetProfile retrieves user profile based on their role
// @Summary Get user profile
// @Description Get current user profile information with role-specific data
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{success=bool,data=ProfileResponse}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id").(string)

	// Get base user profile with roles
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
			u.last_login_at,
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

	var profile ProfileResponse
	var rolesJSON string

	err := db.QueryRow(query, userID).Scan(
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
		log.Printf("Error fetching user profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถดึงข้อมูลโปรไฟล์ได้",
		})
	}

	profile.Roles = json.RawMessage(rolesJSON)

	// Check if user is a student and get student-specific data
	studentData, err := h.getStudentData(db, userID)
	if err == nil && studentData != nil {
		profile.StudentData = studentData
	}

	// For staff roles (admin/officer/interviewer), could add staff-specific data here
	// For now, we'll leave it as optional future enhancement

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
	})
}

// getStudentData retrieves student-specific profile data
func (h *ProfileHandler) getStudentData(db *sql.DB, userID string) (*StudentProfileData, error) {
	query := `
		SELECT
			s.student_id,
			s.faculty_code,
			s.department_code,
			s.year_level,
			s.current_gpa,
			s.admission_year,
			s.student_status,
			CONCAT(u.first_name, ' ', u.last_name) as advisor_name
		FROM students s
		LEFT JOIN users u ON s.advisor_id = u.user_id
		WHERE s.user_id = $1
	`

	var data StudentProfileData
	err := db.QueryRow(query, userID).Scan(
		&data.StudentID,
		&data.FacultyCode,
		&data.DepartmentCode,
		&data.YearLevel,
		&data.CurrentGPA,
		&data.AdmissionYear,
		&data.StudentStatus,
		&data.AdvisorName,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not a student
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateProfile updates user profile information
// @Summary Update user profile
// @Description Update current user profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body object{first_name=string,last_name=string,phone=string,email=string} true "Profile update data"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id").(string)

	var req struct {
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Phone     *string `json:"phone"`
		Email     string  `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "กรุณากรอกข้อมูลให้ครบถ้วน",
		})
	}

	// Check if email is already used by another user
	var existingUserID string
	err := db.QueryRow("SELECT user_id FROM users WHERE email = $1 AND user_id != $2", req.Email, userID).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "อีเมลนี้ถูกใช้งานแล้ว",
		})
	}

	// Update user profile
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, phone = $3, email = $4, updated_at = NOW()
		WHERE user_id = $5
	`

	_, err = db.Exec(query, req.FirstName, req.LastName, req.Phone, req.Email, userID)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอัพเดทโปรไฟล์ได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "อัพเดทโปรไฟล์สำเร็จ",
	})
}

// ChangePassword changes user password
// @Summary Change password
// @Description Change current user password
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body object{current_password=string,new_password=string,confirm_password=string} true "Password change data"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /profile/password [put]
func (h *ProfileHandler) ChangePassword(c *fiber.Ctx) error {
	db := database.DB
	userID := c.Locals("user_id").(string)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ข้อมูลไม่ถูกต้อง",
		})
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "กรุณากรอกข้อมูลให้ครบถ้วน",
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านใหม่ไม่ตรงกัน",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร",
		})
	}

	// Get current password hash
	var currentHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE user_id = $1", userID).Scan(&currentHash)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถดำเนินการได้",
		})
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "รหัสผ่านปัจจุบันไม่ถูกต้อง",
		})
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถเปลี่ยนรหัสผ่านได้",
		})
	}

	// Update password
	_, err = db.Exec(`
		UPDATE users
		SET password_hash = $1, password_changed_at = NOW(), updated_at = NOW()
		WHERE user_id = $2
	`, string(newHash), userID)

	if err != nil {
		log.Printf("Error updating password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถเปลี่ยนรหัสผ่านได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เปลี่ยนรหัสผ่านสำเร็จ",
	})
}
