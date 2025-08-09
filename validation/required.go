package validation

import (
	"fmt"
	"reflect"
)

// RequiredValidator validates that a field is not empty
type RequiredValidator struct{}

func (v *RequiredValidator) Validate(value any) error {
	if value == nil {
		return fmt.Errorf("field is required")
	}
	val := reflect.ValueOf(value)
	if val.IsZero() {
		return fmt.Errorf("field is required")
	}
	return nil
}

// New creates a new RequiredValidator from parameters
func (v *RequiredValidator) New(params map[string]string) (Validator, error) {
	return &RequiredValidator{}, nil
}

// Key returns the registration key for this validator
func (v *RequiredValidator) Key() string {
	return "required"
}
