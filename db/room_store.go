package db

import (
	"context"
	"fmt"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *models.Room) (*models.Room, error)
}

type MongoRoomStore struct {
	connection *MongoConnection
	coll       *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(conn *MongoConnection, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		connection: conn,
		coll:       conn.Database.Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("dropping rooms collection")
	return s.coll.Database().Drop(ctx)
}

// ------------------
// ROOM CRUD METHODS
// ------------------

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelId}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err = s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
