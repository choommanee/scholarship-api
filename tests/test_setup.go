package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
)

// TestSuite provides common test setup
type TestSuite struct {
	suite.Suite
	App *fiber.App
	DB  *sql.DB
	Config *config.Config
}

// SetupSuite runs before all tests
func (s *TestSuite) SetupSuite() {
	// Load test configuration
	s.Config = &config.Config{
		DBHost:     getEnv("TEST_DB_HOST", "localhost"),
		DBPort:     getEnv("TEST_DB_PORT", "5432"),
		DBUser:     getEnv("TEST_DB_USER", "postgres"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:     getEnv("TEST_DB_NAME", "scholarship_test"),
		JWTSecret:  "test_jwt_secret_key_for_testing_only",
		Port:       "8081", // Different port for testing
	}

	// Connect to test database
	err := database.ConnectDB(s.Config)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	s.DB = database.DB

	// Create Fiber app
	s.App = fiber.New(fiber.Config{
		Testing: true,
	})
}

// TearDownSuite runs after all tests
func (s *TestSuite) TearDownSuite() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// SetupTest runs before each test
func (s *TestSuite) SetupTest() {
	// Clean up test data before each test
	s.cleanupTestData()
}

// TearDownTest runs after each test
func (s *TestSuite) TearDownTest() {
	// Clean up test data after each test
	s.cleanupTestData()
}

// cleanupTestData removes test data from database
func (s *TestSuite) cleanupTestData() {
	tables := []string{
		"interview_results",
		"interview_appointments", 
		"interview_schedules",
		"application_documents",
		"scholarship_allocations",
		"scholarship_applications",
		"notifications",
		"scholarship_budgets",
		"scholarships",
		"scholarship_sources",
		"students",
		"user_roles",
		"users",
	}

	for _, table := range tables {
		_, err := s.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE 1=1", table))
		if err != nil {
			log.Printf("Warning: Failed to clean table %s: %v", table, err)
		}
	}
}

// CreateTestUser creates a test user and returns user ID
func (s *TestSuite) CreateTestUser(email, role string) string {
	userID := "test-user-" + email
	
	// Insert user
	_, err := s.DB.Exec(`
		INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, is_active)
		VALUES ($1, $2, $3, '$2a$10$test.hash', 'Test', 'User', true)
	`, userID, email, email)
	s.Require().NoError(err)

	// Get role ID
	var roleID int
	err = s.DB.QueryRow("SELECT role_id FROM roles WHERE role_name = $1", role).Scan(&roleID)
	if err != nil {
		// Create role if it doesn't exist
		err = s.DB.QueryRow(`
			INSERT INTO roles (role_name, role_description, permissions) 
			VALUES ($1, $1, '{}') RETURNING role_id
		`, role).Scan(&roleID)
		s.Require().NoError(err)
	}

	// Assign role
	_, err = s.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id, assigned_by, is_active)
		VALUES ($1, $2, $1, true)
	`, userID, roleID)
	s.Require().NoError(err)

	return userID
}

// CreateTestStudent creates a test student profile
func (s *TestSuite) CreateTestStudent(userID string) string {
	studentID := "STU" + userID[len(userID)-6:]
	
	_, err := s.DB.Exec(`
		INSERT INTO students (student_id, user_id, faculty_code, department_code, 
			year_level, gpa, admission_year, student_status)
		VALUES ($1, $2, 'ENG', 'CSE', 3, 3.25, 2022, 'active')
	`, studentID, userID)
	s.Require().NoError(err)

	return studentID
}

// CreateTestScholarship creates a test scholarship
func (s *TestSuite) CreateTestScholarship(createdBy string) int {
	// Create source first
	var sourceID int
	err := s.DB.QueryRow(`
		INSERT INTO scholarship_sources (source_name, source_type, contact_person, contact_email)
		VALUES ('Test Foundation', 'foundation', 'Test Contact', 'test@example.com')
		RETURNING source_id
	`).Scan(&sourceID)
	s.Require().NoError(err)

	// Create scholarship
	var scholarshipID int
	err = s.DB.QueryRow(`
		INSERT INTO scholarships (source_id, scholarship_name, scholarship_type, amount, 
			total_quota, available_quota, academic_year, application_start_date, 
			application_end_date, is_active, created_by)
		VALUES ($1, 'Test Scholarship', 'merit', 10000.00, 10, 10, '2024', 
			'2024-01-01', '2024-12-31', true, $2)
		RETURNING scholarship_id
	`, sourceID, createdBy).Scan(&scholarshipID)
	s.Require().NoError(err)

	return scholarshipID
}

// GenerateTestJWT generates a test JWT token
func (s *TestSuite) GenerateTestJWT(userID, role string) string {
	// This is a simplified version - in real implementation you'd use proper JWT library
	return fmt.Sprintf("test.jwt.token.%s.%s", userID, role)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// TestMain sets up test environment
func TestMain(m *testing.M) {
	// Setup test database if needed
	setupTestDB()
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	teardownTestDB()
	
	os.Exit(code)
}

func setupTestDB() {
	// This would create test database if it doesn't exist
	// For now, assume test database exists
}

func teardownTestDB() {
	// This would clean up test database after all tests
}