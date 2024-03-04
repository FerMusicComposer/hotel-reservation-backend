package main

import (
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/db/fixtures"
)

func main() {

	fmt.Println("seeding databases...")
	fmt.Println("=============================")
	fmt.Println("")

	fmt.Println("main database...")
	fmt.Println("")
	fixtures.SeedData(db.DBURI, db.DBNAME)

	fmt.Println("=============================")
	fmt.Println("disconnecting main database...")
	fmt.Println("")

	fmt.Println("test database...")
	fmt.Println("")
	fixtures.SeedData(db.DBURI, db.TestDBNAME)

	fmt.Println("=============================")
	fmt.Println("disconnecting test database...")
	fmt.Println("")

}
