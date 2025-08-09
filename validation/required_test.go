package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequiredValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"non-empty string", "test", false},
		{"non-empty int", 123, false},
		{"non-empty float", 123.45, false},
		{"non-empty bool", true, false},
		{"non-empty slice", []string{"test"}, false},
		{"non-empty map", map[string]string{"key": "value"}, false},
		{"non-empty struct", struct{ Name string }{"test"}, false},
		{"empty string", "", true},
		{"zero int", 0, true},
		{"zero float", 0.0, true},
		{"false bool", false, true},
		{"nil slice", []string(nil), true},
		{"empty slice", []string{}, false},
		{"nil map", map[string]string(nil), true},
		{"empty map", map[string]string{}, false},
		{"nil interface", nil, true},
		{"zero struct", struct{}{}, true},
		{"zero time", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &RequiredValidator{}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "field is required")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRequiredValidator_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"space string", " ", false},
		{"tab string", "\t", false},
		{"newline string", "\n", false},
		{"unicode string", "测试", false},
		{"special chars", "!@#$%^&*()", false},
		{"negative int", -1, false},
		{"negative float", -123.45, false},
		{"small positive float", 0.001, false},
		{"large positive float", 999999.999, false},
		{"pointer to zero", &struct{}{}, false},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &RequiredValidator{}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "field is required")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRequiredValidator_New(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
	}{
		{
			name:    "empty params",
			params:  map[string]string{},
			wantErr: false,
		},
		{
			name:    "with params",
			params:  map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name:    "nil params",
			params:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &RequiredValidator{}
			result, err := validator.New(tt.params)

			require.NoError(t, err)
			require.NotNil(t, result)

			requiredValidator, ok := result.(*RequiredValidator)
			require.True(t, ok)
			assert.NotNil(t, requiredValidator)
		})
	}
}

func TestRequiredValidator_Key(t *testing.T) {
	validator := &RequiredValidator{}
	result := validator.Key()
	assert.Equal(t, "required", result)
}
