package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelId primitive.ObjectID `bson:"hotelId" json:"hotelId"`
	Size    string             `bson:"size" json:"size"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
}
