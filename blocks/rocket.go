package blocks

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
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

// LinksColored marks links to existing hyphae as existing. V3
func (lp LaunchPad) LinksColored(ctx mycocontext.Context) LaunchPad {
	var probes []func(string)
	for _, rocket := range lp.Rockets {
		if rocket.IsEmpty {
			continue
		}
		if probe := rocket.Link.HyphaProbe(ctx); probe != nil {
			probes = append(probes, probe)
		}
	}
	ctx.Options().IterateHyphaNamesWith(func(hyphaName string) {
		for _, probe := range probes {
			probe(hyphaName)
		}
	})
	return lp
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
