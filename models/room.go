package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelId     primitive.ObjectID `bson:"hotelId" json:"hotelId"`
	Size        string             `bson:"size" json:"size"`
	Seaside     bool               `bson:"seaside" json:"seaside"`
	Price       float64            `bson:"price" json:"price"`
	MaxCapacity int                `bson:"maxCapacity" json:"maxCapacity"`
	Status      []RoomStatus       `bson:"status" json:"status"`
}

type RoomStatus struct {
	Status     string             `bson:"status" json:"status"`
	BookingID  primitive.ObjectID `bson:"bookingId,omitempty" json:"bookingId"`
	BookedFrom time.Time          `bson:"bookedFrom" json:"bookedFrom"`
	BookedTo   time.Time          `bson:"bookedTo" json:"bookedTo"`
}
