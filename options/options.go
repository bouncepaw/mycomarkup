// Package options provides a configuration data structure to pass when invoking the Mycomarkup parser.
package options

// Options is what you pass when invoking the parser.
type Options struct {
	// Canonical hypha name that is being parsed now.
	HyphaName string

	// Canonical URL (including the protocol) of your website. Used for OpenGraph. Example: https://mycorrhiza.wiki.
	WebSiteURL string

	TransclusionSupported bool
}
