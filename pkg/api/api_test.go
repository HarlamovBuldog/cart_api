package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HarlamovBuldog/cart_api/pkg/mocks"
	"github.com/HarlamovBuldog/cart_api/pkg/mongo"
	"github.com/HarlamovBuldog/cart_api/pkg/service"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_createCart(t *testing.T) {
	type addCartOut struct {
		cart *service.Cart
		err  error
	}
	cartObjIDSet := generatePrimObjIDSet(1)
	tt := []struct {
		name             string
		method           string
		request          string
		expectedResponse string
		expectedStatus   int
		contentType      string
		addCrtOut        *addCartOut
	}{
		{
			name:             "correct test",
			method:           http.MethodPost,
			request:          `{}`,
			expectedResponse: fmt.Sprintf(`{"id":"%s","items":[]}`, cartObjIDSet[0].Hex()),
			expectedStatus:   http.StatusOK,
			addCrtOut: &addCartOut{
				cart: &service.Cart{
					ID:    cartObjIDSet[0],
					Items: []service.CartItem{},
				},
				err: nil,
			},
		},
		{
			name:           "incorrect method",
			method:         http.MethodPatch,
			request:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:             "db error",
			method:           http.MethodPost,
			request:          `{}`,
			expectedResponse: "could not add cart: no carts: not found",
			expectedStatus:   http.StatusInternalServerError,
			addCrtOut: &addCartOut{
				cart: nil,
				err:  errors.Wrap(mongo.ErrNotFound, "no carts"),
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockService(ctrl)
	s := New(mock)

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			if tc.addCrtOut != nil {
				mock.EXPECT().AddCart(gomock.Any()).Times(1).Return(tc.addCrtOut.cart, tc.addCrtOut.err)
			}
			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/carts", server.URL), strings.NewReader(tc.request))
			require.NoError(t, err, "could not create request")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "could not get response")
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Two status codes should be the same")
			if tc.expectedStatus == http.StatusOK {
				respBody := string(bytes.TrimSpace(b))
				assert.Equal(t, tc.expectedResponse, respBody, "To response bodies should be the same")
			}
		})
	}
}

func Test_addToCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(1)
	cartItemObjIDSet := generatePrimObjIDSet(1)
	type addToCartIn struct {
		cartID      string
		productName string
		quantity    float64
	}
	type addToCartOut struct {
		cartItem *service.CartItem
		err      error
	}
	tt := []struct {
		name             string
		method           string
		request          string
		expectedResponse string
		expectedStatus   int
		reqCartID        string
		contentType      string
		addToCrtIn       *addToCartIn
		addToCrtOut      *addToCartOut
	}{
		{
			name:    "correct test",
			method:  http.MethodPost,
			request: `{"product":"product_1", "quantity":10.0}`,
			expectedResponse: fmt.Sprintf(`{"id":"%s","cart_id":"%s","product":"product_1","quantity":10}`,
				cartItemObjIDSet[0].Hex(), cartObjIDSet[0].Hex()),
			expectedStatus: http.StatusOK,
			reqCartID:      cartObjIDSet[0].Hex(),
			addToCrtIn: &addToCartIn{
				cartID:      cartObjIDSet[0].Hex(),
				productName: "product_1",
				quantity:    10.0,
			},
			addToCrtOut: &addToCartOut{
				cartItem: &service.CartItem{
					ID:          cartItemObjIDSet[0],
					CartID:      cartObjIDSet[0],
					ProductName: "product_1",
					Quantity:    10.0,
				},
				err: nil,
			},
		},
		{
			name:           "incorrect method",
			method:         http.MethodPatch,
			reqCartID:      cartObjIDSet[0].Hex(),
			request:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:             "bad request body: could not decode",
			method:           http.MethodPost,
			reqCartID:        cartObjIDSet[0].Hex(),
			request:          `{11:"product_1", 22:10.0}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "could not decode request body: invalid character '1' looking for beginning of object key string",
		},
		{
			name:             "data from request body is not valid",
			method:           http.MethodPost,
			reqCartID:        cartObjIDSet[0].Hex(),
			request:          `{"product":"", "quantity":-10.0}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "data from request body is not valid",
		},
		{
			name:             "db error",
			method:           http.MethodPost,
			request:          `{"product":"product_1", "quantity":10.0}`,
			reqCartID:        cartObjIDSet[0].Hex(),
			expectedResponse: "could not add item to cart: no carts: not found",
			expectedStatus:   http.StatusInternalServerError,
			addToCrtIn: &addToCartIn{
				cartID:      cartObjIDSet[0].Hex(),
				productName: "product_1",
				quantity:    10.0,
			},
			addToCrtOut: &addToCartOut{
				cartItem: nil,
				err:      errors.Wrap(mongo.ErrNotFound, "no carts"),
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockService(ctrl)
	s := New(mock)

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.addToCrtOut != nil {
				mock.EXPECT().AddItemToCart(gomock.Any(), tc.addToCrtIn.cartID, tc.addToCrtIn.productName, tc.addToCrtIn.quantity).
					Times(1).Return(tc.addToCrtOut.cartItem, tc.addToCrtOut.err)
			}
			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/carts/%s/items", server.URL, tc.reqCartID), strings.NewReader(tc.request))
			require.NoError(t, err, "could not create request")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "could not get response")
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Two status codes should be the same")
			respBody := string(bytes.TrimSpace(b))
			assert.Equal(t, tc.expectedResponse, respBody, "Two response bodies should be the same")
		})
	}
}

