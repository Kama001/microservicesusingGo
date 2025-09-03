package handlers

import (
	"net/http"
	"restapi/data"
)

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.log.Info("get all records")
	cur := r.URL.Query().Get("currency")
	prods, err := p.pdb.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.log.Error(err, "Unable to Marshal products")
	}
}

func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.log.Info("Details of record id", id)
	cur := r.URL.Query().Get("currency")
	prod, err := p.pdb.GetProductById(id, cur)
	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.log.Error(err, "fetching product")
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found"}, rw)
		return
	default:
		p.log.Error(err, "fetching product")

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// prod.Price *= resp.Rate

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.log.Error(err, "Unable to Marshal products")
	}
}
