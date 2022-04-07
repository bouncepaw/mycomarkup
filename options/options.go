// Package options provides a configuration data structure to pass when invoking the Mycomarkup parser.
package options

import "errors"

// Options is what you pass when invoking the parser.
type Options struct {
	// Canonical hypha name that is being parsed now.
	HyphaName string

	// Canonical URL (including the protocol) of your website. Used for OpenGraph. Example: https://mycorrhiza.wiki.
	WebSiteURL string

	TransclusionSupported bool

	HyphaExists           func(string) bool
	IterateHyphaNamesWith func(func(string))
	HyphaHTMLData         func(string) (rawText, binaryHtml string, err error)
}

func (opts Options) FillTheRest() Options {
	if opts.HyphaExists == nil {
		opts.HyphaExists = func(hyphaName string) bool {
			return true
		}
	}
	if opts.IterateHyphaNamesWith == nil {
		opts.IterateHyphaNamesWith = func(func(string)) {}
	}
	if opts.HyphaHTMLData == nil {
		opts.HyphaHTMLData = func(_ string) (string, string, error) {
			return "", "", errors.New("HyphaHTMLData not set")
		}
	}
	return opts
}
