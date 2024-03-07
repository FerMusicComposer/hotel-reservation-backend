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
		return Internal.with(http.StatusInternalServerError).Err

	}

	if bookingUserID != user.ID {
		return Unauthorized.with(http.StatusUnauthorized).Err
	}

	return nil
}

func invalidCredentials(ctx *fiber.Ctx) error {
	fmt.Println("unauthorized")
	return Unauthorized.from("invalid credentials", fmt.Errorf("Error %d: Unauthorized", ctx.Response().StatusCode())).Err
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
		internal := Internal.from("createTokenFromUser", err)
		fmt.Printf("Error signing token: %v", internal.Err)
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
			return InvalidID.from("validateRoomParams => GetHotelByID", fmt.Errorf("%d Bad Request: Invalid Hotel ID.\nAdditional Info: %v", http.StatusBadRequest, err)).Err
		}

		return Internal.from("validateRoomParams => GetHotelByID", err).Err
	}

	if params.Size == "" {
		return BadRequest.from("validateRoomParams => Check Size", fmt.Errorf("%d Bad Request: Size cannot be empty", http.StatusBadRequest)).Err
	}

	if params.Price <= 0 {
		return BadRequest.from("validateRoomParams => Check Price", fmt.Errorf("%d Bad Request: Price cannot be negative or zero", http.StatusBadRequest)).Err
	}

	if params.MaxCap <= 0 {
		return BadRequest.from("validateRoomParams => Check MaxCap", fmt.Errorf("%d Bad Request: MaxCapacity cannot be negative or zero", http.StatusBadRequest)).Err
	}

	return nil
}

func (params BookingParams) validateBookingParams(ctx *fiber.Ctx, roomstore db.RoomStore, roomID primitive.ObjectID) error {
	now := time.Now()

	room, err := roomstore.GetRoomByID(ctx.Context(), roomID.Hex())
	if err != nil {
		return InvalidID.from("validateBookingParams => GetRoomByID", err).Err
	}

	if params.Checkin.Before(now) {
		return BadRequest.from("validateBookingParams => Check Checkin date", fmt.Errorf("%d Bad Request: Checkin date cannot be in the past", http.StatusBadRequest)).Err
	}

	if params.Checkin.After(params.Checkout) {
		return BadRequest.from("validateBookingParams => Check Checkin date", fmt.Errorf("%d Bad Request: Checkin date cannot be after checkout date", http.StatusBadRequest)).Err
	}

	if params.NumPeople <= 0 {
		return BadRequest.from("validateBookingParams => Check NumPeople", fmt.Errorf("%d Bad Request: NumPeople cannot be negative or zero", http.StatusBadRequest)).Err
	}

	if room.MaxCapacity < params.NumPeople {
		return BadRequest.from("validateBookingParams => Check MaxCapacity", fmt.Errorf("%d Bad Request: NumPeople cannot be greater than MaxCapacity", http.StatusBadRequest)).Err
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
