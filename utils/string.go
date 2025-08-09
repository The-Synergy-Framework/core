package utils

import "strings"

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(str string) string {
	var result string
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}
