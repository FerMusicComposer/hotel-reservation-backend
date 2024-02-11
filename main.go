package main

import (
	"flag"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/api"
	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/gofiber/fiber/v2"
)

// custom fiber config for custom error handling
var fiberConfig = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "API Server listen address")
	flag.Parse()

	conn, err := db.NewMongoConnection(db.DBURI, db.DBNAME)
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(conn))
	hotelHandler := api.NewHotelHandler(db.NewMongoHotelStore(conn))

	app := fiber.New(fiberConfig)
	apiV1 := app.Group("/api/v1")

	// USER ROUTES
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// HOTEL ROUTES
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)
}
