package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type ScholarshipRepository struct {
	db *sql.DB
}

func NewScholarshipRepository() *ScholarshipRepository {
	return &ScholarshipRepository{
		db: database.DB,
	}
}

// Scholarship Source Methods
func (r *ScholarshipRepository) CreateSource(source *models.ScholarshipSource) error {
	query := `
		INSERT INTO scholarship_sources (source_name, source_type, contact_person, contact_email, contact_phone, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING source_id
	`

	now := time.Now()
	source.CreatedAt = now
	source.UpdatedAt = now

	err := r.db.QueryRow(query,
		source.SourceName,
		source.SourceType,
		source.ContactPerson,
		source.ContactEmail,
		source.ContactPhone,
		source.Description,
		source.IsActive,
		source.CreatedAt,
		source.UpdatedAt,
	).Scan(&source.SourceID)

	return err
}

func (r *ScholarshipRepository) GetSourceByID(sourceID uint) (*models.ScholarshipSource, error) {
	query := `
		SELECT source_id, source_name, source_type, contact_person, contact_email, contact_phone, description, is_active, created_at, updated_at
		FROM scholarship_sources WHERE source_id = $1
	`

	source := &models.ScholarshipSource{}
	err := r.db.QueryRow(query, sourceID).Scan(
		&source.SourceID,
		&source.SourceName,
		&source.SourceType,
		&source.ContactPerson,
		&source.ContactEmail,
		&source.ContactPhone,
		&source.Description,
		&source.IsActive,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	return source, err
}

func (r *ScholarshipRepository) ListSources(limit, offset int, search string) ([]models.ScholarshipSource, int, error) {
	var sources []models.ScholarshipSource
	var totalCount int

	// Count total records
	countQuery := `SELECT COUNT(*) FROM scholarship_sources WHERE is_active = true`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		countQuery += fmt.Sprintf(" AND (source_name ILIKE $%d OR source_type ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get sources
	query := `
		SELECT source_id, source_name, source_type, contact_person, contact_email, contact_phone, description, is_active, created_at, updated_at
		FROM scholarship_sources WHERE is_active = true
	`

	if search != "" {
		query += fmt.Sprintf(" AND (source_name ILIKE $%d OR source_type ILIKE $%d)", argIndex-1, argIndex)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var source models.ScholarshipSource
		err := rows.Scan(
			&source.SourceID,
			&source.SourceName,
			&source.SourceType,
			&source.ContactPerson,
			&source.ContactEmail,
			&source.ContactPhone,
			&source.Description,
			&source.IsActive,
			&source.CreatedAt,
			&source.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		sources = append(sources, source)
	}

	return sources, totalCount, nil
}

// Scholarship Methods
func (r *ScholarshipRepository) Create(scholarship *models.Scholarship) error {
	query := `
		INSERT INTO scholarships (source_id, name, type, amount, total_quota, available_quota,
		                         academic_year, semester, eligibility_criteria, required_documents,
		                         application_start_date, application_end_date, interview_required,
		                         is_active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING scholarship_id
	`

	now := time.Now()
	scholarship.CreatedAt = now
	scholarship.UpdatedAt = now

	// Convert string fields to proper JSON format for JSONB columns
	var eligibilityCriteriaJSON interface{}
	var requiredDocumentsJSON interface{}

	if scholarship.EligibilityCriteria != nil && *scholarship.EligibilityCriteria != "" {
		// Convert string to JSON array format
		criteriaArray := []string{*scholarship.EligibilityCriteria}
		eligibilityCriteriaJSON, _ = json.Marshal(criteriaArray)
	} else {
		eligibilityCriteriaJSON = nil
	}

	if scholarship.RequiredDocuments != nil && *scholarship.RequiredDocuments != "" {
		// Convert comma-separated string to JSON array
		docs := []string{}
		if *scholarship.RequiredDocuments != "" {
			// Split by comma and trim spaces
			parts := strings.Split(*scholarship.RequiredDocuments, ",")
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" {
					docs = append(docs, trimmed)
				}
			}
		}
		requiredDocumentsJSON, _ = json.Marshal(docs)
	} else {
		requiredDocumentsJSON = nil
	}

	err := r.db.QueryRow(query,
		scholarship.SourceID,
		scholarship.ScholarshipName,
		scholarship.ScholarshipType,
		scholarship.Amount,
		scholarship.TotalQuota,
		scholarship.AvailableQuota,
		scholarship.AcademicYear,
		scholarship.Semester,
		eligibilityCriteriaJSON,
		requiredDocumentsJSON,
		scholarship.ApplicationStartDate,
		scholarship.ApplicationEndDate,
		scholarship.InterviewRequired,
		scholarship.IsActive,
		scholarship.CreatedBy,
		scholarship.CreatedAt,
		scholarship.UpdatedAt,
	).Scan(&scholarship.ScholarshipID)

	return err
}

func (r *ScholarshipRepository) GetByID(scholarshipID uint) (*models.Scholarship, error) {
	query := `
		SELECT s.scholarship_id, s.source_id, s.name, s.type, s.amount,
		       s.total_quota, s.available_quota, s.academic_year, s.semester, s.eligibility_criteria,
		       s.required_documents, s.application_start_date, s.application_end_date,
		       s.interview_required, s.is_active, s.created_by, s.created_at, s.updated_at,
		       src.source_id, src.source_name, src.source_type, src.contact_person,
		       src.contact_email, src.contact_phone, src.description, src.is_active,
		       src.created_at, src.updated_at
		FROM scholarships s
		LEFT JOIN scholarship_sources src ON s.source_id = src.source_id
		WHERE s.scholarship_id = $1
	`

	scholarship := &models.Scholarship{}
	source := &models.ScholarshipSource{}

	var eligibilityCriteriaBytes []byte
	var requiredDocumentsBytes []byte

	err := r.db.QueryRow(query, scholarshipID).Scan(
		&scholarship.ScholarshipID,
		&scholarship.SourceID,
		&scholarship.ScholarshipName,
		&scholarship.ScholarshipType,
		&scholarship.Amount,
		&scholarship.TotalQuota,
		&scholarship.AvailableQuota,
		&scholarship.AcademicYear,
		&scholarship.Semester,
		&eligibilityCriteriaBytes,
		&requiredDocumentsBytes,
		&scholarship.ApplicationStartDate,
		&scholarship.ApplicationEndDate,
		&scholarship.InterviewRequired,
		&scholarship.IsActive,
		&scholarship.CreatedBy,
		&scholarship.CreatedAt,
		&scholarship.UpdatedAt,
		&source.SourceID,
		&source.SourceName,
		&source.SourceType,
		&source.ContactPerson,
		&source.ContactEmail,
		&source.ContactPhone,
		&source.Description,
		&source.IsActive,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convert JSONB back to string
	if eligibilityCriteriaBytes != nil {
		var criteriaArray []string
		if json.Unmarshal(eligibilityCriteriaBytes, &criteriaArray) == nil && len(criteriaArray) > 0 {
			criteriaStr := strings.Join(criteriaArray, "; ")
			scholarship.EligibilityCriteria = &criteriaStr
		}
	}

	if requiredDocumentsBytes != nil {
		var docsArray []string
		if json.Unmarshal(requiredDocumentsBytes, &docsArray) == nil && len(docsArray) > 0 {
			docsStr := strings.Join(docsArray, ", ")
			scholarship.RequiredDocuments = &docsStr
		}
	}

	scholarship.Source = source
	return scholarship, nil
}

func (r *ScholarshipRepository) List(limit, offset int, search, scholarshipType, academicYear string, activeOnly bool) ([]models.Scholarship, int, error) {
	var scholarships []models.Scholarship
	var totalCount int

	// Build WHERE conditions
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if activeOnly {
		whereConditions = append(whereConditions, "s.is_active = true")
	}

	if search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(s.name ILIKE $%d)", argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if scholarshipType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("s.type = $%d", argIndex))
		args = append(args, scholarshipType)
		argIndex++
	}

	if academicYear != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("s.academic_year = $%d", argIndex))
		args = append(args, academicYear)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", whereConditions[0])
		for i := 1; i < len(whereConditions); i++ {
			whereClause += " AND " + whereConditions[i]
		}
	}

	// Count total records
	countQuery := `SELECT COUNT(*) FROM scholarships s ` + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get scholarships with sources
	query := `
		SELECT s.scholarship_id, s.source_id, s.name, s.type, s.amount,
		       s.total_quota, s.available_quota, s.academic_year, s.semester, s.eligibility_criteria,
		       s.required_documents, s.application_start_date, s.application_end_date,
		       s.interview_required, s.is_active, s.created_by, s.created_at, s.updated_at,
		       src.source_id, src.source_name, src.source_type, src.contact_person,
		       src.contact_email, src.contact_phone, src.description, src.is_active,
		       src.created_at, src.updated_at
		FROM scholarships s
		LEFT JOIN scholarship_sources src ON s.source_id = src.source_id ` +
		whereClause +
		fmt.Sprintf(" ORDER BY s.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var scholarship models.Scholarship
		var source models.ScholarshipSource
		var eligibilityCriteriaBytes []byte
		var requiredDocumentsBytes []byte

		err := rows.Scan(
			&scholarship.ScholarshipID,
			&scholarship.SourceID,
			&scholarship.ScholarshipName,
			&scholarship.ScholarshipType,
			&scholarship.Amount,
			&scholarship.TotalQuota,
			&scholarship.AvailableQuota,
			&scholarship.AcademicYear,
			&scholarship.Semester,
			&eligibilityCriteriaBytes,
			&requiredDocumentsBytes,
			&scholarship.ApplicationStartDate,
			&scholarship.ApplicationEndDate,
			&scholarship.InterviewRequired,
			&scholarship.IsActive,
			&scholarship.CreatedBy,
			&scholarship.CreatedAt,
			&scholarship.UpdatedAt,
			&source.SourceID,
			&source.SourceName,
			&source.SourceType,
			&source.ContactPerson,
			&source.ContactEmail,
			&source.ContactPhone,
			&source.Description,
			&source.IsActive,
			&source.CreatedAt,
			&source.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		// Convert JSONB back to string
		if eligibilityCriteriaBytes != nil {
			var criteriaArray []string
			if json.Unmarshal(eligibilityCriteriaBytes, &criteriaArray) == nil && len(criteriaArray) > 0 {
				criteriaStr := strings.Join(criteriaArray, "; ")
				scholarship.EligibilityCriteria = &criteriaStr
			}
		}

		if requiredDocumentsBytes != nil {
			var docsArray []string
			if json.Unmarshal(requiredDocumentsBytes, &docsArray) == nil && len(docsArray) > 0 {
				docsStr := strings.Join(docsArray, ", ")
				scholarship.RequiredDocuments = &docsStr
			}
		}

		scholarship.Source = &source
		scholarships = append(scholarships, scholarship)
	}

	return scholarships, totalCount, nil
}

