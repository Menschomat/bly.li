package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Menschomat/bly.li/shared/config"
	m "github.com/Menschomat/bly.li/shared/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database = config.MongoConfig().Database

// ShortExists Check if a short URL exists in MongoDB
func ShortExists(short string) bool {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return false
	}
	collection := _client.Database(database).Collection("urls")

	filter := bson.M{"short": short}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Println("Error checking short URL in MongoDB:", err)
		return false
	}
	return count > 0
}

func StoreShortURL(shortURL m.ShortURL) (interface{}, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}
	collection := _client.Database(database).Collection("urls")

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

func UpdateShortUrl(shortURL m.ShortURL) (interface{}, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}
	collection := _client.Database(database).Collection("urls")

	filter := bson.D{{Key: "short", Value: shortURL.Short}}
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"url":   shortURL.URL,
			"count": shortURL.Count,
			// Add other fields you want to update
		}},
	}

	opts := options.Update().SetUpsert(false) // Change to true if you want to insert if not found

	updateResult, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}

	log.Printf("Matched %d document(s) and modified %d document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return updateResult.UpsertedID, nil
}

func GetShortURLByShort(short string) (*m.ShortURL, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}
	collection := _client.Database(database).Collection("urls")
	var result m.ShortURL
	filter := bson.M{"short": short}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no document found with the given short value: %v", short)
		}
		return nil, fmt.Errorf("failed to find document: %v", err)
	}
	log.Println(result.URL)

	return &result, nil
}

func DeleteShortURLByShort(short string) error {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return _err
	}
	collection := _client.Database(database).Collection("urls")
	filter := bson.M{"short": short}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("no document found with the given short value: %v", short)
		}
		return fmt.Errorf("failed to find document: %v", err)
	}
	return nil
}
