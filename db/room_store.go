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

	GetRooms(context.Context) ([]*models.Room, error)
	GetRoomsByHotelID(context.Context, bson.M) ([]*models.Room, error)
	GetRoomByID(context.Context, string) (*models.Room, error)
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

func (s *MongoRoomStore) GetRooms(ctx context.Context) ([]*models.Room, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var rooms []*models.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, filter bson.M) ([]*models.Room, error) {
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*models.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var room models.Room
	if err := s.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

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
