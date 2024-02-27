package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID      primitive.ObjectID `bson:"roomID" json:"roomID"`
	FromDate    time.Time          `bson:"fromDate" json:"fromDate"`
	ToDate      time.Time          `bson:"toDate" json:"toDate"`
	NumPeople   int                `bson:"numPeople" json:"numPeople"`
	IsCancelled bool               `bson:"isCancelled" json:"isCancelled"`
}
