package main

import (
	"context"
	"fmt"
	faker2 "github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"

	"github.com/mohamedramadan14/hotel-reservation/db"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	ctx        context.Context
	err        error
	faker      = faker2.New()
	count      = 5
)

func init() {
	ctx := context.Background()
	fmt.Println("Seeding database...")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	userStore = db.NewMongoUserStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}

func seedUser(fName string, lName string, email string, password string, isAdmin bool) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fName,
		LastName:  lName,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	_, err = userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}
func seedHotel(name string, location string, rating float64) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:    "Single",
			SeeSide: false,
			Price:   100,
		},
		{
			Size:    "Double",
			SeeSide: false,
			Price:   200,
		},
		{
			Size:    "Single",
			SeeSide: true,
			Price:   500,
		},
		{
			Size:    "Deluxe",
			Price:   500,
			SeeSide: false,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		room.ID = insertedRoom.ID
	}

}
func main() {

	for i := 0; i < count; i++ {

		seedHotel(faker.Company().Name(), faker.Address().City(), float64(rand.Intn(51))/10.0)
	}
	seedUser("mohamed", "ramadan", "mo@gmail.com", "testPass", false)
	seedUser("Maged", "ramadan", "maged@gmail.com", "admin", true)
	fmt.Println("Seeding completed...")
}
