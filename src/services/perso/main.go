package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/Menschomat/bly.li/shared/scheduler"
)

func main() {
	log.Println("*_-_-_-BlyLi-Perso-_-_-_*")

	// Start the scheduler (runs "Hello World" every 30 seconds)
	persistUnsaved := func() {
		keys, err := redis.GetUnsavedKeys()
		if err != nil {
			log.Fatalln("There's an error with the server", err)
		}
		for _, key := range keys {
			// element is the element from someSlice for where we are
			short, err := redis.GetShort(key)
			if err != nil {
				log.Fatalln("There's an error with the server", err)
			}
			slog.Info("Storing changed short: " + short.Short)
			mongo.StoreShortURL(short)
			s, err := mongo.GetShortURLByShort(short.Short)
			if err != nil {
				log.Fatalln("There's an error with the server", err)
			}
			log.Println(short)
			log.Println(s)
			increment := short.Count - s.Count
			s.Count = short.Count
			log.Println(s)
			mongo.UpdateShortUrl(*s)
			log.Println(increment)
			for i := 0; i < increment; i++ {
				err := mongo.RecordClick(s.URL)
				if err != nil {
					log.Fatalln("There's an error with the server", err)
				}
			}
			redis.RemoveUnsaved(short.Short)
		}

	}
	scheduler := scheduler.NewScheduler(10*time.Second, persistUnsaved)

	// Handle shutdown signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stopChan:
		log.Println("Shutdown signal received. Stopping scheduler...")
		scheduler.Stop() // Stop scheduler safely
	}
	log.Println("Server shut down successfully.")
}
