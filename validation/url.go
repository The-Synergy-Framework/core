package validation

import (
	"fmt"
	"reflect"
	"regexp"
)

// URLValidator validates URL format
type URLValidator struct{}

func (v *URLValidator) Validate(value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.String {
		return fmt.Errorf("url validation only applies to strings")
	}

	url := val.String()
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].\S*$`)
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}

// New creates a new URLValidator from parameters
func (v *URLValidator) New(params map[string]string) (Validator, error) {
	return &URLValidator{}, nil
}

// Key returns the registration key for this validator
func (v *URLValidator) Key() string {
	return "url"
}
