package model

import "time"

type ShortURL struct {
	Short string `bson:"short"`
	URL   string `bson:"url"`
	Count int    `bson:"count"`
	Owner string `bson:"owner,omitempty"`
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
