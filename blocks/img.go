package blocks

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
)

// Img is an image gallery, consisting of zero or more images.
type Img struct {
	Entries []ImgEntry
	layout  ImgLayout
}

// ID returns the gallery's id which is img- and a number.
func (img Img) ID(counter *IDCounter) string {
	counter.imgs++
	return fmt.Sprintf("img-%d", counter.imgs)
}

// NewImg returns a new Img.
func NewImg(entries []ImgEntry, layout ImgLayout) Img {
	return Img{
		Entries: entries,
		layout:  layout,
	}
}

// Layout returns Img's layout.
func (img Img) Layout() ImgLayout { return img.layout }

// HasOneImage returns true if img has exactly one image. The image may have a description.
//
// TODO: get rid of the concept.
func (img Img) HasOneImage() bool {
	return len(img.Entries) == 1
}

// WithExistingTargetsMarked returns a new Img with its ImgEntries colored according to their existence.
//
// This functions iterates over hyphae once.
func (img Img) WithExistingTargetsMarked(ctx mycocontext.Context) Img {
	var probes []func(string)
	for _, entry := range img.Entries {
		if probe := entry.Target.HyphaProbe(ctx); probe != nil {
			probes = append(probes, probe)
		}
	}
	ctx.Options().IterateHyphaNamesWith(func(hyphaName string) {
		for _, probe := range probes {
			probe(hyphaName)
		}
	})
	return img
}
