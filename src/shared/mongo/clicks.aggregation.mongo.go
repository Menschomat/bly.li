package mongo

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	aggregationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "blyli_click_aggregation_duration_seconds",
			Help:    "Duration of click aggregation per resolution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"resolution"},
	)
	aggregationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blyli_click_aggregation_errors_total",
			Help: "Number of aggregation errors",
		},
		[]string{"resolution"},
	)
	aggregationBucketsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "blyli_click_aggregation_buckets_total",
			Help: "Number of successfully processed buckets",
		},
		[]string{"resolution"},
	)
)

func initMetrics() {
	prometheus.MustRegister(aggregationDuration)
	prometheus.MustRegister(aggregationErrors)
	prometheus.MustRegister(aggregationBucketsProcessed)
}

// AggregateClicksByResolution aggregates raw clicks into a single time bucket.
func AggregateClicksByResolution(
	ctx context.Context,
	client *mongo.Client,
	resolution string, // e.g., "10min"
	unit string, // e.g., "minute"
	binSize int, // e.g., 10
	from, to time.Time,
) error {
	raw := client.Database(database).Collection(CollectionClicks)
	timer := prometheus.NewTimer(aggregationDuration.WithLabelValues(resolution))
	defer timer.ObserveDuration()
	pipeline := mongo.Pipeline{
		// 1. Zeitraum filtern
		{{
			Key: "$match",
			Value: bson.D{
				{Key: "timestamp", Value: bson.D{
					{Key: "$gte", Value: from},
					{Key: "$lt", Value: to},
				}},
			},
		}},
		// 2. Gruppieren nach short + bucketed timestamp
		{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "shortUrl", Value: "$short"},
					{Key: "timestamp", Value: bson.D{
						{Key: "$dateTrunc", Value: bson.D{
							{Key: "date", Value: "$timestamp"},
							{Key: "unit", Value: unit},
							{Key: "binSize", Value: binSize},
						}},
					}},
					{Key: "resolution", Value: resolution},
				}},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		}},
		// 3. _id auflösen in echte Felder
		{{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "shortUrl", Value: "$_id.shortUrl"},
				{Key: "timestamp", Value: "$_id.timestamp"},
				{Key: "resolution", Value: "$_id.resolution"},
				{Key: "count", Value: "$count"},
			},
		}},
		// 4. Merge → Ersetze bestehende Werte (idempotent)
		{{
			Key: "$merge",
			Value: bson.D{
				{Key: "into", Value: CollectionClicksAggregated},
				{Key: "on", Value: bson.A{"shortUrl", "timestamp", "resolution"}},
				{Key: "whenMatched", Value: "replace"},
				{Key: "whenNotMatched", Value: "insert"},
			},
		}},
	}
	_, err := raw.Aggregate(ctx, pipeline)
	if err != nil {
		aggregationErrors.WithLabelValues(resolution).Inc()
		return err
	}
	aggregationBucketsProcessed.WithLabelValues(resolution).Inc()
	return nil
}

// RunClickAggregation triggers aggregation for the previous *and* current bucket (e.g. last 20 min).
func RunClickAggregation() {
	ctx := context.Background()
	client, _ := GetMongoClient()
	binSize := 10
	unit := "minute"
	resolution := "10min"
	now := time.Now().UTC()
	// Bucket-Startzeiten berechnen
	currentBucketStart := now.Truncate(time.Duration(binSize) * time.Minute)
	prevBucketStart := currentBucketStart.Add(-time.Duration(binSize) * time.Minute)
	nextBucketStart := currentBucketStart.Add(time.Duration(binSize) * time.Minute)
	// 1. Vorheriger Bucket: prevBucketStart bis currentBucketStart
	err := AggregateClicksByResolution(ctx, client, resolution, unit, binSize, prevBucketStart, currentBucketStart)
	if err != nil {
		logger.Error("Failed aggregation for previous bucket",
			"resolution", resolution,
			"bucket_start", prevBucketStart.Format(time.RFC3339),
			"bucket_end", currentBucketStart.Format(time.RFC3339),
			"error", err)
	}
	// 2. Aktueller Bucket: currentBucketStart bis nextBucketStart
	err = AggregateClicksByResolution(ctx, client, resolution, unit, binSize, currentBucketStart, nextBucketStart)
	if err != nil {
		logger.Error("Failed aggregation for current bucket",
			"resolution", resolution,
			"bucket_start", currentBucketStart.Format(time.RFC3339),
			"bucket_end", nextBucketStart.Format(time.RFC3339),
			"error", err)
	}
}
