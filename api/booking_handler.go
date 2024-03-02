package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	bookingStore db.BookingStore
	roomStore    db.RoomStore
}

func NewBookingHandler(bookingStore db.BookingStore, roomStore db.RoomStore) *BookingHandler {
	return &BookingHandler{bookingStore: bookingStore, roomStore: roomStore}
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

	updateBooking := bson.M{
		"isCancelled": true,
	}

	err = bookingHandler.bookingStore.UpdateBooking(ctx.Context(), id, updateBooking)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":              booking.RoomID,
		"status.bookingId": booking.ID,
	}

	updateRoom := bson.M{
		"$set": bson.M{
			"status.$.status": "cancelled",
		},
	}

	err = bookingHandler.roomStore.UpdateRoom(ctx.Context(), filter, updateRoom)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "booking cancelled",
	})
}