func Test_removeFromCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(1)
	itemObjIDSet := generatePrimObjIDSet(1)
	type removeFromCartIn struct {
		cartID string
		itemID string
	}
	tt := []struct {
		name                string
		method              string
		request             string
		requestCartID       string
		requestItemID       string
		expectedResponse    string
		expectedStatus      int
		contentType         string
		removeFromCrtIn     *removeFromCartIn
		removeFromCrtOutErr error
	}{
		{
			name:             "correct test",
			method:           http.MethodDelete,
			request:          `{}`,
			requestCartID:    cartObjIDSet[0].Hex(),
			requestItemID:    itemObjIDSet[0].Hex(),
			expectedResponse: "",
			expectedStatus:   http.StatusOK,
			removeFromCrtIn: &removeFromCartIn{
				cartID: cartObjIDSet[0].Hex(),
				itemID: itemObjIDSet[0].Hex(),
			},
			removeFromCrtOutErr: nil,
		},
		{
			name:           "incorrect method",
			method:         http.MethodPatch,
			request:        `{}`,
			requestCartID:  cartObjIDSet[0].Hex(),
			requestItemID:  itemObjIDSet[0].Hex(),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:             "db error",
			method:           http.MethodDelete,
			request:          `{}`,
			requestCartID:    cartObjIDSet[0].Hex(),
			requestItemID:    itemObjIDSet[0].Hex(),
			expectedResponse: "could not remove item from cart: no carts: not found",
			expectedStatus:   http.StatusInternalServerError,
			removeFromCrtIn: &removeFromCartIn{
				cartID: cartObjIDSet[0].Hex(),
				itemID: itemObjIDSet[0].Hex(),
			},
			removeFromCrtOutErr: errors.Wrap(mongo.ErrNotFound, "no carts"),
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockService(ctrl)
	s := New(mock)

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.removeFromCrtIn != nil {
				mock.EXPECT().RemoveItemFromCart(gomock.Any(), tc.removeFromCrtIn.cartID, tc.removeFromCrtIn.itemID).Times(1).Return(tc.removeFromCrtOutErr)
			}
			req, err := http.NewRequest(
				tc.method,
				fmt.Sprintf("%s/carts/%s/items/%s", server.URL, tc.requestCartID, tc.requestItemID),
				strings.NewReader(tc.request))
			require.NoError(t, err, "could not create request")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "could not get response")
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Two status codes should be the same")
			respBody := string(bytes.TrimSpace(b))
			assert.Equal(t, tc.expectedResponse, respBody, "To response bodies should be the same")
		})
	}
}

func Test_viewCart(t *testing.T) {
	cartObjIDSet := generatePrimObjIDSet(1)
	itemObjIDSet := generatePrimObjIDSet(3)
	type viewCartOut struct {
		cart *service.Cart
		err  error
	}
	tt := []struct {
		name             string
		method           string
		request          string
		requestCartID    string
		expectedResponse string
		expectedStatus   int
		contentType      string
		viewCrtIn        string
		viewCrtOut       *viewCartOut
	}{
		{
			name:          "correct test",
			method:        http.MethodGet,
			request:       `{}`,
			requestCartID: cartObjIDSet[0].Hex(),
			expectedResponse: fmt.Sprintf(`{"id":"%[1]s","items":[{"id":"%[2]s","cart_id":"%[1]s","product":"product_1","quantity":10},`+
				`{"id":"%[3]s","cart_id":"%[1]s","product":"product_2","quantity":15}]}`,
				cartObjIDSet[0].Hex(), itemObjIDSet[0].Hex(), itemObjIDSet[1].Hex()),
			expectedStatus: http.StatusOK,
			viewCrtIn:      cartObjIDSet[0].Hex(),
			viewCrtOut: &viewCartOut{
				cart: &service.Cart{
					ID: cartObjIDSet[0],
					Items: []service.CartItem{
						{
							ID:          itemObjIDSet[0],
							CartID:      cartObjIDSet[0],
							ProductName: "product_1",
							Quantity:    10,
						},
						{
							ID:          itemObjIDSet[1],
							CartID:      cartObjIDSet[0],
							ProductName: "product_2",
							Quantity:    15,
						},
					},
				},
				err: nil,
			},
		},
		{
			name:           "incorrect method",
			method:         http.MethodPatch,
			request:        `{}`,
			requestCartID:  cartObjIDSet[0].Hex(),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:             "db error",
			method:           http.MethodGet,
			request:          `{}`,
			requestCartID:    cartObjIDSet[0].Hex(),
			expectedResponse: "could not get cart: no carts: not found",
			expectedStatus:   http.StatusInternalServerError,
			viewCrtIn:        cartObjIDSet[0].Hex(),
			viewCrtOut: &viewCartOut{
				cart: nil,
				err:  errors.Wrap(mongo.ErrNotFound, "no carts"),
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockService(ctrl)
	s := New(mock)

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.viewCrtOut != nil {
				mock.EXPECT().Cart(gomock.Any(), tc.viewCrtIn).Times(1).Return(tc.viewCrtOut.cart, tc.viewCrtOut.err)
			}
			req, err := http.NewRequest(
				tc.method,
				fmt.Sprintf("%s/carts/%s", server.URL, tc.requestCartID),
				strings.NewReader(tc.request))
			require.NoError(t, err, "could not create request")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "could not get response")
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Two status codes should be the same")
			respBody := string(bytes.TrimSpace(b))
			assert.Equal(t, tc.expectedResponse, respBody, "To response bodies should be the same")
		})
	}
}

func generatePrimObjIDSet(n int) []primitive.ObjectID {
	primObjIDSet := make([]primitive.ObjectID, n)
	for i := 0; i < n; i++ {
		primObjIDSet[i] = primitive.NewObjectID()
	}

	return primObjIDSet
}
