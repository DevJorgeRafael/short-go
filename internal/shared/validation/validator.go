package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	// Usar los nombres del tag json en lugar de los nombres de campo
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}

		if idx := strings.Index(name, ","); idx != -1 {
			return name[:idx]
		}
		return name
	})

	return v
}