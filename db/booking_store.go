package db

import (
	"context"
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	Dropper

	GetBookings(ctx context.Context, filter bson.M) ([]*models.Booking, error)
	InsertBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error)
}

type MongoBookingStore struct {
	connection *MongoConnection
	coll       *mongo.Collection
}

func NewMongoBookingStore(conn *MongoConnection) *MongoBookingStore {
	return &MongoBookingStore{
		connection: conn,
		coll:       conn.Database.Collection(bookingColl),
	}
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("dropping bookings collection")
	return s.coll.Database().Drop(ctx)
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*models.Booking, error) {
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*models.Booking

	if err := res.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}