func (r *ScholarshipRepository) Update(scholarship *models.Scholarship) error {
	query := `
		UPDATE scholarships
		SET source_id = $2, name = $3, type = $4, amount = $5,
		    total_quota = $6, available_quota = $7, academic_year = $8, semester = $9,
		    eligibility_criteria = $10, required_documents = $11, application_start_date = $12,
		    application_end_date = $13, interview_required = $14, is_active = $15, updated_at = $16
		WHERE scholarship_id = $1
	`

	scholarship.UpdatedAt = time.Now()

	// Convert string fields to proper JSON format for JSONB columns
	var eligibilityCriteriaJSON interface{}
	var requiredDocumentsJSON interface{}

	if scholarship.EligibilityCriteria != nil && *scholarship.EligibilityCriteria != "" {
		criteriaArray := []string{*scholarship.EligibilityCriteria}
		eligibilityCriteriaJSON, _ = json.Marshal(criteriaArray)
	} else {
		eligibilityCriteriaJSON = nil
	}

	if scholarship.RequiredDocuments != nil && *scholarship.RequiredDocuments != "" {
		docs := []string{}
		parts := strings.Split(*scholarship.RequiredDocuments, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				docs = append(docs, trimmed)
			}
		}
		requiredDocumentsJSON, _ = json.Marshal(docs)
	} else {
		requiredDocumentsJSON = nil
	}

	_, err := r.db.Exec(query,
		scholarship.ScholarshipID,
		scholarship.SourceID,
		scholarship.ScholarshipName,
		scholarship.ScholarshipType,
		scholarship.Amount,
		scholarship.TotalQuota,
		scholarship.AvailableQuota,
		scholarship.AcademicYear,
		scholarship.Semester,
		eligibilityCriteriaJSON,
		requiredDocumentsJSON,
		scholarship.ApplicationStartDate,
		scholarship.ApplicationEndDate,
		scholarship.InterviewRequired,
		scholarship.IsActive,
		scholarship.UpdatedAt,
	)

	return err
}

