package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ========================
// AUTH HELPERS AND PARAMS
// ========================
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

func invalidCredentials(c *fiber.Ctx) error {
	fmt.Println("unauthorized")
	return c.Status(http.StatusUnauthorized).JSON(AuthResponse{
		Status: http.StatusUnauthorized,
		Msg:    "unauthorized",
	})
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

// ================================
// ROOM HANDLER HELPERS AND PARAMS
// ================================
type RoomParams struct {
	HotelID string  `json:"hotelId"`
	Size    string  `json:"size"`
	Seaside bool    `json:"seaside"`
	Price   float64 `json:"price"`
	MaxCap  int     `json:"maxCapacity"`
}

type BookingParams struct {
	Checkin     time.Time `json:"checkin"`
	Checkout    time.Time `json:"checkout"`
	NumPeople   int       `json:"numPeople"`
	IsCancelled bool      `json:"isCancelled"`
}

func (params RoomParams) validateRoomParams(ctx *fiber.Ctx, hotelStore db.HotelStore) error {
	_, err := hotelStore.GetHotelByID(ctx.Context(), params.HotelID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid hotel id")
		}

		return fmt.Errorf("internal server error")
	}

	if params.Size == "" {
		return fmt.Errorf("must specify a room size")
	}

	if params.Price <= 0 {
		return fmt.Errorf("price cannot be negative or zero")
	}

	if params.MaxCap <= 0 {
		return fmt.Errorf("max capacity cannot be negative or zero")
	}

	return nil
}

func (params BookingParams) validateBookingParams(ctx *fiber.Ctx, roomstore db.RoomStore, roomID primitive.ObjectID) error {
	now := time.Now()

	room, err := roomstore.GetRoomByID(ctx.Context(), roomID.Hex())
	if err != nil {
		return err
	}

	if params.Checkin.Before(now) {
		return fmt.Errorf("cannot book a room in the past")
	}

	if params.Checkin.After(params.Checkout) {
		return fmt.Errorf("from date cannot be superior to end date")
	}

	if params.NumPeople <= 0 {
		return fmt.Errorf("invalid number of people")
	}

	if room.MaxCapacity < params.NumPeople {
		return fmt.Errorf("room capacity exedeced")
	}

	for _, status := range room.Status {
		if datesAreWithinRange(params.Checkin, params.Checkout, status.BookedFrom, status.BookedTo) &&
			status.Status != "cancelled" {
			return fmt.Errorf("room is already booked")
		}
	}

	return nil
}

func datesAreWithinRange(checkin time.Time, checkout time.Time, bookedFrom time.Time, bookedTo time.Time) bool {
	if checkin.Equal(bookedFrom) ||
		checkout.Equal(bookedTo) ||
		checkin.After(bookedFrom) && checkout.Before(bookedTo) ||
		checkin.Before(bookedFrom) && checkout.After(bookedTo) ||
		checkin.Before(bookedTo) && (checkout.Equal(bookedTo) || checkout.After(bookedTo)) ||
		checkin.Before(bookedFrom) && (checkout.Equal(bookedFrom) || (checkout.After(bookedFrom) && checkout.Before(bookedTo))) {

		return true
	}

	return false
}

// ===================================
// BOOKING HANDLER HELPERS AND PARAMS
// ===================================
type BookingQueryParams struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}
