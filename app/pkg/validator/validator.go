package validator

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

// Validate validates a struct and returns validation errors.
func Validate(s any) error {
	return v.Struct(s)
}
