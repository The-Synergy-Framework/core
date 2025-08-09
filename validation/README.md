# Validation Package

Package `validation` provides a flexible, tag-based validation system for Go structs, similar to Java annotations. It allows defining validation rules using struct tags and provides a comprehensive validation framework for the Synergy Framework.

**Import Path:** `core/validation`

## Overview

The validation package offers a declarative approach to data validation using struct tags. It provides a rich set of built-in validators and supports custom validators for specialized validation needs. The package is designed to be:

- **Simple to use**: Define validation rules directly in struct tags
- **Extensible**: Register custom validators for domain-specific validation
- **Type-safe**: Leverages Go's reflection system for robust validation
- **Comprehensive**: Supports nested structs, embedded structs, and complex validation scenarios

The package automatically handles different data types and provides clear, actionable error messages when validation fails.

## Quick Start

```go
package main

import (
    "fmt"
    "core/validation"
)

type User struct {
    Name     string `validate:"required"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"min:18"`
    Username string `validate:"required,min:3,max:20"`
}

func main() {
    user := User{
        Name:     "",
        Email:    "invalid-email",
        Age:      16,
        Username: "ab",
    }
    
    result := validation.Validate(user)
    if !result.IsValid {
        for _, err := range result.Errors {
            fmt.Println(err.Error())
        }
    }
}
```

**Output:**
```
validation failed for field 'Name': field is required (value: )
validation failed for field 'Email': invalid email format (value: invalid-email)
validation failed for field 'Age': value must be at least 18 (value: 16)
validation failed for field 'Username': string length must be at least 3 (value: ab)
```

## API Reference

### Core Functions

#### `Validate(targetStruct any) *Result`

Validates a struct using validation tags and the default registry.

**Parameters:**
- `targetStruct`: The struct to validate (can be a pointer or value)

**Returns:**
- `*Result`: Validation result containing success status and any errors

**Example:**
```go
type Product struct {
    Name  string `validate:"required"`
    Price float64 `validate:"min:0"`
}

product := Product{Name: "Widget", Price: -10}
result := validation.Validate(product)
```

#### `ValidateWithCustomValidators(targetStruct any, customValidators ...Validator) *Result`

Validates a struct with additional custom validators without permanently registering them with the default registry.

**Parameters:**
- `targetStruct`: The struct to validate
- `customValidators`: Variable number of custom validators to use for this validation, alongside the default validators

**Returns:**
- `*Result`: Validation result

**Example:**
```go
customValidator := &MyCustomValidator{}
result := validation.ValidateWithCustomValidators(product, customValidator)
```

#### `RegisterCustomValidator(validator Validator)`

Registers a custom validator with the default registry for use in all subsequent validations.

**Parameters:**
- `validator`: Custom validator implementing the Validator interface

#### `RegisterCustomValidators(validators ...Validator)`

Registers multiple custom validators with the default registry.

**Parameters:**
- `validators`: Variable number of custom validators

#### `GetRegisteredValidators() []string`

Returns all registered validator names.

**Returns:**
- `[]string`: Slice of validator names

#### `HasValidator(name string) bool`

Checks if a validator exists in the default registry.

**Parameters:**
- `name`: Validator name to check

**Returns:**
- `bool`: True if validator exists, false otherwise

### Core Types

#### `Validator` Interface

Defines the contract for validation rules.

```go
type Validator interface {
    Validate(value any) error
    New(params map[string]string) (Validator, error)
    Key() string
}
```

**Methods:**
- `Validate(value any) error`: Validates a value and returns an error if validation fails
- `New(params map[string]string) (Validator, error)`: Creates a new validator instance from parameters
- `Key() string`: Returns the registration key for this validator

#### `Rule` Struct

Represents a single validation rule.

```go
type Rule struct {
    Name   string
    Params map[string]string
}
```

**Fields:**
- `Name`: The name of the validation rule
- `Params`: Parameters for the validation rule

#### `Result` Struct

