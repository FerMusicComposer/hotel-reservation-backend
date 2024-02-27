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
	GetBookingByID(ctx context.Context, id string) (*models.Booking, error)
	InsertBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error)
	UpdateBooking(ctx context.Context, id string, update bson.M) error
}

type MongoBookingStore struct {
	connection *MongoConnection
	collection *mongo.Collection
}

func NewMongoBookingStore(conn *MongoConnection) *MongoBookingStore {
	return &MongoBookingStore{
		connection: conn,
		collection: conn.Database.Collection(bookingColl),
	}
}

func (store *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("dropping bookings collection")
	return store.collection.Database().Drop(ctx)
}

func (store *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*models.Booking, error) {
	res, err := store.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*models.Booking

	if err := res.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (store *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking *models.Booking

	res := store.collection.FindOne(ctx, bson.M{"_id": objID})
	if err := res.Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (store *MongoBookingStore) InsertBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	res, err := store.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (store *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = store.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}
