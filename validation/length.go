package validation

import (
	"fmt"
	"reflect"
	"strconv"
)

// LenValidator validates exact length
type LenValidator struct {
	ExpectedLen int
}

func (v *LenValidator) Validate(value any) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		if len(val.String()) != v.ExpectedLen {
			return fmt.Errorf("string length must be exactly %d", v.ExpectedLen)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() != v.ExpectedLen {
			return fmt.Errorf("slice length must be exactly %d", v.ExpectedLen)
		}
	}
	return nil
}

// New creates a new LenValidator from parameters
func (v *LenValidator) New(params map[string]string) (Validator, error) {
	lenStr := params["value"]
	if lenStr == "" {
		return nil, fmt.Errorf("len validation requires a value parameter")
	}
	expectedLen, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid len value: %s", lenStr)
	}
	return &LenValidator{ExpectedLen: expectedLen}, nil
}

// Key returns the registration key for this validator
func (v *LenValidator) Key() string {
	return "len"
}
