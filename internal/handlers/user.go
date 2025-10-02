package handlers

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type UserHandler struct {
	cfg *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{cfg: cfg}
}

// GetUsers retrieves all users with pagination and filters
// @Summary Get users
// @Description Get paginated list of users with search and filters (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role"
// @Param search query string false "Search by name, email, or username"
// @Param is_active query string false "Filter by active status (true/false)"
// @Success 200 {object} object{data=[]object,pagination=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	role := c.Query("role")
	search := c.Query("search")
	isActive := c.Query("is_active")
	
	offset := (page - 1) * limit

	query := `SELECT u.user_id, u.username, u.email, u.first_name, u.last_name,
		u.phone, u.is_active, u.sso_provider, u.created_at, u.last_login,
		COALESCE(STRING_AGG(r.role_name, ', '), '') as roles
		FROM users u
		LEFT JOIN user_roles ur ON u.user_id = ur.user_id AND ur.is_active = true
		LEFT JOIN roles r ON ur.role_id = r.role_id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if search != "" {
		argCount++
		query += " AND (u.first_name ILIKE $" + strconv.Itoa(argCount) + 
			" OR u.last_name ILIKE $" + strconv.Itoa(argCount) + 
			" OR u.email ILIKE $" + strconv.Itoa(argCount) + 
			" OR u.username ILIKE $" + strconv.Itoa(argCount) + ")"
		args = append(args, "%"+search+"%")
	}

	if isActive != "" {
		argCount++
		query += " AND u.is_active = $" + strconv.Itoa(argCount)
		args = append(args, isActive == "true")
	}

	query += " GROUP BY u.user_id, u.username, u.email, u.first_name, u.last_name, u.phone, u.is_active, u.sso_provider, u.created_at, u.last_login"

	if role != "" {
		query += " HAVING STRING_AGG(r.role_name, ', ') ILIKE '%" + role + "%'"
	}

	query += " ORDER BY u.created_at DESC"
	
	argCount++
	query += " LIMIT $" + strconv.Itoa(argCount)
	args = append(args, limit)
	
	argCount++
	query += " OFFSET $" + strconv.Itoa(argCount)
	args = append(args, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user map[string]interface{} = make(map[string]interface{})
		var userID, username, email, firstName, lastName string
		var phone, ssoProvider, roles *string
		var isActive bool
		var createdAt, lastLogin *string

		err := rows.Scan(
			&userID, &username, &email, &firstName, &lastName,
			&phone, &isActive, &ssoProvider, &createdAt, &lastLogin, &roles,
		)
		if err != nil {
			continue
		}

		user["user_id"] = userID
		user["username"] = username
		user["email"] = email
		user["first_name"] = firstName
		user["last_name"] = lastName
		user["phone"] = phone
		user["is_active"] = isActive
		user["sso_provider"] = ssoProvider
		user["created_at"] = createdAt
		user["last_login"] = lastLogin
		user["roles"] = roles

		users = append(users, user)
	}

	// Get total count
	countQuery := `SELECT COUNT(DISTINCT u.user_id) FROM users u
		LEFT JOIN user_roles ur ON u.user_id = ur.user_id AND ur.is_active = true
		LEFT JOIN roles r ON ur.role_id = r.role_id WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 0

	if search != "" {
		countArgCount++
		countQuery += " AND (u.first_name ILIKE $" + strconv.Itoa(countArgCount) + 
			" OR u.last_name ILIKE $" + strconv.Itoa(countArgCount) + 
			" OR u.email ILIKE $" + strconv.Itoa(countArgCount) + 
			" OR u.username ILIKE $" + strconv.Itoa(countArgCount) + ")"
		countArgs = append(countArgs, "%"+search+"%")
	}

	if isActive != "" {
		countArgCount++
		countQuery += " AND u.is_active = $" + strconv.Itoa(countArgCount)
		countArgs = append(countArgs, isActive == "true")
	}

	var total int
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	return c.JSON(fiber.Map{
		"data": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUser retrieves a specific user by ID
// @Summary Get user by ID
// @Description Get detailed information about a specific user (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	query := `SELECT u.user_id, u.username, u.email, u.first_name, u.last_name,
		u.phone, u.is_active, u.sso_provider, u.sso_user_id, u.created_at, u.last_login
		FROM users u WHERE u.user_id = $1`

	var user models.User
	err := database.DB.QueryRow(query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.IsActive, &user.SSOProvider, &user.SSOUserID,
		&user.CreatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Get user roles
	roleQuery := `SELECT r.role_id, r.role_name, r.role_description, ur.assigned_at
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.role_id
		WHERE ur.user_id = $1 AND ur.is_active = true`

	roleRows, err := database.DB.Query(roleQuery, userID)
	if err == nil {
		defer roleRows.Close()
		var roles []map[string]interface{}
		for roleRows.Next() {
			var roleID int
			var roleName, roleDescription string
			var assignedAt string
			
			if err := roleRows.Scan(&roleID, &roleName, &roleDescription, &assignedAt); err == nil {
				roles = append(roles, map[string]interface{}{
					"role_id": roleID,
					"role_name": roleName,
					"role_description": roleDescription,
					"assigned_at": assignedAt,
				})
			}
		}
		user.Roles = roles
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

// CreateUser creates a new user (admin only)
// @Summary Create user
// @Description Create a new user account (Admin only)
// @Tags User Administration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.User true "User data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.FirstName == "" || user.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, first name, and last name are required",
		})
	}

	// Hash password if provided
	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Insert user
	query := `INSERT INTO users (username, email, password_hash, first_name, last_name, phone, sso_provider)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`

	err := database.DB.QueryRow(query,
		user.Username, user.Email, user.PasswordHash, user.FirstName,
		user.LastName, user.Phone, user.SSOProvider,
	).Scan(&user.UserID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser updates an existing user
// @Summary Update user
// @Description Update user information (Admin only)
// @Tags User Administration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body object{username=string,email=string,first_name=string,last_name=string,phone=string,is_active=bool} true "User update data"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	
	var updateData struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 0

	if updateData.Username != "" {
		argCount++
		setParts = append(setParts, "username = $"+strconv.Itoa(argCount))
		args = append(args, updateData.Username)
	}
	if updateData.Email != "" {
		argCount++
		setParts = append(setParts, "email = $"+strconv.Itoa(argCount))
		args = append(args, updateData.Email)
	}
	if updateData.FirstName != "" {
		argCount++
		setParts = append(setParts, "first_name = $"+strconv.Itoa(argCount))
		args = append(args, updateData.FirstName)
	}
	if updateData.LastName != "" {
		argCount++
		setParts = append(setParts, "last_name = $"+strconv.Itoa(argCount))
		args = append(args, updateData.LastName)
	}
	if updateData.Phone != "" {
		argCount++
		setParts = append(setParts, "phone = $"+strconv.Itoa(argCount))
		args = append(args, updateData.Phone)
	}
	if updateData.IsActive != nil {
		argCount++
		setParts = append(setParts, "is_active = $"+strconv.Itoa(argCount))
		args = append(args, *updateData.IsActive)
	}

	if len(setParts) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	argCount++
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	
	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Description Assign a role to a specific user (Admin only)
// @Tags User Administration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param role body object{role_id=int} true "Role assignment data"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id}/roles [post]
func (h *UserHandler) AssignRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	assignerID := c.Locals("user_id").(string)

	var request struct {
		RoleID int `json:"role_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	database.DB.QueryRow(checkQuery, userID).Scan(&exists)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if role exists
	checkRoleQuery := `SELECT EXISTS(SELECT 1 FROM roles WHERE role_id = $1)`
	database.DB.QueryRow(checkRoleQuery, request.RoleID).Scan(&exists)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	// Insert or update user role
	query := `INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user_id, role_id) 
		DO UPDATE SET is_active = true, assigned_at = CURRENT_TIMESTAMP, assigned_by = $3`

	_, err := database.DB.Exec(query, userID, request.RoleID, assignerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role assigned successfully",
	})
}

