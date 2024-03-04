package api

import (
	"net/http"
	"testing"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticate(t *testing.T) {
	testDB := setup(db.DBURI, db.TestDBNAME)

	app := fiber.New()
	authHandler := NewAuthHandler(testDB.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	type testCase struct {
		name               string
		params             AuthParams
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			name: "Successful authentication",
			params: AuthParams{
				Email:    "jdoe@me.com",
				Password: "password1123456789",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Failing authentication with incorrect password",
			params: AuthParams{
				Email:    "jdoe@me.com",
				Password: "wrongpassword",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Failing authentication with incorrect email",
			params: AuthParams{
				Email:    "wrongmail@example.com",
				Password: "password1123456789",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAuth(t, app, tc.params)
		})
	}
}
