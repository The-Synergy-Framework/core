package validation

import "fmt"

// validatorRegistry holds all available validators and provides methods for managing them
type validatorRegistry struct {
	validators map[string]Validator
}

// newValidatorRegistry creates a new validator registry with built-in validators (internal use)
func newValidatorRegistry() *validatorRegistry {
	registry := &validatorRegistry{
		validators: make(map[string]Validator),
	}

	registry.registerBuiltInValidators()
	return registry
}

// registerBuiltInValidators registers all the built-in validators
func (r *validatorRegistry) registerBuiltInValidators() {
	r.registerValidator(&RequiredValidator{})
	r.registerValidator(&EmailValidator{})
	r.registerValidator(&URLValidator{})

	r.registerValidator(&MinValidator{})
	r.registerValidator(&MaxValidator{})
	r.registerValidator(&LenValidator{})
	r.registerValidator(&OneOfValidator{})
	r.registerValidator(&RegexpValidator{})

	r.registerValidator(&ComparisonValidator{Operator: ">"})
	r.registerValidator(&ComparisonValidator{Operator: "<"})
	r.registerValidator(&ComparisonValidator{Operator: ">="})
	r.registerValidator(&ComparisonValidator{Operator: "<="})
}

// registerValidator adds a validator to the registry using its own Key() method (internal use)
func (r *validatorRegistry) registerValidator(validator Validator) {
	r.validators[validator.Key()] = validator
}

// hasValidator checks if a validator exists (internal use)
func (r *validatorRegistry) hasValidator(name string) bool {
	_, exists := r.validators[name]
	return exists
}

// listValidators returns all registered validator names (internal use)
func (r *validatorRegistry) listValidators() []string {
	names := make([]string, 0, len(r.validators))
	for name := range r.validators {
		names = append(names, name)
	}
	return names
}

// getValidator creates a validator instance from a rule (internal use)
func (r *validatorRegistry) getValidator(rule Rule) (Validator, error) {
	validator, exists := r.validators[rule.Name]
	if !exists {
		return nil, fmt.Errorf("unknown validation rule: %s", rule.Name)
	}
	return validator.New(rule.Params)
}
