// Package globals provides global variables.
package globals

import (
	"errors"
)

// TODO: get rid of these three ⤵. It requires quite an amount of work.︎

// HyphaExists holds function that checks if the hypha is present. By default, it is set to a function that is always true.
var HyphaExists func(string) bool

// HyphaAccess holds function that accesses a hypha by its name. By default, it is set to a function that always returns an error.
var HyphaAccess func(string) (rawText, binaryHtml string, err error)

// HyphaIterate is a function that iterates all hypha names existing. By default, it is set to a function that does nothing.
var HyphaIterate func(func(string))

func init() {
	HyphaExists = func(_ string) bool {
		return true
	}
	HyphaAccess = func(_ string) (string, string, error) {
		return "", "", errors.New("globals.HyphaAccess not set")
	}
	HyphaIterate = func(_ func(string)) {}
}
