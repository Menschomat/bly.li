package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DATABASE string = "short_url_db"

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
