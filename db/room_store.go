package db

import (
	"context"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, interface{}) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, store HotelStore) *MongoRoomStore {
	dbName := os.Getenv("MONGO_DB_NAME")
	return &MongoRoomStore{
		client: client,
		coll:   client.Database(dbName).Collection("rooms"),

		HotelStore: store,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter interface{}) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)
	// Update the hotel ID in the room
	filter := bson.M{"_id": room.HotelID}
	values := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.UpdateHotel(ctx, filter, values); err != nil {
		return nil, err
	}
	return room, nil
}
