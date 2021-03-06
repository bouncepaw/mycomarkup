package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v2/globals"
	"github.com/bouncepaw/mycomarkup/v2/util"
)

// Img is an image gallery, consisting of zero or more images.
type Img struct {
	// All entries
	Entries   []ImgEntry
	HyphaName string
}

func (img Img) isBlock() {}

// ID returns the gallery's id which is img- and a number.
func (img Img) ID(counter *IDCounter) string {
	counter.imgs++
	return fmt.Sprintf("img-%d", counter.imgs)
}

// HasOneImage returns true if img has exactly one image. The image may have a description.
func (img *Img) HasOneImage() bool {
	return len(img.Entries) == 1
}

// MarkExistenceOfSrcLinks effectively checks if the links in the gallery are blue or red.
func (img *Img) MarkExistenceOfSrcLinks() {
	globals.HyphaIterate(func(hn string) {
		for _, entry := range img.Entries {
			if hn == util.CanonicalName(entry.Srclink.TargetHypha()) {
				entry.Srclink.MarkAsExisting()
			}
		}
	})
}
