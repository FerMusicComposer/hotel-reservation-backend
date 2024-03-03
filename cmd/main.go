package main

import (
	"flag"
	"log"
	"os"

	"github.com/FerMusicComposer/hotel-reservation-backend/api"
	"github.com/FerMusicComposer/hotel-reservation-backend/api/middleware"
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// custom fiber config for custom error handling
var fiberConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listenAddr := flag.String("listenAddr", os.Getenv("HTTP_LISTEN_ADDRESS"), "API Server listen address")
	flag.Parse()

	conn, err := db.NewMongoConnection(db.DBURI, db.DBNAME)
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userStore := db.NewMongoUserStore(conn)
	hotelStore := db.NewMongoHotelStore(conn)
	roomStore := db.NewMongoRoomStore(conn)
	bookingStore := db.NewMongoBookingStore(conn)

	authHandler := api.NewAuthHandler(userStore)
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelStore)
	roomHandler := api.NewRoomHandler(roomStore, hotelStore, bookingStore)
	bookingHandler := api.NewBookingHandler(bookingStore, roomStore)

	app := fiber.New(fiberConfig)
	auth := app.Group("/api")
	apiV1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiV1.Group("/admin", middleware.AdminAuth)

	// AUTH ROUTES
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// ADMIN ONLY ROUTES
	admin.Get("/booking", bookingHandler.HandleGetAllBookings)
	admin.Get("/booking", bookingHandler.HandleGetAllBookingsWithinDateRange)
	admin.Get("/user", userHandler.HandleGetUsers)
	admin.Post("/room", roomHandler.HandlePostRoom)

	// USER ROUTES
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// HOTEL ROUTES
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotelById)
	apiV1.Get("/hotel/:id/rooms", roomHandler.HandleGetRoomsByHotelID)

	// ROOM ROUTES
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// BOOKING ROUTES
	apiV1.Get("/booking/:id", bookingHandler.HandleGetUserBooking)
	apiV1.Post("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	app.Listen(*listenAddr)
}
