package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HarlamovBuldog/cart_api/pkg/api"
	"github.com/HarlamovBuldog/cart_api/pkg/config"
	"github.com/HarlamovBuldog/cart_api/pkg/mongo"
)

const (
	dbName  = "cart_api"
	connStr = "mongodb://localhost:27018"
)

func main() {
	dbConfig := new(config.DatabaseConfig)
	dbConfig.Load(config.SERVICENAME)

	var (
		db  *mongo.DB
		err error
	)

	if dbConfig.DBName != "" && dbConfig.ConnectionString != "" {
		db, err = mongo.Connect(context.Background(), dbConfig.ConnectionString, dbConfig.DBName)
	} else {
		db, err = mongo.Connect(context.Background(), connStr, dbName)
	}

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
