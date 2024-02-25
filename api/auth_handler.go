package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User   *models.User `json:"user"`
	Token  string       `json:"token"`
	Status int          `json:"status"`
	Msg    string       `json:"msg"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

func invalidCredentials(c *fiber.Ctx) error {
	fmt.Println("unauthorized")
	return c.Status(http.StatusUnauthorized).JSON(AuthResponse{
		Status: http.StatusUnauthorized,
		Msg:    "unauthorized",
	})
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

func createTokenFromUser(user *models.User) string {
	expires := time.Now().Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Printf("Error signing token: %v", err)
	}

	return tokenStr
}
