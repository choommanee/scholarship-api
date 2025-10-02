package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type DocumentHandler struct {
	cfg *config.Config
}

func NewDocumentHandler(cfg *config.Config) *DocumentHandler {
	return &DocumentHandler{cfg: cfg}
}

// UploadDocument handles file upload for application documents
// @Summary Upload document
// @Description Upload supporting document for scholarship application (Student only)
// @Tags Document Management
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param application_id path string true "Application ID"
// @Param file formData file true "Document file (PDF, JPEG, PNG, DOC, DOCX - max 10MB)"
// @Param document_type formData string true "Document type (id_card, transcript, income_certificate, etc.)"
// @Success 201 {object} object{message=string,document_id=int,filename=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /documents/applications/{application_id}/upload [post]
func (h *DocumentHandler) UploadDocument(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")
	documentType := c.FormValue("document_type")
	
	if documentType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Document type is required",
		})
	}

	// Check if application exists and belongs to current user
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := userIDValue.(uuid.UUID).String()
	var studentID string
	checkQuery := `SELECT s.student_id FROM scholarship_applications sa
		JOIN students s ON sa.student_id = s.student_id
		JOIN users u ON s.user_id = u.user_id
		WHERE sa.application_id = $1 AND u.user_id = $2`
	
	err := database.DB.QueryRow(checkQuery, applicationID, userID).Scan(&studentID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Application not found or access denied",
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file size (10MB limit)
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size exceeds 10MB limit",
		})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf":  true,
		"image/jpeg":       true,
		"image/jpg":        true,
		"image/png":        true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: PDF, JPEG, PNG, DOC, DOCX",
		})
	}

	// Create upload directory if not exists
	uploadDir := filepath.Join("./uploads", "documents", applicationID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", documentType, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveFile(file, filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Save document metadata to database
	query := `INSERT INTO application_documents 
		(application_id, document_type, document_name, file_path, file_size, mime_type)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING document_id`

	var documentID int
	err = database.DB.QueryRow(query,
		applicationID, documentType, file.Filename, filepath,
		file.Size, file.Header.Get("Content-Type"),
	).Scan(&documentID)

	if err != nil {
		// Delete uploaded file if database insert fails
		os.Remove(filepath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save document metadata",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Document uploaded successfully",
		"document_id": documentID,
		"filename": filename,
	})
}

// GetDocuments retrieves documents for an application
// @Summary Get application documents
// @Description Get list of documents for a specific application
// @Tags Document Management
// @Produce json
// @Security BearerAuth
// @Param application_id path string true "Application ID"
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /documents/applications/{application_id} [get]
func (h *DocumentHandler) GetDocuments(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("user_role").(string)

	// Check access permissions
	if userRole == "student" {
		var studentID string
		checkQuery := `SELECT s.student_id FROM scholarship_applications sa
			JOIN students s ON sa.student_id = s.student_id
			JOIN users u ON s.user_id = u.user_id
			WHERE sa.application_id = $1 AND u.user_id = $2`
		
		err := database.DB.QueryRow(checkQuery, applicationID, userID).Scan(&studentID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
	}

	query := `SELECT document_id, document_type, document_name, file_path,
		file_size, mime_type, upload_status, verification_notes,
		uploaded_at, verified_at
		FROM application_documents
		WHERE application_id = $1
		ORDER BY uploaded_at DESC`

	rows, err := database.DB.Query(query, applicationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch documents",
		})
	}
	defer rows.Close()

	var documents []models.ApplicationDocument
	for rows.Next() {
		var doc models.ApplicationDocument
		err := rows.Scan(
			&doc.DocumentID, &doc.DocumentType, &doc.DocumentName,
			&doc.FilePath, &doc.FileSize, &doc.MimeType,
			&doc.UploadStatus, &doc.VerificationNotes,
			&doc.UploadedAt, &doc.VerifiedAt,
		)
		if err != nil {
			continue
		}
		
		// Don't expose full file path to client
		doc.FilePath = ""
		documents = append(documents, doc)
	}

	return c.JSON(fiber.Map{
		"data": documents,
	})
}

