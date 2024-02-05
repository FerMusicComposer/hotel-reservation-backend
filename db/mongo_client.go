package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client       *mongo.Client
	Database     *mongo.Database
	DatabaseName string
}

func NewMongoConnection(uri, dbName string) (*MongoConnection, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoConnection{
		Client:       client,
		Database:     db,
		DatabaseName: dbName,
	}, nil
}
