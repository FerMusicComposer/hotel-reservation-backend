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

	GetHotels(context.Context, bson.M) ([]*models.Hotel, error)
	GetHotelByID(context.Context, string) (*models.Hotel, error)
	InsertHotel(context.Context, *models.Hotel) (*models.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	DeleteHotel(context.Context, string) error
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

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*models.Hotel, error) {
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*models.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*models.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel models.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

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

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}
