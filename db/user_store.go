package db

import (
	"context"

	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*models.User, error)
}

type MongoUserStore struct {
	connection *MongoConnection
	coll       *mongo.Collection
}

func NewMongoUserStore(conn *MongoConnection) *MongoUserStore {
	return &MongoUserStore{
		connection: conn,
		coll:       conn.Database.Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
