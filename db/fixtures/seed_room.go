package fixtures

import (
	"context"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func seedRoom(size string, seaside bool, price float64, maxCapacity int, hotelID primitive.ObjectID, roomStore db.RoomStore, ctx context.Context) primitive.ObjectID {

	room := models.Room{
		HotelId:     hotelID,
		Size:        size,
		Seaside:     seaside,
		Price:       price,
		MaxCapacity: maxCapacity,
		Status:      []models.RoomStatus{},
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom.ID

}
