package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Items []CartItem         `json:"items" bson:"items"`
}

type CartItem struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	CartID      primitive.ObjectID `json:"cart_id" bson:"cart_id"`
	ProductName string             `json:"product" bson:"product"`
	Quantity    float64            `json:"quantity" bson:"quantity"`
}

type Service interface {
	AddCart(ctx context.Context) (*Cart, error)

	Cart(ctx context.Context, id string) (*Cart, error)

	AddItemToCart(ctx context.Context, cartID, productName string, quantity float64) (*CartItem, error)

	RemoveItemFromCart(ctx context.Context, cartID, cartItemID string) error
}
