package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinValidator_Validate_NumericTypes(t *testing.T) {
	tests := []struct {
		name    string
		min     float64
		value   interface{}
		wantErr bool
	}{
		{"int valid", 5, 10, false},
		{"int equal", 10, 10, false},
		{"int invalid", 15, 10, true},
		{"int8 valid", 5, int8(10), false},
		{"int8 invalid", 15, int8(10), true},
		{"int16 valid", 5, int16(10), false},
		{"int16 invalid", 15, int16(10), true},
		{"int32 valid", 5, int32(10), false},
		{"int32 invalid", 15, int32(10), true},
		{"int64 valid", 5, int64(10), false},
		{"int64 invalid", 15, int64(10), true},
		{"uint valid", 5, uint(10), false},
		{"uint invalid", 15, uint(10), true},
		{"uint8 valid", 5, uint8(10), false},
		{"uint8 invalid", 15, uint8(10), true},
		{"uint16 valid", 5, uint16(10), false},
		{"uint16 invalid", 15, uint16(10), true},
		{"uint32 valid", 5, uint32(10), false},
		{"uint32 invalid", 15, uint32(10), true},
		{"uint64 valid", 5, uint64(10), false},
		{"uint64 invalid", 15, uint64(10), true},
		{"float32 valid", 5.0, float32(10.0), false},
		{"float32 equal", 10.0, float32(10.0), false},
		{"float32 invalid", 15.0, float32(10.0), true},
		{"float64 valid", 5.0, float64(10.0), false},
		{"float64 equal", 10.0, float64(10.0), false},
		{"float64 invalid", 15.0, float64(10.0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MinValidator{Min: tt.min}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "value must be at least")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMinValidator_Validate_String(t *testing.T) {
	tests := []struct {
		name    string
		min     float64
		value   string
		wantErr bool
	}{
		{"valid length", 3, "hello", false},
		{"equal length", 5, "hello", false},
		{"invalid length", 10, "hi", true},
		{"empty string", 0, "", false},
		{"empty string invalid", 5, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MinValidator{Min: tt.min}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "string length must be at least")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMinValidator_Validate_UnsupportedTypes(t *testing.T) {
	validator := &MinValidator{Min: 10}

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

func TestMinValidator_New(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid parameters",
			params:  map[string]string{"value": "10.5"},
			wantErr: false,
		},
		{
			name:    "missing value parameter",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "min validation requires a value parameter",
		},
		{
			name:    "invalid value parameter",
			params:  map[string]string{"value": "invalid"},
			wantErr: true,
			errMsg:  "invalid min value: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MinValidator{}
			result, err := validator.New(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				minValidator, ok := result.(*MinValidator)
				require.True(t, ok)
				assert.Equal(t, 10.5, minValidator.Min)
			}
		})
	}
}

func TestMinValidator_Key(t *testing.T) {
	validator := &MinValidator{}
	result := validator.Key()
	assert.Equal(t, "min", result)
}
