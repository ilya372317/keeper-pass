package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var requiredValidator *validator.Validate

func ValidateDTOWithRequired(dto interface{}) error {
	if requiredValidator == nil {
		requiredValidator = validator.New(validator.WithRequiredStructEnabled())
	}

	if err := requiredValidator.Struct(dto); err != nil {
		return fmt.Errorf("invalid struct given: %w", err)
	}

	return nil
}
