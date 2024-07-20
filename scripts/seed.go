package main

import (
	"context"
	"fmt"
	faker2 "github.com/jaswdr/faker"
	"github.com/mohamedramadan14/hotel-reservation/api"
	"github.com/mohamedramadan14/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fmt.Println("Seeding database...")

	faker := faker2.New()
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "doe", false)
	fmt.Println("user ---> ", user)
	fmt.Println("user token ---> ", api.GenerateTokenFromUser(user))
	admin := fixtures.AddUser(store, "mo", "ramadan", true)
	fmt.Println("admin ---> ", admin)
	fmt.Println("admin token ---> ", api.GenerateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, faker.Company().Name(), faker.Address().Address(), 3.5, nil)
	fmt.Println("hotel ---> ", hotel)
	room := fixtures.AddRoom(store, "large", true, 88.4, hotel.ID)
	fmt.Println("room ---> ", room)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ---> ", booking)

	fmt.Println("Seeding completed...")
}
