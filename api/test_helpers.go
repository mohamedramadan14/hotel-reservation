package api

import (
	"context"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

const (
	testDBUri  = "mongodb://localhost:27017"
	testDBName = "hotel-reservation-test"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.client.Database(testDBName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(_ *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDBUri))
	if err != nil {
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
		}}
}
