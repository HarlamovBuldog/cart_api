package mongo

import (
	"context"

	"github.com/HarlamovBuldog/cart_api/pkg/service"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddItemToCart adds item to item list of a cart with a specified ID.
// Func returns ErrNotFound if no cart was found.
func (db *DB) AddItemToCart(ctx context.Context, cartID, productName string, quantity float64) (*service.CartItem, error) {
	cartObjID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert %s to ObjectID", cartID)
	}
	cartItemID := primitive.NewObjectID()
	updateResult, err := db.Carts.UpdateOne(
		ctx,
		bson.M{"_id": cartObjID},
		bson.D{
			bson.E{Key: "$addToSet", Value: bson.D{
				bson.E{Key: "items", Value: bson.M{
					"id":       cartItemID,
					"cart_id":  cartObjID,
					"product":  productName,
					"quantity": quantity,
				}},
			}},
		})
	switch {
	case err != nil:
		return nil, errors.Wrap(err, "could not add item to cart")
	case updateResult.MatchedCount == 0:
		return nil, errors.Wrap(ErrNotFound, "no carts")
	case updateResult.ModifiedCount == 0:
		return nil, errors.New("could not add item")
	default:
		return &service.CartItem{
			ID:          cartItemID,
			CartID:      cartObjID,
			ProductName: productName,
			Quantity:    quantity,
		}, nil
	}
}

// RemoveItemFromCart removes an item with a specified ID from a cart with a specified ID.
// Func returns ErrNotFound if no cart was found or item.
func (db *DB) RemoveItemFromCart(ctx context.Context, cartID, cartItemID string) error {
	cartObjID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to ObjectID", cartID)
	}
	cartItemObjID, err := primitive.ObjectIDFromHex(cartItemID)
	if err != nil {
		return errors.Wrapf(err, "could not convert %s to ObjectID", cartItemID)
	}

	updateResult, err := db.Carts.UpdateOne(
		ctx,
		bson.M{"_id": cartObjID},
		bson.M{"$pull": bson.M{"items": bson.M{"id": cartItemObjID}}})
	switch {
	case err != nil:
		return errors.Wrap(err, "could not delete item from cart")
	case updateResult.MatchedCount == 0:
		return errors.Wrap(ErrNotFound, "no carts")
	case updateResult.ModifiedCount == 0:
		return errors.Wrap(ErrNotFound, "no items")
	default:
		return nil
	}
}

// ItemFromCart get an item with a specified ID from a cart with a specified ID.
// Func returns ErrNotFound if no cart was found or item.
func (db *DB) ItemFromCart(ctx context.Context, cartID, cartItemID string) (*service.CartItem, error) {
	cartObjID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert %s to ObjectID", cartID)
	}
	cartItemObjID, err := primitive.ObjectIDFromHex(cartItemID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not convert %s to ObjectID", cartItemID)
	}
	var cart service.Cart
	err = db.Carts.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: cartObjID},
			bson.E{Key: "items", Value: bson.M{
				"$elemMatch": bson.M{
					"id": cartItemObjID},
			}},
		}).Decode(&cart)

	switch {
	case err == mongo.ErrNoDocuments:
		return nil, errors.Wrap(ErrNotFound, "no carts or items")
	case err != nil:
		return nil, errors.Wrap(err, "could not decode document")
	default:
		var cartItemToReturn service.CartItem
		for _, cartItem := range cart.Items {
			if cartItem.ID == cartItemObjID {
				cartItemToReturn = cartItem
				break
			}
		}
		return &cartItemToReturn, nil
	}
}
