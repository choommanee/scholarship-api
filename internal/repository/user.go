package repository

import (
	"database/sql"
	"fmt"
	"time"

	"scholarship-system/internal/database"
	"scholarship-system/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, phone, is_active, sso_provider, sso_user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(query,
		user.UserID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.SSOProvider,
		user.SSOUserID,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, first_name, last_name, phone, is_active, 
		       sso_provider, sso_user_id, created_at, updated_at, last_login
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.SSOProvider,
		&user.SSOUserID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, first_name, last_name, phone, is_active, 
		       sso_provider, sso_user_id, created_at, updated_at, last_login
		FROM users 
		WHERE user_id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.SSOProvider,
		&user.SSOUserID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, first_name, last_name, phone, is_active, 
		       sso_provider, sso_user_id, created_at, updated_at, last_login
		FROM users 
		WHERE username = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.SSOProvider,
		&user.SSOUserID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, first_name = $4, last_name = $5, phone = $6, 
		    is_active = $7, updated_at = $8
		WHERE user_id = $1
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		user.UserID,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) UpdatePassword(userID uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE user_id = $3`
	_, err := r.db.Exec(query, passwordHash, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateLastLogin(userID uuid.UUID) error {
	query := `UPDATE users SET last_login = $1 WHERE user_id = $2`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *UserRepository) Delete(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) GetUserRoles(userID uuid.UUID) ([]models.UserRole, error) {
	query := `
		SELECT ur.user_id, ur.role_id, ur.assigned_at, ur.assigned_by, ur.is_active,
		       r.role_id, r.role_name, r.role_description, r.permissions, r.created_at
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.role_id
		WHERE ur.user_id = $1 AND ur.is_active = true
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []models.UserRole

	for rows.Next() {
		var userRole models.UserRole
		var role models.Role

		err := rows.Scan(
			&userRole.UserID,
			&userRole.RoleID,
			&userRole.AssignedAt,
			&userRole.AssignedBy,
			&userRole.IsActive,
			&role.RoleID,
			&role.RoleName,
			&role.RoleDescription,
			&role.Permissions,
			&role.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		userRole.Role = &role
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (r *UserRepository) GetUserWithRoles(userID uuid.UUID) (*models.User, error) {
	user, err := r.GetByID(userID)
	if err != nil {
		return nil, err
	}

	userRoles, err := r.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoles
	return user, nil
}

func (r *UserRepository) GetUserWithRolesByEmail(email string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	userRoles, err := r.GetUserRoles(user.UserID)
	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoles
	return user, nil
}

func (r *UserRepository) GetUserWithRolesByUsername(username string) (*models.User, error) {
	user, err := r.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	userRoles, err := r.GetUserRoles(user.UserID)
	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoles
	return user, nil
}

func (r *UserRepository) AssignRole(userID uuid.UUID, roleID uint, assignedBy *uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, role_id) 
		DO UPDATE SET is_active = $5, assigned_at = $3, assigned_by = $4
	`

	_, err := r.db.Exec(query, userID, roleID, time.Now(), assignedBy, true)
	return err
}

func (r *UserRepository) RemoveRole(userID uuid.UUID, roleID uint) error {
	query := `UPDATE user_roles SET is_active = false WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(query, userID, roleID)
	return err
}

func (r *UserRepository) GetRoleByName(roleName string) (*models.Role, error) {
	query := `SELECT role_id, role_name, role_description, permissions, created_at FROM roles WHERE role_name = $1`

	role := &models.Role{}
	err := r.db.QueryRow(query, roleName).Scan(
		&role.RoleID,
		&role.RoleName,
		&role.RoleDescription,
		&role.Permissions,
		&role.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *UserRepository) List(limit, offset int, search string) ([]models.User, int, error) {
	var users []models.User
	var totalCount int

	// Count total records
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		countQuery += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR username ILIKE $%d)", argIndex, argIndex+1, argIndex+2, argIndex+3)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT user_id, username, email, first_name, last_name, phone, is_active, 
		       sso_provider, sso_user_id, created_at, updated_at, last_login
		FROM users 
		WHERE 1=1
	`

	if search != "" {
		query += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR username ILIKE $%d)", argIndex-3, argIndex-2, argIndex-1, argIndex)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.IsActive,
			&user.SSOProvider,
			&user.SSOUserID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		)

		if err != nil {
			return nil, 0, err
		}

		// Don't include password hash in list response
		user.PasswordHash = ""
		users = append(users, user)
	}

	return users, totalCount, nil
}
