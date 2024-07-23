package api

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

//func init() {
//	if err := godotenv.Load("../.env"); err != nil {
//		log.Fatal(err)
//	}
//}

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := godotenv.Load("..\\.env"); err != nil {
		log.Fatal(err)
	}
	dbName := os.Getenv(db.MongoDBNameKey)
	if err := tdb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(_ *testing.T) *testdb {
	if err := godotenv.Load("..\\.env"); err != nil {
		log.Fatal(err)
	}
	dbUri := os.Getenv("MONGO_DB_URL")
	dbName := os.Getenv(db.MongoDBNameKey)
	fmt.Println("-->", dbUri)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(dbName).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			Hotel:   hotelStore,
		},
	}
}
