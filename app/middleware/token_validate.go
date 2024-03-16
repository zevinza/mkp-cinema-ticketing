package middleware

import (
	"api/app/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// TokenValidator middleware
func TokenValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(viper.GetString("HEADER_TOKEN_KEY")) // Get the value of the "x-Token" header
		// Validate the token
		if token != viper.GetString("VALUE_TOKEN_KEY") {
			return lib.ErrorUnauthorized(c, "Wrong x-Token header")
		}
		// Token is valid, continue to the next middleware or route handler
		return c.Next()
	}
}
