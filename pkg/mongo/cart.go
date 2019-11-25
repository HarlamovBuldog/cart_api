package mongo

import (
	"context"

	"github.com/HarlamovBuldog/cart_api/pkg/service"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddCart inserts cart to collection with primitiveObjectID generated by mongo.
func (db *DB) AddCart(ctx context.Context) (*service.Cart, error) {
	insertResult, err := db.Carts.InsertOne(ctx,
		service.Cart{
			Items: []service.CartItem{},
		})
	if err != nil {
		return nil, errors.Wrap(err, "could not insert cart")
	}
	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not convert to primitive.ObjectID")
	}

	return &service.Cart{
		ID:    insertedID,
		Items: []service.CartItem{},
	}, nil
}

// Cart returns cart with a specified id.
// Func returns ErrNotFound if no carts were found.
func (db *DB) Cart(ctx context.Context, id string) (*service.Cart, error) {
	cartID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert %s to ObjectID", id)
	}

	var cart service.Cart
	err = db.Carts.FindOne(ctx, bson.M{"_id": cartID}).Decode(&cart)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, errors.Wrap(ErrNotFound, "no carts")
	case err != nil:
		return nil, errors.Wrap(err, "could not decode document")
	default:
		return &cart, nil
	}
}