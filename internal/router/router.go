// Package router provides HTTP route configuration for the scholarship management system
package router

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "scholarship-system/docs" // Swagger docs
	"scholarship-system/internal/config"
	"scholarship-system/internal/handlers"
	"scholarship-system/internal/middleware"
)

// SetupRoutes configures all HTTP routes for the application
// @title Scholarship Management System API
// @version 1.0
// @description API for managing scholarship applications, allocations, and administration
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize all handlers
	authHandler := handlers.NewAuthHandler(cfg)
	scholarshipHandler := handlers.NewScholarshipHandler(cfg)
	applicationHandler := handlers.NewApplicationHandler(cfg)
	interviewHandler := handlers.NewInterviewHandler(cfg)
	allocationHandler := handlers.NewAllocationHandler(cfg)
	notificationHandler := handlers.NewNotificationHandler(cfg)
	reportHandler := handlers.NewReportHandler(cfg)
	documentHandler := handlers.NewDocumentHandler(cfg)
	userHandler := handlers.NewUserHandler(cfg)
	studentHandler := handlers.NewStudentHandler(cfg)
	adminHandler := handlers.NewAdminHandler(cfg)
	newsHandler := handlers.NewNewsHandler(cfg)
	profileHandler := handlers.NewProfileHandler(cfg)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Scholarship Management System API",
			"version": "1.0.0",
		})
	})

	// API routes - public routes
	api := app.Group("/api/v1")

	// Setup auth routes
	setupAuthRoutes(api, authHandler, cfg)

	// Setup public scholarship routes (no authentication)
	setupPublicScholarshipRoutes(api, scholarshipHandler)

	// Setup public news routes (no authentication)
	setupPublicNewsRoutes(api, newsHandler)

	// Protected routes with JWT middleware
	protected := api.Group("", middleware.JWTMiddleware(cfg))

	// Setup protected routes
	setupProtectedRoutes(protected,
		authHandler, scholarshipHandler, applicationHandler,
		interviewHandler, allocationHandler, notificationHandler,
		reportHandler, documentHandler, userHandler, studentHandler, adminHandler, newsHandler, profileHandler, cfg)

	// Setup admin application routes
	setupAdminApplicationRoutes(protected, applicationHandler)

	// Enhanced Application routes (protected)
	applicationEnhanced := handlers.NewApplicationEnhancedHandler()

	// Multi-step application
	protected.Post("/applications/multi-step", applicationEnhanced.StartMultiStepApplication)
	protected.Get("/applications/steps-config", applicationEnhanced.GetStepsConfiguration)

	// Payment routes (admin/officer only)
	paymentHandler := handlers.NewPaymentHandler(cfg)
	setupPaymentRoutes(protected, paymentHandler)

	// Analytics routes (admin/officer only)
	analyticsHandler := handlers.NewAnalyticsHandler(cfg)
	setupAnalyticsRoutes(protected, analyticsHandler)

	// Draft management
	protected.Post("/applications/draft", applicationEnhanced.SaveDraft)
	protected.Get("/applications/draft", applicationEnhanced.LoadDraft)
	protected.Delete("/applications/draft", applicationEnhanced.DeleteDraft)

	// Document management
	protected.Post("/documents/bulk-upload", applicationEnhanced.StartBulkUpload)
	protected.Post("/documents/bulk-upload/files", applicationEnhanced.UploadBulkFiles)
	protected.Get("/documents/bulk-upload/progress", applicationEnhanced.GetUploadProgress)

	// Validation
	protected.Post("/applications/validate", applicationEnhanced.ValidateApplication)
	protected.Get("/applications/validation-rules", applicationEnhanced.GetValidationRules)

	// Preview
	protected.Post("/applications/preview", applicationEnhanced.PreviewApplication)

	// Interview & Review Management Routes
	interviewReviewHandler := handlers.NewInterviewReviewHandler(cfg)

	// Interview Slot Management (Officer/Admin only)
	protected.Post("/interview/slots", middleware.RequireRole("scholarship_officer", "admin"), interviewReviewHandler.CreateInterviewSlot)
	protected.Get("/interview/slots", middleware.RequireRole("scholarship_officer", "admin", "interviewer"), interviewReviewHandler.GetInterviewSlots)
	protected.Put("/interview/slots/:id", middleware.RequireRole("scholarship_officer", "admin"), interviewReviewHandler.UpdateInterviewSlot)
	protected.Delete("/interview/slots/:id", middleware.RequireRole("scholarship_officer", "admin"), interviewReviewHandler.DeleteInterviewSlot)

	// Interview Booking (Student & Officer)
	protected.Get("/interview/availability", middleware.JWTMiddleware(cfg), interviewReviewHandler.GetAvailableSlots)
	protected.Post("/interview/book", middleware.RequireRole("student"), interviewReviewHandler.BookInterviewSlot)

	// Interview Booking Management (Admin/Officer)
	protected.Get("/interview/bookings", middleware.RequireRole("scholarship_officer", "admin", "interviewer"), interviewReviewHandler.GetAllBookings)
	protected.Get("/interview/bookings/:id", middleware.JWTMiddleware(cfg), interviewReviewHandler.GetBookingByID)
	protected.Put("/interview/bookings/:id", middleware.RequireRole("scholarship_officer", "admin"), interviewReviewHandler.UpdateBooking)
	protected.Delete("/interview/bookings/:id", middleware.JWTMiddleware(cfg), interviewReviewHandler.CancelBooking)

	// Interview Booking Actions
	protected.Post("/interview/bookings/:id/reschedule", middleware.JWTMiddleware(cfg), interviewReviewHandler.RescheduleBooking)
	protected.Post("/interview/bookings/:id/confirm", middleware.JWTMiddleware(cfg), interviewReviewHandler.ConfirmBooking)
	protected.Post("/interview/bookings/:id/checkin", middleware.RequireRole("scholarship_officer", "admin", "interviewer"), interviewReviewHandler.CheckInBooking)
	protected.Post("/interview/bookings/:id/checkout", middleware.RequireRole("scholarship_officer", "admin", "interviewer"), interviewReviewHandler.CheckOutBooking)

	// Interview Statistics
	protected.Get("/interview/statistics", middleware.RequireRole("scholarship_officer", "admin"), interviewReviewHandler.GetStatistics)
	
	// Add catch-all route for undefined endpoints
	addCatchAllRoute(app)
}

