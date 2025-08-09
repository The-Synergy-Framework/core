package validation

import (
	"fmt"
	"strings"
)

// OneOfValidator validates that a value is one of the allowed values
type OneOfValidator struct {
	AllowedValues []string
}

func (v *OneOfValidator) Validate(value any) error {
	valueStr := fmt.Sprintf("%v", value)

	for _, allowed := range v.AllowedValues {
		if strings.TrimSpace(allowed) == valueStr {
			return nil
		}
	}

	return fmt.Errorf("value must be one of: %s", strings.Join(v.AllowedValues, "|"))
}

// New creates a new OneOfValidator from parameters
func (v *OneOfValidator) New(params map[string]string) (Validator, error) {
	valuesStr := params["values"]
	if valuesStr == "" {
		return nil, fmt.Errorf("oneof validation requires a values parameter")
	}
	allowedValues := strings.Split(valuesStr, "|")
	return &OneOfValidator{AllowedValues: allowedValues}, nil
}

// Key returns the registration key for this validator
func (v *OneOfValidator) Key() string {
	return "oneof"
}
