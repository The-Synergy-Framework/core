package validation

import (
	"fmt"
	"reflect"
	"regexp"
)

// EmailValidator validates email format
type EmailValidator struct{}

func (v *EmailValidator) Validate(value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.String {
		return fmt.Errorf("email validation only applies to strings")
	}

	email := val.String()
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// New creates a new EmailValidator from parameters
func (v *EmailValidator) New(params map[string]string) (Validator, error) {
	// Email validator doesn't need any parameters
	return &EmailValidator{}, nil
}

// Key returns the registration key for this validator
func (v *EmailValidator) Key() string {
	return "email"
}
