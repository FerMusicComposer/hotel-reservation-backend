package api

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticate(t *testing.T) {
	tdb := setup(t)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
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
				Email:    "H3XK1@example.com",
				Password: "password8978",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Failing authentication with incorrect password",
			params: AuthParams{
				Email:    "H3XK1@example.com",
				Password: "wrongpassword",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Failing authentication with incorrect email",
			params: AuthParams{
				Email:    "wrongmail@example.com",
				Password: "password8978",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAuth(t, app, tc.params)
		})
	}

	tdb.teardown(t)
}
