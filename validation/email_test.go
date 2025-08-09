package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailValidator_Validate(t *testing.T) {
	validator := &EmailValidator{}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with subdomain", "user@sub.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "user.name@example.com", false},
		{"valid with underscore", "user_name@example.com", false},
		{"empty string", "", true},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
		{"missing tld", "user@example", true},
		{"invalid chars", "user@example!.com", true},
		{"multiple @", "user@example@com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.email)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid email format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEmailValidator_Validate_NonStringTypes(t *testing.T) {
	validator := &EmailValidator{}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 123},
		{"bool", true},
		{"slice", []string{"test"}},
		{"map", map[string]string{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "email validation only applies to strings")
		})
	}
}

func TestEmailValidator_New(t *testing.T) {
	validator := &EmailValidator{}

	result, err := validator.New(map[string]string{})
	require.NoError(t, err)
	require.NotNil(t, result)

	emailValidator, ok := result.(*EmailValidator)
	require.True(t, ok)
	assert.NotNil(t, emailValidator)
}

func TestEmailValidator_Key(t *testing.T) {
	validator := &EmailValidator{}
	result := validator.Key()
	assert.Equal(t, "email", result)
}
