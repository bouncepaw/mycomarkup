package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/util"
	"strings"
)

// MatchesImg is true if the line starts with img {.
func MatchesImg(line string) bool {
	return strings.HasPrefix(line, `^ img {`)
}

type imgState int

const (
	InRoot imgState = iota
	InName
	InDimensionsW
	InDimensionsH
	InDescription
)

// Img is an image gallery, consisting of zero or more images.
type Img struct {
	// All entries but the last one
	Entries []ImgEntry
	// The last entry
	CurrEntry ImgEntry
	HyphaName string
	// Parsing state. TODO: move to a different place.
	State imgState
}

func (img Img) isBlock() {}

// ID returns the gallery's id which is img- and a number.
func (img Img) ID(counter *IDCounter) string {
	counter.imgs++
	return fmt.Sprintf("img-%d", counter.imgs)
}

// HasOneImage returns true if img has exactly one image and that images has no description.
func (img *Img) HasOneImage() bool {
	return len(img.Entries) == 1 && img.Entries[0].Desc.Len() == 0
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
