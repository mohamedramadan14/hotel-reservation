package db

import (
	"context"
	"fmt"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, m interface{}) ([]*types.Booking, error)
	GetBookingByID(ctx context.Context, id string) (*types.Booking, error)
	UpdateBooking(ctx context.Context, id string, update interface{}) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbName := os.Getenv(MongoDBNameKey)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter interface{}) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	fmt.Println(bookings)
	return bookings, nil
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	bookingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking

	if err := s.coll.FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update = bson.M{
		"$set": update,
	}
	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}
