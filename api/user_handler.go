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

func (userHandler *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := userHandler.userStore.GetUsers(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(users)
}

func (userHandler *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := userHandler.userStore.GetUserByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"error": fmt.Sprintf("user with ID %v not found", id)})
		}
		return Internal.from("HandleGetUser => error obtaining user by ID", err).Err
	}

	return ctx.JSON(user)
}

func (userHandler *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params models.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return BadRequest.from("HandlePostUser => error parsing user params", err).Err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}

	user, err := models.NewUserFromParams(params)
	if err != nil {
		return Internal.from("HandlePostUser => error creating user from params", err).Err
	}

	insertedUser, err := userHandler.userStore.InsertUser(ctx.Context(), user)
	if err != nil {
		return Internal.from("HandlePostUser => error inserting user", err).Err
	}

	return ctx.JSON(insertedUser)
}

func (userHandler *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		params models.UpdateUserParams
		userId = ctx.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return InvalidID.from("HandlePutUser => error converting ID(string) to ObjectID", err).Err
	}

	if err := ctx.BodyParser(&params); err != nil {
		return BadRequest.from("HandlePutUser => error parsing user params", err).Err
	}

	filter := bson.M{"_id": oid}
	if err := userHandler.userStore.UpdateUser(ctx.Context(), filter, params); err != nil {
		return ctx.JSON(map[string]string{"error": fmt.Sprintf("failed to update userID %v: %v", userId, err.Error())})
	}
	return ctx.JSON(map[string]string{"message": fmt.Sprintf("successfully updated userID %v", userId)})
}

func (userHandler *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")

	if err := userHandler.userStore.DeleteUser(ctx.Context(), userId); err != nil {
		return ctx.JSON(map[string]string{"error": fmt.Sprintf("failed to delete userID %v: %v", userId, err.Error())})
	}

	return ctx.JSON(map[string]string{"message": fmt.Sprintf("successfully deleted userID %v", userId)})
}
