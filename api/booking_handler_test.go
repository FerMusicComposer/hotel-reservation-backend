package api

import (
	"net/http"
	"testing"

	"github.com/FerMusicComposer/hotel-reservation-backend/api/middleware"
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	testDB := setup(db.DBURI, db.TestDBNAME)

	app := fiber.New()
	admin := app.Group("/admin", middleware.JWTAuthentication(testDB.UserStore), middleware.AdminAuth)
	authHandler := NewAuthHandler(testDB.UserStore)
	bookingHandler := NewBookingHandler(testDB.BookingStore, testDB.RoomStore)

	app.Post("/auth", authHandler.HandleAuthenticate)
	admin.Get("/bookings", bookingHandler.HandleGetAllBookings)

	type testCase struct {
		name               string
		params             AuthParams
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name: "Successful retrieval of all bookings",
			params: AuthParams{
				Email:    "miller@me.com",
				Password: "password1123456789",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Failed retrieval of all bookings for non-admin user",
			params: AuthParams{
				Email:    "jdoe@me.com",
				Password: "password1123456789",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAdminGetAllBookings(t, app, tc.params)
		})
	}
}
