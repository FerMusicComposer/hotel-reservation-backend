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

func (bookingHandler *BookingHandler) HandleGetAllBookings(ctx *fiber.Ctx) error {
	bookings, err := bookingHandler.bookingStore.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}

	return ctx.JSON(bookings)
}

func (bookingHandler *BookingHandler) HandleGetAllBookingsWithinDateRange(ctx *fiber.Ctx) error {
	var params BookingQueryParams
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}

	where := bson.M{
		"fromDate": bson.M{"$gte": params.FromDate},
		"toDate":   bson.M{"$lte": params.ToDate},
	}

	bookings, err := bookingHandler.bookingStore.GetBookings(ctx.Context(), where)
	if err != nil {
		return err
	}

	return ctx.JSON(bookings)
}

// -----------
// USER ROUTES
// -----------
func (bookingHandler *BookingHandler) HandleGetUserBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	booking, err := bookingHandler.bookingStore.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	err = isUserToken(ctx, booking.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(booking)
}

func (bookingHandler *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	booking, err := bookingHandler.bookingStore.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	err = isUserToken(ctx, booking.UserID)
	if err != nil {
		return err
	}

	update := bson.M{
		"isCancelled": true,
	}

	err = bookingHandler.bookingStore.UpdateBooking(ctx.Context(), id, update)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "booking cancelled",
	})
}
