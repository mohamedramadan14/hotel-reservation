package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrUnAuthorized()
	}
	if !user.IsAdmin {
		return ErrUnAuthorized()
	}

	return c.Next()
}
