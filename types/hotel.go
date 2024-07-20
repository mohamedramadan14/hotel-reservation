package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string               `json:"name" bson:"name"`
	Location string               `json:"location" bson:"location"`
	Rooms    []primitive.ObjectID `json:"rooms" bson:"rooms"`
	Rating   float64              `json:"rating" bson:"rating"`
}

type Room struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SeeSide bool               `json:"seeSide" bson:"seeSide"`
	Size    string             `json:"size" bson:"size"`
	Price   float64            `json:"price" bson:"price"`
	HotelID primitive.ObjectID `json:"hotelID" bson:"hotelID"`
}
