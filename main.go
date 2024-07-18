package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/api"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	portNumber := flag.String("portNumber", ":5000", "The Port Number that server listen to")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore   = db.NewMongoUserStore(client)
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		userHandler = api.NewUserHandler(userStore)
		store       = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		app          = fiber.New(config)
		apiV1        = app.Group("/api/v1")
	)

	// User
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	// Hotel
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	err = app.Listen(*portNumber)
	if err != nil {
		log.Fatal(err)
	}
}
