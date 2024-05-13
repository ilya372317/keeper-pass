package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func CardExpValidation(fl validator.FieldLevel) bool {
	expDateStr := fl.Field().String()
	if _, err := time.Parse("01/06", expDateStr); err != nil {
		return false
	}

	return true
}
