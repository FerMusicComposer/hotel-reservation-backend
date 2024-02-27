package api

import (
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isUserToken(ctx *fiber.Ctx, bookingUserID primitive.ObjectID) error {
	user, ok := ctx.Context().UserValue("user").(*models.User)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	if bookingUserID != user.ID {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthorized",
			},
		)
	}

	return nil
}
