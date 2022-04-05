// Package globals provides global variables.
package globals

import (
	"errors"
)

// TODO: get rid of these functions ⤵. It requires quite an amount of work.︎

// HyphaAccess holds function that accesses a hypha by its name. By default, it is set to a function that always returns an error.
var HyphaAccess func(string) (rawText, binaryHtml string, err error)

func init() {
	HyphaAccess = func(_ string) (string, string, error) {
		return "", "", errors.New("globals.HyphaAccess not set")
	}
}
