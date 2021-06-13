package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/util"
	"strings"

	"github.com/bouncepaw/mycomarkup/links"
)

// LaunchPad is a container for RocketLinks.
type LaunchPad struct {
	Rockets []RocketLink
}

// ID returns the launchpad's id which is rocket- and a number. Note that it does not say launchpad.
func (lp LaunchPad) ID(counter *IDCounter) string {
	counter.launchpads++
	return fmt.Sprintf("rocket-%d", counter.launchpads)
}

func (lp LaunchPad) isBlock() {}

// MakeLaunchPad returns an empty launchpad. Add rocket links there using AddRocket.
func MakeLaunchPad() LaunchPad {
	return LaunchPad{[]RocketLink{}}
}

// AddRocket stores the rocket link in the launchpad.
func (lp *LaunchPad) AddRocket(rl RocketLink) {
	lp.Rockets = append(lp.Rockets, rl)
}

// ColorRockets marks links to existing hyphae as existing.
func (lp *LaunchPad) ColorRockets() {
	globals.HyphaIterate(func(hyphaName string) {
		for _, rocket := range lp.Rockets {
			// TODO: do not canonize every time
			if util.CanonicalName(rocket.TargetHypha()) == hyphaName {
				rocket.Link.MarkAsExisting()
			}
		}
	})
}

// RocketLink is a rocket link which is meant to be nested inside LaunchPad.
type RocketLink struct {
	IsEmpty bool
	links.Link
}

func (r RocketLink) isBlock() {}

// ID returns an empty string because rocket links do not have ids on their own.
func (r RocketLink) ID(_ *IDCounter) string {
	return ""
}

// MakeRocketLink parses the rocket link on the given line and returns it.
func MakeRocketLink(line, hyphaName string) RocketLink {
	line = strings.TrimSpace(line[2:])
	if line == "" {
		return RocketLink{IsEmpty: true}
	}

	var (
		// Address is text after => till first whitespace
		addr = strings.Fields(line)[0]
		// Display is what is left
		display = strings.TrimPrefix(line, addr)
		rl      = RocketLink{
			IsEmpty: false,
			Link:    *links.From(addr, display, hyphaName),
		}
	)

	return rl
}
