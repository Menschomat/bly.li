package mongo

import (
	"context"
	"fmt"
	"time"

	m "github.com/Menschomat/bly.li/shared/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateTimeSeriesCollection creates the time-series collection and necessary indexes.
func CreateTimeSeriesCollection(client *mongo.Client, dbName, collectionName string) error {
	db := client.Database(dbName)

	// Create the time-series collection if it doesn't exist
	cmd := bson.D{
		{Key: "create", Value: collectionName},
		{Key: "timeseries", Value: bson.D{
			{Key: "timeField", Value: "timestamp"},
			{Key: "metaField", Value: "short"},
			{Key: "granularity", Value: "minutes"},
		}},
	}

	err := db.RunCommand(context.Background(), cmd).Err()
	if err != nil && !isNamespaceExistsError(err) {
		return fmt.Errorf("failed to create time-series collection: %v", err)
	}

	// Create a unique compound index on (short, timestamp)
	//coll := db.Collection(collectionName)
	//indexModel := mongo.IndexModel{
	//	Keys: bson.D{
	//		{Key: "short", Value: 1},
	//		{Key: "timestamp", Value: 1},
	//	},
	//	Options: options.Index().SetUnique(true),
	//}
	//
	//_, err = coll.Indexes().CreateOne(context.Background(), indexModel)
	//if err != nil {
	//	// Check if the error is about the index already existing; if not, return error
	//	if !isIndexAlreadyExistsError(err) {
	//		return fmt.Errorf("failed to create unique index: %v", err)
	//	}
	//}

	return nil
}

// isNamespaceExistsError checks if the error is due to the collection already existing.
func isNamespaceExistsError(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
		return true
	}
	return false
}

// isIndexAlreadyExistsError checks if the error is due to the index already existing.
func isIndexAlreadyExistsError(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok && (cmdErr.Code == 85 || cmdErr.Code == 86) {
		return true
	}
	return false
}

// IncrementShortClickCount handles click counting using insert and update on duplicate.
func IncrementShortClickCount(client *mongo.Client, dbName, colName, shortID string, clickTime time.Time) error {
	coll := client.Database(dbName).Collection(colName)
	roundedTime := roundTo5Min(clickTime)

	// Attempt to insert a new document
	data := m.ShortClickCount{
		Short:     shortID,
		Timestamp: roundedTime,
		Count:     1,
	}
	_, err := coll.InsertOne(context.Background(), data)
	if err == nil {
		return nil // Successfully inserted
	}

	// Check for duplicate key error
	if !isDuplicateKeyError(err) {
		return fmt.Errorf("failed to insert click count: %v", err)
	}

	// Update existing document to increment count
	filter := bson.D{
		{Key: "short", Value: shortID},
		{Key: "timestamp", Value: roundedTime},
	}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}},
	}
	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update click count after duplicate: %v", err)
	}

	return nil
}

// isDuplicateKeyError checks if the error is a MongoDB duplicate key error.
func isDuplicateKeyError(err error) bool {
	if writeErr, ok := err.(mongo.WriteException); ok {
		for _, e := range writeErr.WriteErrors {
			if e.Code == 11000 { // MongoDB duplicate key error code
				return true
			}
		}
	}
	return false
}

// Existing helper functions remain unchanged
func roundTo5Min(t time.Time) time.Time {
	return t.Truncate(5 * time.Minute)
}

// GetClicksForShort queries all time-series documents for a given short URL.
// Optionally, you can add date-range filters.
func GetClicksForShort(client *mongo.Client, dbName, colName, shortID string) ([]m.ShortClickCount, error) {
	coll := client.Database(dbName).Collection(colName)

	cursor, err := coll.Find(
		context.Background(),
		bson.M{"short": shortID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query click counts: %v", err)
	}
	defer cursor.Close(context.Background())

	var results []m.ShortClickCount
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

// RecordClick is a high-level function that increments the time-series
// click counter for a given short ID at the current time.
func RecordClick(shortID string) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	return IncrementShortClickCount(client, database, "click_counts", shortID, time.Now())
}
