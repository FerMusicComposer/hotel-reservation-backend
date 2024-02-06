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
	GetUsers(context.Context) ([]*models.User, error)
	InsertUser(context.Context, *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, filter bson.M, params models.UpdateUserParams) error
	DeleteUser(context.Context, string) error
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

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*models.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*models.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err // return empty slice of users. err is nill otherwise we get the err
	}

	return users, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params models.UpdateUserParams) error {
	values := params.ToBSON()
	update := bson.D{{Key: "$set", Value: values}}

	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
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
