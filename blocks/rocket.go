package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v4/links"
	"github.com/bouncepaw/mycomarkup/v4/mycocontext"
	"github.com/bouncepaw/mycomarkup/v4/util"
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

// NewLaunchPad returns a launchpad with the given RocketLinks inside
func NewLaunchPad(rockets []RocketLink) LaunchPad {
	return LaunchPad{Rockets: rockets}
}

// ColorRockets marks links to existing hyphae as existing. V3
func (lp *LaunchPad) ColorRockets(ctx mycocontext.Context) {
	mycocontext.IterateHyphaNamesWith(ctx, func(hyphaName string) {
		for i, rocket := range lp.Rockets {
			// TODO: do not canonize every time
			if util.CanonicalName(rocket.TargetHypha()) == hyphaName {
				rocket.Link = rocket.Link.CopyMarkedAsExisting()
			}
			lp.Rockets[i] = rocket
		}
	})
}

// RocketLink is a rocket link which is meant to be nested inside LaunchPad.
type RocketLink struct {
	IsEmpty bool
	links.Link
}

// ID returns an empty string because rocket links do not have ids on their own.
func (r RocketLink) ID(_ *IDCounter) string {
	return ""
}
