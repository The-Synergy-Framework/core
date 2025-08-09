package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLenValidator_Validate_String(t *testing.T) {
	tests := []struct {
		name        string
		expectedLen int
		value       string
		wantErr     bool
	}{
		{"exact length", 5, "hello", false},
		{"too short", 5, "hi", true},
		{"too long", 5, "hello world", true},
		{"empty string", 0, "", false},
		{"zero length", 0, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &LenValidator{ExpectedLen: tt.expectedLen}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), fmt.Sprintf("string length must be exactly %d", tt.expectedLen))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLenValidator_Validate_Slice(t *testing.T) {
	tests := []struct {
		name        string
		expectedLen int
		value       interface{}
		wantErr     bool
	}{
		{"exact length", 3, []int{1, 2, 3}, false},
		{"too short", 3, []int{1, 2}, true},
		{"too long", 3, []int{1, 2, 3, 4}, true},
		{"empty slice", 0, []string{}, false},
		{"zero length", 0, []int{1, 2, 3}, true},
		{"string slice", 2, []string{"a", "b"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &LenValidator{ExpectedLen: tt.expectedLen}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), fmt.Sprintf("slice length must be exactly %d", tt.expectedLen))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLenValidator_Validate_Array(t *testing.T) {
	tests := []struct {
		name        string
		expectedLen int
		value       interface{}
		wantErr     bool
	}{
		{"exact length", 3, [3]int{1, 2, 3}, false},
		{"wrong length", 3, [2]int{1, 2}, true},
		{"empty array", 0, [0]int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &LenValidator{ExpectedLen: tt.expectedLen}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), fmt.Sprintf("slice length must be exactly %d", tt.expectedLen))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLenValidator_Validate_UnsupportedTypes(t *testing.T) {
	validator := &LenValidator{ExpectedLen: 5}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 123},
		{"bool", true},
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

func TestLenValidator_New(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid parameters",
			params:  map[string]string{"value": "10"},
			wantErr: false,
		},
		{
			name:    "missing value parameter",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "len validation requires a value parameter",
		},
		{
			name:    "invalid value parameter",
			params:  map[string]string{"value": "invalid"},
			wantErr: true,
			errMsg:  "invalid len value: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &LenValidator{}
			result, err := validator.New(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				lenValidator, ok := result.(*LenValidator)
				require.True(t, ok)
				assert.Equal(t, 10, lenValidator.ExpectedLen)
			}
		})
	}
}

func TestLenValidator_Key(t *testing.T) {
	validator := &LenValidator{}
	result := validator.Key()
	assert.Equal(t, "len", result)
}
