package repository

import (
	"database/sql"

	"scholarship-system/internal/models"
)

// AnalyticsRepository handles analytics-related database operations
type AnalyticsRepository struct {
	db *sql.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// CreateStatistics creates scholarship statistics
func (r *AnalyticsRepository) CreateStatistics(stats *models.ScholarshipStatistics) error {
	query := `
		INSERT INTO scholarship_statistics (
			stat_id, academic_year, scholarship_round, total_applications,
			approved_applications, rejected_applications, total_budget,
			allocated_budget, remaining_budget, average_amount, success_rate,
			processing_time_avg, total_faculties, most_popular_scholarship
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (academic_year, scholarship_round)
		DO UPDATE SET
			total_applications = EXCLUDED.total_applications,
			approved_applications = EXCLUDED.approved_applications,
			rejected_applications = EXCLUDED.rejected_applications,
			total_budget = EXCLUDED.total_budget,
			allocated_budget = EXCLUDED.allocated_budget,
			remaining_budget = EXCLUDED.remaining_budget,
			average_amount = EXCLUDED.average_amount,
			success_rate = EXCLUDED.success_rate,
			processing_time_avg = EXCLUDED.processing_time_avg,
			total_faculties = EXCLUDED.total_faculties,
			most_popular_scholarship = EXCLUDED.most_popular_scholarship
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		stats.StatID, stats.AcademicYear, stats.ScholarshipRound, stats.TotalApplications,
		stats.ApprovedApplications, stats.RejectedApplications, stats.TotalBudget,
		stats.AllocatedBudget, stats.RemainingBudget, stats.AverageAmount, stats.SuccessRate,
		stats.ProcessingTimeAvg, stats.TotalFaculties, stats.MostPopularScholarship,
	).Scan(&stats.CreatedAt)
}

// GetStatistics retrieves statistics for a specific year and round
func (r *AnalyticsRepository) GetStatistics(year, round string) (*models.ScholarshipStatistics, error) {
	query := `
		SELECT stat_id, academic_year, scholarship_round, total_applications,
			   approved_applications, rejected_applications, total_budget,
			   allocated_budget, remaining_budget, average_amount, success_rate,
			   processing_time_avg, total_faculties, most_popular_scholarship, created_at
		FROM scholarship_statistics
		WHERE academic_year = $1 AND scholarship_round = $2`

	stats := &models.ScholarshipStatistics{}
	err := r.db.QueryRow(query, year, round).Scan(
		&stats.StatID, &stats.AcademicYear, &stats.ScholarshipRound, &stats.TotalApplications,
		&stats.ApprovedApplications, &stats.RejectedApplications, &stats.TotalBudget,
		&stats.AllocatedBudget, &stats.RemainingBudget, &stats.AverageAmount, &stats.SuccessRate,
		&stats.ProcessingTimeAvg, &stats.TotalFaculties, &stats.MostPopularScholarship, &stats.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// GetAllStatistics retrieves all statistics
func (r *AnalyticsRepository) GetAllStatistics() ([]models.ScholarshipStatistics, error) {
	query := `
		SELECT stat_id, academic_year, scholarship_round, total_applications,
			   approved_applications, rejected_applications, total_budget,
			   allocated_budget, remaining_budget, average_amount, success_rate,
			   processing_time_avg, total_faculties, most_popular_scholarship, created_at
		FROM scholarship_statistics
		ORDER BY academic_year DESC, scholarship_round`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statistics []models.ScholarshipStatistics
	for rows.Next() {
		var s models.ScholarshipStatistics
		err := rows.Scan(
			&s.StatID, &s.AcademicYear, &s.ScholarshipRound, &s.TotalApplications,
			&s.ApprovedApplications, &s.RejectedApplications, &s.TotalBudget,
			&s.AllocatedBudget, &s.RemainingBudget, &s.AverageAmount, &s.SuccessRate,
			&s.ProcessingTimeAvg, &s.TotalFaculties, &s.MostPopularScholarship, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		statistics = append(statistics, s)
	}
	return statistics, nil
}

// CreateApplicationAnalytics creates application analytics
func (r *AnalyticsRepository) CreateApplicationAnalytics(analytics *models.ApplicationAnalytics) error {
	query := `
		INSERT INTO application_analytics (
			analytics_id, application_id, processing_time, total_steps,
			completed_steps, bottleneck_step, time_in_each_step,
			document_upload_time, review_time, interview_score, final_score
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at`

	return r.db.QueryRow(
		query,
		analytics.AnalyticsID, analytics.ApplicationID, analytics.ProcessingTime,
		analytics.TotalSteps, analytics.CompletedSteps, analytics.BottleneckStep,
		analytics.TimeInEachStep, analytics.DocumentUploadTime, analytics.ReviewTime,
		analytics.InterviewScore, analytics.FinalScore,
	).Scan(&analytics.CreatedAt)
}

// GetApplicationAnalytics retrieves analytics for a specific application
func (r *AnalyticsRepository) GetApplicationAnalytics(applicationID int) (*models.ApplicationAnalytics, error) {
	query := `
		SELECT analytics_id, application_id, processing_time, total_steps,
			   completed_steps, bottleneck_step, time_in_each_step,
			   document_upload_time, review_time, interview_score, final_score, created_at
		FROM application_analytics
		WHERE application_id = $1`

	analytics := &models.ApplicationAnalytics{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&analytics.AnalyticsID, &analytics.ApplicationID, &analytics.ProcessingTime,
		&analytics.TotalSteps, &analytics.CompletedSteps, &analytics.BottleneckStep,
		&analytics.TimeInEachStep, &analytics.DocumentUploadTime, &analytics.ReviewTime,
		&analytics.InterviewScore, &analytics.FinalScore, &analytics.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return analytics, nil
}

// GetAverageProcessingTime gets average processing time
func (r *AnalyticsRepository) GetAverageProcessingTime() (float64, error) {
	query := `SELECT AVG(processing_time) FROM application_analytics`

	var avg sql.NullFloat64
	err := r.db.QueryRow(query).Scan(&avg)
	if err != nil {
		return 0, err
	}

	if !avg.Valid {
		return 0, nil
	}
	return avg.Float64, nil
}

// GetBottleneckSteps identifies common bottleneck steps
func (r *AnalyticsRepository) GetBottleneckSteps() (map[string]int, error) {
	query := `
		SELECT bottleneck_step, COUNT(*) as count
		FROM application_analytics
		WHERE bottleneck_step IS NOT NULL
		GROUP BY bottleneck_step
		ORDER BY count DESC
		LIMIT 10`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bottlenecks := make(map[string]int)
	for rows.Next() {
		var step string
		var count int
		err := rows.Scan(&step, &count)
		if err != nil {
			return nil, err
		}
		bottlenecks[step] = count
	}
	return bottlenecks, nil
}
