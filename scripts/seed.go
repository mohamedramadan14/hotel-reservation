package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mohamedramadan14/hotel-reservation/db"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	fmt.Println("Seeding database...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Grand Hotel",
		Location: "New York",
	}

	roomA := types.Room{
		Type:      types.Single,
		BasePrice: 100,
	}

	roomB := types.Room{
		Type:      types.SeaSide,
		BasePrice: 500,
	}

	_ = roomB

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	roomA.HotelID = insertedHotel.ID

	insertedRoom, err := roomStore.InsertRoom(ctx, &roomA)
	if err != nil {
		log.Fatal(err)
	}

	roomA.ID = insertedRoom.ID

	fmt.Println(insertedHotel)
	fmt.Println(insertedRoom)
}
