package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func seedHotel(name, location string, rating float64, hotelStore db.HotelStore, roomStore db.RoomStore, ctx context.Context) {
	hotel := models.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted hotel: ", insertedHotel)

	for _, room := range rooms {
		insertedRoomID := seedRoom(room.Size, room.Seaside, room.Price, room.MaxCapacity, insertedHotel.ID, roomStore, ctx)

		filter := bson.M{"_id": insertedHotel.ID}
		update := bson.M{"$push": bson.M{"rooms": insertedRoomID}}
		hotelStore.UpdateHotel(ctx, filter, update)

		fmt.Println("Added to hotel room: ", insertedRoomID)
		fmt.Println("------------------------------")
	}
}