Contains the result of a validation operation.

```go
type Result struct {
    IsValid bool
    Errors  []Error
}
```

**Fields:**
- `IsValid`: True if validation passed, false otherwise
- `Errors`: Slice of validation errors

#### `Error` Struct

Represents a single validation error.

```go
type Error struct {
    Field   string
    Rule    string
    Message string
    Value   any
}
```

**Fields:**
- `Field`: The field name that failed validation
- `Rule`: The validation rule that failed
- `Message`: Human-readable error message
- `Value`: The actual value that failed validation

**Methods:**
- `Error() string`: Returns a formatted error message

## Built-in Validators

### Required Validator

Validates that a field is not empty or nil.

**Tag:** `required`

**Supported Types:** All types

**Example:**
```go
type User struct {
    Name string `validate:"required"`
}
```

### Email Validator

Validates email format using a comprehensive regex pattern.

**Tag:** `email`

**Supported Types:** String

**Example:**
```go
type User struct {
    Email string `validate:"required,email"`
}
```

### URL Validator

Validates URL format.

**Tag:** `url`

**Supported Types:** String

**Example:**
```go
type Website struct {
    URL string `validate:"required,url"`
}
```

### Min Validator

Validates minimum values for numbers and minimum length for strings.

**Tag:** `min:value`

**Parameters:**
- `value`: Minimum value (numeric) or minimum length (string)

**Supported Types:** Numeric types, String

**Example:**
```go
type Product struct {
    Price    float64 `validate:"min:0"`
    Name     string  `validate:"min:3"`
    Quantity int     `validate:"min:1"`
}
```

### Max Validator

Validates maximum values for numbers and maximum length for strings.

**Tag:** `max:value`

**Parameters:**
- `value`: Maximum value (numeric) or maximum length (string)

**Supported Types:** Numeric types, String

**Example:**
```go
type User struct {
    Username string `validate:"max:20"`
    Age      int    `validate:"max:120"`
}
```

### Length Validator

Validates exact string length.

**Tag:** `len:value`

**Parameters:**
- `value`: Exact length required

**Supported Types:** String

**Example:**
```go
type Code struct {
    PIN string `validate:"len:4"`
}
```

### OneOf Validator

Validates that a value is one of the allowed values.

**Tag:** `oneof:values`

**Parameters:**
- `values`: Pipe-separated list of allowed values

**Supported Types:** All types (converted to string for comparison)

**Example:**
```go
type Status struct {
    State string `validate:"oneof:active|inactive|pending"`
}
```

### Regexp Validator

Validates string against a regex pattern.

**Tag:** `regexp:pattern`

**Parameters:**
- `pattern`: Regular expression pattern

**Supported Types:** String

**Example:**
```go
type Phone struct {
    Number string `validate:"regexp:^\\+?[1-9]\\d{1,14}$"`
}
```

### Comparison Validators

Validates numeric comparisons.

**Tags:** `>`, `<`, `>=`, `<=`

**Parameters:**
- `value`: Comparison value

**Supported Types:** Numeric types

**Example:**
```go
type Score struct {
    Points int `validate:">:0,<:100"`
}
```

## Examples

### Basic User Validation

```go
type User struct {
    ID       int     `validate:"min:1"`
    Username string  `validate:"required,min:3,max:20"`
    Email    string  `validate:"required,email"`
    Age      int     `validate:"min:13,max:120"`
    Status   string  `validate:"oneof:active|inactive|suspended"`
}

func validateUser(user User) error {
    result := validation.Validate(user)
    if !result.IsValid {
        return fmt.Errorf("user validation failed: %v", result.Errors)
    }
    return nil
}
```

### Nested Struct Validation

