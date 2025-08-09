package validation

import (
	"fmt"
	"reflect"
	"regexp"
)

// RegexpValidator validates string against a regex pattern
type RegexpValidator struct {
	Pattern *regexp.Regexp
}

func (v *RegexpValidator) Validate(value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.String {
		return fmt.Errorf("regexp validation only applies to strings")
	}

	if !v.Pattern.MatchString(val.String()) {
		return fmt.Errorf("value does not match pattern: %s", v.Pattern.String())
	}
	return nil
}

// New creates a new RegexpValidator from parameters
func (v *RegexpValidator) New(params map[string]string) (Validator, error) {
	pattern := params["pattern"]
	if pattern == "" {
		return nil, fmt.Errorf("regexp validation requires a pattern parameter")
	}
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp pattern: %s", pattern)
	}
	return &RegexpValidator{Pattern: regex}, nil
}

// Key returns the registration key for this validator
func (v *RegexpValidator) Key() string {
	return "regexp"
}
