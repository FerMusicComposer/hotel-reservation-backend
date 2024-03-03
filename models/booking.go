package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID      primitive.ObjectID `bson:"roomID" json:"roomID"`
	Checkin     time.Time          `bson:"checkin" json:"checkin"`
	Checkout    time.Time          `bson:"checkout" json:"checkout"`
	NumPeople   int                `bson:"numPeople" json:"numPeople"`
	IsCancelled bool               `bson:"isCancelled" json:"isCancelled"`
}
