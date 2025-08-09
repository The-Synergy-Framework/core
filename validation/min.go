package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

// MinValidator validates minimum values
type MinValidator struct {
	Min float64
}

func (v *MinValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(val.Int()) < v.Min {
			return fmt.Errorf("value must be at least %v", v.Min)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(val.Uint()) < v.Min {
			return fmt.Errorf("value must be at least %v", v.Min)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() < v.Min {
			return fmt.Errorf("value must be at least %v", v.Min)
		}
	case reflect.String:
		if float64(len(val.String())) < v.Min {
			return fmt.Errorf("string length must be at least %v", v.Min)
		}
	}
	return nil
}

// New creates a new MinValidator from parameters
func (v *MinValidator) New(params map[string]string) (Validator, error) {
	minStr := params["value"]
	if minStr == "" {
		return nil, fmt.Errorf("min validation requires a value parameter")
	}
	minValue, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid min value: %s", minStr)
	}
	return &MinValidator{Min: minValue}, nil
}

// Key returns the registration key for this validator
func (v *MinValidator) Key() string {
	return "min"
}
