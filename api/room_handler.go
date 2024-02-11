package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	roomStore db.RoomStore
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
