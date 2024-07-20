package api

import (
	"fmt"
	"github.com/mohamedramadan14/hotel-reservation/db/fixtures"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	user := fixtures.AddUser(db.Store, "mo", "gmail", false)
	hotel := fixtures.AddHotel(db.Store, "bar Hotel", "New York", 4.0, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 261.5, hotel.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 3)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)

	fmt.Println(booking)
}
