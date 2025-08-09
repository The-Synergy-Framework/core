// Package validation provides a tag-based validation system similar to Java annotations.
// It allows defining validation rules using struct tags and provides a flexible
// validation framework for the Synergy Framework.
package validation

// Validator defines the contract for validation rules
type Validator interface {
	Validate(value any) error
	// New creates a new instance of this validator from parameters
	New(params map[string]string) (Validator, error)
	// Key returns the registration key for this validator
	Key() string
}

// Rule represents a single validation rule
type Rule struct {
	Name   string
	Params map[string]string
}

// The Result contains the result of a validation operation
type Result struct {
	IsValid bool
	Errors  []Error
}
