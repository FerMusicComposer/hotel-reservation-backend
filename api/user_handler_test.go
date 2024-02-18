package api

import (
	"testing"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := models.CreateUserParams{
		FirstName: "Leonardo",
		LastName:  "da Vinci",
		Email:     "H3XK1@example.com",
		Password:  "password8978",
	}

	testPostUser(t, app, params)

}
