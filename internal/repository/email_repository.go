package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"scholarship-system/internal/models"
)

// EmailRepository handles email-related database operations
type EmailRepository struct {
	db *sql.DB
}

// NewEmailRepository creates a new email repository
func NewEmailRepository(db *sql.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

// CreateEmailQueue adds an email to the queue
func (r *EmailRepository) CreateEmailQueue(email *models.EmailQueue) error {
	query := `
		INSERT INTO email_queue (
			queue_id, recipient_email, recipient_name, sender_email,
			subject, body, template_id, priority, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		email.QueueID, email.RecipientEmail, email.RecipientName, email.SenderEmail,
		email.Subject, email.Body, email.TemplateID, email.Priority, email.Status,
	).Scan(&email.CreatedAt)
}

// GetPendingEmails retrieves pending emails from queue
func (r *EmailRepository) GetPendingEmails(limit int) ([]models.EmailQueue, error) {
	query := `
		SELECT queue_id, recipient_email, recipient_name, sender_email,
			   subject, body, template_id, priority, status, sent_at, error_message, created_at
		FROM email_queue
		WHERE status = 'pending'
		ORDER BY priority DESC, created_at ASC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []models.EmailQueue
	for rows.Next() {
		var e models.EmailQueue
		err := rows.Scan(
			&e.QueueID, &e.RecipientEmail, &e.RecipientName, &e.SenderEmail,
			&e.Subject, &e.Body, &e.TemplateID, &e.Priority, &e.Status,
			&e.SentAt, &e.ErrorMessage, &e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		emails = append(emails, e)
	}
	return emails, nil
}

// UpdateEmailStatus updates email status
func (r *EmailRepository) UpdateEmailStatus(queueID uuid.UUID, status string, errorMsg *string) error {
	var query string
	if status == "sent" {
		query = `
			UPDATE email_queue
			SET status = $1, sent_at = CURRENT_TIMESTAMP
			WHERE queue_id = $2`
		_, err := r.db.Exec(query, status, queueID)
		return err
	}

	query = `
		UPDATE email_queue
		SET status = $1, error_message = $2
		WHERE queue_id = $3`
	_, err := r.db.Exec(query, status, errorMsg, queueID)
	return err
}

// GetTemplateByType retrieves email template by type
func (r *EmailRepository) GetTemplateByType(templateType string) (*models.EmailTemplate, error) {
	query := `
		SELECT template_id, template_name, subject, body, variables, template_type, is_active, created_at
		FROM email_templates
		WHERE template_type = $1 AND is_active = true
		LIMIT 1`

	template := &models.EmailTemplate{}
	err := r.db.QueryRow(query, templateType).Scan(
		&template.TemplateID, &template.TemplateName, &template.Subject, &template.Body,
		&template.Variables, &template.TemplateType, &template.IsActive, &template.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// GetTemplateByID retrieves email template by ID
func (r *EmailRepository) GetTemplateByID(id uuid.UUID) (*models.EmailTemplate, error) {
	query := `
		SELECT template_id, template_name, subject, body, variables, template_type, is_active, created_at
		FROM email_templates
		WHERE template_id = $1`

	template := &models.EmailTemplate{}
	err := r.db.QueryRow(query, id).Scan(
		&template.TemplateID, &template.TemplateName, &template.Subject, &template.Body,
		&template.Variables, &template.TemplateType, &template.IsActive, &template.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return template, nil
}

// GetAllTemplates retrieves all email templates
func (r *EmailRepository) GetAllTemplates() ([]models.EmailTemplate, error) {
	query := `
		SELECT template_id, template_name, subject, body, variables, template_type, is_active, created_at
		FROM email_templates
		ORDER BY template_type, template_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.EmailTemplate
	for rows.Next() {
		var t models.EmailTemplate
		err := rows.Scan(
			&t.TemplateID, &t.TemplateName, &t.Subject, &t.Body,
			&t.Variables, &t.TemplateType, &t.IsActive, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

// CreateTemplate creates a new email template
func (r *EmailRepository) CreateTemplate(template *models.EmailTemplate) error {
	query := `
		INSERT INTO email_templates (
			template_id, template_name, subject, body, variables, template_type, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		template.TemplateID, template.TemplateName, template.Subject, template.Body,
		template.Variables, template.TemplateType, template.IsActive,
	).Scan(&template.CreatedAt)
}

// UpdateTemplate updates an email template
func (r *EmailRepository) UpdateTemplate(template *models.EmailTemplate) error {
	query := `
		UPDATE email_templates
		SET template_name = $1, subject = $2, body = $3, variables = $4, is_active = $5
		WHERE template_id = $6`

	_, err := r.db.Exec(
		query,
		template.TemplateName, template.Subject, template.Body,
		template.Variables, template.IsActive, template.TemplateID,
	)
	return err
}

// DeleteTemplate deletes an email template
func (r *EmailRepository) DeleteTemplate(id uuid.UUID) error {
	query := `DELETE FROM email_templates WHERE template_id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
