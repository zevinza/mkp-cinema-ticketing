package account

import (
	"api/app/lib"

	"github.com/gofiber/fiber/v2"
)

func PostAccountLogout(c *fiber.Ctx) error {
	return lib.OK(c)
}
