package api

import (
	"testing"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	testDB := setup(db.DBURI, db.TestDBNAME)

	app := fiber.New()
	userHandler := NewUserHandler(testDB.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := models.CreateUserParams{
		FirstName: "Leonardo",
		LastName:  "da Vinci",
		Email:     "ldavinci@me.com",
		Password:  "password1123456789",
	}

	testPostUser(t, app, params)

}
