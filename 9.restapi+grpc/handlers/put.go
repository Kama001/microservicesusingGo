package handlers

import (
	"net/http"
	"restapi/data"
)

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	id := getProductID(r)

	p.log.Info("updating record id", id)
	prod.ID = id
	err := data.UpdateProduct(prod, id)
	if err == data.ErrProductNotFound {
		p.log.Error(err, "[ERROR] product not found")
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}
	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