// setupAuthRoutes configures authentication-related routes
func setupAuthRoutes(api fiber.Router, authHandler *handlers.AuthHandler, cfg *config.Config) {
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/refresh", authHandler.RefreshToken) // Refresh token endpoint doesn't need JWT middleware

	// Protected auth routes that require JWT token
	authProtected := auth.Use(middleware.JWTMiddleware(cfg))
	authProtected.Get("/me", authHandler.GetProfile) // Current user profile
	authProtected.Post("/logout", func(c *fiber.Ctx) error {
		// Simple logout endpoint - client handles token removal
		return c.JSON(fiber.Map{"message": "Logged out successfully"})
	})
}

// setupPublicScholarshipRoutes configures public scholarship routes that don't require authentication
func setupPublicScholarshipRoutes(api fiber.Router, scholarshipHandler *handlers.ScholarshipHandler) {
	publicScholarships := api.Group("/scholarships")
	publicScholarships.Get("/", scholarshipHandler.GetScholarships)
	publicScholarships.Get("/available", scholarshipHandler.GetAvailableScholarships)
	publicScholarships.Get("/:id", scholarshipHandler.GetScholarship)
}

// setupPublicNewsRoutes configures public news routes that don't require authentication
func setupPublicNewsRoutes(api fiber.Router, newsHandler *handlers.NewsHandler) {
	publicNews := api.Group("/news")
	publicNews.Get("/", newsHandler.ListNews)   // Public listing
	publicNews.Get("/:id", newsHandler.GetNews) // Public reading
}

// setupProtectedRoutes configures routes that require authentication
func setupProtectedRoutes(protected fiber.Router,
	authHandler *handlers.AuthHandler,
	scholarshipHandler *handlers.ScholarshipHandler,
	applicationHandler *handlers.ApplicationHandler,
	interviewHandler *handlers.InterviewHandler,
	allocationHandler *handlers.AllocationHandler,
	notificationHandler *handlers.NotificationHandler,
	reportHandler *handlers.ReportHandler,
	documentHandler *handlers.DocumentHandler,
	userHandler *handlers.UserHandler,
	studentHandler *handlers.StudentHandler,
	adminHandler *handlers.AdminHandler,
	newsHandler *handlers.NewsHandler,
	profileHandler *handlers.ProfileHandler,
	cfg *config.Config) {

	// Universal profile routes (all roles)
	setupProfileRoutes(protected, profileHandler)

	// User profile routes
	setupUserProfileRoutes(protected, authHandler)

	// Protected scholarship routes
	setupProtectedScholarshipRoutes(protected, scholarshipHandler)

	// Protected news routes
	setupProtectedNewsRoutes(protected, newsHandler)

	// Application routes
	setupApplicationRoutes(protected, applicationHandler, cfg)

	// Interview routes
	setupInterviewRoutes(protected, interviewHandler)

	// Document routes
	setupDocumentRoutes(protected, documentHandler)

	// Allocation routes
	setupAllocationRoutes(protected, allocationHandler)

	// Notification routes
	setupNotificationRoutes(protected, notificationHandler)

	// Student profile routes
	setupStudentRoutes(protected, studentHandler)

	// User management routes
	setupUserManagementRoutes(protected, userHandler)

	// Admin and report routes
	setupAdminRoutes(protected, adminHandler, reportHandler, scholarshipHandler)
	setupReportRoutes(protected, reportHandler)
}

