package mongo

import (
	"context"
	"testing"

	"github.com/HarlamovBuldog/cart_api/pkg/service"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddItemToCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(3)
	tt := []struct {
		name                string
		initColParams       initCollectionParams
		cartID              string
		productName         string
		quantity            float64
		expectedErr         error
		isCustomErrExpected bool
	}{
		{
			name:        "correct test",
			cartID:      cartObjIDSet[0].Hex(),
			productName: "product_1",
			quantity:    10.0,
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID:    cartObjIDSet[0],
						Items: []service.CartItem{},
					},
					service.Cart{
						ID:    cartObjIDSet[1],
						Items: []service.CartItem{},
					},
				},
				Opts: nil,
			},
			isCustomErrExpected: false,
			expectedErr:         nil,
		},
		{
			name:   "incorrect test: bad cartID provided",
			cartID: "bad_id",
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID:    cartObjIDSet[0],
						Items: []service.CartItem{},
					},
					service.Cart{
						ID:    cartObjIDSet[1],
						Items: []service.CartItem{},
					},
				},
				Opts: nil,
			},
			isCustomErrExpected: true,
			expectedErr:         errors.New("could not convert"),
		},
		{
			name:   "incorrect test: ErrNotFound",
			cartID: cartObjIDSet[2].Hex(),
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID:    cartObjIDSet[0],
						Items: []service.CartItem{},
					},
					service.Cart{
						ID:    cartObjIDSet[1],
						Items: []service.CartItem{},
					},
				},
				Opts: nil,
			},
			isCustomErrExpected: false,
			expectedErr:         ErrNotFound,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			connTest, err := Connect(context.Background(), dbTestConnString, dbTestName)
			require.NoError(t, err, "could not create db instance")

			defer func() {
				err = cleanUpCollection(connTest, tc.initColParams.CollectionName)
				assert.NoError(t, err, "cleanUpCollection")
			}()

			err = initCollection(connTest, tc.initColParams)
			require.NoError(t, err, "initCollection")

			expectedCartItem, err := connTest.AddItemToCart(context.Background(), tc.cartID, tc.productName, tc.quantity)
			switch {
			case tc.isCustomErrExpected && tc.expectedErr != nil && err != nil:
				assert.Contains(t, err.Error(), tc.expectedErr.Error(), "Actual error should contain text from expected error")
			default:
				assert.Equal(t, tc.expectedErr, errors.Cause(err), "Two errors should be the same")
				if expectedCartItem != nil {
					actualCartItem, cartErr := connTest.ItemFromCart(context.Background(), tc.cartID, expectedCartItem.ID.Hex())
					assert.Equal(t, expectedCartItem, actualCartItem, "Two objects should be the same")
					assert.NoError(t, cartErr)
				}
			}
		})
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(3)
	cartItemObjIDSet := generatePrimObjIDSet(4)
	tt := []struct {
		name                          string
		initColParams                 initCollectionParams
		cartID                        string
		cartItemID                    string
		expectedCart                  *service.Cart
		expectedRemoveItemErr         error
		isCustomRemoveItemErrExpected bool
		expectedGetCartErr            error
	}{
		{
			name:       "correct test",
			cartID:     cartObjIDSet[0].Hex(),
			cartItemID: cartItemObjIDSet[0].Hex(),
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID: cartObjIDSet[0],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    10.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    20.0,
							},
						},
					},
					service.Cart{
						ID: cartObjIDSet[1],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    5.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    15.0,
							},
						},
					},
				},
				Opts: nil,
			},
			expectedCart: &service.Cart{
				ID: cartObjIDSet[0],
				Items: []service.CartItem{
					{
						ID:          cartItemObjIDSet[1],
						CartID:      cartObjIDSet[0],
						ProductName: "product_2",
						Quantity:    20.0,
					},
				},
			},
			isCustomRemoveItemErrExpected: false,
			expectedGetCartErr:            nil,
		},
		{
			name:       "incorrect test: bad cartID provided",
			cartID:     "bad_id",
			cartItemID: cartItemObjIDSet[0].Hex(),
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID: cartObjIDSet[0],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    10.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    20.0,
							},
						},
					},
					service.Cart{
						ID: cartObjIDSet[1],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    5.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    15.0,
							},
						},
					},
				},
				Opts: nil,
			},
			isCustomRemoveItemErrExpected: true,
			expectedRemoveItemErr:         errors.New("could not convert"),
		},
		{
			name:       "incorrect test: bad cartItemID provided",
			cartID:     cartObjIDSet[0].Hex(),
			cartItemID: "bad_id",
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID: cartObjIDSet[0],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    10.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    20.0,
							},
						},
					},
					service.Cart{
						ID: cartObjIDSet[1],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    5.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    15.0,
							},
						},
					},
				},
				Opts: nil,
			},
			isCustomRemoveItemErrExpected: true,
			expectedRemoveItemErr:         errors.New("could not convert")},
		{
			name:       "incorrect test: ErrNotFound: no carts",
			cartID:     cartObjIDSet[2].Hex(),
			cartItemID: cartItemObjIDSet[0].Hex(),
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID: cartObjIDSet[0],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    10.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    20.0,
							},
						},
					},
					service.Cart{
						ID: cartObjIDSet[1],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    5.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    15.0,
							},
						},
					},
				},
				Opts: nil,
			},
			isCustomRemoveItemErrExpected: false,
			expectedRemoveItemErr:         ErrNotFound,
			expectedCart:                  nil,
			expectedGetCartErr:            ErrNotFound,
		},
		{
			name:       "incorrect test: ErrNotFound: no items",
			cartID:     cartObjIDSet[0].Hex(),
			cartItemID: cartItemObjIDSet[2].Hex(),
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						ID: cartObjIDSet[0],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    10.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    20.0,
							},
						},
					},
					service.Cart{
						ID: cartObjIDSet[1],
						Items: []service.CartItem{
							{
								ID:          cartItemObjIDSet[0],
								CartID:      cartObjIDSet[0],
								ProductName: "product_1",
								Quantity:    5.0,
							},
							{
								ID:          cartItemObjIDSet[1],
								CartID:      cartObjIDSet[0],
								ProductName: "product_2",
								Quantity:    15.0,
							},
						},
					},
				},
				Opts: nil,
			},
			isCustomRemoveItemErrExpected: false,
			expectedRemoveItemErr:         ErrNotFound,
			expectedCart: &service.Cart{
				ID: cartObjIDSet[0],
				Items: []service.CartItem{
					{
						ID:          cartItemObjIDSet[0],
						CartID:      cartObjIDSet[0],
						ProductName: "product_1",
						Quantity:    10.0,
					},
					{
						ID:          cartItemObjIDSet[1],
						CartID:      cartObjIDSet[0],
						ProductName: "product_2",
						Quantity:    20.0,
					},
				},
			},
			expectedGetCartErr: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			connTest, err := Connect(context.Background(), dbTestConnString, dbTestName)
			require.NoError(t, err, "could not create db instance")

			defer func() {
				err = cleanUpCollection(connTest, tc.initColParams.CollectionName)
				assert.NoError(t, err, "cleanUpCollection")
			}()

			err = initCollection(connTest, tc.initColParams)
			require.NoError(t, err, "initCollection")

			err = connTest.RemoveItemFromCart(context.Background(), tc.cartID, tc.cartItemID)
			switch {
			case tc.isCustomRemoveItemErrExpected && tc.expectedRemoveItemErr != nil && err != nil:
				assert.Contains(t, err.Error(), tc.expectedRemoveItemErr.Error(), "Actual error should contain text from expected error")
			default:
				assert.Equal(t, tc.expectedRemoveItemErr, errors.Cause(err), "Two errors should be the same")

				actualCart, cartErr := connTest.Cart(context.Background(), tc.cartID)
				assert.Equal(t, tc.expectedCart, actualCart, "Two objects should be the same")
				assert.Equal(t, tc.expectedGetCartErr, errors.Cause(cartErr), "Two errors should be the same")
			}
		})
	}
}
