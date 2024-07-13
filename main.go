package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/api"
)

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Working and Running"})
}
func main() {
	portNumber := flag.String("portNumber", ":5000", "The Port Number that server listen to")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*portNumber)
}