```go
type Address struct {
    Street  string `validate:"required"`
    City    string `validate:"required"`
    Country string `validate:"required"`
}

type Customer struct {
    Name    string  `validate:"required"`
    Email   string  `validate:"required,email"`
    Address Address `validate:"required"`
}

func validateCustomer(customer Customer) error {
    result := validation.Validate(customer)
    if !result.IsValid {
        for _, err := range result.Errors {
            fmt.Printf("Field: %s, Error: %s\n", err.Field, err.Message)
        }
    }
    return nil
}
```

### Custom Validator

```go
type StrongPasswordValidator struct{}

func (v *StrongPasswordValidator) Validate(value any) error {
    password, ok := value.(string)
    if !ok {
        return fmt.Errorf("password must be a string")
    }
    
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    
    if !hasUpper || !hasLower || !hasDigit {
        return fmt.Errorf("password must contain uppercase, lowercase, and digit")
    }
    
    return nil
}

func (v *StrongPasswordValidator) New(params map[string]string) (Validator, error) {
    return &StrongPasswordValidator{}, nil
}

func (v *StrongPasswordValidator) Key() string {
    return "strong_password"
}

// Register the custom validator
func init() {
    validation.RegisterCustomValidator(&StrongPasswordValidator{})
}

// Use in struct
type Account struct {
    Password string `validate:"required,strong_password"`
}
```

### Multiple Validation Rules

```go
type Product struct {
    Name        string  `validate:"required,min:3,max:100"`
    Price       float64 `validate:"min:0"`
    Category    string  `validate:"oneof:electronics|clothing|books|food"`
    SKU         string  `validate:"required,regexp:^[A-Z]{2}-[0-9]{6}$"`
    Description string  `validate:"max:500"`
}

func validateProduct(product Product) *validation.Result {
    return validation.Validate(product)
}
```

## Best Practices

### 1. Use Descriptive Field Names

Field names in error messages are derived from struct field names. Use clear, descriptive names:

```go
// Good
type User struct {
    EmailAddress string `validate:"required,email"`
}

// Avoid
type User struct {
    Email string `validate:"required,email"` // Less descriptive
}
```

### 2. Combine Related Validations

Group related validation rules together:

```go
type User struct {
    // Basic info
    Name     string `validate:"required,min:2,max:50"`
    Email    string `validate:"required,email"`
    
    // Security
    Password string `validate:"required,min:8"`
    
    // Business rules
    Age      int    `validate:"min:13,max:120"`
    Status   string `validate:"oneof:active|inactive|suspended"`
}
```

### 3. Handle Validation Results Properly

Always check the validation result and handle errors appropriately:

```go
func processUser(user User) error {
    result := validation.Validate(user)
    if !result.IsValid {
        // Log validation errors for debugging
        for _, err := range result.Errors {
            log.Printf("Validation error: %s", err.Error())
        }
        return fmt.Errorf("user validation failed")
    }
    return nil
}
```

### 4. Use Custom Validators for Complex Logic

For complex validation logic that can't be expressed with built-in validators:

```go
type DateRangeValidator struct {
    StartDate time.Time
    EndDate   time.Time
}

func (v *DateRangeValidator) Validate(value any) error {
    date, ok := value.(time.Time)
    if !ok {
        return fmt.Errorf("value must be a time.Time")
    }
    
    if date.Before(v.StartDate) || date.After(v.EndDate) {
        return fmt.Errorf("date must be between %v and %v", v.StartDate, v.EndDate)
    }
    
    return nil
}
```

### 5. Validate Early and Often

Validate data as early as possible in your application flow:

```go
func createUser(userData map[string]interface{}) (*User, error) {
    // Validate input data first
    user := User{
        Name:  userData["name"].(string),
        Email: userData["email"].(string),
    }
    
    if result := validation.Validate(user); !result.IsValid {
        return nil, fmt.Errorf("invalid user data: %v", result.Errors)
    }
    
    // Process valid data
    return &user, nil
}
```

## Performance Considerations

### 1. Validation is Reflection-Based

The validation system uses Go's reflection package, which has some performance overhead. For high-performance applications:

- Cache validation results when possible
- Use validation sparingly in hot paths
- Consider pre-validating data structures

