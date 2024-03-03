package api

import (
	"fmt"
	"net/http"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	roomStore    db.RoomStore
	bookingStore db.BookingStore
	hotelStore   db.HotelStore
}

func NewRoomHandler(roomStore db.RoomStore, hotelStore db.HotelStore, bookingStore db.BookingStore) *RoomHandler {
	return &RoomHandler{
		roomStore:    roomStore,
		hotelStore:   hotelStore,
		bookingStore: bookingStore,
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
		fmt.Println("error validating booking params: ", err)
		return err
	}

	booking := &models.Booking{
		UserID:    user.ID,
		RoomID:    roomId,
		Checkin:   params.Checkin,
		Checkout:  params.Checkout,
		NumPeople: params.NumPeople,
	}

	bkng, err := roomHandler.bookingStore.InsertBooking(ctx.Context(), booking)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": roomId}
	update := bson.M{"$push": bson.M{
		"status": models.RoomStatus{
			Status:     "booked",
			BookingID:  bkng.ID,
			BookedTo:   bkng.Checkout,
			BookedFrom: bkng.Checkin,
		},
	}}
	if err = roomHandler.roomStore.UpdateRoom(ctx.Context(), filter, update); err != nil {
		fmt.Println("error updating room: ", err)
		return err
	}

	return ctx.JSON(bkng)
}
