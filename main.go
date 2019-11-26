package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HarlamovBuldog/cart_api/pkg/api"
	"github.com/HarlamovBuldog/cart_api/pkg/mongo"
)

func main() {
	//dbConfig := new(config.DatabaseConfig)
	//dbConfig.Load(config.SERVICENAME)

	db, err := mongo.Connect(context.Background(), "mongodb://localhost:27018", "cart_api")
	if err != nil {
		log.Fatal("could not connect to mongo")
	}
	srv := &http.Server{
		Addr:    ":27000",
		Handler: api.New(db),
	}

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe(): %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Print("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error shutdown server: %s", err)
	}

	log.Print("Server stopped")
}
