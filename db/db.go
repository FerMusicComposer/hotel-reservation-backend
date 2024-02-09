package db

import "context"

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

type Dropper interface {
	Drop(context.Context) error
}
