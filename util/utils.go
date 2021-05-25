// Package util provides simple utility functions usable in all modules.
package util

import (
	"strings"
	"unicode"
)

// StringID sanitizes the string and makes it more suitable for the id attribute in HTML.
func StringID(s string) string {
	var (
		ret            strings.Builder
		usedUnderscore bool
	)
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			ret.WriteRune(r)
			usedUnderscore = false
		} else if !usedUnderscore {
			ret.WriteRune('_')
			usedUnderscore = true
		}
	}
	return strings.Trim(ret.String(), "_")
}

// BeautifulName makes the ugly name beautiful by replacing _ with spaces and using title case
func BeautifulName(uglyName string) string {
	// What other transformations can we apply for a better beautifying process?
	if uglyName == "" {
		return uglyName
	}
	return strings.Title(strings.ReplaceAll(uglyName, "_", " "))
}

// CanonicalName returns the canonical form of the name. A name is canonical if it is lowercase, all left and right whitespace is trimmed and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return strings.ToLower(
		strings.ReplaceAll(
			strings.TrimRight(
				strings.TrimLeft(name, "_"),
				"_",
			), " ", "_"))
}

// Remover returns a function that can strip prefix and trim whitespace when called.
func Remover(prefix string) func(string) string {
	return func(l string) string {
		return strings.TrimSpace(strings.TrimPrefix(l, prefix))
	}
}
