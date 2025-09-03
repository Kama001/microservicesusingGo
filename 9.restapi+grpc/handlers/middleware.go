package handlers

import (
	"context"
	"net/http"
	"restapi/data"
)

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var prod data.Product
		err := data.FromJSON(&prod, r.Body)
		if err != nil {
			p.log.Error(err, "deserializing product")
			http.Error(rw, "Unable to UnMarshal", http.StatusBadRequest)
			return
		}
		if err = prod.ValidateJSON(); err != nil {
			p.log.Error(err, "validating product")
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
