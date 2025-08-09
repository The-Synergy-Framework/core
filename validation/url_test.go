package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid http url", "http://example.com", false},
		{"valid https url", "https://example.com", false},
		{"valid url with path", "https://example.com/path", false},
		{"valid url with query", "https://example.com?param=value", false},
		{"valid url with fragment", "https://example.com#section", false},
		{"valid url with port", "https://example.com:8080", false},
		{"valid url with subdomain", "https://sub.example.com", false},
		{"valid url with multiple subdomains", "https://a.b.c.example.com", false},
		{"valid url with ip address", "https://192.168.1.1", false},
		{"valid url with localhost", "http://localhost:3000", false},
		{"valid url with special chars", "https://example.com/path-with-spaces", false},
		{"invalid url missing protocol", "example.com", true},
		{"invalid url wrong protocol", "ftp://example.com", true},
		{"invalid url empty", "", true},
		{"invalid url just protocol", "https://", true},
		{"invalid url with spaces in host", "https://ex ample.com", true},
		{"invalid url with spaces before protocol", " https://example.com", true},
		{"invalid url with spaces after protocol", "https:// example.com", true},
		{"invalid url with invalid chars", "https://example.com/\n", true},
		{"invalid url with control chars", "https://example.com/\t", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &URLValidator{}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid URL format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestURLValidator_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"url with unicode", "https://example.com/test", false},
		{"url with encoded chars", "https://example.com/path%20with%20spaces", false},
		{"url with complex query", "https://example.com?param1=value1&param2=value2", false},
		{"url with hash", "https://example.com/page#section", false},
		{"url with user info", "https://user:pass@example.com", false},
		{"url with trailing slash", "https://example.com/", false},
		{"url with multiple slashes", "https://example.com//path", false},
		{"url with dots in path", "https://example.com/path/../other", false},
		{"url with numbers in host", "https://example123.com", false},
		{"url with hyphens in host", "https://my-example.com", false},
		{"url with underscores in host", "https://my_example.com", false},
		{"url with very long host", "https://very-long-subdomain-name-that-exceeds-normal-length-limits.example.com", false},
		{"url with very long path", "https://example.com/" + string(make([]byte, 1000)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &URLValidator{}
			err := validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid URL format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestURLValidator_Validate_NonStringTypes(t *testing.T) {
	validator := &URLValidator{}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"int", 123},
		{"float", 123.45},
		{"bool", true},
		{"slice", []string{"https://example.com"}},
		{"map", map[string]string{"url": "https://example.com"}},
		{"struct", struct{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "url validation only applies to strings")
		})
	}
}

func TestURLValidator_New(t *testing.T) {
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
			validator := &URLValidator{}
			result, err := validator.New(tt.params)

			require.NoError(t, err)
			require.NotNil(t, result)

			urlValidator, ok := result.(*URLValidator)
			require.True(t, ok)
			assert.NotNil(t, urlValidator)
		})
	}
}

func TestURLValidator_Key(t *testing.T) {
	validator := &URLValidator{}
	result := validator.Key()
	assert.Equal(t, "url", result)
}
