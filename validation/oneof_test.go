package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneOfValidator_Validate(t *testing.T) {
	tests := []struct {
		name          string
		allowedValues []string
		value         interface{}
		wantErr       bool
	}{
		{
			name:          "valid string",
			allowedValues: []string{"red", "green", "blue"},
			value:         "red",
			wantErr:       false,
		},
		{
			name:          "valid int",
			allowedValues: []string{"1", "2", "3"},
			value:         1,
			wantErr:       false,
		},
		{
			name:          "valid float",
			allowedValues: []string{"1.5", "2.5", "3.5"},
			value:         2.5,
			wantErr:       false,
		},
		{
			name:          "valid bool",
			allowedValues: []string{"true", "false"},
			value:         true,
			wantErr:       false,
		},
		{
			name:          "invalid value",
			allowedValues: []string{"red", "green", "blue"},
			value:         "yellow",
			wantErr:       true,
		},
		{
			name:          "case sensitive",
			allowedValues: []string{"Red", "Green", "Blue"},
			value:         "red",
			wantErr:       true,
		},
		{
			name:          "whitespace handling",
			allowedValues: []string{" red ", " green ", " blue "},
			value:         "red",
			wantErr:       false,
		},
		{
			name:          "empty allowed values",
			allowedValues: []string{},
			value:         "test",
			wantErr:       true,
		},
		{
			name:          "single allowed value",
			allowedValues: []string{"test"},
			value:         "test",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &OneOfValidator{AllowedValues: tt.allowedValues}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "value must be one of:")
				// Check that all allowed values are in the error message
				for _, allowed := range tt.allowedValues {
					assert.Contains(t, err.Error(), strings.TrimSpace(allowed))
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOneOfValidator_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		allowedValues []string
		value         interface{}
		wantErr       bool
	}{
		{
			name:          "empty string value",
			allowedValues: []string{"", "test"},
			value:         "",
			wantErr:       false,
		},
		{
			name:          "zero value",
			allowedValues: []string{"0", "1", "2"},
			value:         0,
			wantErr:       false,
		},
		{
			name:          "negative value",
			allowedValues: []string{"-1", "0", "1"},
			value:         -1,
			wantErr:       false,
		},
		{
			name:          "complex string",
			allowedValues: []string{"hello world", "test"},
			value:         "hello world",
			wantErr:       false,
		},
		{
			name:          "special characters",
			allowedValues: []string{"test@example.com", "user+tag@domain.com"},
			value:         "test@example.com",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &OneOfValidator{AllowedValues: tt.allowedValues}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "value must be one of:")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOneOfValidator_New(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid parameters",
			params:  map[string]string{"values": "red|green|blue"},
			wantErr: false,
		},
		{
			name:    "single value",
			params:  map[string]string{"values": "test"},
			wantErr: false,
		},
		{
			name:    "empty values",
			params:  map[string]string{"values": ""},
			wantErr: true,
			errMsg:  "oneof validation requires a values parameter",
		},
		{
			name:    "missing values parameter",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "oneof validation requires a values parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &OneOfValidator{}
			result, err := validator.New(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				oneOfValidator, ok := result.(*OneOfValidator)
				require.True(t, ok)
				assert.NotNil(t, oneOfValidator.AllowedValues)

				// Verify the values were parsed correctly
				if tt.params["values"] != "" {
					expectedValues := strings.Split(tt.params["values"], "|")
					assert.Equal(t, expectedValues, oneOfValidator.AllowedValues)
				}
			}
		})
	}
}

func TestOneOfValidator_Key(t *testing.T) {
	validator := &OneOfValidator{}
	result := validator.Key()
	assert.Equal(t, "oneof", result)
}
