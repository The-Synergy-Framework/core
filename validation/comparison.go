package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

// ComparisonValidator validates numeric comparisons
type ComparisonValidator struct {
	Operator     string
	CompareValue float64
}

func (v *ComparisonValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	var actualValue float64
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualValue = float64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualValue = float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		actualValue = val.Float()
	default:
		return fmt.Errorf("comparison validation only applies to numeric types")
	}

	var isValid bool
	switch v.Operator {
	case ">=":
		isValid = actualValue >= v.CompareValue
	case "<=":
		isValid = actualValue <= v.CompareValue
	case ">":
		isValid = actualValue > v.CompareValue
	case "<":
		isValid = actualValue < v.CompareValue
	}

	if !isValid {
		return fmt.Errorf("value must be %s %v", v.Operator, v.CompareValue)
	}
	return nil
}

// New creates a new ComparisonValidator from parameters
func (v *ComparisonValidator) New(params map[string]string) (Validator, error) {
	valueStr := params["value"]
	if valueStr == "" {
		return nil, fmt.Errorf("comparison validation requires a value parameter")
	}
	compareValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid comparison value: %s", valueStr)
	}
	return &ComparisonValidator{Operator: v.Operator, CompareValue: compareValue}, nil
}

// Key returns the registration key for this validator
func (v *ComparisonValidator) Key() string {
	return v.Operator
}
