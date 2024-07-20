package main

import (
	"context"
	"fmt"
	faker2 "github.com/jaswdr/faker"
	"github.com/mohamedramadan14/hotel-reservation/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"time"

	"github.com/mohamedramadan14/hotel-reservation/db"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          context.Context
	err          error
	faker        = faker2.New()
	count        = 5
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
	bookingStore = db.NewMongoBookingStore(client)
}

func seedUser(fName string, lName string, email string, password string, isAdmin bool) *types.User {
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

	insertedUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s", user.Email, api.GenerateTokenFromUser(user))
	return insertedUser

}
func seedHotel(name string, location string, rating float64) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		SeaSide: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(roomID primitive.ObjectID, userID primitive.ObjectID, from time.Time, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	booking, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Booking: ", booking.ID)
}

func main() {

	mo := seedUser("mohamed", "ramadan", "mo@gmail.com", "testPass", false)
	seedUser("Maged", "ramadan", "maged@gmail.com", "admin", true)

	for i := 0; i < count; i++ {
		h := seedHotel(faker.Company().Name(), faker.Address().City(), float64(rand.Intn(51))/10.0)
		seedRoom("small", true, 89.99, h.ID)
		seedRoom("normal", false, 129.99, h.ID)
		room := seedRoom("King", true, 389.99, h.ID)
		seedBooking(room.ID, mo.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	}
	fmt.Println("Seeding completed...")
}
