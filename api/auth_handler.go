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

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return invalidCredentials(c)
	}

	if !models.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	token := createTokenFromUser(user)

	response := AuthResponse{
		User:   user,
		Token:  token,
		Status: http.StatusOK,
		Msg:    "OK",
	}

	fmt.Printf("authenticated --> %s ; role: %s", user.Email, user.Role)

	return c.JSON(response)
}
