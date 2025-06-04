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
	collection := _client.Database(database).Collection(CollectionShorts)

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
	collection := _client.Database(database).Collection(CollectionShorts)

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
	collection := _client.Database(database).Collection(CollectionShorts)

	filter := bson.D{{Key: "short", Value: shortURL.Short}}
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"url":       shortURL.URL,
			"count":     shortURL.Count,
			"updatedAt": shortURL.UpdatedAt,
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
	collection := _client.Database(database).Collection(CollectionShorts)
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

func GetShortsByOwner(owner string) *[]m.ShortURL {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Fatal(_err)
		return &[]m.ShortURL{}
	}
	collection := _client.Database(database).Collection(CollectionShorts)
	filter := bson.M{"owner": owner}

	//err := collection.Find(context.Background(), filter).Decode(&result)
	// Retrieves documents that match the query filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	// Unpacks the cursor into a slice
	var results []m.ShortURL
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return &results
}

func DeleteShortURLByShort(short string) error {
	_client, _err := GetMongoClient()
	if _err != nil {
		log.Printf("Failed to get Mongo client: %v", _err)
		return _err
	}
	database := _client.Database(database)
	collectionShorts := database.Collection(CollectionShorts)
	clicksSeries := database.Collection(CollectionClicks)
	aggClicksSeries := database.Collection(CollectionClicksAggregated)
	filter := bson.M{"short": short}

	if _, err := clicksSeries.DeleteMany(context.Background(), filter); err != nil {
		log.Printf("Failed to delete from clicks series: %v", err)
		return fmt.Errorf("failed to delete clicks: %v", err)
	}

	if _, err := aggClicksSeries.DeleteMany(context.Background(), filter); err != nil {
		log.Printf("Failed to delete from aggregated clicks series: %v", err)
		return fmt.Errorf("failed to delete aggregated clicks: %v", err)
	}

	if _, err := collectionShorts.DeleteOne(context.Background(), filter); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("no document found with the given short value: %v", short)
		}
		log.Printf("Failed to delete short URL: %v", err)
		return fmt.Errorf("failed to delete short URL: %v", err)
	}

	return nil
}