// setupUserProfileRoutes configures user profile management routes
func setupUserProfileRoutes(protected fiber.Router, authHandler *handlers.AuthHandler) {
	user := protected.Group("/user")
	user.Get("/profile", authHandler.GetProfile)
	user.Put("/profile", authHandler.UpdateProfile)
	user.Put("/password", authHandler.ChangePassword)
}

// setupProtectedScholarshipRoutes configures protected scholarship management routes
func setupProtectedScholarshipRoutes(protected fiber.Router, scholarshipHandler *handlers.ScholarshipHandler) {
	// Protected scholarship routes (authenticated users)
	scholarships := protected.Group("/scholarships")

	// Admin/Officer scholarship management routes
	scholarshipAdmin := scholarships.Use(middleware.RequireRole("admin", "scholarship_officer"))
	scholarshipAdmin.Post("/", scholarshipHandler.CreateScholarship)
	scholarshipAdmin.Put("/:id", scholarshipHandler.UpdateScholarship)
	scholarshipAdmin.Delete("/:id", scholarshipHandler.DeleteScholarship)
	scholarshipAdmin.Post("/:id/toggle", scholarshipHandler.ToggleScholarshipStatus)
	scholarshipAdmin.Post("/:id/publish", scholarshipHandler.PublishScholarship)
	scholarshipAdmin.Post("/:id/close", scholarshipHandler.CloseScholarship)
	scholarshipAdmin.Post("/:id/suspend", scholarshipHandler.SuspendScholarship)
	scholarshipAdmin.Post("/:id/duplicate", scholarshipHandler.DuplicateScholarship)

	// Scholarship sources routes
	sources := protected.Group("/scholarship-sources")
	sourcesAdmin := sources.Use(middleware.RequireRole("admin", "scholarship_officer"))
	sourcesAdmin.Post("/", scholarshipHandler.CreateSource)
	sourcesAdmin.Get("/", scholarshipHandler.GetSources)
}

// setupApplicationRoutes configures application management routes
func setupApplicationRoutes(protected fiber.Router, applicationHandler *handlers.ApplicationHandler, cfg *config.Config) {
	applications := protected.Group("/applications")

	// Student application routes
	applications.Post("/", middleware.RequireRole("student"), applicationHandler.CreateApplication)
	applications.Get("/my", middleware.RequireRole("student"), applicationHandler.GetMyApplications)
	applications.Get("/:id", applicationHandler.GetApplication)
	applications.Put("/:id", applicationHandler.UpdateApplication)
	applications.Post("/:id/submit", middleware.RequireRole("student"), applicationHandler.SubmitApplication)
	applications.Delete("/:id", applicationHandler.DeleteApplication)

	// Application Details routes (Student only)
	setupApplicationDetailsRoutes(applications, middleware.RequireRole("student"), cfg)

	// Admin/Officer application management routes
	applicationsAdmin := applications.Use(middleware.RequireRole("admin", "scholarship_officer", "interviewer"))
	applicationsAdmin.Get("/", applicationHandler.GetApplications)
	applicationsAdmin.Post("/:id/review", middleware.RequireRole("admin", "scholarship_officer"), applicationHandler.ReviewApplication)
}