### 2. Regex Compilation

Regexp validators compile patterns on each validation. For frequently used patterns:

```go
// Pre-compile regex patterns
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type CustomEmailValidator struct{}

func (v *CustomEmailValidator) Validate(value any) error {
    email, ok := value.(string)
    if !ok {
        return fmt.Errorf("email must be a string")
    }
    
    if !emailRegex.MatchString(email) {
        return fmt.Errorf("invalid email format")
    }
    
    return nil
}
```

### 3. Struct Validation Order

Validation processes struct fields in declaration order. For optimal performance:

- Put most commonly failing validations first
- Use `required` validation early to avoid unnecessary processing

## Advanced Usage

### Embedded Structs

The validation system automatically handles embedded structs:

```go
type BaseEntity struct {
    ID        int       `validate:"min:1"`
    CreatedAt time.Time `validate:"required"`
}

type User struct {
    BaseEntity // Embedded struct
    Name       string `validate:"required"`
    Email      string `validate:"required,email"`
}
```

### Conditional Validation

For conditional validation, use custom validators:

```go
type ConditionalValidator struct {
    Condition func(any) bool
    Validator Validator
}

func (v *ConditionalValidator) Validate(value any) error {
    if v.Condition(value) {
        return v.Validator.Validate(value)
    }
    return nil
}
```

### Validation Context

For context-aware validation, pass additional data through custom validators:

```go
type ContextualValidator struct {
    Context map[string]interface{}
    Rule    string
}

func (v *ContextualValidator) Validate(value any) error {
    // Use context to make validation decisions
    if v.Context["user_role"] == "admin" {
        // Skip certain validations for admins
        return nil
    }
    
    // Apply normal validation rules
    return applyRule(v.Rule, value)
}
```

## Error Handling

### Understanding Error Messages

Error messages follow this format:
```
validation failed for field 'FieldName': Error message (value: actual_value)
```

### Custom Error Messages

To provide custom error messages, create custom validators:

```go
type CustomRequiredValidator struct{}

func (v *CustomRequiredValidator) Validate(value any) error {
    if value == nil {
        return fmt.Errorf("this field cannot be empty")
    }
    // ... rest of validation logic
}
```

### Error Aggregation

The validation system collects all errors rather than stopping at the first failure:

```go
result := validation.Validate(user)
if !result.IsValid {
    fmt.Printf("Found %d validation errors:\n", len(result.Errors))
    for i, err := range result.Errors {
        fmt.Printf("%d. %s: %s\n", i+1, err.Field, err.Message)
    }
}
```

## Thread Safety

The validation package is designed to be thread-safe:

- The default registry is safe for concurrent access
- Validator instances are stateless and can be safely shared
- Validation operations don't modify shared state

However, when registering custom validators:

```go
// Do this once during initialization
func init() {
    validation.RegisterCustomValidator(&MyCustomValidator{})
}

// Don't register validators during concurrent validation
func validateConcurrently() {
    // This is NOT thread-safe
    validation.RegisterCustomValidator(&MyCustomValidator{})
    result := validation.Validate(data)
}
```

## Integration with Other Packages

### HTTP Request Validation

```go
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    result := validation.Validate(user)
    if !result.IsValid {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "errors": result.Errors,
        })
        return
    }
    
    // Process valid user
}
```

### Database Model Validation

```go
type UserModel struct {
    ID       int       `db:"id" validate:"min:1"`
    Username string    `db:"username" validate:"required,min:3,max:20"`
    Email    string    `db:"email" validate:"required,email"`
    Created  time.Time `db:"created_at" validate:"required"`
}

func (u *UserModel) Validate() error {
    result := validation.Validate(u)
    if !result.IsValid {
        return fmt.Errorf("model validation failed: %v", result.Errors)
    }
    return nil
}
```

This comprehensive documentation provides everything needed to effectively use the validation package, from basic usage to advanced customization patterns.