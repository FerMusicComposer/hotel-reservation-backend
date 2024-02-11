package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelData struct {
	name     string
	location string
	rating   float64
}

func main() {
	hotels := []hotelData{
		{
			name:     "Hilton Rome",
			location: "Rome",
			rating:   4.5,
		},
		{
			name:     "Melia London",
			location: "London",
			rating:   4.0,
		},
		{
			name:     "Marriott Paris",
			location: "Paris",
			rating:   3.5,
		},
		{
			name:     "Palacio Real",
			location: "Madrid",
			rating:   5.0,
		},
	}
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

	for _, hotel := range hotels {
		seedHotel(hotel.name, hotel.location, hotel.rating, hotelStore, roomStore, ctx)
		fmt.Println("------------------------------")
	}

	conn.Database.Client().Disconnect(ctx)

}

func seedHotel(name, location string, rating float64, hotelStore db.HotelStore, roomStore db.RoomStore, ctx context.Context) {
	hotel := models.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []models.Room{
		{
			Size:    "single",
			Seaside: true,
			Price:   99.9,
		},
		{
			Size:    "double",
			Seaside: false,
			Price:   199.9,
		},
		{
			Size:    "king",
			Seaside: false,
			Price:   299.9,
		},
		{
			Size:    "king deluxe",
			Seaside: true,
			Price:   399.9,
		},
	}

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
