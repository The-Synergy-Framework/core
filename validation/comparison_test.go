package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComparisonValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		validator *ComparisonValidator
		value     interface{}
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "int greater than or equal - valid",
			validator: &ComparisonValidator{Operator: ">=", CompareValue: 10},
			value:     15,
			wantErr:   false,
		},
		{
			name:      "int greater than or equal - equal value",
			validator: &ComparisonValidator{Operator: ">=", CompareValue: 10},
			value:     10,
			wantErr:   false,
		},
		{
			name:      "int greater than or equal - invalid",
			validator: &ComparisonValidator{Operator: ">=", CompareValue: 10},
			value:     5,
			wantErr:   true,
			errMsg:    "value must be >= 10",
		},
		{
			name:      "int less than or equal - valid",
			validator: &ComparisonValidator{Operator: "<=", CompareValue: 10},
			value:     5,
			wantErr:   false,
		},
		{
			name:      "int less than or equal - equal value",
			validator: &ComparisonValidator{Operator: "<=", CompareValue: 10},
			value:     10,
			wantErr:   false,
		},
		{
			name:      "int less than or equal - invalid",
			validator: &ComparisonValidator{Operator: "<=", CompareValue: 10},
			value:     15,
			wantErr:   true,
			errMsg:    "value must be <= 10",
		},
		{
			name:      "int greater than - valid",
			validator: &ComparisonValidator{Operator: ">", CompareValue: 10},
			value:     15,
			wantErr:   false,
		},
		{
			name:      "int greater than - equal value",
			validator: &ComparisonValidator{Operator: ">", CompareValue: 10},
			value:     10,
			wantErr:   true,
			errMsg:    "value must be > 10",
		},
		{
			name:      "int less than - valid",
			validator: &ComparisonValidator{Operator: "<", CompareValue: 10},
			value:     5,
			wantErr:   false,
		},
		{
			name:      "int less than - equal value",
			validator: &ComparisonValidator{Operator: "<", CompareValue: 10},
			value:     10,
			wantErr:   true,
			errMsg:    "value must be < 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate(tt.value)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestComparisonValidator_Validate_DifferentTypes(t *testing.T) {
	validator := &ComparisonValidator{Operator: ">=", CompareValue: 5}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"int8", int8(10)},
		{"int16", int16(10)},
		{"int32", int32(10)},
		{"int64", int64(10)},
		{"uint", uint(10)},
		{"uint8", uint8(10)},
		{"uint16", uint16(10)},
		{"uint32", uint32(10)},
		{"uint64", uint64(10)},
		{"float32", float32(10.0)},
		{"float64", float64(10.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			require.NoError(t, err)
		})
	}
}

func TestComparisonValidator_Validate_NonNumericTypes(t *testing.T) {
	validator := &ComparisonValidator{Operator: ">=", CompareValue: 5}

	tests := []struct {
		name  string
		value interface{}
	}{
		{"string", "test"},
		{"bool", true},
		{"slice", []int{1, 2, 3}},
		{"map", map[string]int{"a": 1}},
		{"struct", struct{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.value)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "comparison validation only applies to numeric types")
		})
	}
}

func TestComparisonValidator_New(t *testing.T) {
	validator := &ComparisonValidator{Operator: ">="}

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
			errMsg:  "comparison validation requires a value parameter",
		},
		{
			name:    "invalid value parameter",
			params:  map[string]string{"value": "invalid"},
			wantErr: true,
			errMsg:  "invalid comparison value: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.New(tt.params)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				comparisonValidator, ok := result.(*ComparisonValidator)
				require.True(t, ok)
				assert.Equal(t, ">=", comparisonValidator.Operator)
				assert.Equal(t, 10.5, comparisonValidator.CompareValue)
			}
		})
	}
}

func TestComparisonValidator_Key(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected string
	}{
		{"greater than or equal", ">=", ">="},
		{"less than or equal", "<=", "<="},
		{"greater than", ">", ">"},
		{"less than", "<", "<"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &ComparisonValidator{Operator: tt.operator}
			result := validator.Key()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComparisonValidator_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		validator   *ComparisonValidator
		value       interface{}
		description string
	}{
		{
			name:        "zero value comparison",
			validator:   &ComparisonValidator{Operator: ">=", CompareValue: 0},
			value:       0,
			description: "should accept zero values",
		},
		{
			name:        "negative value comparison",
			validator:   &ComparisonValidator{Operator: ">=", CompareValue: -10},
			value:       -5,
			description: "should handle negative values",
		},
		{
			name:        "large value comparison",
			validator:   &ComparisonValidator{Operator: ">=", CompareValue: 1e6},
			value:       2e6,
			description: "should handle large values",
		},
		{
			name:        "decimal value comparison",
			validator:   &ComparisonValidator{Operator: ">=", CompareValue: 3.14159},
			value:       3.14159,
			description: "should handle decimal values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate(tt.value)
			require.NoError(t, err, tt.description)
		})
	}
}
