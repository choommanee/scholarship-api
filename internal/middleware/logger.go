package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format: "${time} ${status} - ${method} ${path} ${latency}\n",
	})
}