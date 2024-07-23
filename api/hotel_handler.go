package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	resp, err := h.store.Room.GetRooms(c.Context(), bson.M{"hotelID": oid})
	if err != nil {
		return NewError(http.StatusNotFound, "Room not found")
	}
	return c.JSON(resp)
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return NewError(http.StatusBadRequest, "Bad Request")
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &params.Pagination)
	if err != nil {
		return NewError(http.StatusNotFound, "hotel not found")
	}
	resp := ResourceResp{Data: hotels, Results: len(hotels), Page: int(params.Page)}

	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, "hotel not found")
	}
	return c.JSON(hotel)
}
