// Package validation provides a tag-based validation system similar to Java annotations.
// It allows defining validation rules using struct tags and provides a flexible
// validation framework for the Synergy Framework.
package validation

import (
	"reflect"
	"strings"
)

const (
	defaultRuleCount = 3
)

// Global validator registry with built-in validators
var defaultRegistry = newValidatorRegistry()

// Validate validates a struct using validation tags and the default registry
func Validate(targetStruct any) *Result {
	return validateWithRegistry(targetStruct, defaultRegistry)
}

// ValidateWithCustomValidators validates a struct with additional custom validators
// without permanently registering them with the default registry
func ValidateWithCustomValidators(targetStruct any, customValidators ...Validator) *Result {
	registry := newValidatorRegistry()
	for _, validator := range customValidators {
		registry.registerValidator(validator)
	}
	return validateWithRegistry(targetStruct, registry)
}

// RegisterCustomValidator registers a custom validator with the default registry
func RegisterCustomValidator(validator Validator) {
	defaultRegistry.registerValidator(validator)
}

// RegisterCustomValidators registers multiple custom validators with the default registry
func RegisterCustomValidators(validators ...Validator) {
	for _, validator := range validators {
		defaultRegistry.registerValidator(validator)
	}
}

// GetRegisteredValidators returns all registered validator names
func GetRegisteredValidators() []string {
	return defaultRegistry.listValidators()
}

// HasValidator checks if a validator exists in the default registry
func HasValidator(name string) bool {
	return defaultRegistry.hasValidator(name)
}

func validateWithRegistry(targetStruct any, registry *validatorRegistry) *Result {
	result := &Result{
		IsValid: true,
		Errors:  []Error{},
	}

	val := reflect.ValueOf(targetStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if !isValidStruct(val) {
		result.IsValid = false
		result.Errors = append(result.Errors, NewValidationError("root", "type", "object must be a struct", targetStruct))
		return result
	}

	validateStruct(val, "", result, registry)
	return result
}

func isValidStruct(val reflect.Value) bool {
	return val.Kind() == reflect.Struct
}

func validateStruct(val reflect.Value, prefix string, result *Result, registry *validatorRegistry) {
	valType := val.Type()

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		fieldValue := val.Field(i)

		validationTag := field.Tag.Get("validate")
		if validationTag == "" {
			continue
		}

		fieldName := buildFieldName(prefix, field.Name)
		validateField(fieldValue, fieldName, validationTag, result, registry)

		if isEmbeddedStruct(field, fieldValue) {
			validateStruct(fieldValue, fieldName, result, registry)
		}
	}
}

func buildFieldName(prefix, fieldName string) string {
	if prefix == "" {
		return fieldName
	}
	return prefix + "." + fieldName
}

func validateField(fieldValue reflect.Value, fieldName, validationTag string, result *Result, registry *validatorRegistry) {
	rules := parseValidationRules(validationTag)

	for _, rule := range rules {
		if err := applyValidationRule(fieldValue, rule, registry); err != nil {
			result.IsValid = false
			result.Errors = append(result.Errors, NewValidationError(fieldName, rule.Name, err.Error(), fieldValue.Interface()))
		}
	}
}

func applyValidationRule(fieldValue reflect.Value, rule Rule, registry *validatorRegistry) error {
	validator, err := registry.getValidator(rule)
	if err != nil {
		return err
	}

	return validator.Validate(fieldValue.Interface())
}

func isEmbeddedStruct(field reflect.StructField, fieldValue reflect.Value) bool {
	return field.Anonymous && fieldValue.Kind() == reflect.Struct
}

func parseValidationRules(tag string) []Rule {
	rules := make([]Rule, 0, defaultRuleCount)

	ruleStrings := strings.Split(tag, ",")
	for _, ruleString := range ruleStrings {
		ruleString = strings.TrimSpace(ruleString)
		if ruleString == "" {
			continue
		}

		rule := parseSingleRule(ruleString)
		rules = append(rules, rule)
	}

	return rules
}

func parseSingleRule(ruleString string) Rule {
	parts := strings.SplitN(ruleString, ":", 2)
	ruleName := strings.TrimSpace(parts[0])
	params := make(map[string]string)

	if len(parts) > 1 {
		params = parseRuleParameters(parts[1])
	}

	return Rule{
		Name:   ruleName,
		Params: params,
	}
}

func parseRuleParameters(paramString string) map[string]string {
	params := make(map[string]string)

	paramStrings := strings.Split(paramString, ",")
	for _, paramString := range paramStrings {
		paramString = strings.TrimSpace(paramString)
		if paramString == "" {
			continue
		}

		key, value := parseParameterKeyValue(paramString)
		if key != "" {
			params[key] = value
		}
	}

	return params
}

func parseParameterKeyValue(paramString string) (key, value string) {
	keyValue := strings.SplitN(paramString, "=", 2)
	if len(keyValue) == 2 {
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		return key, value
	}

	return "value", paramString
}
