package blocks

import (
	"strings"

	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/links"
)

// LaunchPad is a container for RocketLinks.
type LaunchPad struct {
	Rockets []RocketLink
}

func MakeLaunchPad() LaunchPad {
	return LaunchPad{[]RocketLink{}}
}

func (lp *LaunchPad) AddRocket(rl RocketLink) {
	lp.Rockets = append(lp.Rockets, rl)
}

// RocketLink is a rocket link which is meant to be nested inside LaunchPad.
type RocketLink struct {
	IsEmpty bool
	links.Link
}

func MakeRocketLink(line, hyphaName string) RocketLink {
	line = strings.TrimSpace(line[2:])
	if line == "" {
		return RocketLink{IsEmpty: true}
	}

	var (
		// Address is text after => till first whitespace
		addr = strings.Fields(line)[0]
		// Display is what is left
		display = strings.TrimPrefix(addr, addr)
		rl      = RocketLink{
			IsEmpty: false,
			Link:    *links.From(addr, display, hyphaName),
		}
	)

	if rl.OfKind(links.LinkLocalHypha) && globals.HyphaExists(rl.Address()) {
		rl.DestinationUnknown = false
	}

	return rl
}

// LinkParts determines what href, text and class should resulting <a> have based on mycomarkup's addr, display and hypha Target.
//
// => addr display
// [[addr|display]]
// TODO: deprecate
func LinkParts(addr, display, hyphaName string) (href, text, class string) {
	l := links.From(addr, display, hyphaName)
	if l.OfKind(links.LinkLocalHypha) && !globals.HyphaExists(l.Address()) {
		l.DestinationUnknown = true
	}
	return l.Href(), l.Display(), l.Classes()
}
