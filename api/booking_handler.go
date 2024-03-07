package api

import (
	"net/http"

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
		return Internal.from("HandleGetAllBookings => error obtaining bookings", err).Err
	}

	return ctx.JSON(bookings)
}

func (bookingHandler *BookingHandler) HandleGetAllBookingsWithinDateRange(ctx *fiber.Ctx) error {
	var params BookingQueryParams
	if err := ctx.QueryParser(&params); err != nil {
		return Internal.from("HandleGetAllBookingsWithinDateRange => error parsing query params", err).Err
	}

	where := bson.M{
		"fromDate": bson.M{"$gte": params.FromDate},
		"toDate":   bson.M{"$lte": params.ToDate},
	}

	bookings, err := bookingHandler.bookingStore.GetBookings(ctx.Context(), where)
	if err != nil {
		return Internal.from("HandleGetAllBookingsWithinDateRange => error obtaining bookings", err).Err
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
		return InvalidID.from("HandleGetUserBooking => error obtaining booking by ID", err).Err
	}

	err = isUserToken(ctx, booking.UserID)
	if err != nil {
		return Internal.with(http.StatusUnauthorized).Err
	}

	return ctx.JSON(booking)
}

func (bookingHandler *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	booking, err := bookingHandler.bookingStore.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return InvalidID.from("HandleCancelBooking => error obtaining booking by ID", err).Err
	}

	err = isUserToken(ctx, booking.UserID)
	if err != nil {
		return Internal.with(http.StatusUnauthorized).Err
	}

	updateBooking := bson.M{
		"isCancelled": true,
	}

	err = bookingHandler.bookingStore.UpdateBooking(ctx.Context(), id, updateBooking)
	if err != nil {
		return Internal.from("HandleCancelBooking => error updating booking", err).Err
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
		return Internal.from("HandleCancelBooking => error updating room", err).Err
	}

	return ctx.JSON(fiber.Map{
		"message": "booking cancelled",
	})
}
