package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Mohamed",
		LastName:  "Ramadan",
		ID:        "kf85qr-fksjs5",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON("James" + id)
}
