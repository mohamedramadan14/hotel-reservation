package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"userID,omitempty"`
	RoomID     primitive.ObjectID `bson:"room_id,omitempty" json:"roomID,omitempty"`
	NumPersons int                `bson:"num_persons" json:"numPersons"`
	FromDate   time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
}
