package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDb struct {
	client *mongo.Client
}

var mongodb MongoDb

func (db *MongoDb) connect() (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	db.client = client
	return client, err
}

func (db *MongoDb) disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.client.Disconnect(ctx); err != nil {
    return err
	}
  return nil
}

func (db *MongoDb) Ping() error {
	client, err := db.connect()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer db.disconnect()
	err = client.Ping(ctx, readpref.Primary())
	return err
}

func (db *MongoDb) EntryCollection() *mongo.Collection {
	db.connect()
  defer db.disconnect()
	collection := db.client.Database("calorie").Collection("entity")
	return collection
}
