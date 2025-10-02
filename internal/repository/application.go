package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository() *ApplicationRepository {
	return &ApplicationRepository{
		db: database.DB,
	}
}

func (r *ApplicationRepository) Create(application *models.ScholarshipApplication) error {
	query := `
		INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, application_data, 
		                                    family_income, monthly_expenses, siblings_count, special_abilities, 
		                                    activities_participation, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING application_id
	`
	
	now := time.Now()
	application.CreatedAt = now
	application.UpdatedAt = now
	
	err := r.db.QueryRow(query,
		application.StudentID,
		application.ScholarshipID,
		application.ApplicationStatus,
		application.ApplicationData,
		application.FamilyIncome,
		application.MonthlyExpenses,
		application.SiblingsCount,
		application.SpecialAbilities,
		application.ActivitiesParticipation,
		application.CreatedAt,
		application.UpdatedAt,
	).Scan(&application.ApplicationID)
	
	return err
}

func (r *ApplicationRepository) GetByID(applicationID uint) (*models.ScholarshipApplication, error) {
	query := `
		SELECT sa.application_id, sa.student_id, sa.scholarship_id, sa.application_status, 
		       sa.application_data, sa.family_income, sa.monthly_expenses, sa.siblings_count, 
		       sa.special_abilities, sa.activities_participation, sa.submitted_at, sa.reviewed_by, 
		       sa.reviewed_at, sa.review_notes, sa.priority_score, sa.created_at, sa.updated_at,
		       s.name as scholarship_name, s.type as scholarship_type, s.amount,
		       st.user_id, u.first_name, u.last_name, u.email
		FROM scholarship_applications sa
		LEFT JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		LEFT JOIN students st ON sa.student_id = st.student_id
		LEFT JOIN users u ON st.user_id = u.user_id
		WHERE sa.application_id = $1
	`
	
	application := &models.ScholarshipApplication{}
	var scholarshipName, scholarshipType sql.NullString
	var amount sql.NullFloat64
	var userID sql.NullString
	var firstName, lastName, email sql.NullString
	
	err := r.db.QueryRow(query, applicationID).Scan(
		&application.ApplicationID,
		&application.StudentID,
		&application.ScholarshipID,
		&application.ApplicationStatus,
		&application.ApplicationData,
		&application.FamilyIncome,
		&application.MonthlyExpenses,
		&application.SiblingsCount,
		&application.SpecialAbilities,
		&application.ActivitiesParticipation,
		&application.SubmittedAt,
		&application.ReviewedBy,
		&application.ReviewedAt,
		&application.ReviewNotes,
		&application.PriorityScore,
		&application.CreatedAt,
		&application.UpdatedAt,
		&scholarshipName,
		&scholarshipType,
		&amount,
		&userID,
		&firstName,
		&lastName,
		&email,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Load scholarship info if available
	if scholarshipName.Valid {
		application.Scholarship = &models.Scholarship{
			ScholarshipID:   application.ScholarshipID,
			ScholarshipName: scholarshipName.String,
			ScholarshipType: scholarshipType.String,
			Amount:          amount.Float64,
		}
	}
	
	// Load student info if available
	if userID.Valid {
		uuid, _ := uuid.Parse(userID.String)
		application.Student = &models.Student{
			StudentID: application.StudentID,
			UserID:    uuid,
		}
	}
	
	return application, nil
}

func (r *ApplicationRepository) GetByStudentAndScholarship(studentID string, scholarshipID uint) (*models.ScholarshipApplication, error) {
	query := `
		SELECT application_id, student_id, scholarship_id, application_status, application_data, 
		       family_income, monthly_expenses, siblings_count, special_abilities, 
		       activities_participation, submitted_at, reviewed_by, reviewed_at, review_notes, 
		       priority_score, created_at, updated_at
		FROM scholarship_applications 
		WHERE student_id = $1 AND scholarship_id = $2
	`
	
	application := &models.ScholarshipApplication{}
	err := r.db.QueryRow(query, studentID, scholarshipID).Scan(
		&application.ApplicationID,
		&application.StudentID,
		&application.ScholarshipID,
		&application.ApplicationStatus,
		&application.ApplicationData,
		&application.FamilyIncome,
		&application.MonthlyExpenses,
		&application.SiblingsCount,
		&application.SpecialAbilities,
		&application.ActivitiesParticipation,
		&application.SubmittedAt,
		&application.ReviewedBy,
		&application.ReviewedAt,
		&application.ReviewNotes,
		&application.PriorityScore,
		&application.CreatedAt,
		&application.UpdatedAt,
	)
	
	return application, err
}

func (r *ApplicationRepository) ListByStudent(studentID string, limit, offset int) ([]models.ScholarshipApplication, int, error) {
	var applications []models.ScholarshipApplication
	var totalCount int
	
	// Count total records
	countQuery := `SELECT COUNT(*) FROM scholarship_applications WHERE student_id = $1`
	err := r.db.QueryRow(countQuery, studentID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	
	// Get applications with scholarship info
	query := `
		SELECT sa.application_id, sa.student_id, sa.scholarship_id, sa.application_status, 
		       sa.application_data, sa.family_income, sa.monthly_expenses, sa.siblings_count, 
		       sa.special_abilities, sa.activities_participation, sa.submitted_at, sa.reviewed_by, 
		       sa.reviewed_at, sa.review_notes, sa.priority_score, sa.created_at, sa.updated_at,
		       s.name as scholarship_name, s.type as scholarship_type, s.amount, s.academic_year
		FROM scholarship_applications sa
		LEFT JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		WHERE sa.student_id = $1
		ORDER BY sa.created_at DESC LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var application models.ScholarshipApplication
		var scholarshipName, scholarshipType, academicYear sql.NullString
		var amount sql.NullFloat64
		
		err := rows.Scan(
			&application.ApplicationID,
			&application.StudentID,
			&application.ScholarshipID,
			&application.ApplicationStatus,
			&application.ApplicationData,
			&application.FamilyIncome,
			&application.MonthlyExpenses,
			&application.SiblingsCount,
			&application.SpecialAbilities,
			&application.ActivitiesParticipation,
			&application.SubmittedAt,
			&application.ReviewedBy,
			&application.ReviewedAt,
			&application.ReviewNotes,
			&application.PriorityScore,
			&application.CreatedAt,
			&application.UpdatedAt,
			&scholarshipName,
			&scholarshipType,
			&amount,
			&academicYear,
		)
		
		if err != nil {
			return nil, 0, err
		}
		
		// Load scholarship info if available
		if scholarshipName.Valid {
			application.Scholarship = &models.Scholarship{
				ScholarshipID:   application.ScholarshipID,
				ScholarshipName: scholarshipName.String,
				ScholarshipType: scholarshipType.String,
				Amount:          amount.Float64,
				AcademicYear:    academicYear.String,
			}
		}
		
		applications = append(applications, application)
	}
	
	return applications, totalCount, nil
}

func (r *ApplicationRepository) List(limit, offset int, status, scholarshipType string, scholarshipID *uint) ([]models.ScholarshipApplication, int, error) {
	var applications []models.ScholarshipApplication
	var totalCount int
	
	// Build WHERE conditions
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1
	
	if status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.application_status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}
	
	if scholarshipType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("s.type = $%d", argIndex))
		args = append(args, scholarshipType)
		argIndex++
	}
	
	if scholarshipID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.scholarship_id = $%d", argIndex))
		args = append(args, *scholarshipID)
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
	countQuery := `
		SELECT COUNT(*) 
		FROM scholarship_applications sa
		LEFT JOIN scholarships s ON sa.scholarship_id = s.scholarship_id ` + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	
	// Get applications
	query := `
		SELECT sa.application_id, sa.student_id, sa.scholarship_id, sa.application_status, 
		       sa.application_data, sa.family_income, sa.monthly_expenses, sa.siblings_count, 
		       sa.special_abilities, sa.activities_participation, sa.submitted_at, sa.reviewed_by, 
		       sa.reviewed_at, sa.review_notes, sa.priority_score, sa.created_at, sa.updated_at,
		       s.name as scholarship_name, s.type as scholarship_type, s.amount, s.academic_year,
		       st.user_id, u.first_name, u.last_name, u.email
		FROM scholarship_applications sa
		LEFT JOIN scholarships s ON sa.scholarship_id = s.scholarship_id
		LEFT JOIN students st ON sa.student_id = st.student_id
		LEFT JOIN users u ON st.user_id = u.user_id ` + 
		whereClause + 
		fmt.Sprintf(" ORDER BY sa.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var application models.ScholarshipApplication
		var scholarshipName, scholarshipType, academicYear sql.NullString
		var amount sql.NullFloat64
		var userID sql.NullString
		var firstName, lastName, email sql.NullString
		
		err := rows.Scan(
			&application.ApplicationID,
			&application.StudentID,
			&application.ScholarshipID,
			&application.ApplicationStatus,
			&application.ApplicationData,
			&application.FamilyIncome,
			&application.MonthlyExpenses,
			&application.SiblingsCount,
			&application.SpecialAbilities,
			&application.ActivitiesParticipation,
			&application.SubmittedAt,
			&application.ReviewedBy,
			&application.ReviewedAt,
			&application.ReviewNotes,
			&application.PriorityScore,
			&application.CreatedAt,
			&application.UpdatedAt,
			&scholarshipName,
			&scholarshipType,
			&amount,
			&academicYear,
			&userID,
			&firstName,
			&lastName,
			&email,
		)
		
		if err != nil {
			return nil, 0, err
		}
		
		// Load scholarship info if available
		if scholarshipName.Valid {
			application.Scholarship = &models.Scholarship{
				ScholarshipID:   application.ScholarshipID,
				ScholarshipName: scholarshipName.String,
				ScholarshipType: scholarshipType.String,
				Amount:          amount.Float64,
				AcademicYear:    academicYear.String,
			}
		}
		
		// Load student info if available
		if userID.Valid {
			uuid, _ := uuid.Parse(userID.String)
			application.Student = &models.Student{
				StudentID: application.StudentID,
				UserID:    uuid,
			}
		}
		
		applications = append(applications, application)
	}
	
	return applications, totalCount, nil
}

func (r *ApplicationRepository) Update(application *models.ScholarshipApplication) error {
	query := `
		UPDATE scholarship_applications 
		SET application_status = $2, application_data = $3, family_income = $4, monthly_expenses = $5, 
		    siblings_count = $6, special_abilities = $7, activities_participation = $8, 
		    submitted_at = $9, reviewed_by = $10, reviewed_at = $11, review_notes = $12, 
		    priority_score = $13, updated_at = $14
		WHERE application_id = $1
	`
	
	application.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		application.ApplicationID,
		application.ApplicationStatus,
		application.ApplicationData,
		application.FamilyIncome,
		application.MonthlyExpenses,
		application.SiblingsCount,
		application.SpecialAbilities,
		application.ActivitiesParticipation,
		application.SubmittedAt,
		application.ReviewedBy,
		application.ReviewedAt,
		application.ReviewNotes,
		application.PriorityScore,
		application.UpdatedAt,
	)
	
	return err
}

func (r *ApplicationRepository) UpdateStatus(applicationID uint, status string, reviewedBy *uuid.UUID, reviewNotes *string) error {
	query := `
		UPDATE scholarship_applications 
		SET application_status = $2, reviewed_by = $3, reviewed_at = $4, review_notes = $5, updated_at = $6
		WHERE application_id = $1
	`
	
	now := time.Now()
	_, err := r.db.Exec(query, applicationID, status, reviewedBy, &now, reviewNotes, now)
	return err
}

func (r *ApplicationRepository) Submit(applicationID uint) error {
	query := `
		UPDATE scholarship_applications 
		SET application_status = 'submitted', submitted_at = $2, updated_at = $3
		WHERE application_id = $1
	`
	
	now := time.Now()
	_, err := r.db.Exec(query, applicationID, now, now)
	return err
}

func (r *ApplicationRepository) Delete(applicationID uint) error {
	query := `DELETE FROM scholarship_applications WHERE application_id = $1`
	_, err := r.db.Exec(query, applicationID)
	return err
}

// Document methods
func (r *ApplicationRepository) AddDocument(doc *models.ApplicationDocument) error {
	query := `
		INSERT INTO application_documents (application_id, document_type, document_name, file_path, 
		                                 file_size, mime_type, upload_status, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING document_id
	`
	
	doc.UploadedAt = time.Now()
	
	err := r.db.QueryRow(query,
		doc.ApplicationID,
		doc.DocumentType,
		doc.DocumentName,
		doc.FilePath,
		doc.FileSize,
		doc.MimeType,
		doc.UploadStatus,
		doc.UploadedAt,
	).Scan(&doc.DocumentID)
	
	return err
}

func (r *ApplicationRepository) GetDocuments(applicationID uint) ([]models.ApplicationDocument, error) {
	query := `
		SELECT document_id, application_id, document_type, document_name, file_path, 
		       file_size, mime_type, upload_status, verification_notes, uploaded_at, 
		       verified_by, verified_at
		FROM application_documents 
		WHERE application_id = $1
		ORDER BY uploaded_at DESC
	`
	
	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var documents []models.ApplicationDocument
	
	for rows.Next() {
		var doc models.ApplicationDocument
		
		err := rows.Scan(
			&doc.DocumentID,
			&doc.ApplicationID,
			&doc.DocumentType,
			&doc.DocumentName,
			&doc.FilePath,
			&doc.FileSize,
			&doc.MimeType,
			&doc.UploadStatus,
			&doc.VerificationNotes,
			&doc.UploadedAt,
			&doc.VerifiedBy,
			&doc.VerifiedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		documents = append(documents, doc)
	}
	
	return documents, nil
}