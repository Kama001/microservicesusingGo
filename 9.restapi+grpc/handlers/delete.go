package handlers

import (
	"net/http"
	"restapi/data"
)

// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.log.Info("deleting record id", id)
	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.log.Info("[INFO] deleting record id does not exist", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	if err != nil {
		p.log.Error(err, "[ERROR] deleting record")
		http.Error(rw, "Product not deleted", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
