package validation

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexpValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid email",
			pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email",
			pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
			value:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "valid phone number",
			pattern: `^\d{3}-\d{3}-\d{4}$`,
			value:   "123-456-7890",
			wantErr: false,
		},
		{
			name:    "invalid phone number",
			pattern: `^\d{3}-\d{3}-\d{4}$`,
			value:   "1234567890",
			wantErr: true,
		},
		{
			name:    "valid alphanumeric",
			pattern: `^[a-zA-Z0-9]+$`,
			value:   "abc123",
			wantErr: false,
		},
		{
			name:    "invalid alphanumeric",
			pattern: `^[a-zA-Z0-9]+$`,
			value:   "abc-123",
			wantErr: true,
		},
		{
			name:    "valid date format",
			pattern: `^\d{4}-\d{2}-\d{2}$`,
			value:   "2023-12-25",
			wantErr: false,
		},
		{
			name:    "invalid date format",
			pattern: `^\d{4}-\d{2}-\d{2}$`,
			value:   "2023/12/25",
			wantErr: true,
		},
		{
			name:    "empty string",
			pattern: `^.*$`,
			value:   "",
			wantErr: false,
		},
		{
			name:    "empty string with required pattern",
			pattern: `^.+$`,
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := regexp.Compile(tt.pattern)
			require.NoError(t, err)

			validator := &RegexpValidator{Pattern: pattern}
			err = validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "value does not match pattern:")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRegexpValidator_Validate_NonStringTypes(t *testing.T) {
	pattern := regexp.MustCompile(`^test$`)

	validator := &RegexpValidator{Pattern: pattern}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 123},
		{"float", 123.45},
		{"bool", true},
		{"slice", []string{"test"}},
		{"map", map[string]string{"key": "value"}},
		{"struct", struct{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "regexp validation only applies to strings")
		})
	}
}

func TestRegexpValidator_New(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid pattern",
			params:  map[string]string{"pattern": `^[a-z]+$`},
			wantErr: false,
		},
		{
			name:    "missing pattern parameter",
			params:  map[string]string{},
			wantErr: true,
			errMsg:  "regexp validation requires a pattern parameter",
		},
		{
			name:    "invalid pattern",
			params:  map[string]string{"pattern": `[invalid`},
			wantErr: true,
			errMsg:  "invalid regexp pattern: [invalid",
		},
		{
			name:    "complex pattern",
			params:  map[string]string{"pattern": `^https?://[^\s/$.?#].[^\s]*$`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &RegexpValidator{}
			result, err := validator.New(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				regexpValidator, ok := result.(*RegexpValidator)
				require.True(t, ok)
				assert.NotNil(t, regexpValidator.Pattern)

				// Test that the pattern works
				if tt.params["pattern"] != "" {
					// Test with a simple string that should match
					testValue := "test"
					if tt.params["pattern"] == `^[a-z]+$` {
						err = regexpValidator.Validate(testValue)
						require.NoError(t, err)
					}
				}
			}
		})
	}
}

func TestRegexpValidator_Key(t *testing.T) {
	validator := &RegexpValidator{}
	result := validator.Key()
	assert.Equal(t, "regexp", result)
}
