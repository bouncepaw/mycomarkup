// Package util provides simple utility functions usable in all modules.
package util

import (
	"regexp"
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

var badCharactersRe = regexp.MustCompile(`[?!:#@<>*|'"&%{}\\]`)

// SanitizedName is unsafeName with all unsafe characters removed.
func SanitizedName(unsafeName string) string {
	return badCharactersRe.ReplaceAllString(unsafeName, "")
}

// BeautifulName makes the ugly name beautiful by replacing _ with spaces and using title case
func BeautifulName(uglyName string) string {
	uglyName = SanitizedName(uglyName)
	if uglyName == "" {
		return uglyName
	}
	// What other transformations can we apply for a better beautifying process?
	return strings.Title(strings.ReplaceAll(uglyName, "_", " "))
}

// CanonicalName returns the canonical form of the name. A name is canonical if it is lowercase, all left and right whitespace is trimmed and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	var (
		spaceless = strings.ReplaceAll(SanitizedName(name), " ", "_")
		trimmed   = strings.Trim(spaceless, "_")
	)
	return strings.ToLower(trimmed)
}

// DefaultString returns d if s is an empty string, s otherwise.
func DefaultString(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// TernaryConditionString is an approximation of an expression like (cond ? thenBranch : elseBranch) from other programming languages. There is no real need for this function, I just felt the urge to add it.
func TernaryConditionString(cond bool, thenBranch, elseBranch string) string {
	if cond {
		return thenBranch
	}
	return elseBranch
}
