// reference to validation
// https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go

package data

import (
	"regexp"

	"github.com/go-playground/validator"
)

func validateSKU(fl validator.FieldLevel) bool {
	ptm := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	return ptm.MatchString(fl.Field().String())
}

func (p *Product) ValidateJSON() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}
