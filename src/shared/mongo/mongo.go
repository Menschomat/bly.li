package mongo

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Menschomat/bly.li/shared/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

func GetMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var err error
		clientInstance, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoConfig().MongoServerUrl))
		if err != nil {
			log.Fatal(err)
		}
		// check the connection
		err = clientInstance.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		InitMongoCollections(clientInstance)
	})
	return clientInstance, nil
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

// InitMongoCollections sets up indexes/collections. Call once during startup.
func InitMongoCollections(mongo_client *mongo.Client) {
	// 1) Ensure "urls" collection has a unique index on "short".
	urlsColl := mongo_client.Database(database).Collection("urls")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"short": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := urlsColl.Indexes().CreateOne(context.Background(), indexModel); err != nil && !mongo.IsDuplicateKeyError(err) {
		log.Fatalf("Failed to create index on 'urls.short': %v", err)
	}

	// 2) Create or validate the time-series "click_counts" collection.
	err := CreateTimeSeriesCollection(mongo_client, database, "click_counts")
	if err != nil {
		log.Fatalf("Could not create time-series collection: %v", err)
	}
}
