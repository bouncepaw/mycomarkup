// Package options provides a configuration data structure to pass when invoking the Mycomarkup parser.
package options

import (
	"errors"
)

// Options is what you pass when invoking the parser.
type Options struct {
	// Canonical hypha name that is being parsed now.
	HyphaName string

	// Canonical URL (including the protocol) of your website. Used for OpenGraph. Example: https://mycorrhiza.wiki.
	WebSiteURL string

	TransclusionSupported bool
	RedLinksSupported     bool

	HyphaExists           func(string) bool
	IterateHyphaNamesWith func(func(string))
	HyphaHTMLData         func(string) (rawText, binaryHtml string, err error)

	LocalLinkHref                    func(string) string
	LocalImgSrc                      func(string) string
	LinkHrefFormatForInterwikiPrefix func(string) string
	ImgSrcFormatForInterwikiPrefix   func(string) string
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
	if opts.LocalLinkHref == nil {
		opts.LocalLinkHref = func(hyphaName string) string {
			return hyphaName
		}
	}
	if opts.LocalImgSrc == nil {
		opts.LocalImgSrc = func(hyphaName string) string {
			return hyphaName
		}
	}
	if opts.LinkHrefFormatForInterwikiPrefix == nil {
		opts.LinkHrefFormatForInterwikiPrefix = func(prefix string) string {
			return "{NAME}"
		}
	}
	if opts.ImgSrcFormatForInterwikiPrefix == nil {
		opts.ImgSrcFormatForInterwikiPrefix = func(prefix string) string {
			return "{NAME}"
		}
	}
	return opts
}
