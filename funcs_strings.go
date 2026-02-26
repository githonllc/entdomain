package entdomain

import (
	"strings"
	"unicode"
)

// hasPrefix checks if a string has a prefix.
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// camelCase converts a snake_case or PascalCase string to camelCase.
// Examples: "phone_number" → "phoneNumber", "PhoneNumber" → "phoneNumber", "name" → "name".
func camelCase(s string) string {
	if s == "" {
		return s
	}

	// Handle snake_case: split on underscores and capitalize each part after the first.
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		for i := range parts {
			if i == 0 {
				parts[i] = strings.ToLower(parts[i])
			} else if len(parts[i]) > 0 {
				runes := []rune(parts[i])
				runes[0] = unicode.ToUpper(runes[0])
				parts[i] = string(runes[:1]) + strings.ToLower(string(runes[1:]))
			}
		}
		return strings.Join(parts, "")
	}

	// Handle PascalCase: lowercase the first letter.
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
