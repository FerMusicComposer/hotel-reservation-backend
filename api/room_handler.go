package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	roomStore db.RoomStore
	// bookingStore db.BookingStore
}

type BookingParams struct {
	FromDate  time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
	NumPeople int       `json:"numPeople"`
}

func NewRoomHandler(roomStore db.RoomStore) *RoomHandler {
	return &RoomHandler{roomStore: roomStore}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.roomStore.GetRooms(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
func (h *RoomHandler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelId": oid}
	rooms, err := h.roomStore.GetRoomsByHotelID(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleGetRoomByID(c *fiber.Ctx) error {
	id := c.Params("id")

	room, err := h.roomStore.GetRoomByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(room)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookingParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*models.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(AuthResponse{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
		})
	}

	booking := models.Booking{
		UserID:    user.ID,
		RoomID:    roomId,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
		NumPeople: params.NumPeople,
	}

	fmt.Printf("%+v\n", booking)

	return nil
}
