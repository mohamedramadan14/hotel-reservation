package api

import (
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
		return err
	}
	return ctx.JSON(&bookings)
}

// HandleGetBooking This need to be user authorized
func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	return ctx.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	if err := h.store.Booking.UpdateBooking(ctx.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}
	return ctx.JSON(genericResp{
		Type: "msg",
		Msg:  "Booking has been canceled successfully",
	})
}
