package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"scholarship-system/internal/models"
)

// FileRepository handles file-related database operations
type FileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// CreateFile creates a new file entry
func (r *FileRepository) CreateFile(file *models.FileStorage) error {
	query := `
		INSERT INTO file_storage (
			file_id, original_name, stored_name, stored_path, file_size,
			mime_type, file_hash, uploaded_by, related_table, related_id,
			storage_type, access_level
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(
		query,
		file.FileID, file.OriginalName, file.StoredName, file.StoredPath, file.FileSize,
		file.MimeType, file.FileHash, file.UploadedBy, file.RelatedTable, file.RelatedID,
		file.StorageType, file.AccessLevel,
	).Scan(&file.CreatedAt, &file.UpdatedAt)
}

// GetFileByID retrieves a file by ID
func (r *FileRepository) GetFileByID(id uuid.UUID) (*models.FileStorage, error) {
	query := `
		SELECT file_id, original_name, stored_name, stored_path, file_size,
			   mime_type, file_hash, uploaded_by, related_table, related_id,
			   storage_type, access_level, created_at, updated_at
		FROM file_storage
		WHERE file_id = $1`

	file := &models.FileStorage{}
	err := r.db.QueryRow(query, id).Scan(
		&file.FileID, &file.OriginalName, &file.StoredName, &file.StoredPath, &file.FileSize,
		&file.MimeType, &file.FileHash, &file.UploadedBy, &file.RelatedTable, &file.RelatedID,
		&file.StorageType, &file.AccessLevel, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetFilesByRelated retrieves files by related entity
func (r *FileRepository) GetFilesByRelated(relatedTable string, relatedID uuid.UUID) ([]models.FileStorage, error) {
	query := `
		SELECT file_id, original_name, stored_name, stored_path, file_size,
			   mime_type, file_hash, uploaded_by, related_table, related_id,
			   storage_type, access_level, created_at, updated_at
		FROM file_storage
		WHERE related_table = $1 AND related_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, relatedTable, relatedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.FileStorage
	for rows.Next() {
		var f models.FileStorage
		err := rows.Scan(
			&f.FileID, &f.OriginalName, &f.StoredName, &f.StoredPath, &f.FileSize,
			&f.MimeType, &f.FileHash, &f.UploadedBy, &f.RelatedTable, &f.RelatedID,
			&f.StorageType, &f.AccessLevel, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

// DeleteFile deletes a file
func (r *FileRepository) DeleteFile(id uuid.UUID) error {
	query := `DELETE FROM file_storage WHERE file_id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CreateFileVersion creates a file version
func (r *FileRepository) CreateFileVersion(version *models.FileVersion) error {
	query := `
		INSERT INTO document_versions (
			version_id, file_id, version_number, change_description,
			uploaded_by, file_size, is_current
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		version.VersionID, version.FileID, version.VersionNumber, version.ChangeDescription,
		version.UploadedBy, version.FileSize, version.IsCurrent,
	).Scan(&version.CreatedAt)
}

// GetFileVersions retrieves all versions of a file
func (r *FileRepository) GetFileVersions(fileID uuid.UUID) ([]models.FileVersion, error) {
	query := `
		SELECT version_id, file_id, version_number, change_description,
			   uploaded_by, file_size, is_current, created_at, replaced_at
		FROM document_versions
		WHERE file_id = $1
		ORDER BY version_number DESC`

	rows, err := r.db.Query(query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.FileVersion
	for rows.Next() {
		var v models.FileVersion
		err := rows.Scan(
			&v.VersionID, &v.FileID, &v.VersionNumber, &v.ChangeDescription,
			&v.UploadedBy, &v.FileSize, &v.IsCurrent, &v.CreatedAt, &v.ReplacedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

// LogFileAccess logs file access
func (r *FileRepository) LogFileAccess(log *models.FileAccessLog) error {
	query := `
		INSERT INTO file_access_logs (
			access_id, file_id, user_id, access_time, action,
			ip_address, user_agent, success
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(
		query,
		log.AccessID, log.FileID, log.UserID, log.AccessTime, log.Action,
		log.IPAddress, log.UserAgent, log.Success,
	)
	return err
}

// GetFileAccessLogs retrieves access logs for a file
func (r *FileRepository) GetFileAccessLogs(fileID uuid.UUID, limit int) ([]models.FileAccessLog, error) {
	query := `
		SELECT access_id, file_id, user_id, access_time, action,
			   ip_address, user_agent, success
		FROM file_access_logs
		WHERE file_id = $1
		ORDER BY access_time DESC
		LIMIT $2`

	rows, err := r.db.Query(query, fileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.FileAccessLog
	for rows.Next() {
		var l models.FileAccessLog
		err := rows.Scan(
			&l.AccessID, &l.FileID, &l.UserID, &l.AccessTime, &l.Action,
			&l.IPAddress, &l.UserAgent, &l.Success,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
