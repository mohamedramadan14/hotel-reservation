package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}

	return c.Next()
}
