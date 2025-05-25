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
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	coll := client.Database(database).Collection("click_counts")
	roundedTime := clickTime //roundTo5Min(clickTime)
	// Attempt to insert a new document
	data := m.ShortClickCount{
		Short:     shortID,
		Timestamp: roundedTime,
		Count:     count,
	}
	_, err = coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return nil // Successfully inserted
}

//func InsetTimeseriesDoc(shortID string, count int, clickTime time.Time) error {
//	client, err := GetMongoClient()
//	if err != nil {
//		log.Printf("[ERROR] [%s] Failed to get MongoDB client: %v", time.Now().Format(time.RFC3339), err)
//		return err
//	}
//	coll := client.Database(database).Collection("click_counts")
//
//	roundedTime := roundTo5Min(clickTime)
//
//	filter := bson.M{
//		"short":     shortID,
//		"timestamp": roundedTime,
//	}
//	update := bson.M{
//		"$inc": bson.M{"count": count},
//	}
//
//	// Versuche Update ohne Upsert (wegen TimeSeries-Einschränkung)
//	result, err := coll.UpdateMany(context.Background(), filter, update)
//	if err != nil {
//		log.Printf("[ERROR] [%s] UpdateMany failed for shortID=%s at %s: %v",
//			time.Now().Format(time.RFC3339), shortID, roundedTime.Format(time.RFC3339), err)
//		return err
//	}
//
//	if result.MatchedCount == 0 {
//		// Kein bestehender Block -> Insert
//		doc := bson.M{
//			"short":     shortID,
//			"timestamp": roundedTime,
//			"count":     count,
//		}
//		_, err = coll.InsertOne(context.Background(), doc)
//		if err != nil {
//			log.Printf("[ERROR] [%s] InsertOne failed for shortID=%s at %s: %v",
//				time.Now().Format(time.RFC3339), shortID, roundedTime.Format(time.RFC3339), err)
//			return err
//		}
//		log.Printf("[INFO] [%s] Inserted new document for shortID=%s at %s with count=%d",
//			time.Now().Format(time.RFC3339), shortID, roundedTime.Format(time.RFC3339), count)
//	} else {
//		log.Printf("[INFO] [%s] Updated count for shortID=%s at %s (+%d)",
//			time.Now().Format(time.RFC3339), shortID, roundedTime.Format(time.RFC3339), count)
//	}
//
//	return nil
//}

// Existing helper functions remain unchanged
func roundTo5Min(t time.Time) time.Time {
	return t.Truncate(5 * time.Minute)
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

// RecordClick is a high-level function that increments the time-series
// click counter for a given short ID at the current time.
func RecordClick(shortID string) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	return IncrementShortClickCount(client, database, "click_counts", shortID)
}
