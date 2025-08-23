package data

import (
	"testing"
)

func TestProductValidation(t *testing.T) {
	product := &Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		SKU:         "test-sku-sku",
	}

	err := product.ValidateJSON()
	if err != nil {
		t.Errorf("Product validation failed: %v", err)
	}

	err = product.ValidateJSON()
	if err != nil {
		t.Errorf("Product validation failed: %v", err)
	}
}
