package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	bookingStore db.BookingStore
}

func NewBookingHandler(bookingStore db.BookingStore) *BookingHandler {
	return &BookingHandler{bookingStore: bookingStore}
}

// this is an Admin only route
func (h *BookingHandler) HandleGetAllBookings(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetUserBooking(c *fiber.Ctx) error {
	return nil
}
