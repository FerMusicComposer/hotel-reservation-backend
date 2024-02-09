package db

import (
	"context"
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	Dropper

	InsertHotel(context.Context, *models.Hotel) (*models.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	connection *MongoConnection
	coll       *mongo.Collection
}

func NewMongoHotelStore(conn *MongoConnection) *MongoHotelStore {
	return &MongoHotelStore{
		connection: conn,
		coll:       conn.Database.Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("dropping hotels collection")
	return s.coll.Database().Drop(ctx)
}

// ------------------
// HOTEL CRUD METHODS
// ------------------

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
