package authentication

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JwtProtect(jwtSecret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret), 
		ContextKey: "user",                 
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println("JWT validation failed:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized or invalid token",
			})
		},
	})
}
