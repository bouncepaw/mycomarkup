package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v4/links"
	"github.com/bouncepaw/mycomarkup/v4/mycocontext"
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
	// bouncepaw: I'm so sorry this function is this complex.

	// We create this structure to keep track of what targets we have ‘ticked‘ ✅.
	// We do not compare hypha names with ticked targets.
	// Important: the structure retains the same order as the original img.Entries.
	type check struct {
		shouldCheck bool
		target      links.Link
	}
	var entryCheckList []check
	for _, entry := range img.Entries {
		entryCheckList = append(entryCheckList, check{
			shouldCheck: entry.Target.OfKind(links.LinkLocalHypha), // Other kinds are blue by definition
			target:      entry.Target,
		})
	}

	mycocontext.IterateHyphaNamesWith(ctx, func(hn string) {
		// Go through every entry and mark them accordingly.
		for i, entryCheck := range entryCheckList {
			shouldCheck, target := entryCheck.shouldCheck, entryCheck.target
			if shouldCheck && hn == target.TargetHypha() {
				entryCheckList[i] = check{
					shouldCheck: false,
					target:      target.CopyMarkedAsExisting(),
				}
			}
		}
	})

	// Collect the results. Some entries are left unmarked. It means they are red.
	var entries []ImgEntry
	for i, entry := range img.Entries {
		// Indices of entryCheckList and img.Entries are the same for the corresponding elements.
		entry.Target = entryCheckList[i].target
		entries = append(entries, entry)
	}

	return NewImg(entries, img.Layout())
}
