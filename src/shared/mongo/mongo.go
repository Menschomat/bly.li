package mongo

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

const DATABASE string = "short_url_db"

func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
		if err != nil {
			log.Fatal(err)
		}
		// check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	return client, nil
}

func CloseClientDB() {
	if clientInstance != nil {
		log.Println("Closing MongoDB connection...")

		// Create a timeout of 5 seconds for disconnection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Disconnect the client
		if err := clientInstance.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}

		log.Println("MongoDB connection closed.")
	}
}
