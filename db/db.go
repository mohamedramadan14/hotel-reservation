package db

import (
	"github.com/joho/godotenv"
	"log"
)

const MongoDBNameKey = "MONGO_DB_NAME"

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
