package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortURL struct {
	Name      string    `bson:"name"`
	Short     string    `bson:"short"`
	URL       string    `bson:"url"`
	Count     int       `bson:"count"`
	Owner     string    `bson:"owner,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type QrCode struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB primary key
	Short     string             `bson:"short"`         // Reference to short code or object
	Owner     string             `bson:"owner"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type ShortClick struct {
	Short     string
	Timestamp time.Time
	Ip        string
	UsrAgent  string
}

type ShortClickCount struct {
	Short     string    `bson:"short"`
	Timestamp time.Time `bson:"timestamp"`
	Count     int       `bson:"count"`
}

type ShortnReq struct {
	Url string
	//Short string
}

type ShortnRes struct {
	Url   string
	Short string
}
