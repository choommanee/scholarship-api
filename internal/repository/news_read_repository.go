package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

// NewsReadRepository handles database operations for news read tracking
type NewsReadRepository struct {
	db *sql.DB
}

// NewNewsReadRepository creates a new NewsReadRepository
func NewNewsReadRepository() *NewsReadRepository {
	return &NewsReadRepository{
		db: database.DB,
	}
}

// MarkNewsAsRead marks a news article as read by a user
func (r *NewsReadRepository) MarkNewsAsRead(userID string, newsID uuid.UUID) error {
	// Check if the user has already read this news
	var count int
	row := r.db.QueryRow(`
		SELECT COUNT(*) FROM news_read 
		WHERE user_id = $1 AND news_id = $2
	`, userID, newsID)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	// If the user has already read this news, do nothing
	if count > 0 {
		return nil
	}

	// Otherwise, insert a new record
	now := time.Now()
	readID := uuid.New()

	_, err = r.db.Exec(`
		INSERT INTO news_read (id, user_id, news_id, read_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, readID, userID, newsID, now, now, now)

	return err
}

// GetUnreadNewsCount returns the count of unread news for a user
func (r *NewsReadRepository) GetUnreadNewsCount(userID string) (int, error) {
	var count int
	row := r.db.QueryRow(`
		SELECT COUNT(*) FROM news n
		WHERE n.is_published = true
		AND n.publish_date <= NOW()
		AND (n.expire_date IS NULL OR n.expire_date > NOW())
		AND NOT EXISTS (
			SELECT 1 FROM news_read nr
			WHERE nr.news_id = n.id AND nr.user_id = $1
		)
	`, userID)

	err := row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, nil
}

// GetUnreadNews returns the list of unread news for a user
func (r *NewsReadRepository) GetUnreadNews(userID string, limit, offset int) ([]*models.News, error) {
	rows, err := r.db.Query(`
		SELECT n.id, n.title, n.content, n.summary, n.image_url, n.publish_date, n.expire_date, 
		n.category, n.tags, n.is_published, n.created_by, n.created_at, n.updated_at 
		FROM news n
		WHERE n.is_published = true
		AND n.publish_date <= NOW()
		AND (n.expire_date IS NULL OR n.expire_date > NOW())
		AND NOT EXISTS (
			SELECT 1 FROM news_read nr
			WHERE nr.news_id = n.id AND nr.user_id = $1
		)
		ORDER BY n.publish_date DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := []*models.News{}
	for rows.Next() {
		n := &models.News{}
		var tags []byte
		err := rows.Scan(
			&n.ID, &n.Title, &n.Content, &n.Summary, &n.ImageURL, &n.PublishDate, &n.ExpireDate,
			&n.Category, &tags, &n.IsPublished, &n.CreatedBy, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse tags from JSON
		if len(tags) > 0 {
			if err := json.Unmarshal(tags, &n.Tags); err != nil {
				return nil, err
			}
		}

		news = append(news, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}
