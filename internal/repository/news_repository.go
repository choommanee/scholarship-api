package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{
		db: database.DB,
	}
}

// CreateNews creates a new news article
func (r *NewsRepository) CreateNews(news *models.News) error {
	news.ID = uuid.New()
	news.CreatedAt = time.Now()
	news.UpdatedAt = time.Now()

	query := `
		INSERT INTO news (
			id, title, content, summary, image_url, publish_date, expire_date, 
			category, tags, is_published, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id
	`

	return r.db.QueryRow(
		query,
		news.ID, news.Title, news.Content, news.Summary, news.ImageURL,
		news.PublishDate, news.ExpireDate, news.Category, pq.Array(news.Tags),
		news.IsPublished, news.CreatedBy, news.CreatedAt, news.UpdatedAt,
	).Scan(&news.ID)
}

// GetNewsByID retrieves a news article by ID
func (r *NewsRepository) GetNewsByID(id uuid.UUID) (*models.News, error) {
	var news models.News
	query := `
		SELECT 
			id, title, content, summary, image_url, publish_date, expire_date, 
			category, tags, is_published, created_by, created_at, updated_at
		FROM news
		WHERE id = $1
	`

	var tags []string
	err := r.db.QueryRow(query, id).Scan(
		&news.ID, &news.Title, &news.Content, &news.Summary, &news.ImageURL,
		&news.PublishDate, &news.ExpireDate, &news.Category, pq.Array(&tags),
		&news.IsPublished, &news.CreatedBy, &news.CreatedAt, &news.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("news not found")
		}
		return nil, err
	}

	news.Tags = tags
	return &news, nil
}

// UpdateNews updates an existing news article
func (r *NewsRepository) UpdateNews(news *models.News) error {
	news.UpdatedAt = time.Now()

	query := `
		UPDATE news SET
			title = $1,
			content = $2,
			summary = $3,
			image_url = $4,
			publish_date = $5,
			expire_date = $6,
			category = $7,
			tags = $8,
			is_published = $9,
			updated_at = $10
		WHERE id = $11
	`

	_, err := r.db.Exec(
		query,
		news.Title, news.Content, news.Summary, news.ImageURL,
		news.PublishDate, news.ExpireDate, news.Category, pq.Array(news.Tags),
		news.IsPublished, news.UpdatedAt, news.ID,
	)

	return err
}

// DeleteNews deletes a news article
func (r *NewsRepository) DeleteNews(id uuid.UUID) error {
	query := "DELETE FROM news WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

// ListNews retrieves a paginated list of news articles with optional filtering
func (r *NewsRepository) ListNews(limit, offset int, category, search string, publishedOnly bool) ([]*models.News, int, error) {
	var args []interface{}
	argCount := 1

	// Base query
	countQuery := "SELECT COUNT(*) FROM news WHERE 1=1"
	query := `
		SELECT 
			id, title, content, summary, image_url, publish_date, expire_date, 
			category, tags, is_published, created_by, created_at, updated_at
		FROM news
		WHERE 1=1
	`

	// Add filters
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		countQuery += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, category)
		argCount++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argCount, argCount)
		countQuery += fmt.Sprintf(" AND (title ILIKE $%d OR content ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	if publishedOnly {
		query += fmt.Sprintf(" AND is_published = $%d", argCount)
		countQuery += fmt.Sprintf(" AND is_published = $%d", argCount)
		args = append(args, true)
		argCount++
	}

	// Add pagination
	query += " ORDER BY publish_date DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, args[:argCount-1]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Execute main query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var newsList []*models.News
	for rows.Next() {
		var news models.News
		var tags []string

		err := rows.Scan(
			&news.ID, &news.Title, &news.Content, &news.Summary, &news.ImageURL,
			&news.PublishDate, &news.ExpireDate, &news.Category, pq.Array(&tags),
			&news.IsPublished, &news.CreatedBy, &news.CreatedAt, &news.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		news.Tags = tags
		newsList = append(newsList, &news)
	}

	return newsList, total, nil
}
