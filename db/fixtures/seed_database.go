package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
)

//=================
// HOTEL SEED DATA
//=================

type hotelData struct {
	name     string
	location string
	rating   float64
}

var hotels = []hotelData{
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

// ================
// ROOM SEED DATA
// ================
var rooms = []models.Room{
	{
		Size:        "single",
		Seaside:     true,
		Price:       99.9,
		MaxCapacity: 2,
	},
	{
		Size:        "double",
		Seaside:     false,
		Price:       199.9,
		MaxCapacity: 4,
	},
	{
		Size:        "king",
		Seaside:     false,
		Price:       299.9,
		MaxCapacity: 6,
	},
	{
		Size:        "king deluxe",
		Seaside:     true,
		Price:       399.9,
		MaxCapacity: 8,
	},
}

// ================
// USER SEED DATA
// ================
type userData struct {
	fname    string
	lname    string
	email    string
	password string
	role     string
}

var users = []userData{
	{
		fname:    "John",
		lname:    "Doe",
		email:    "jdoe@me.com",
		password: "password1123456789",
		role:     "user",
	},
	{
		fname:    "Jane",
		lname:    "Doe",
		email:    "jane@me.com",
		password: "password1123456789",
		role:     "user",
	},
	{
		fname:    "Mike",
		lname:    "Miller",
		email:    "miller@me.com",
		password: "password1123456789",
		role:     "admin",
	},
	{
		fname:    "Sarah",
		lname:    "Smith",
		email:    "smith@me.com",
		password: "password1123456789",
		role:     "user",
	},
}

func SeedData() {
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

	for _, hotel := range hotels {
		seedHotel(hotel.name, hotel.location, hotel.rating, hotelStore, roomStore, ctx)
		fmt.Println("=============================")
	}

	for _, user := range users {
		seedUser(user.fname, user.lname, user.email, user.password, user.role, userStore, ctx)
		fmt.Println("------------------------------")
	}

	conn.Database.Client().Disconnect(ctx)

}
