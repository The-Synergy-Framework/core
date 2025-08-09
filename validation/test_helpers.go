package validation

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testValidatorNew is a helper function to test the New method of validators
func testValidatorNew(t *testing.T, validator Validator, params map[string]string, expectedValue interface{}, expectedField string) {
	result, err := validator.New(params)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Use reflection to get the field value
	val := reflect.ValueOf(result)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(expectedField)
	require.True(t, field.IsValid())
	assert.Equal(t, expectedValue, field.Interface())
}

// testValidatorNewError is a helper function to test error cases in the New method
func testValidatorNewError(t *testing.T, validator Validator, params map[string]string, expectedErrMsg string) {
	result, err := validator.New(params)

	require.Error(t, err)
	require.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErrMsg)
}

// testValidatorKey is a helper function to test the Key method of validators
func testValidatorKey(t *testing.T, validator Validator, expectedKey string) {
	result := validator.Key()
	assert.Equal(t, expectedKey, result)
}
