package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/Menschomat/bly.li/shared/mongo"
	r "github.com/Menschomat/bly.li/shared/redis"
)

func persistUnsaved() {
	keys, err := r.GetUnsavedKeys()
	if err != nil {
		log.Fatalln("There's an error with the server:", err)
	}
	for _, key := range keys {
		short, err := r.GetShort(key)
		if err != nil {
			log.Fatalln("There's an error with the server:", err)
		}
		slog.Info("Storing changed short: " + short.Short)
		mongo.StoreShortURL(*short)
		s, err := mongo.GetShortURLByShort(short.Short)
		if err != nil {
			log.Fatalln("There's an error with the server:", err)
		}
		log.Println(short)
		log.Println(s)
		increment := short.Count - s.Count
		s.Count = short.Count
		log.Println(s)
		mongo.UpdateShortUrl(*s)
		log.Println(increment)
		mongo.InsetTimeseriesDoc(s.Short, increment, time.Now())
		r.RemoveUnsaved(short.Short)
	}
}
