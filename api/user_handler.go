package api

import (
	"strconv"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

var users = []models.User{
	{FirstName: "James", LastName: "Bond"},
	{FirstName: "Miss", LastName: "Moneypenny"},
	{FirstName: "M", LastName: "Hmmmm"},
	{FirstName: "Dr", LastName: "No"},
}

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}

func HandleGetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	user := users[id]
	return c.JSON(user)
}
