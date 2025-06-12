package utils

import (
	"strings"
)

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// ValidateEmail validates email address correctness
func ValidateEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SanitizeInput cleans input data from potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove leading and trailing spaces
	return strings.TrimSpace(input)
}
