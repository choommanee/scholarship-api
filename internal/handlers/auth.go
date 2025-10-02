package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"scholarship-system/internal/config"
	"scholarship-system/internal/middleware"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type AuthHandler struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg:      cfg,
		userRepo: repository.NewUserRepository(),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone"`
	StudentID string `json:"student_id,omitempty"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Try to get user by email first, then by username
	var user *models.User
	var err error

	// Debug logging for request
	fmt.Printf("Debug - Login request - Email: %s, Username: %s\n", req.Email, req.Username)

	if req.Email != "" {
		fmt.Printf("Debug - Searching user by email: %s\n", req.Email)
		user, err = h.userRepo.GetUserWithRolesByEmail(req.Email)
	} else if req.Username != "" {
		fmt.Printf("Debug - Searching user by username: %s\n", req.Username)
		user, err = h.userRepo.GetUserWithRolesByUsername(req.Username)
	} else {
		fmt.Printf("Debug - No email or username provided\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email or username is required",
		})
	}

	fmt.Printf("Debug - Database query error: %v\n", err)
	if user != nil {
		fmt.Printf("Debug - User found - ID: %s, Email: %s, Username: %s, Active: %t\n",
			user.UserID.String(), user.Email, user.Username, user.IsActive)
	} else {
		fmt.Printf("Debug - No user found\n")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Debug - User not found in database (sql.ErrNoRows)\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials xx",
			})
		}
		fmt.Printf("Debug - Database error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if !user.IsActive {
		fmt.Printf("Debug - User account is disabled\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Account is disabled",
		})
	}

	// Debug logging for password comparison
	fmt.Printf("Debug - Password provided: '%s' (length: %d)\n", req.Password, len(req.Password))
	fmt.Printf("Debug - Password hash from DB: %s\n", user.PasswordHash)
	fmt.Printf("Debug - Starting bcrypt comparison...\n")

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		fmt.Printf("Debug - bcrypt comparison error: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	fmt.Printf("Debug - Password comparison successful\n")

	// Extract user roles
	var roles []string
	for _, userRole := range user.UserRoles {
		if userRole.IsActive && userRole.Role != nil {
			roles = append(roles, userRole.Role.RoleName)
		}
	}

	// Generate JWT token with longer expiration for testing
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7 days instead of 24 hours
	
	fmt.Printf("Debug - Token generation:\n")
	fmt.Printf("  Current time: %s\n", now.Format(time.RFC3339))
	fmt.Printf("  Expires at: %s\n", expiresAt.Format(time.RFC3339))
	fmt.Printf("  Duration: 7 days\n")
	
	claims := middleware.Claims{
		UserID:   user.UserID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.UserID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Update last login
	h.userRepo.UpdateLastLogin(user.UserID)

	user.PasswordHash = "" // Don't send password hash in response

	// Create response with primary role
	primaryRole := "student" // default
	if len(roles) > 0 {
		primaryRole = roles[0] // Use first role as primary
	}

	// Create user response with role field
	userResponse := map[string]interface{}{
		"user_id":    user.UserID,
		"id":         user.UserID,
		"username":   user.Username,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"last_login": user.LastLogin,
		"role":       primaryRole,
		"roles":      roles,
		"user_roles": user.UserRoles,
	}

	return c.JSON(map[string]interface{}{
		"token":      tokenString,
		"user":       userResponse,
		"expires_at": expiresAt,
		"success":    true,
	})
}

// Register creates a new user account
// @Summary User registration
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Registration data"
// @Success 201 {object} object{message=string,user=object}
// @Failure 400 {object} object{error=string}
// @Failure 409 {object} object{error=string}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user already exists
	if existingUser, err := h.userRepo.GetByEmail(req.Email); err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	if existingUser, err := h.userRepo.GetByUsername(req.Username); err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := &models.User{
		UserID:       uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        &req.Phone,
		IsActive:     true,
	}

	if err := h.userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Assign default student role
	studentRole, err := h.userRepo.GetRoleByName("student")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find student role",
		})
	}

	if err := h.userRepo.AssignRole(user.UserID, studentRole.RoleID, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign user role",
		})
	}

	// TODO: Create student record if StudentID is provided
	// This would require a StudentRepository

	user.PasswordHash = "" // Don't send password hash in response

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

// GetProfile retrieves current user's profile
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags User Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /user/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := h.userRepo.GetUserWithRoles(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	user.PasswordHash = "" // Don't send password hash in response

	return c.JSON(user)
}

// UpdateProfile updates current user's profile
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body object{first_name=string,last_name=string,phone=string} true "Profile data"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Router /user/profile [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = &req.Phone

	if err := h.userRepo.Update(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	user.PasswordHash = "" // Don't send password hash in response

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// RefreshToken refreshes the JWT token for an authenticated user
// @Summary Refresh JWT token
// @Description Refresh JWT token using the current valid token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} LoginResponse
// @Failure 401 {object} object{error=string}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Extract token from Bearer prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTSecret), nil
	})

	// If token is invalid but not expired, still reject
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Extract claims from token
	claims, ok := token.Claims.(*middleware.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Get user with roles from database using claims
	user, err := h.userRepo.GetUserWithRoles(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Create new token claims
	newClaims := middleware.Claims{
		UserID:   user.UserID,
		Email:    user.Email,
		Username: user.Username,
		Roles:    extractRoles(user),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.UserID.String(),
		},
	}

	// Create and sign token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		fmt.Println(newTokenString)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Calculate token expiry time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	return c.JSON(LoginResponse{
		Token:     tokenString,
		User:      user,
		ExpiresAt: expiresAt,
	})
}

// extractRoles extracts role names from user roles
func extractRoles(user *models.User) []string {
	roles := make([]string, 0)
	for _, ur := range user.UserRoles {
		if ur.Role != nil && ur.IsActive {
			roles = append(roles, ur.Role.RoleName)
		}
	}
	return roles
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash new password",
		})
	}

	// Update password
	if err := h.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}
