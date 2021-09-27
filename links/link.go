// Package links provides a data type for links.
package links

import (
	"fmt"
	"html"
	"path"
	"strings"

	"github.com/bouncepaw/mycomarkup/v2/util"
)

// LinkType tells what type the given link is.
type LinkType int

const (
	// LinkLocalRoot is a link like "/list", "/user-list", etc.
	LinkLocalRoot LinkType = iota
	// LinkLocalHypha is a link like "test", "../test", etc.
	LinkLocalHypha
	// LinkExternal is an external link with specified protocol.
	LinkExternal
)

// Link is an abstraction for universal representation of links, be they links in mycomarkup links or whatever.
type Link struct {
	// Parsed stuff
	kind     LinkType
	protocol string
	address  string //
	anchor   string // # and everything after it

	display string
	// Settable stuff
	destinationKnown bool // set to true when you have //checked// that the target hypha exists. It might be false for non-hypha links.
}

// From makes a link from the given source address and display text on the given hypha. The arguments are stripped of whitespace on both sides before further processing.
func From(srcAddress, srcDisplay, srcHypha string) Link {
	srcAddress = strings.TrimSpace(srcAddress)
	srcDisplay = strings.TrimSpace(srcDisplay)
	srcHypha = strings.TrimSpace(srcHypha)

	link := Link{
		address:          srcAddress,
		destinationKnown: false,
	}

	// If there is a hash sign in the address, move everything starting from the sign to the end of the address to the anchor field and truncate the address.
	if pos := strings.IndexRune(srcAddress, '#'); pos != -1 {
		link.anchor = srcAddress[pos:]
		link.address = link.address[:pos]
	}

	// NOTE: This part will need some extending with introduction of interwiki.

	switch {
	// If is an external link
	case strings.ContainsRune(link.address, ':'):
		pos := strings.IndexRune(link.address, ':')
		link.kind = LinkExternal
		link.protocol = link.address[:pos+1]
		link.address = link.address[pos+1:]
		if strings.HasPrefix(link.address, "//") && len(link.address) > 2 {
			link.protocol += "//"
			link.address = link.address[2:]
		}
		link.display = link.address + link.anchor
	case strings.HasPrefix(link.address, "/"):
		link.kind = LinkLocalRoot
		link.display = link.address + link.anchor
	case strings.HasPrefix(link.address, "./"):
		link.kind = LinkLocalHypha
		link.display = link.address + link.anchor
		link.address = util.CanonicalName(path.Join(srcHypha, link.address[2:]))
	case link.address == "..":
		link.kind = LinkLocalHypha
		link.address = util.CanonicalName(path.Dir(srcHypha))
		link.display = ".."
	case strings.HasPrefix(link.address, "../"):
		link.kind = LinkLocalHypha
		link.display = link.address + link.anchor
		link.address = util.CanonicalName(path.Join(path.Dir(srcHypha), link.address[3:]))
	case strings.HasPrefix(link.address, "#"):
		link.kind = LinkLocalHypha
		link.anchor = link.address
		link.address = util.CanonicalName(srcHypha)
		link.display = link.anchor
	default:
		link.kind = LinkLocalHypha
		link.display = link.address + link.anchor
		link.address = util.CanonicalName(link.address)
	}

	if srcDisplay != "" {
		link.display = srcDisplay
	}

	return link
}

// IsBlueLink is true if the link should be blue, not red. Red links are links to hyphae that do not exist, all other links are blue.
func (link Link) IsBlueLink() bool {
	return !(link.OfKind(LinkLocalHypha) && !link.destinationKnown)
}

// CopyMarkedAsExisting returns a copy of the link that is marked as existing, i/e colored in blue.
func (link Link) CopyMarkedAsExisting() Link {
	link.destinationKnown = true
	return link
}

// Classes returns CSS class string for given link. It is not wrapped in any quotes, wrap yourself.
func (link Link) Classes() (classes string) {
	classes = "wikilink"
	switch link.kind {
	case LinkLocalRoot, LinkLocalHypha:
		classes += " wikilink_internal"
		if !link.IsBlueLink() {
			classes += " wikilink_new"
		}
	case LinkExternal:
		classes += fmt.Sprintf(
			" wikilink_external wikilink_%s",
			strings.TrimSuffix(strings.TrimSuffix(link.protocol, "://"), ":"),
		)
	}
	return classes
}

// Href returns escaped content for the href attribute for HTML link. You should always use it.
func (link Link) Href() string {
	switch link.kind {
	case LinkExternal, LinkLocalRoot:
		return html.EscapeString(link.protocol + link.address + link.anchor)
	default:
		// TODO: configure the path
		return "/hypha/" + html.EscapeString(link.address+link.anchor)
	}
}

// ImgSrc returns escaped content for src attribute of img tag. Used with `img{}`.
func (link Link) ImgSrc() string {
	switch link.kind {
	case LinkExternal, LinkLocalRoot:
		return html.EscapeString(link.protocol + link.address + link.anchor)
	default:
		// TODO: configure the path
		return "/binary/" + html.EscapeString(link.address)
	}
}

// Display returns the display text of the given link. It is not escaped, escape by yourself.
func (link Link) Display() string {
	return link.display
}

// TargetHypha returns the canonical name of the target hypha. Use for hypha links.
func (link Link) TargetHypha() string {
	return util.CanonicalName(link.address)
}

// OfKind is true if the given link is of the given kind, i/e the kinds are equal.
func (link Link) OfKind(kind LinkType) bool {
	return link.kind == kind
}
