package fixtures

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SeedBooking(userStore db.UserStore, hotelStore db.HotelStore, roomStore db.RoomStore, bookingStore db.BookingStore, ctx context.Context) {
	users, err := userStore.GetUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	hotels, err := hotelStore.GetHotels(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		selectedHotels := randomSubsetOfHotels(hotels, 1, len(hotels)/2)

		for hIndex, hotel := range selectedHotels {
			rooms, err := roomStore.GetRoomsByHotelID(ctx, bson.M{"hotelId": hotel.ID})
			if err != nil {
				log.Fatal(err)
			}

			selectedRooms := randomSubsetOfRooms(rooms, 1, len(rooms)/2)

			for rIndex, room := range selectedRooms {
				dayOffset := (hIndex * len(rooms)) + rIndex
				checkin := time.Now().Add(time.Hour * 24 * time.Duration(dayOffset))
				checkout := checkin.Add(time.Hour * 24)

				booking := &models.Booking{
					UserID:    user.ID,
					RoomID:    room.ID,
					Checkin:   checkin,
					Checkout:  checkout,
					NumPeople: 2,
				}

				_, err = bookingStore.InsertBooking(ctx, booking)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func randomSubsetOfHotels(hotels []*models.Hotel, min, max int) []*models.Hotel {
	n := rand.Intn(max-min+1) + min
	rand.Shuffle(len(hotels), func(i, j int) {
		hotels[i], hotels[j] = hotels[j], hotels[i]
	})
	return hotels[:n]
}

// Implement randomSubsetOfRooms to select a random subset of rooms
func randomSubsetOfRooms(rooms []*models.Room, min, max int) []*models.Room {
	n := rand.Intn(max-min+1) + min
	rand.Shuffle(len(rooms), func(i, j int) {
		rooms[i], rooms[j] = rooms[j], rooms[i]
	})
	return rooms[:n]
}

// func SeedBooking(userStore db.UserStore, hotelStore db.HotelStore, roomStore db.RoomStore, bookingStore db.BookingStore, ctx context.Context) {
// 	user, err := userStore.GetUserByEmail(ctx, "jdoe@me.com")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	hotels, err := hotelStore.GetHotels(ctx, bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	hotelID := hotels[0].ID

// 	rooms, err := roomStore.GetRoomsByHotelID(ctx, bson.M{"hotelId": hotelID})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, room := range rooms {
// 		booking := &models.Booking{
// 			UserID:    user.ID,
// 			RoomID:    room.ID,
// 			Checkin:   time.Now(),
// 			Checkout:  time.Now().Add(time.Hour * 24),
// 			NumPeople: 2,
// 		}

// 		_, err = bookingStore.InsertBooking(ctx, booking)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