// DownloadDocument handles document download
// @Summary Download document
// @Description Download a specific document file
// @Tags Document Management
// @Produce application/octet-stream
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {file} file "Document file"
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /documents/{document_id}/download [get]
func (h *DocumentHandler) DownloadDocument(c *fiber.Ctx) error {
	documentID := c.Params("document_id")
	userID := c.Locals("user_id").(string)
	userRole := c.Locals("user_role").(string)

	// Get document info and check permissions
	var doc models.ApplicationDocument
	var applicationID int
	query := `SELECT ad.document_id, ad.application_id, ad.document_type,
		ad.document_name, ad.file_path, ad.mime_type
		FROM application_documents ad`

	if userRole == "student" {
		query += ` JOIN scholarship_applications sa ON ad.application_id = sa.application_id
			JOIN students s ON sa.student_id = s.student_id
			JOIN users u ON s.user_id = u.user_id
			WHERE ad.document_id = $1 AND u.user_id = $2`
		err := database.DB.QueryRow(query, documentID, userID).Scan(
			&doc.DocumentID, &applicationID, &doc.DocumentType,
			&doc.DocumentName, &doc.FilePath, &doc.MimeType,
		)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Document not found or access denied",
			})
		}
	} else {
		query += " WHERE ad.document_id = $1"
		err := database.DB.QueryRow(query, documentID).Scan(
			&doc.DocumentID, &applicationID, &doc.DocumentType,
			&doc.DocumentName, &doc.FilePath, &doc.MimeType,
		)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Document not found",
			})
		}
	}

	// Check if file exists
	if _, err := os.Stat(doc.FilePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found on server",
		})
	}

	// Set headers for file download
	c.Set("Content-Type", doc.MimeType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.DocumentName))

	return c.SendFile(doc.FilePath)
}

// VerifyDocument allows officers to verify uploaded documents
// @Summary Verify document
// @Description Verify or reject an uploaded document (Admin/Officer only)
// @Tags Document Verification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Param verification body object{status=string,notes=string} true "Verification data (status: verified/rejected)"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /documents/{document_id}/verify [post]
func (h *DocumentHandler) VerifyDocument(c *fiber.Ctx) error {
	documentID := c.Params("document_id")
	userID := c.Locals("user_id").(string)

	var verification struct {
		Status string `json:"status"` // "verified" or "rejected"
		Notes  string `json:"notes"`
	}

	if err := c.BodyParser(&verification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if verification.Status != "verified" && verification.Status != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status must be 'verified' or 'rejected'",
		})
	}

	query := `UPDATE application_documents 
		SET upload_status = $1, verification_notes = $2, 
		    verified_by = $3, verified_at = CURRENT_TIMESTAMP
		WHERE document_id = $4`

	result, err := database.DB.Exec(query, verification.Status, verification.Notes, userID, documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify document",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Document verification updated successfully",
	})
}

// DeleteDocument allows students to delete their own documents (if not verified)
// @Summary Delete document
// @Description Delete own uploaded document (only if not verified, Student only)
// @Tags Document Management
// @Produce json
// @Security BearerAuth
// @Param document_id path string true "Document ID"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /documents/{document_id} [delete]
func (h *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	documentID := c.Params("document_id")
	userID := c.Locals("user_id").(string)

	// Check if document belongs to current user and is not verified
	var filePath string
	var uploadStatus string
	checkQuery := `SELECT ad.file_path, ad.upload_status
		FROM application_documents ad
		JOIN scholarship_applications sa ON ad.application_id = sa.application_id
		JOIN students s ON sa.student_id = s.student_id
		JOIN users u ON s.user_id = u.user_id
		WHERE ad.document_id = $1 AND u.user_id = $2`

	err := database.DB.QueryRow(checkQuery, documentID, userID).Scan(&filePath, &uploadStatus)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document not found or access denied",
		})
	}

	if uploadStatus == "verified" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete verified document",
		})
	}

	// Delete from database
	deleteQuery := `DELETE FROM application_documents WHERE document_id = $1`
	_, err = database.DB.Exec(deleteQuery, documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete document record",
		})
	}

	// Delete physical file
	if err := os.Remove(filePath); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to delete file %s: %v\n", filePath, err)
	}

	return c.JSON(fiber.Map{
		"message": "Document deleted successfully",
	})
}