// setupAdminApplicationRoutes configures admin application management routes
func setupAdminApplicationRoutes(protected fiber.Router, applicationHandler *handlers.ApplicationHandler) {
	adminApplications := protected.Group("/admin/applications", middleware.RequireRole("admin", "scholarship_officer"))

	// Admin application endpoints
	adminApplications.Get("/", applicationHandler.GetApplications)
	adminApplications.Get("/stats", applicationHandler.GetApplicationStats)
	adminApplications.Get("/:id", applicationHandler.GetApplication)
	adminApplications.Post("/:id/review", applicationHandler.ReviewApplication)
	adminApplications.Put("/:id/status", applicationHandler.UpdateApplicationStatus)
	adminApplications.Delete("/:id", applicationHandler.DeleteApplication)
}

// setupInterviewRoutes configures interview management routes
func setupInterviewRoutes(protected fiber.Router, interviewHandler *handlers.InterviewHandler) {
	interviews := protected.Group("/interviews")

	// Interview management (officers/admins)
	interviewAdmin := interviews.Use(middleware.RequireRole("admin", "scholarship_officer"))
	interviewAdmin.Post("/schedules", interviewHandler.CreateSchedule)
	interviewAdmin.Get("/schedules", interviewHandler.GetSchedules)

	// Student interview routes
	interviews.Post("/applications/:application_id/schedules/:schedule_id/book",
		middleware.RequireRole("student"), interviewHandler.BookInterview)
	interviews.Post("/appointments/:appointment_id/confirm",
		middleware.RequireRole("student"), interviewHandler.ConfirmInterview)

	// Interviewer routes
	interviews.Post("/appointments/:appointment_id/result",
		middleware.RequireRole("interviewer", "admin", "scholarship_officer"),
		interviewHandler.SubmitInterviewResult)
	interviews.Get("/my", interviewHandler.GetMyInterviews)
}

// setupDocumentRoutes configures document management routes
func setupDocumentRoutes(protected fiber.Router, documentHandler *handlers.DocumentHandler) {
	documents := protected.Group("/documents")
	documents.Post("/applications/:application_id/upload",
		middleware.RequireRole("student"), documentHandler.UploadDocument)
	documents.Get("/applications/:application_id", documentHandler.GetDocuments)
	documents.Get("/:document_id/download", documentHandler.DownloadDocument)
	documents.Delete("/:document_id", middleware.RequireRole("student"), documentHandler.DeleteDocument)
	documents.Get("/types", documentHandler.GetDocumentTypes)
	documents.Get("/stats", documentHandler.GetDocumentStats)

	// Document verification (officers/admins)
	documentAdmin := documents.Use(middleware.RequireRole("admin", "scholarship_officer"))
	documentAdmin.Post("/:document_id/verify", documentHandler.VerifyDocument)
	documentAdmin.Post("/bulk-verify", documentHandler.BulkVerifyDocuments)
}

// setupAllocationRoutes configures allocation management routes
func setupAllocationRoutes(protected fiber.Router, allocationHandler *handlers.AllocationHandler) {
	allocations := protected.Group("/allocations", middleware.RequireRole("admin", "scholarship_officer"))
	allocations.Post("/", allocationHandler.CreateAllocation)
	allocations.Get("/", allocationHandler.GetAllocations)
	allocations.Get("/:id", allocationHandler.GetAllocationDetails)
	allocations.Post("/:id/approve", allocationHandler.ApproveAllocation)
	allocations.Post("/:id/disburse", allocationHandler.DisburseAllocation)
	allocations.Get("/budget/summary", allocationHandler.GetBudgetSummary)
}

// setupNotificationRoutes configures notification management routes
func setupNotificationRoutes(protected fiber.Router, notificationHandler *handlers.NotificationHandler) {
	notifications := protected.Group("/notifications")
	notifications.Get("/", notificationHandler.GetNotifications)
	notifications.Post("/:id/read", notificationHandler.MarkAsRead)
	notifications.Post("/mark-all-read", notificationHandler.MarkAllAsRead)
	notifications.Get("/unread-count", notificationHandler.GetUnreadCount)
	notifications.Delete("/:id", notificationHandler.DeleteNotification)
	notifications.Get("/types", notificationHandler.GetNotificationTypes)

	// Notification management (admins/officers)
	notificationAdmin := notifications.Use(middleware.RequireRole("admin", "scholarship_officer"))
	notificationAdmin.Post("/send", notificationHandler.SendNotification)
	notificationAdmin.Post("/bulk-send", notificationHandler.SendBulkNotification)
}

