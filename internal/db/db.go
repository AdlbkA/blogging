package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB *mongo.Client

func InitDb(uri string) error {
	clientdb := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientdb)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")

	DB = client

	return nil
}

func CloseDb() {
	if err := DB.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to close MongoDB connection: %v", err)
	}
}
