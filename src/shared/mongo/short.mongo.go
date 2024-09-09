package mongo

import (
	"context"
	"errors"
	"fmt"
	m "github.com/Menschomat/bly.li/shared/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// ShortExists Check if a short URL exists in MongoDB
func ShortExists(short string) bool {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return false
	}
	collection := _client.Database(DATABASE).Collection("urls")

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
	collection := _client.Database(DATABASE).Collection("urls")

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
	collection := _client.Database(DATABASE).Collection("urls")
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

func GetShortsByUsr(userId string) ([]m.ShortURL, error) {
	// Get MongoDB client
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return nil, _err
	}

	// Select the "urls" collection from the database
	collection := _client.Database(DATABASE).Collection("urls")

	// Create a filter to match the owner (userId)
	filter := bson.M{"owner": userId}

	// Find all documents matching the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())

	var results []m.ShortURL

	// Iterate through the cursor and decode each document
	for cursor.Next(context.Background()) {
		var result m.ShortURL
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
		results = append(results, result)
	}

	// Check for any errors encountered during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during cursor iteration: %v", err)
	}

	// If no results were found, return an error
	if len(results) == 0 {
		return nil, fmt.Errorf("no documents found for user: %s", userId)
	}

	// Return the results slice
	return results, nil
}

func DeleteShortURLByShort(short string) error {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return _err
	}
	collection := _client.Database(DATABASE).Collection("urls")
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
