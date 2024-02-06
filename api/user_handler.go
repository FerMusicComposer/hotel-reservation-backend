package api

import (
	"errors"
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": fmt.Sprintf("user with ID %v not found", id)})
		}
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

func (u *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params models.UpdateUserParams
		userId = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := u.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return c.JSON(map[string]string{"error": fmt.Sprintf("failed to update userID %v: %v", userId, err.Error())})
	}
	return c.JSON(map[string]string{"message": fmt.Sprintf("successfully updated userID %v", userId)})
}

func (u *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	if err := u.userStore.DeleteUser(c.Context(), userId); err != nil {
		return c.JSON(map[string]string{"error": fmt.Sprintf("failed to delete userID %v: %v", userId, err.Error())})
	}

	return c.JSON(map[string]string{"message": fmt.Sprintf("successfully deleted userID %v", userId)})
}