func (r *ScholarshipRepository) UpdateQuota(scholarshipID uint, availableQuota int) error {
	query := `UPDATE scholarships SET available_quota = $1, updated_at = $2 WHERE scholarship_id = $3`
	_, err := r.db.Exec(query, availableQuota, time.Now(), scholarshipID)
	return err
}

func (r *ScholarshipRepository) Delete(scholarshipID uint) error {
	query := `DELETE FROM scholarships WHERE scholarship_id = $1`
	result, err := r.db.Exec(query, scholarshipID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ScholarshipRepository) GetAvailableScholarships() ([]models.Scholarship, error) {
	query := `
		SELECT s.scholarship_id, s.source_id, s.name, s.type, s.amount,
		       s.total_quota, s.available_quota, s.academic_year, s.semester, s.eligibility_criteria,
		       s.required_documents, s.application_start_date, s.application_end_date,
		       s.interview_required, s.is_active, s.created_by, s.created_at, s.updated_at,
		       src.source_name, src.source_type
		FROM scholarships s
		LEFT JOIN scholarship_sources src ON s.source_id = src.source_id
		WHERE s.is_active = true
		AND s.available_quota > 0
		AND s.application_start_date <= CURRENT_DATE
		AND s.application_end_date >= CURRENT_DATE
		ORDER BY s.application_end_date ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scholarships []models.Scholarship

	for rows.Next() {
		var scholarship models.Scholarship
		var sourceName, sourceType sql.NullString
		var eligibilityCriteriaBytes []byte
		var requiredDocumentsBytes []byte

		err := rows.Scan(
			&scholarship.ScholarshipID,
			&scholarship.SourceID,
			&scholarship.ScholarshipName,
			&scholarship.ScholarshipType,
			&scholarship.Amount,
			&scholarship.TotalQuota,
			&scholarship.AvailableQuota,
			&scholarship.AcademicYear,
			&scholarship.Semester,
			&eligibilityCriteriaBytes,
			&requiredDocumentsBytes,
			&scholarship.ApplicationStartDate,
			&scholarship.ApplicationEndDate,
			&scholarship.InterviewRequired,
			&scholarship.IsActive,
			&scholarship.CreatedBy,
			&scholarship.CreatedAt,
			&scholarship.UpdatedAt,
			&sourceName,
			&sourceType,
		)

		if err != nil {
			return nil, err
		}

		// Convert JSONB back to string
		if eligibilityCriteriaBytes != nil {
			var criteriaArray []string
			if json.Unmarshal(eligibilityCriteriaBytes, &criteriaArray) == nil && len(criteriaArray) > 0 {
				criteriaStr := strings.Join(criteriaArray, "; ")
				scholarship.EligibilityCriteria = &criteriaStr
			}
		}

		if requiredDocumentsBytes != nil {
			var docsArray []string
			if json.Unmarshal(requiredDocumentsBytes, &docsArray) == nil && len(docsArray) > 0 {
				docsStr := strings.Join(docsArray, ", ")
				scholarship.RequiredDocuments = &docsStr
			}
		}

		if sourceName.Valid && sourceType.Valid {
			scholarship.Source = &models.ScholarshipSource{
				SourceName: sourceName.String,
				SourceType: sourceType.String,
			}
		}

		scholarships = append(scholarships, scholarship)
	}

	return scholarships, nil
}
