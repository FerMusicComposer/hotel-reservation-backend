package main

import (
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/db/fixtures"
)

func main() {

	fmt.Println("seeding database...")
	fmt.Println("=============================")

	fixtures.SeedData()

	fmt.Println("=============================")
	fmt.Println("disconnecting...")

}