// GetDocumentTypes returns available document types
// @Summary Get document types
// @Description Get list of available document types for uploads
// @Tags Document Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]object}
// @Failure 401 {object} object{error=string}
// @Router /documents/types [get]
func (h *DocumentHandler) GetDocumentTypes(c *fiber.Ctx) error {
    types := []map[string]interface{}{
        {"type": "id_card", "description": "Identity Card", "required": true},
        {"type": "transcript", "description": "Academic Transcript", "required": true},
        {"type": "income_certificate", "description": "Family Income Certificate", "required": true},
        {"type": "residence_photo", "description": "Residence Photo", "required": true},
        {"type": "activity_certificate", "description": "Activity Certificate", "required": false},
        {"type": "recommendation_letter", "description": "Recommendation Letter", "required": false},
        {"type": "medical_certificate", "description": "Medical Certificate", "required": false},
        {"type": "scholarship_certificate", "description": "Previous Scholarship Certificate", "required": false},
        {"type": "bank_account", "description": "Bank Account Statement", "required": false},
        {"type": "other", "description": "Other Supporting Documents", "required": false},
    }

    return c.JSON(fiber.Map{
        "data": types,
    })
}

// GetDocumentStats returns document verification statistics
// @Summary Get document statistics
// @Description Get document verification statistics, optionally filtered by application
// @Tags Document Management
// @Produce json
// @Security BearerAuth
// @Param application_id query string false "Filter by application ID"
// @Success 200 {object} object{data=object}
// @Failure 401 {object} object{error=string}
// @Router /documents/stats [get]
func (h *DocumentHandler) GetDocumentStats(c *fiber.Ctx) error {
	applicationID := c.Query("application_id")
	
	baseQuery := `SELECT 
		COUNT(*) as total_documents,
		COUNT(CASE WHEN upload_status = 'pending' THEN 1 END) as pending_docs,
		COUNT(CASE WHEN upload_status = 'verified' THEN 1 END) as verified_docs,
		COUNT(CASE WHEN upload_status = 'rejected' THEN 1 END) as rejected_docs
		FROM application_documents`
	
	args := []interface{}{}
	if applicationID != "" {
		baseQuery += " WHERE application_id = $1"
		args = append(args, applicationID)
	}

	var stats struct {
		TotalDocuments    int `json:"total_documents"`
		PendingDocuments  int `json:"pending_documents"`
		VerifiedDocuments int `json:"verified_documents"`
		RejectedDocuments int `json:"rejected_documents"`
	}

	err := database.DB.QueryRow(baseQuery, args...).Scan(
		&stats.TotalDocuments, &stats.PendingDocuments,
		&stats.VerifiedDocuments, &stats.RejectedDocuments,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch document statistics",
		})
	}

	return c.JSON(fiber.Map{
		"data": stats,
	})
}

// BulkVerifyDocuments allows officers to verify multiple documents at once
// @Summary Bulk verify documents
// @Description Verify or reject multiple documents at once (Admin/Officer only)
// @Tags Document Verification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param verification body object{document_ids=[]int,status=string,notes=string} true "Bulk verification data"
// @Success 200 {object} object{message=string,updated_count=int}
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Router /documents/bulk-verify [post]
func (h *DocumentHandler) BulkVerifyDocuments(c *fiber.Ctx) error {
	var request struct {
		DocumentIDs []int  `json:"document_ids"`
		Status      string `json:"status"`
		Notes       string `json:"notes"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(request.DocumentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Document IDs are required",
		})
	}

	if request.Status != "verified" && request.Status != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status must be 'verified' or 'rejected'",
		})
	}

	userID := c.Locals("user_id").(string)
	
	// Build query with placeholders
	placeholders := make([]string, len(request.DocumentIDs))
	args := []interface{}{request.Status, request.Notes, userID}
	
	for i, id := range request.DocumentIDs {
		placeholders[i] = "$" + strconv.Itoa(i+4)
		args = append(args, id)
	}

	query := fmt.Sprintf(`UPDATE application_documents 
		SET upload_status = $1, verification_notes = $2, 
		    verified_by = $3, verified_at = CURRENT_TIMESTAMP
		WHERE document_id IN (%s)`, strings.Join(placeholders, ","))

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify documents",
		})
	}

	rowsAffected, _ := result.RowsAffected()

	return c.JSON(fiber.Map{
		"message": "Documents verified successfully",
		"updated_count": rowsAffected,
	})
}