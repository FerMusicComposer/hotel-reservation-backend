package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelData struct {
	name     string
	location string
	rating   float64
}
type userData struct {
	fname    string
	lname    string
	email    string
	password string
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

	users := []userData{
		{
			fname:    "John",
			lname:    "Doe",
			email:    "jdoe@me.com",
			password: "password1123456789",
		},
		{
			fname:    "Jane",
			lname:    "Doe",
			email:    "jane@me.com",
			password: "password1123456789",
		},
		{
			fname:    "Mike",
			lname:    "Miller",
			email:    "miller@me.com",
			password: "password1123456789",
		},
		{
			fname:    "Sarah",
			lname:    "Smith",
			email:    "smith@me.com",
			password: "password1123456789",
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
	roomStore := db.NewMongoRoomStore(conn)
	userStore := db.NewMongoUserStore(conn)

	fmt.Println("seeding database...")

	for _, hotel := range hotels {
		seedHotel(hotel.name, hotel.location, hotel.rating, hotelStore, roomStore, ctx)
		fmt.Println("------------------------------")
	}

	for _, user := range users {
		seedUser(user.fname, user.lname, user.email, user.password, userStore, ctx)
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
			Size:        "single",
			Seaside:     true,
			Price:       99.9,
			MaxCapacity: 2,
			Status: []models.RoomStatus{
				{Status: "available"},
			},
		},
		{
			Size:        "double",
			Seaside:     false,
			Price:       199.9,
			MaxCapacity: 4,
			Status: []models.RoomStatus{
				{Status: "available"},
			},
		},
		{
			Size:        "king",
			Seaside:     false,
			Price:       299.9,
			MaxCapacity: 6,
			Status: []models.RoomStatus{
				{Status: "available"},
			},
		},
		{
			Size:        "king deluxe",
			Seaside:     true,
			Price:       399.9,
			MaxCapacity: 8,
			Status: []models.RoomStatus{
				{Status: "available"},
			},
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

		room.ID = insertedRoom.ID
		filter := bson.M{"_id": insertedHotel.ID}
		update := bson.M{"$push": bson.M{"rooms": room.ID}}
		hotelStore.UpdateHotel(ctx, filter, update)

		// _ = append(insertedHotel.Rooms, room.ID)

		fmt.Println("Added to hotel room: ", room.ID)
	}
}

func seedUser(fname, lname, email, password string, userStore db.UserStore, ctx context.Context) {
	user, err := models.NewUserFromParams(models.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted user: ", insertedUser)

}
