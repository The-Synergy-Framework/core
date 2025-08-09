package validation

import "fmt"

// Error represents a single validation error
type Error struct {
	Field   string
	Rule    string
	Message string
	Value   any
}

// NewValidationError creates a new ValidationError instance
func NewValidationError(field, rule, message string, value any) Error {
	return Error{
		Field:   field,
		Rule:    rule,
		Message: message,
		Value:   value,
	}
}

// Error returns the error message for the ValidationError
func (e Error) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}
