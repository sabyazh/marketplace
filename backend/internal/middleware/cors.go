package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns a Fiber middleware that configures Cross-Origin Resource
// Sharing headers. It allows all origins, common HTTP methods, and the
// standard headers needed for authenticated API requests.
// Note: AllowCredentials is false because JWT auth uses the Authorization
// header (not cookies), and wildcard origins cannot be used with credentials.
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length,Content-Range",
	})
}
