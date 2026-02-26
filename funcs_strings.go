package entdomain

import (
	"strings"
)

// hasPrefix checks if a string has a prefix.
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
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
