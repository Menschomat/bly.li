package mongo

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/Menschomat/bly.li/shared/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Collection names
	CollectionShorts           = "shorts"
	CollectionClicks           = "clicks"
	CollectionClicksCounts     = "clicks_counts"
	CollectionClicksCounty     = "clicks_country"
	CollectionClicksAggregated = "clicks_aggregated"
)

var (
	logger         *slog.Logger
	once           sync.Once
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

func InitMongoPackage(_logger *slog.Logger) {
	once.Do(func() {
		initMetrics()
		logger = _logger
	})
}

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
func InitMongoCollections(mongoClient *mongo.Client) {
	ctx := context.Background()

	// 1) Index für Kurz-URLs
	urlsColl := mongoClient.Database(database).Collection(CollectionShorts)

	shortIndex := mongo.IndexModel{
		Keys:    bson.M{"short": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := urlsColl.Indexes().CreateOne(ctx, shortIndex); err != nil && !mongo.IsDuplicateKeyError(err) {
		log.Fatalf("Failed to create index on 'short': %v", err)
	}

	userIndex := mongo.IndexModel{
		Keys: bson.M{"userID": 1},
	}
	if _, err := urlsColl.Indexes().CreateOne(ctx, userIndex); err != nil {
		log.Fatalf("Failed to create index on 'userID': %v", err)
	}

	// 2) Time-series clicks collection
	if err := CreateTimeSeriesCollection(mongoClient, database, CollectionClicks); err != nil {
		log.Fatalf("Could not create time-series collection: %v", err)
	}

	// 3) Time-series click_counts collection (z. B. für sekundäre Zählungen)
	if err := CreateTimeSeriesCollection(mongoClient, database, CollectionClicksCounts); err != nil {
		log.Fatalf("Could not create time-series collection: %v", err)
	}

	// 4) Aggregierte Klickdaten
	aggColl := mongoClient.Database(database).Collection(CollectionClicksAggregated)

	// Neue Struktur: Index auf echte Felder statt auf _id
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "shortUrl", Value: 1},
				{Key: "resolution", Value: 1},
				{Key: "timestamp", Value: 1},
			},
			Options: options.Index().SetName("short_resolution_timestamp_idx").SetUnique(true),
		},
		// Optional: Für spätere Filter / Drilldowns
		// {
		// 	Keys: bson.D{
		// 		{Key: "shortUrl", Value: 1},
		// 		{Key: "timestamp", Value: 1},
		// 		{Key: "browser", Value: 1},
		// 		{Key: "country", Value: 1},
		// 	},
		// 	Options: options.Index().
		// 		SetName("detailed_access_idx").
		// 		SetSparse(true),
		// },
	}

	if _, err := aggColl.Indexes().CreateMany(ctx, indexes); err != nil {
		log.Fatalf("Failed to create indexes for clicks_aggregated: %v", err)
	}
}