// RemoveRole removes a role from a user
// @Summary Remove role from user
// @Description Remove a specific role from a user (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id}/roles/{role_id} [delete]
func (h *UserHandler) RemoveRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("role_id")

	query := `UPDATE user_roles SET is_active = false 
		WHERE user_id = $1 AND role_id = $2`

	result, err := database.DB.Exec(query, userID, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove role",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User role assignment not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role removed successfully",
	})
}

// GetRoles retrieves all available roles
// @Summary Get all roles
// @Description Get list of all available roles in the system (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /users/roles [get]
func (h *UserHandler) GetRoles(c *fiber.Ctx) error {
	query := `SELECT role_id, role_name, role_description, permissions, created_at
		FROM roles ORDER BY role_name`

	rows, err := database.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.RoleID, &role.RoleName, &role.RoleDescription,
			&role.Permissions, &role.CreatedAt,
		)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return c.JSON(fiber.Map{
		"data": roles,
	})
}

// DeactivateUser deactivates a user account
// @Summary Deactivate user
// @Description Deactivate a user account (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1`

	result, err := database.DB.Exec(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deactivate user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}

// ReactivateUser reactivates a user account
// @Summary Reactivate user
// @Description Reactivate a user account (Admin only)
// @Tags User Administration
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} object{message=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /users/{id}/reactivate [post]
func (h *UserHandler) ReactivateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	query := `UPDATE users SET is_active = true, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1`

	result, err := database.DB.Exec(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reactivate user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User reactivated successfully",
	})
}