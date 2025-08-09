package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxValidator_Validate_NumericTypes(t *testing.T) {
	tests := []struct {
		name    string
		max     float64
		value   interface{}
		wantErr bool
	}{
		{"int valid", 10, 5, false},
		{"int equal", 10, 10, false},
		{"int invalid", 10, 15, true},
		{"int8 valid", 10, int8(5), false},
		{"int8 invalid", 10, int8(15), true},
		{"int16 valid", 10, int16(5), false},
		{"int16 invalid", 10, int16(15), true},
		{"int32 valid", 10, int32(5), false},
		{"int32 invalid", 10, int32(15), true},
		{"int64 valid", 10, int64(5), false},
		{"int64 invalid", 10, int64(15), true},
		{"uint valid", 10, uint(5), false},
		{"uint invalid", 10, uint(15), true},
		{"uint8 valid", 10, uint8(5), false},
		{"uint8 invalid", 10, uint8(15), true},
		{"uint16 valid", 10, uint16(5), false},
		{"uint16 invalid", 10, uint16(15), true},
		{"uint32 valid", 10, uint32(5), false},
		{"uint32 invalid", 10, uint32(15), true},
		{"uint64 valid", 10, uint64(5), false},
		{"uint64 invalid", 10, uint64(15), true},
		{"float32 valid", 10.0, float32(5.0), false},
		{"float32 equal", 10.0, float32(10.0), false},
		{"float32 invalid", 10.0, float32(15.0), true},
		{"float64 valid", 10.0, float64(5.0), false},
		{"float64 equal", 10.0, float64(10.0), false},
		{"float64 invalid", 10.0, float64(15.0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MaxValidator{Max: tt.max}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "value must be at most")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMaxValidator_Validate_String(t *testing.T) {
	tests := []struct {
		name    string
		max     float64
		value   string
		wantErr bool
	}{
		{"valid length", 10, "hello", false},
		{"equal length", 5, "hello", false},
		{"invalid length", 3, "hello", true},
		{"empty string", 5, "", false},
		{"zero max", 0, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MaxValidator{Max: tt.max}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "string length must be at most")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMaxValidator_Validate_UnsupportedTypes(t *testing.T) {
	validator := &MaxValidator{Max: 10}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"bool", true},
		{"slice", []int{1, 2, 3}},
		{"map", map[string]int{"a": 1}},
		{"struct", struct{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			// Should not error, just return nil for unsupported types
			require.NoError(t, err)
		})
	}
}

func TestMaxValidator_New(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]string
		wantErr       bool
		errMsg        string
		expectedValue interface{}
		expectedField string
	}{
		{
			name:          "valid parameters",
			params:        map[string]string{"value": "10.5"},
			wantErr:       false,
			expectedValue: 10.5,
			expectedField: "Max",
		},
		{
			name:    "missing value parameter",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "max validation requires a value parameter",
		},
		{
			name:    "invalid value parameter",
			params:  map[string]string{"value": "invalid"},
			wantErr: true,
			errMsg:  "invalid max value: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MaxValidator{}

			if tt.wantErr {
				testValidatorNewError(t, validator, tt.params, tt.errMsg)
			} else {
				testValidatorNew(t, validator, tt.params, tt.expectedValue, tt.expectedField)
			}
		})
	}
}

func TestMaxValidator_Key(t *testing.T) {
	validator := &MaxValidator{}
	result := validator.Key()
	assert.Equal(t, "max", result)
}
