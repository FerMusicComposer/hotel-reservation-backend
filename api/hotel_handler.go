package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore db.HotelStore
}

type HotelQueryParams struct {
	Rooms  bool
	Rating float64
}

func NewHotelHandler(hotelStore db.HotelStore) *HotelHandler {
	return &HotelHandler{hotelStore: hotelStore}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	// var q HotelQueryParams
	// if err := c.QueryParser(q); err != nil {
	// 	return err
	// }

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
