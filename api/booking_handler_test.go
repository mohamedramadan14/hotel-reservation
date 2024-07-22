package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/api/middleware"
	"github.com/mohamedramadan14/hotel-reservation/db/fixtures"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "maged", "gmail", true)
		user           = fixtures.AddUser(db.Store, "mo", "gmail", false)
		hotel          = fixtures.AddHotel(db.Store, "bar Hotel", "New York", 4.0, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 261.5, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.Store.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-token", GenerateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 Response Got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("Expected one booking Got %d", len(bookings))
	}

	haveBooking := bookings[0]

	if haveBooking.ID != booking.ID {
		t.Fatalf("Expecting booking with ID: %s but Got %s", booking.ID, haveBooking.ID)
	}
	if haveBooking.UserID != booking.UserID {
		t.Fatalf("Expecting booking with UserID: %s but Got %s", booking.UserID, haveBooking.UserID)
	}

	// Test non-admin cannot access the bookings

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", GenerateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 Got %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		userNotOwnedBooking = fixtures.AddUser(db.Store, "koko", "gmail", false)
		user                = fixtures.AddUser(db.Store, "mo", "gmail", false)
		hotel               = fixtures.AddHotel(db.Store, "bar Hotel", "New York", 4.0, nil)
		room                = fixtures.AddRoom(db.Store, "small", true, 261.5, hotel.ID)
		from                = time.Now()
		till                = time.Now().AddDate(0, 0, 3)
		booking             = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app                 = fiber.New()
		route               = app.Group("/", middleware.JWTAuthentication(db.Store.User))
		bookingHandler      = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("x-api-token", GenerateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	var bookingResult *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResult); err != nil {
		fmt.Println(err)
	}

	if bookingResult.ID != booking.ID {
		t.Fatalf("Expecting booking with ID: %s but Got %s", booking.ID, bookingResult.ID)
	}
	if bookingResult.UserID != booking.UserID {
		t.Fatalf("Expecting booking with UserID: %s but Got %s", booking.UserID, bookingResult.UserID)
	}

	// test non-owned user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", GenerateTokenFromUser(userNotOwnedBooking))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 Got %d", resp.StatusCode)
	}
}
