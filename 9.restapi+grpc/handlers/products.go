package handlers

import (
	"net/http"
	"restapi/data"
	"strconv"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
)

type Products struct {
	log *klog.Logger
	pdb *data.ProductsDB
}

func NewProducts(l *klog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, pdb}
}

type KeyProduct struct{}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	ids := mux.Vars(r)
	id, err := strconv.Atoi(ids["id"])
	if err != nil {
		panic(err)
	}
	return id
}
