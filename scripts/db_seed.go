package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	hotel = models.Hotel{
		Name:     "Hilton",
		Location: "London",
		Rooms:    []primitive.ObjectID{},
	}

	rooms = []models.Room{
		{
			Type:      models.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      models.DoubleRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      models.SeaSideRoomType,
			BasePrice: 299.9,
		},
		{
			Type:      models.DeluxeRoomType,
			BasePrice: 399.9,
		},
	}
)

func main() {
	ctx := context.Background()

	conn, err := db.NewMongoConnection(db.DBURI, db.DBNAME)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Database.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(conn)
	roomStore := db.NewMongoRoomStore(conn, hotelStore)

	fmt.Println("seeding database...")

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted hotel: ", insertedHotel)

	for _, room := range rooms {
		room.HotelId = insertedHotel.ID

		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		_ = append(insertedHotel.Rooms, insertedRoom.ID)

		fmt.Println("inserted room: ", insertedRoom.ID)
	}

}
