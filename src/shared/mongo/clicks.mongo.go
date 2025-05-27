package mongo

import (
	"context"
	"fmt"
	"time"

	m "github.com/Menschomat/bly.li/shared/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTimeSeriesCollection creates the time-series collection and necessary indexes.
func CreateTimeSeriesCollection(client *mongo.Client, dbName, collectionName string) error {
	db := client.Database(dbName)

	cmd := bson.D{
		{Key: "create", Value: collectionName},
		{Key: "timeseries", Value: bson.D{
			{Key: "timeField", Value: "timestamp"},
			{Key: "metaField", Value: "short"},
			//({Key: "granularity", Value: "minutes"},
			{Key: "bucketMaxSpanSeconds", Value: 300},  // 5 minutes
			{Key: "bucketRoundingSeconds", Value: 300}, // 5-minute alignment
		}},
	}

	err := db.RunCommand(context.Background(), cmd).Err()
	if err != nil && !isNamespaceExistsError(err) {
		return fmt.Errorf("failed to create time-series collection: %v", err)
	}
	return nil
}

// isNamespaceExistsError checks if the error is due to the collection already existing.
func isNamespaceExistsError(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
		return true
	}
	return false
}

// IncrementShortClickCount handles atomic increments with proper time-series constraints
func IncrementShortClickCount(client *mongo.Client, dbName, colName, shortID string) error {
	coll := client.Database(dbName).Collection(colName)
	now := time.Now().UTC()

	// Magic happens in the document ID structure
	docID := bson.D{
		{Key: "short", Value: shortID},
		{Key: "timestamp", Value: now.Truncate(5 * time.Minute)},
	}

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}},
		{Key: "$setOnInsert", Value: docID}, // For new documents
	}

	// Leverage upsert through native driver capabilities
	_, err := coll.UpdateMany(
		context.Background(),
		docID, // Exact match on compound ID
		update,
		options.Update().SetUpsert(false),
	)

	return err
}

func InsetTimeseriesDoc(shortID string, count int, clickTime time.Time) error {
	data := m.ShortClickCount{
		Short:     shortID,
		Timestamp: clickTime,
		Count:     count,
	}
	// Use a slice literal to wrap 'data' as []interface{}
	return InsetTimeseriesData("click_counts", []m.ShortClickCount{data})
}

func InsetTimeseriesData[T any](colname string, data []T) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	coll := client.Database(database).Collection(colname)

	// Convert []T to []interface{} only inside this function
	interfaceSlice := make([]interface{}, len(data))
	for i, v := range data {
		interfaceSlice[i] = v
	}

	_, err = coll.InsertMany(context.Background(), interfaceSlice)
	if err != nil {
		return err
	}
	return nil // Successfully inserted
}

// GetClicksForShort queries all time-series documents for a given short URL.
// Optionally, you can add date-range filters.
func GetClicksForShort(client *mongo.Client, dbName, colName, shortID string) ([]m.ShortClickCount, error) {
	coll := client.Database(dbName).Collection(colName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "short", Value: shortID}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$dateTrunc", Value: bson.D{
					{Key: "date", Value: "$timestamp"},
					{Key: "unit", Value: "minute"},
					{Key: "binSize", Value: 5},
				}},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: "$count"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "timestamp", Value: "$_id"},
			{Key: "count", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate click counts: %v", err)
	}
	defer cursor.Close(context.Background())

	var results []m.ShortClickCount
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
