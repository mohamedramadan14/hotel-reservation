package main

import (
	"context"
	"flag"
	"github.com/mohamedramadan14/hotel-reservation/api/middleware"
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
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		roomHandler    = api.NewRoomHandler(store)
		hotelHandler   = api.NewHotelHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// Auth
	auth.Post("/auth/login", authHandler.HandleAuthentication)

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

	// Room
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Booking
	// TODO: Cancel A Booking
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Put("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// Admin Routes
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	err = app.Listen(*portNumber)
	if err != nil {
		log.Fatal(err)
	}
}