// setupUserManagementRoutes configures user administration routes
func setupUserManagementRoutes(protected fiber.Router, userHandler *handlers.UserHandler) {
	users := protected.Group("/users", middleware.RequireRole("admin"))
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Post("/", userHandler.CreateUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Post("/:id/roles", userHandler.AssignRole)
	users.Delete("/:id/roles/:role_id", userHandler.RemoveRole)
	users.Post("/:id/deactivate", userHandler.DeactivateUser)
	users.Post("/:id/reactivate", userHandler.ReactivateUser)
	users.Get("/roles", userHandler.GetRoles)
}

// setupStudentRoutes configures student profile routes
func setupStudentRoutes(protected fiber.Router, studentHandler *handlers.StudentHandler) {
	student := protected.Group("/student", middleware.RequireRole("student"))
	student.Get("/profile", studentHandler.GetStudentProfile)
	student.Post("/profile", studentHandler.CreateStudentProfile)
	student.Get("/application-history", studentHandler.GetStudentApplicationHistory)

	// Priority Score Calculator
	student.Post("/calculate-score", studentHandler.CalculatePriorityScore)

	// Eligible Scholarships
	student.Get("/eligible-scholarships", studentHandler.GetEligibleScholarships)
}

// setupAdminRoutes configures administrative routes
func setupAdminRoutes(protected fiber.Router, adminHandler *handlers.AdminHandler, reportHandler *handlers.ReportHandler, scholarshipHandler *handlers.ScholarshipHandler) {
	admin := protected.Group("/admin", middleware.RequireRole("admin"))
	admin.Get("/dashboard", reportHandler.GetDashboardSummary)
	admin.Get("/stats", adminHandler.GetSystemStats)
	admin.Get("/activity-log", adminHandler.GetActivityLog)

	// System Configuration
	admin.Get("/config", adminHandler.GetSystemConfig)
	admin.Put("/config", adminHandler.UpdateSystemConfig)

	// System Testing
	admin.Post("/test-email", adminHandler.TestEmailConnection)
	admin.Post("/test-database", adminHandler.TestDatabaseConnection)

	// Dashboard APIs
	adminDashboard := admin.Group("/dashboard")
	adminDashboard.Get("/stats", adminHandler.GetDashboardStats)
	adminDashboard.Get("/alerts", adminHandler.GetDashboardAlerts)
	adminDashboard.Get("/activities", adminHandler.GetDashboardActivities)
	adminDashboard.Get("/resources", adminHandler.GetDashboardResources)

	// Admin Scholarship Management
	adminScholarships := admin.Group("/scholarships")
	adminScholarships.Get("/", scholarshipHandler.GetScholarships)
	adminScholarships.Get("/stats", scholarshipHandler.GetScholarshipStats)
	adminScholarships.Get("/:id", scholarshipHandler.GetScholarship)
	adminScholarships.Post("/", scholarshipHandler.CreateScholarship)
	adminScholarships.Put("/:id", scholarshipHandler.UpdateScholarship)
	adminScholarships.Delete("/:id", scholarshipHandler.DeleteScholarship)
	adminScholarships.Post("/:id/publish", scholarshipHandler.PublishScholarship)
	adminScholarships.Post("/:id/close", scholarshipHandler.CloseScholarship)
	adminScholarships.Post("/:id/suspend", scholarshipHandler.SuspendScholarship)
	adminScholarships.Post("/:id/duplicate", scholarshipHandler.DuplicateScholarship)

	// Admin Profile Management
	admin.Get("/profile", adminHandler.GetAdminProfile)
	admin.Put("/profile", adminHandler.UpdateAdminProfile)
	admin.Put("/profile/password", adminHandler.ChangeAdminPassword)
}

// setupProtectedNewsRoutes configures protected news management routes
func setupProtectedNewsRoutes(protected fiber.Router, newsHandler *handlers.NewsHandler) {
	// News routes for all authenticated users
	newsUser := protected.Group("/news")
	newsUser.Get("/unread/count", newsHandler.GetUnreadNewsCount)
	newsUser.Post("/:id/read", newsHandler.MarkNewsAsRead)

	// Protected news routes (admin/editor only)
	newsAdmin := protected.Group("/news", middleware.RequireRole("admin", "superadmin"))
	newsAdmin.Post("/", newsHandler.CreateNews)
	newsAdmin.Put("/:id", newsHandler.UpdateNews)
	newsAdmin.Delete("/:id", newsHandler.DeleteNews)
}

// setupReportRoutes configures reporting routes
func setupReportRoutes(protected fiber.Router, reportHandler *handlers.ReportHandler) {
	reports := protected.Group("/reports", middleware.RequireRole("admin", "scholarship_officer"))
	reports.Get("/dashboard", reportHandler.GetDashboardSummary)
	reports.Get("/applications", reportHandler.GetApplicationReport)
	reports.Get("/scholarships", reportHandler.GetScholarshipReport)
	reports.Get("/budget", reportHandler.GetBudgetReport)
	reports.Get("/students", reportHandler.GetStudentReport)
	reports.Get("/export/:type", reportHandler.ExportReport)
}

// setupPaymentRoutes configures payment-related routes
func setupPaymentRoutes(protected fiber.Router, paymentHandler *handlers.PaymentHandler) {
	payments := protected.Group("/payments", middleware.RequireRole("admin", "scholarship_officer"))

	// Payment transactions
	payments.Post("/transactions", paymentHandler.CreateTransaction)
	payments.Get("/transactions/:id", paymentHandler.GetTransaction)
	payments.Get("/allocations/:allocation_id/transactions", paymentHandler.GetAllocationTransactions)
	payments.Put("/transactions/:id/status", paymentHandler.UpdateTransactionStatus)

	// Disbursement schedules
	payments.Post("/disbursements", paymentHandler.CreateDisbursementSchedule)
	payments.Get("/allocations/:allocation_id/disbursements", paymentHandler.GetDisbursementSchedules)
	payments.Get("/disbursements/pending", paymentHandler.GetPendingDisbursements)

	// Payment methods
	payments.Get("/methods", paymentHandler.GetPaymentMethods)
}

// setupAnalyticsRoutes configures analytics-related routes
func setupAnalyticsRoutes(protected fiber.Router, analyticsHandler *handlers.AnalyticsHandler) {
	analytics := protected.Group("/analytics", middleware.RequireRole("admin", "scholarship_officer"))

	// Scholarship statistics
	analytics.Get("/statistics", analyticsHandler.GetScholarshipStatistics)
	analytics.Get("/statistics/all", analyticsHandler.GetAllStatistics)
	analytics.Post("/statistics", analyticsHandler.CreateStatistics)

	// Application analytics
	analytics.Get("/applications/:application_id", analyticsHandler.GetApplicationAnalytics)
	analytics.Post("/applications", analyticsHandler.CreateApplicationAnalytics)
	analytics.Get("/processing-time", analyticsHandler.GetAverageProcessingTime)
	analytics.Get("/bottlenecks", analyticsHandler.GetBottleneckSteps)
	analytics.Get("/dashboard", analyticsHandler.GetDashboardSummary)
}

// setupApplicationDetailsRoutes configures detailed application form routes
func setupApplicationDetailsRoutes(applications fiber.Router, roleMiddleware fiber.Handler, cfg *config.Config) {
	// Initialize application details handler
	appDetailsHandler := handlers.NewApplicationDetailsHandler(cfg)

	// Application details routes (requires student role)
	details := applications.Group("/:id", roleMiddleware)

	// Personal Information
	details.Post("/personal-info", appDetailsHandler.SavePersonalInfo)

	// Addresses
	details.Post("/addresses", appDetailsHandler.SaveAddresses)

	// Education History
	details.Post("/education", appDetailsHandler.SaveEducation)

	// Family Information (includes family members, guardians, siblings, living situation)
	details.Post("/family", appDetailsHandler.SaveFamily)

	// Financial Information (includes financial info, assets, scholarship history, health info, funding needs)
	details.Post("/financial", appDetailsHandler.SaveFinancial)

	// Activities and References
	details.Post("/activities", appDetailsHandler.SaveActivities)

	// Complete Form - Save all sections at once
	details.Post("/complete-form", appDetailsHandler.SaveCompleteForm)
	details.Get("/complete-form", appDetailsHandler.GetCompleteForm)

	// Submit Application
	details.Put("/submit", appDetailsHandler.SubmitApplication)
}

// Add catch-all route for 404 handling at the end of SetupRoutes function
func addCatchAllRoute(app *fiber.App) {
	// Catch-all route for undefined endpoints
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
			"path":  c.Path(),
			"method": c.Method(),
		})
	})
}

// setupProfileRoutes configures universal profile routes (all roles)
func setupProfileRoutes(protected fiber.Router, profileHandler *handlers.ProfileHandler) {
	// Universal profile routes - accessible by all authenticated users
	protected.Get("/profile", profileHandler.GetProfile)
	protected.Put("/profile", profileHandler.UpdateProfile)
	protected.Put("/profile/password", profileHandler.ChangePassword)
}
