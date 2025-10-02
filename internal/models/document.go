package models

import "time"

type ApplicationDocument struct {
	DocumentID        int        `json:"document_id" db:"document_id"`
	ApplicationID     int        `json:"application_id" db:"application_id"`
	DocumentType      string     `json:"document_type" db:"document_type"`
	DocumentName      string     `json:"document_name" db:"document_name"`
	FilePath          string     `json:"file_path,omitempty" db:"file_path"`
	FileSize          int64      `json:"file_size" db:"file_size"`
	MimeType          string     `json:"mime_type" db:"mime_type"`
	UploadStatus      string     `json:"upload_status" db:"upload_status"`
	VerificationNotes *string    `json:"verification_notes,omitempty" db:"verification_notes"`
	UploadedAt        time.Time  `json:"uploaded_at" db:"uploaded_at"`
	VerifiedBy        *string    `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty" db:"verified_at"`
}