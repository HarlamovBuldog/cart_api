package mongo

import (
	"context"
	"testing"

	"github.com/HarlamovBuldog/cart_api/pkg/service"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCart(t *testing.T) {
	tt := []struct {
		name          string
		initColParams initCollectionParams
		expectedErr   error
	}{
		{
			name: "correct test",
			initColParams: initCollectionParams{
				CollectionName: cartsCollectionName,
				Documents: []interface{}{
					service.Cart{
						Items: []service.CartItem{},
					},
					service.Cart{
						Items: []service.CartItem{},
					},
				},
				Opts: nil,
			},
			expectedErr: nil,
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

			expectedUser, actualErr := connTest.AddCart(context.Background())
			assert.Equal(t, tc.expectedErr, errors.Cause(actualErr), "Two errors should be the same")
			actualUser, err := connTest.Cart(context.Background(), expectedUser.ID.Hex())
			assert.Equal(t, expectedUser, actualUser)
		})
	}
}

func TestCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(3)
	tt := []struct {
		name                string
		initColParams       initCollectionParams
		id                  string
		expectedCart        *service.Cart
		isCustomErrExpected bool
		expectedErr         error
	}{
		{
			name: "correct test",
			id:   cartObjIDSet[0].Hex(),
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
			expectedCart: &service.Cart{
				ID:    cartObjIDSet[0],
				Items: []service.CartItem{},
			},
			isCustomErrExpected: false,
			expectedErr:         nil,
		},
		{
			name: "incorrect test: ErrNotFound",
			id:   cartObjIDSet[2].Hex(),
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
			expectedCart:        nil,
			isCustomErrExpected: false,
			expectedErr:         ErrNotFound,
		},
		{
			name: "incorrect test: bad id provided",
			id:   "bad_id",
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
			expectedCart:        nil,
			isCustomErrExpected: true,
			expectedErr:         errors.New("could not convert"),
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

			actualCart, actualErr := connTest.Cart(context.Background(), tc.id)
			switch {
			case tc.isCustomErrExpected && tc.expectedErr != nil && actualErr != nil:
				assert.Contains(t, actualErr.Error(), tc.expectedErr.Error(), "Actual error should contain text from expected error")
			default:
				assert.Equal(t, tc.expectedErr, errors.Cause(actualErr), "Two errors should be the same")
			}
			assert.Equal(t, tc.expectedCart, actualCart, "Two objects should be the same")
		})
	}
}
