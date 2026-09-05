package validation

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	Validate = validator.New()
	skuRegex = regexp.MustCompile(`^[A-Z0-9_-]+$`)
)

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	if err := Validate.RegisterValidation("sku", validateSKU); err != nil {
		panic(err)
	}
}

func validateSKU(fl validator.FieldLevel) bool {
	sku := fl.Field().String()

	if len(sku) < 2 || len(sku) > 64 {
		return false
	}
	return skuRegex.MatchString(sku)
}
