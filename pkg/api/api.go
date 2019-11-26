package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HarlamovBuldog/cart_api/pkg/service"

	"github.com/gorilla/mux"
)

// Server contains http handler and service interface with database interaction futures.
type Server struct {
	http.Handler
	service service.Service
}

type newItem struct {
	ProductName string  `json:"product"`
	Quantity    float64 `json:"quantity"`
}

// New initializes new api with router and entrypoints.
func New(db service.Service) *Server {
	router := mux.NewRouter()
	s := Server{
		service: db,
		Handler: router,
	}
	router.HandleFunc("/carts", s.createCart).Methods("POST")
	router.HandleFunc("/carts/{cart_id}/items", s.addToCart).Methods("POST")
	router.HandleFunc("/carts/{cart_id}/items/{item_id}", s.removeFromCart).Methods("DELETE")
	router.HandleFunc("/carts/{cart_id}", s.viewCart).Methods("GET")

	return &s
}

func (s *Server) createCart(w http.ResponseWriter, req *http.Request) {
	cart, err := s.service.AddCart(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not add cart: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not encode json: %s", err)
		return
	}
}

func (s *Server) addToCart(w http.ResponseWriter, req *http.Request) {
	var item newItem
	err := json.NewDecoder(req.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not decode request body: %s", err)
		return
	}

	vars := mux.Vars(req)
	cartID, ok := vars["cart_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "cart_id is not provided")
		return
	}
	if valid := isNewItemDataValid(item); !valid {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "data from request body is not valid")
		return
	}

	cartItem, err := s.service.AddItemToCart(req.Context(), cartID, item.ProductName, item.Quantity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not add item to cart: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode(cartItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not encode json: %s", err)
		return
	}
}

func (s *Server) removeFromCart(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cartID, ok := vars["cart_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "cart_id is not provided")
		return
	}
	itemID, ok := vars["item_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "item_id is not provided")
		return
	}

	err := s.service.RemoveItemFromCart(req.Context(), cartID, itemID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not remove item from cart: %s", err)
		return
	}
}

func (s *Server) viewCart(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cartID, ok := vars["cart_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "cart_id is not provided")
		return
	}

	cart, err := s.service.Cart(req.Context(), cartID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not get cart: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not encode json: %s", err)
		return
	}
}

func isNewItemDataValid(item newItem) bool {
	switch {
	case item.ProductName == "" || item.Quantity <= 0:
		return false
	default:
		return true
	}
}
