package api

import (
	"fmt"
	"net/http"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

func (authHandler *AuthHandler) HandleAuthenticate(ctx *fiber.Ctx) error {
	var params AuthParams

	if err := ctx.BodyParser(&params); err != nil {
		return Internal.from("HandleAuthenticate => ctx.BodyParser", err).Err
	}

	user, err := authHandler.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		return invalidCredentials(ctx)
	}

	if !models.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return invalidCredentials(ctx)
	}

	token := createTokenFromUser(user)

	response := AuthResponse{
		User:   user,
		Token:  token,
		Status: http.StatusOK,
		Msg:    "OK",
	}

	fmt.Printf("authenticated --> %s ; role: %s", user.Email, user.Role)

	return ctx.JSON(response)
}
