package data

import (
	"context"
	"fmt"

	protos "github.com/Kama001/microservicesusingGo/8.grpc/currency_converter/protos/currency"
	"k8s.io/klog/v2"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      *klog.Logger
}

func NewProductsDB(cc protos.CurrencyClient, l *klog.Logger) *ProductsDB {
	return &ProductsDB{
		currency: cc,
		log:      l,
	}
}

// GetProducts returns all products from the database
// with price of products in requested currency
func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error(err, fmt.Sprintf("cannot get rate for currency %s", currency))
		return nil, err
	}
	// price converted products list
	convertedPL := []*Product{}
	for _, p := range productList {
		// new product
		// as p is reference to products
		// if we modify p
		// products will be changed
		np := *p
		np.Price *= rate
		convertedPL = append(convertedPL, &np)
	}
	return convertedPL, nil
}

// AddProduct adds a new product to the database
func AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func UpdateProduct(p Product, id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}
	productList[i] = &p
	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}
	productList = append(productList[:i], productList[i+1:]...)
	return nil
}

func (p *ProductsDB) GetProductById(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}
	if currency == "" {
		return productList[i], nil
	}
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error(err, fmt.Sprintf("cannot get rate for currency %s", currency))
		return nil, err
	}
	np := *productList[i]
	np.Price *= rate
	return &np, nil
}

func findIndexByProductID(id int) int {
	for i, product := range productList {
		if product.ID == id {
			return i
		}
	}
	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	rr := protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]).String(),
		Destination: protos.Currencies(protos.Currencies_value[destination]).String(),
	}

	resp, err := p.currency.GetRate(context.Background(), &rr)

	return resp.Rate, err
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
