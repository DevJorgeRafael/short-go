package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("El campo '%s' es requerido", field)
			case "email":
				return "Formato de email inválido"
			case "min":
				return fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres", field, e.Param())
			}
		}
	}
	return "Datos inválidos"
}