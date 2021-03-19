package utils

import (
	"regexp"
	"strings"
)

// Strip hypha name from all ancestor names, replace _ with spaces, title case
func BeautifulName(uglyName string) string {
	if uglyName == "" {
		return uglyName
	}
	return strings.Title(strings.ReplaceAll(uglyName, "_", " "))
}

// CanonicalName makes sure the `name` is canonical. A name is canonical if it is lowercase and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// HyphaPattern is a pattern which all hyphae must match.
var HyphaPattern = regexp.MustCompile(`[^?!:#@><*|"\'&%{}]+`)

// Function that returns a function that can strip `prefix` and trim whitespace when called.
func remover(prefix string) func(string) string {
	return func(l string) string {
		return strings.TrimSpace(strings.TrimPrefix(l, prefix))
	}
}

// Remove #, ## or ### from beginning of `line`.
func removeHeadingOctothorps(line string) string {
	f := remover("#")
	return f(f(f(line)))
}

// Return a canonical representation of a hypha `name`.
func canonicalName(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
}
