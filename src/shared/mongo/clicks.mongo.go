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

// FetchClicksRange queries clicks from 'from' (inclusive) to 'to' (exclusive)
func FetchClicksRange(from, to time.Time, short string) ([]m.ShortClickCount, error) {
	_client, _err := GetMongoClient()
	if _err != nil {
		logger.Error("Error getting Mongo-Client", "error", _err)
		return nil, _err
	}
	collection := _client.Database(database).Collection(CollectionClicksAggregated)
	interval := 10 * time.Minute // or your actual interval

	// ensure UTC and truncation for slots
	from = from.UTC().Truncate(interval)
	to = to.UTC().Truncate(interval)

	filter := bson.M{
		"timestamp": bson.M{
			"$gte": from,
			"$lt":  to,
		},
		"short": short,
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	resultsMap := make(map[time.Time]int)
	for cursor.Next(context.Background()) {
		var result m.ShortClickCount
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		slot := result.Timestamp.UTC().Truncate(interval)
		resultsMap[slot] += result.Count
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	var results []m.ShortClickCount
	for t := from; t.Before(to); t = t.Add(interval) {
		results = append(results, m.ShortClickCount{
			Short:     short,
			Timestamp: t,
			Count:     resultsMap[t],
		})
	}
	return results, nil
}
