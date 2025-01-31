package main

import (
	"context"
	"fmt"
	faker2 "github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"github.com/mohamedramadan14/hotel-reservation/api"
	"github.com/mohamedramadan14/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"

	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}
}
func main() {
	fmt.Println("Seeding database...")

	var (
		faker         = faker2.New()
		ctx           = context.Background()
		mongoEndPoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv(db.MongoDBNameKey)
	)
	fmt.Println("--->", mongoDBName)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndPoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
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

	fmt.Println("Creating 100 Dummy Hotels")

	for i := 0; i < 100; i++ {
		fixtures.AddHotel(store, faker.Company().Name(), faker.Address().Address(), faker.Float64(1, 1, 5), nil)
	}
	fmt.Println("Seeding completed...")
}
