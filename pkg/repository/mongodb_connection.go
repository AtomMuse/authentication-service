package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConnection represents a MongoDB connection.
type MongoDBConnection struct {
	Client *mongo.Client
}

// NewMongoDBConnection creates a new MongoDB connection.
func NewMongoDBConnection(uri string) (*MongoDBConnection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return &MongoDBConnection{
		Client: client,
	}, nil
}

// GetDatabase returns a database object.
func (conn *MongoDBConnection) GetDatabase(dbName string) *mongo.Database {
	return conn.Client.Database(dbName)
}

// Close closes the MongoDB connection.
func (conn *MongoDBConnection) Close() error {
	err := conn.Client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	log.Println("Disconnected from MongoDB!")

	return nil
}
