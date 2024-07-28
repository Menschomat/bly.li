package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	m "github.com/Menschomat/bly.li/shared/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetMongoClient() (*mongo.Client, error) {
	if client != nil {
		return client, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client, nil
}
func CloseClientDB() {
	if client == nil {
		return
	}
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}

func StoreShortURL(shortURL m.ShortURL) (interface{}, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}
	collection := _client.Database("short_url_db").Collection("urls")

	// Create an index on the short field to ensure uniqueness
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"short": 1}, // Index in ascending order
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	insertResult, err := collection.InsertOne(context.Background(), shortURL)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}

	return insertResult.InsertedID, nil
}

func GetShortURLByShort(short string) (*m.ShortURL, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}
	collection := _client.Database("short_url_db").Collection("urls")
	var result m.ShortURL
	filter := bson.M{"short": short}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found with the given short value: %v", short)
		}
		return nil, fmt.Errorf("failed to find document: %v", err)
	}

	return &result, nil
}
