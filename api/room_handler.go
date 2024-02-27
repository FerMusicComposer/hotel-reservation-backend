package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomHandler struct {
	roomStore    db.RoomStore
	bookingStore db.BookingStore
	hotelStore   db.HotelStore
}

type RoomParams struct {
	HotelID string  `json:"hotelId"`
	Size    string  `json:"size"`
	Seaside bool    `json:"seaside"`
	Price   float64 `json:"price"`
	MaxCap  int     `json:"maxCapacity"`
}

type BookingParams struct {
	FromDate    time.Time `json:"fromDate"`
	ToDate      time.Time `json:"toDate"`
	NumPeople   int       `json:"numPeople"`
	IsCancelled bool      `json:"isCancelled"`
}

func (params RoomParams) validateRoomParams(ctx *fiber.Ctx, hotelStore db.HotelStore) error {
	_, err := hotelStore.GetHotelByID(ctx.Context(), params.HotelID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid hotel id")
		}

		return fmt.Errorf("internal server error")
	}

	if params.Size == "" {
		return fmt.Errorf("must specify a room size")
	}

	if params.Price <= 0 {
		return fmt.Errorf("price cannot be negative or zero")
	}

	if params.MaxCap <= 0 {
		return fmt.Errorf("max capacity cannot be negative or zero")
	}

	return nil
}

func (params BookingParams) validateBookingParams(ctx *fiber.Ctx, roomstore db.RoomStore, roomID primitive.ObjectID) error {
	now := time.Now()
	room, err := roomstore.GetRoomByID(ctx.Context(), roomID.Hex())
	if err != nil {
		return err
	}

	if params.FromDate.Before(now) {
		return fmt.Errorf("cannot book a room in the past")
	}

	if params.FromDate.After(params.ToDate) {
		return fmt.Errorf("from date cannot be superior to end date")
	}

	if params.NumPeople <= 0 {
		return fmt.Errorf("invalid number of people")
	}

	if room.MaxCapacity < params.NumPeople {
		return fmt.Errorf("room capacity exedeced")
	}

	for _, status := range room.Status {
		if datesAreWithinRange(params.FromDate, params.ToDate, status.BookedFrom, status.BookedTo) &&
			!params.IsCancelled {
			return fmt.Errorf("room is already booked")
		}
	}

	return nil
}

func datesAreWithinRange(fromDate time.Time, toDate time.Time, bookedFrom time.Time, bookedTo time.Time) bool {
	if fromDate.Equal(bookedFrom) ||
		toDate.Equal(bookedTo) ||
		fromDate.After(bookedFrom) && toDate.Before(bookedTo) ||
		fromDate.Before(bookedFrom) && toDate.After(bookedTo) ||
		fromDate.Before(bookedTo) && (toDate.Equal(bookedTo) || toDate.After(bookedTo)) ||
		fromDate.Before(bookedFrom) && (toDate.Equal(bookedFrom) || (toDate.After(bookedFrom) && toDate.Before(bookedTo))) {

		return true
	}

	return false
}

func NewRoomHandler(roomStore db.RoomStore, bookingStore db.BookingStore, hotelStore db.HotelStore) *RoomHandler {
	return &RoomHandler{
		roomStore:    roomStore,
		bookingStore: bookingStore,
		hotelStore:   hotelStore,
	}
}

func (roomHandler *RoomHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := roomHandler.roomStore.GetRooms(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(rooms)
}
func (roomHandler *RoomHandler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelId": oid}
	rooms, err := roomHandler.roomStore.GetRoomsByHotelID(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (roomHandler *RoomHandler) HandleGetRoomByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	room, err := roomHandler.roomStore.GetRoomByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (roomHandler *RoomHandler) HandlePostRoom(ctx *fiber.Ctx) error {
	var params RoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validateRoomParams(ctx, roomHandler.hotelStore); err != nil {
		return err
	}

	hotelOID, err := primitive.ObjectIDFromHex(params.HotelID)
	if err != nil {
		return err
	}

	room, err := roomHandler.roomStore.InsertRoom(ctx.Context(), &models.Room{
		HotelId:     hotelOID,
		Size:        params.Size,
		Seaside:     params.Seaside,
		Price:       params.Price,
		MaxCapacity: params.MaxCap,
		Status: []models.RoomStatus{
			{
				Status: "available",
			},
		},
	})
	if err != nil {
		return err
	}

	filter := bson.M{"_id": room.HotelId}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err = roomHandler.hotelStore.UpdateHotel(ctx.Context(), filter, update); err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (roomHandler *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	var params BookingParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	roomId, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}

	user, ok := ctx.Context().Value("user").(*models.User)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(AuthResponse{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
		})
	}

	if err := params.validateBookingParams(ctx, roomHandler.roomStore, roomId); err != nil {
		return err
	}

	booking := models.Booking{
		UserID:    user.ID,
		RoomID:    roomId,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
		NumPeople: params.NumPeople,
	}

	fmt.Printf("%+v\n", booking)

	bkng, err := roomHandler.bookingStore.InsertBooking(ctx.Context(), &booking)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": roomId}
	update := bson.M{"$push": bson.M{
		"status": models.RoomStatus{
			Status:     "booked",
			BookingID:  bkng.ID,
			BookedTo:   bkng.ToDate,
			BookedFrom: bkng.FromDate,
		},
	}}
	if err = roomHandler.roomStore.UpdateRoom(ctx.Context(), filter, update); err != nil {
		return err
	}

	return ctx.JSON(bkng)
}
