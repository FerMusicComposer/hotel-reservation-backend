package main

import (
	"flag"

	"github.com/FerMusicComposer/hotel-reservation-backend/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "API Server listen address")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/users", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}
