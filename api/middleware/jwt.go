package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("JWT Authentication")

		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		fmt.Println("token claims: ", claims)

		expires := claims["expires"].(float64)

		if time.Now().After(time.Unix(int64(expires), 0)) {
			fmt.Println("token expired")
			return fmt.Errorf("unauthorized")
		}

		userId := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userId)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")

		return []byte(secret), nil
	})
	if err != nil {
		fmt.Printf("Error parsing token: %v", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("Invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Error mapping claims: %v", err)
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
