package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

// MaxValidator validates maximum values
type MaxValidator struct {
	Max float64
}

func (v *MaxValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if float64(val.Int()) > v.Max {
			return fmt.Errorf("value must be at most %v", v.Max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if float64(val.Uint()) > v.Max {
			return fmt.Errorf("value must be at most %v", v.Max)
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() > v.Max {
			return fmt.Errorf("value must be at most %v", v.Max)
		}
	case reflect.String:
		if float64(len(val.String())) > v.Max {
			return fmt.Errorf("string length must be at most %v", v.Max)
		}
	}
	return nil
}

// New creates a new MaxValidator from parameters
func (v *MaxValidator) New(params map[string]string) (Validator, error) {
	maxStr := params["value"]
	if maxStr == "" {
		return nil, fmt.Errorf("max validation requires a value parameter")
	}
	maxValue, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid max value: %s", maxStr)
	}
	return &MaxValidator{Max: maxValue}, nil
}

// Key returns the registration key for this validator
func (v *MaxValidator) Key() string {
	return "max"
}
