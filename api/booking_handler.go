package api

import (
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	bookingStore db.BookingStore
}

type BookingQueryParams struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

func NewBookingHandler(bookingStore db.BookingStore) *BookingHandler {
	return &BookingHandler{bookingStore: bookingStore}
}

// -----------------
// ADMIN ONLY ROUTES
// -----------------

func (h *BookingHandler) HandleGetAllBookings(c *fiber.Ctx) error {
	bookings, err := h.bookingStore.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBookinByID(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.bookingStore.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleGetAllBookingsWithinDateRange(c *fiber.Ctx) error {
	var params BookingQueryParams
	if err := c.QueryParser(&params); err != nil {
		return err
	}

	where := bson.M{
		"fromDate": bson.M{"$gte": params.FromDate},
		"toDate":   bson.M{"$lte": params.ToDate},
	}

	bookings, err := h.bookingStore.GetBookings(c.Context(), where)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// -----------
// USER ROUTES
// -----------
func (h *BookingHandler) HandleGetUserBooking(c *fiber.Ctx) error {
	return nil
}
