package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbTestConnString = "mongodb://localhost:27018"
	dbTestName       = "cart_api_test_db"
)

type initCollectionParams struct {
	CollectionName string
	Documents      []interface{}
	Opts           *options.InsertManyOptions
}

func generatePrimObjIDSet(n int) []primitive.ObjectID {
	primObjIDSet := make([]primitive.ObjectID, n)
	for i := 0; i < n; i++ {
		primObjIDSet[i] = primitive.NewObjectID()
	}

	return primObjIDSet
}

func initCollection(db *DB, initColParams initCollectionParams) error {
	var err error
	switch initColParams.CollectionName {
	case cartsCollectionName:
		_, err = db.Carts.InsertMany(context.TODO(), initColParams.Documents, initColParams.Opts)
	default:
		return errors.New("no such collection")
	}
	if err != nil {
		return errors.Wrapf(err, "could not insert to %s collection", initColParams.CollectionName)
	}

	return nil
}

func cleanUpCollection(db *DB, colName string) error {
	var err error
	switch colName {
	case cartsCollectionName:
		err = db.Carts.Drop(context.TODO())
	default:
		return errors.New("no such collection")
	}
	if err != nil {
		return errors.Wrapf(err, "could not drop %s collection", colName)
	}

	return nil
}
