package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restapi/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Failed to marshal products", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle POST Product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)
	id, err := strconv.Atoi(ids["id"])
	if err != nil {
		http.Error(rw, "Unable to get id", http.StatusBadRequest)
		return
	}
	fmt.Println("Handle PUT Product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = prod.UpdateProduct(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var prod data.Product
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to UnMarshal", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
