package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore db.HotelStore
}

func NewHotelHandler(hotelStore db.HotelStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return Internal.from("HandleGetHotels => error obtaining hotels", err).Err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return InvalidID.from("HandleGetHotelById => error obtaining hotel by ID", err).Err
	}

	return c.JSON(hotel)
}
