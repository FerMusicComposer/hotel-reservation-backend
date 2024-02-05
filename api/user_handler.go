package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (u *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := u.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (u *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := u.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (u *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params models.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := models.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := u.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}
