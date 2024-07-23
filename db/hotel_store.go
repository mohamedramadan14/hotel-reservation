package db

import (
	"context"
	"fmt"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type HotelStore interface {
	Dropper
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, interface{}, interface{}) error
	GetHotels(context.Context, interface{}, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbName := os.Getenv(MongoDBNameKey)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbName).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping Hotel Collection ---")
	return s.coll.Drop(ctx)
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter interface{}, pag *Pagination) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit).SetLimit(pag.Limit)

	cur, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}
func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter interface{}, values interface{}) error {
	_, err := s.coll.UpdateOne(ctx, filter, values)
	return err
}
