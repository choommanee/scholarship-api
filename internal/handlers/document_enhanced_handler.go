package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type DocumentEnhancedHandler struct {
	cfg             *config.Config
	applicationRepo *repository.ApplicationRepository
	userRepo        *repository.UserRepository
}

func NewDocumentEnhancedHandler(cfg *config.Config) *DocumentEnhancedHandler {
	return &DocumentEnhancedHandler{
		cfg:             cfg,
		applicationRepo: repository.NewApplicationRepository(),
		userRepo:        repository.NewUserRepository(),
	}
}

// UploadDocumentEnhanced handles file upload for application documents
// @Summary Upload document
// @Description Upload a document for an application
// @Tags Document Management
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param document_type formData string true "Document type (id_card, transcript, income_certificate, house_registration, etc.)"
// @Param file formData file true "Document file (PDF, JPEG, PNG - max 10MB)"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/documents/applications/{id}/upload [post]
func (h *DocumentEnhancedHandler) UploadDocumentEnhanced(c *fiber.Ctx) error {
	// Get user_id from context
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	documentType := c.FormValue("document_type")
	if documentType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Document type is required",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only upload documents for your own applications",
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file size based on document type
	maxSize := h.getMaxFileSize(documentType)
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("File size exceeds %dMB limit", maxSize/(1024*1024)),
		})
	}

	// Validate file type
	mimeType := file.Header.Get("Content-Type")
	if !h.isAllowedMimeType(mimeType) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: PDF, JPEG, PNG",
		})
	}

	// Create upload directory if not exists
	uploadDir := filepath.Join("./uploads", "applications", fmt.Sprintf("%d", applicationID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", documentType, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Save document metadata to database
	doc := &models.ApplicationDocument{
		ApplicationID: applicationID,
		DocumentType:  documentType,
		DocumentName:  file.Filename,
		FilePath:      filePath,
		FileSize:      file.Size,
		MimeType:      mimeType,
		UploadStatus:  "pending",
		UploadedAt:    time.Now(),
	}

	if err := h.applicationRepo.AddDocument(doc); err != nil {
		// Delete uploaded file if database insert fails
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save document metadata",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Document uploaded successfully",
		"data": fiber.Map{
			"document_id":         doc.DocumentID,
			"document_type":       doc.DocumentType,
			"file_name":          doc.DocumentName,
			"file_path":          filePath,
			"file_size":          doc.FileSize,
			"mime_type":          doc.MimeType,
			"verification_status": doc.UploadStatus,
			"uploaded_at":        doc.UploadedAt,
		},
	})
}

// DeleteDocument deletes an uploaded document
// @Summary Delete document
// @Description Delete an uploaded document
// @Tags Document Management
// @Produce json
// @Security BearerAuth
// @Param id path int true "Document ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/documents/{id} [delete]
func (h *DocumentEnhancedHandler) DeleteDocument(c *fiber.Ctx) error {
	// Get user_id from context
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	documentID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	// Get document from database
	var doc models.ApplicationDocument
	var filePath string
	query := `
		SELECT d.document_id, d.application_id, d.document_type, d.document_name,
		       d.file_path, d.file_size, d.mime_type, d.upload_status
		FROM application_documents d
		WHERE d.document_id = $1
	`

	err = database.DB.QueryRow(query, documentID).Scan(
		&doc.DocumentID,
		&doc.ApplicationID,
		&doc.DocumentType,
		&doc.DocumentName,
		&filePath,
		&doc.FileSize,
		&doc.MimeType,
		&doc.UploadStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Document not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch document",
		})
	}

	// Verify document belongs to user's application
	if err := h.verifyApplicationOwnership(uint(doc.ApplicationID), userID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own documents",
		})
	}

	// Verify application is still in draft status
	application, err := h.applicationRepo.GetByID(uint(doc.ApplicationID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify application status",
		})
	}

	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete documents for applications that are not in draft status",
		})
	}

	// Delete file from disk
	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			// Log error but continue with database deletion
			fmt.Printf("Warning: Failed to delete file %s: %v\n", filePath, err)
		}
	}

	// Delete database record
	deleteQuery := `DELETE FROM application_documents WHERE document_id = $1`
	_, err = database.DB.Exec(deleteQuery, documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete document record",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Document deleted successfully",
	})
}

// DownloadDocument serves a document for download
// @Summary Download document
// @Description Download a document file
// @Tags Document Management
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path int true "Document ID"
// @Success 200 {file} binary
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/documents/{id}/download [get]
func (h *DocumentEnhancedHandler) DownloadDocument(c *fiber.Ctx) error {
	// Get user_id from context
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	documentID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	// Get document from database
	var doc models.ApplicationDocument
	var filePath string
	query := `
		SELECT d.document_id, d.application_id, d.document_type, d.document_name,
		       d.file_path, d.file_size, d.mime_type
		FROM application_documents d
		WHERE d.document_id = $1
	`

	err = database.DB.QueryRow(query, documentID).Scan(
		&doc.DocumentID,
		&doc.ApplicationID,
		&doc.DocumentType,
		&doc.DocumentName,
		&filePath,
		&doc.FileSize,
		&doc.MimeType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Document not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch document",
		})
	}

	// Verify user has permission (own application or is officer)
	roles := c.Locals("roles")
	isOfficer := false
	if roleList, ok := roles.([]string); ok {
		for _, role := range roleList {
			if role == "admin" || role == "scholarship_officer" {
				isOfficer = true
				break
			}
		}
	}

	if !isOfficer {
		// Student can only download their own documents
		if err := h.verifyApplicationOwnership(uint(doc.ApplicationID), userID); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only download your own documents",
			})
		}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document file not found on server",
		})
	}

	// Set headers for download
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.DocumentName))
	c.Set("Content-Type", doc.MimeType)

	// Send file
	return c.SendFile(filePath)
}

// Helper functions

func (h *DocumentEnhancedHandler) getMaxFileSize(documentType string) int64 {
	// Default 10MB
	maxSize := int64(10 * 1024 * 1024)

	// Specific limits for different document types
	switch documentType {
	case "id_card", "transcript", "income_certificate":
		maxSize = int64(5 * 1024 * 1024) // 5MB
	case "house_photos", "living_situation_photos":
		maxSize = int64(20 * 1024 * 1024) // 20MB
	}

	return maxSize
}

func (h *DocumentEnhancedHandler) isAllowedMimeType(mimeType string) bool {
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
	}

	return allowedTypes[strings.ToLower(mimeType)]
}

func (h *DocumentEnhancedHandler) verifyApplicationOwnership(applicationID uint, userID uuid.UUID) error {
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		return err
	}

	// Get student ID from students table by user_id
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no student record exists, try matching with email
			user, userErr := h.userRepo.GetByID(userID)
			if userErr != nil {
				return userErr
			}

			if application.StudentID != user.Email {
				return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
			}
			return nil
		}
		return err
	}

	if application.StudentID != studentID {
		return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
	}

	return nil
}
