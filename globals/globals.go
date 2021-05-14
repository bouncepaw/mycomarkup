// Package globals provides global variables.
package globals

// UseBatch is true when mycomarkup is invoked as a program.
var UseBatch bool

// HyphaExists holds function that checks that a hypha is present.
var HyphaExists func(string) bool

// HyphaAccess holds function that accesses a hypha by its name.
var HyphaAccess func(string) (rawText, binaryHtml string, err error)

// HyphaIterate is a function that iterates all hypha names existing.
var HyphaIterate func(func(string))
