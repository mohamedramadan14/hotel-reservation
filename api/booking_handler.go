package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// HandleGetBookings This need to be admin authorized
func (h *BookingHandler) HandleGetBookings(ctx *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return NewError(http.StatusNotFound, "no bookings")
	}
	return ctx.JSON(&bookings)
}

// HandleGetBooking This need to be user authorized
func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, fmt.Sprintf("No booking with this Id: %s", id))
	}

	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return ErrUnAuthorized()
	}
	return ctx.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, fmt.Sprintf("No booking with this Id: %s", id))
	}

	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return ErrUnAuthorized()
	}
	if err := h.store.Booking.UpdateBooking(ctx.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}
	return ctx.JSON(genericResp{
		Type: "msg",
		Msg:  "Booking has been canceled successfully",
	})
}
