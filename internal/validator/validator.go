// Package validator for validator
package validator

import (
	"errors"

	"github.com/go-playground/validator"
)

// Validator wraps the go playground validator for the echo framework interface.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

// Validate implements the echo framework validator interface.
func (val *Validator) Validate(i interface{}) error {
	err := val.validator.Struct(i)
	if err == nil {
		return nil
	}
	err = errors.New("some validation error")
	return err
}
