package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"scholarship-system/internal/models"
)

// AuthEnhancedHandler handles enhanced authentication operations
type AuthEnhancedHandler struct {
	// Will add repository/usecase dependencies later
}

func NewAuthEnhancedHandler() *AuthEnhancedHandler {
	return &AuthEnhancedHandler{}
}

// @Summary Register Student
// @Description Register a new student account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.StudentRegistrationRequest true "Student registration data"
// @Success 200 {object} models.RegistrationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/register/student [post]
func (h *AuthEnhancedHandler) RegisterStudent(c *fiber.Ctx) error {
	var req models.StudentRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if req.StudentID == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing required fields",
		})
	}

	// Validate email domain (university email)
	if !strings.HasSuffix(req.Email, "@student.mahidol.ac.th") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email must be a valid university email",
		})
	}

	// Validate password strength
	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 8 characters long",
		})
	}

	// TODO: Call usecase to create student user
	// user, err := h.userUsecase.RegisterStudent(c.Context(), &req)

	// Mock response for now
	response := models.RegistrationResponse{
		Success:          true,
		Message:          "Registration successful. Please check your email for verification.",
		UserID:           "mock-user-id",
		VerificationSent: true,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// @Summary Register Staff
// @Description Register a new staff account (admin only)
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.StaffRegistrationRequest true "Staff registration data"
// @Success 200 {object} models.RegistrationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/register/staff [post]
func (h *AuthEnhancedHandler) RegisterStaff(c *fiber.Ctx) error {
	// Check if user is admin
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	role := userClaims["role"].(string)

	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Only administrators can register staff accounts",
		})
	}

	var req models.StaffRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate required fields
	if req.EmployeeID == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing required fields",
		})
	}

	// Validate role
	validRoles := []string{"scholarship_officer", "interviewer", "admin"}
	roleValid := false
	for _, validRole := range validRoles {
		if req.Role == validRole {
			roleValid = true
			break
		}
	}
	if !roleValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid role specified",
		})
	}

	// TODO: Call usecase to create staff user
	// user, err := h.userUsecase.RegisterStaff(c.Context(), &req)

	response := models.RegistrationResponse{
		Success:          true,
		Message:          "Staff registration successful",
		UserID:           "mock-staff-id",
		VerificationSent: false,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// @Summary Verify Email
// @Description Verify user email with verification token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} models.EmailVerificationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/verify-email [get]
func (h *AuthEnhancedHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Verification token is required",
		})
	}

	// TODO: Call usecase to verify email
	// err := h.userUsecase.VerifyEmail(c.Context(), token)

	response := models.EmailVerificationResponse{
		Success: true,
		Message: "Email verified successfully",
	}

	return c.JSON(response)
}

// @Summary Request Password Reset
// @Description Request password reset via email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.PasswordResetRequest true "Password reset request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/password-reset [post]
func (h *AuthEnhancedHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var req models.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email is required",
		})
	}

	// TODO: Call usecase to send password reset email
	// err := h.userUsecase.RequestPasswordReset(c.Context(), req.Email)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "If the email exists, a password reset link has been sent",
	})
}

// @Summary Reset Password
// @Description Reset password with reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.PasswordResetConfirm true "Password reset confirmation"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/password-reset/confirm [post]
func (h *AuthEnhancedHandler) ConfirmPasswordReset(c *fiber.Ctx) error {
	var req models.PasswordResetConfirm
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Token and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 8 characters long",
		})
	}

	// TODO: Call usecase to reset password
	// err := h.userUsecase.ConfirmPasswordReset(c.Context(), &req)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password reset successfully",
	})
}

// @Summary Change Password
// @Description Change user password (authenticated)
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.PasswordChangeRequest true "Password change request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/change-password [post]
func (h *AuthEnhancedHandler) ChangePassword(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	_ = userClaims["user_id"].(string) // Will be used when implementing usecase

	var req models.PasswordChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Current password and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "New password must be at least 8 characters long",
		})
	}

	// TODO: Call usecase to change password
	// err := h.userUsecase.ChangePassword(c.Context(), userID, &req)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

// @Summary Get User Profile
// @Description Get authenticated user profile with completion status
// @Tags Profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/profile [get]
func (h *AuthEnhancedHandler) GetProfile(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	_ = userClaims["user_id"].(string) // Will be used when implementing usecase

	// TODO: Call usecase to get enhanced user profile
	// user, err := h.userUsecase.GetEnhancedProfile(c.Context(), userID)

	// Mock response for now
	mockUser := &models.EnhancedUser{
		EmailVerified:               false,
		ProfileCompleted:            false,
		ProfileCompletionPercentage: 60,
	}

	response := models.ProfileResponse{
		EnhancedUser:    mockUser,
		CompletionSteps: mockUser.GetCompletionSteps(),
	}

	return c.JSON(response)
}

// @Summary Update User Profile
// @Description Update authenticated user profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.ProfileUpdateRequest true "Profile update data"
// @Success 200 {object} models.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/profile [put]
func (h *AuthEnhancedHandler) UpdateProfile(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	_ = userClaims["user_id"].(string) // Will be used when implementing usecase

	var req models.ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// TODO: Call usecase to update profile
	// user, err := h.userUsecase.UpdateProfile(c.Context(), userID, &req)

	// Mock response
	mockUser := &models.EnhancedUser{
		EmailVerified:               false,
		ProfileCompleted:            true,
		ProfileCompletionPercentage: 85,
	}

	response := models.ProfileResponse{
		EnhancedUser:    mockUser,
		CompletionSteps: mockUser.GetCompletionSteps(),
	}

	return c.JSON(response)
}

// @Summary Upload Avatar
// @Description Upload user avatar image
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/profile/avatar [post]
func (h *AuthEnhancedHandler) UploadAvatar(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	_ = userClaims["user_id"].(string) // Will be used when implementing usecase

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Avatar file is required",
		})
	}

	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	fileType := file.Header.Get("Content-Type")
	typeAllowed := false
	for _, allowedType := range allowedTypes {
		if fileType == allowedType {
			typeAllowed = true
			break
		}
	}

	if !typeAllowed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Only JPEG, PNG, and GIF files are allowed",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "File size must be less than 5MB",
		})
	}

	// TODO: Save file and update user avatar URL
	// avatarURL, err := h.fileService.SaveAvatar(c.Context(), userID, file)
	// err = h.userUsecase.UpdateAvatarURL(c.Context(), userID, avatarURL)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Avatar uploaded successfully",
		"avatar_url": "/uploads/avatars/mock-avatar.jpg",
	})
}

// @Summary Get Login History
// @Description Get user login history
// @Tags Authentication
// @Produce json
// @Security ApiKeyAuth
// @Param limit query int false "Number of records to return" default(10)
// @Success 200 {array} models.LoginHistory
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/login-history [get]
func (h *AuthEnhancedHandler) GetLoginHistory(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	_ = userClaims["user_id"].(string) // Will be used when implementing usecase

	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50 // Max 50 records
	}

	// TODO: Call usecase to get login history
	// history, err := h.userUsecase.GetLoginHistory(c.Context(), userID, limit)

	// Mock response
	mockHistory := []models.LoginHistory{}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    mockHistory,
	})
}

// @Summary Logout
// @Description Logout user and invalidate session
// @Tags Authentication
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *AuthEnhancedHandler) Logout(c *fiber.Ctx) error {
	// Get session token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Authorization header required",
		})
	}

	_ = strings.Replace(authHeader, "Bearer ", "", 1) // Will be used when implementing usecase

	// TODO: Call usecase to invalidate session
	// err := h.userUsecase.InvalidateSession(c.Context(), tokenString)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
