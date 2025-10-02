package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	// Get frontend URL from environment
	frontendURL := os.Getenv("FRONTEND_URL")

	// Build allowed origins list
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
	}

	// Add frontend URL if provided
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	// Add Railway app domains (both production and preview)
	allowedOrigins = append(allowedOrigins, "https://*.railway.app")

	return cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ","),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: true,
	})
}