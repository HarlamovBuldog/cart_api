package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pkg/errors"
)

// DB is the repository, with all of the methods that are required to get info from the db.
type DB struct {
	Carts *mongo.Collection
}

const cartsCollectionName = "carts"

// ErrNotFound is used when result of select statement is empty.
var ErrNotFound = errors.New("not found")

// Connect connects to mongo DB with url, gets database with dbName and returns DB.
func Connect(ctx context.Context, url, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, errors.Wrap(err, "could not create mongo client")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not ping mongo client")
	}

	db := client.Database(dbName)
	carts := db.Collection(cartsCollectionName)

	return &DB{Carts: carts}, nil
}
