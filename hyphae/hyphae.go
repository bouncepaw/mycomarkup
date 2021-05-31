// Package hyphae contains things related to working with Mycorrhiza's hyphae.
package hyphae

// Hypha represents a hypha, duh. Only a subset of hyphal operations is in the interface.
type Hypha interface {
	// NameString returns hypha's canonical name.
	NameString() string

	// DoesExist is true when the hypha does exist.
	DoesExist() bool

	// HasText is true when the hypha has a text part.
	HasText() bool

	// HasAttachment is true when the hypha has an attachment.
	HasAttachment() bool
}

// Visitor is a function that looks at the Hypha and, judging by its content, does something. It does not modify the Hypha. These functions are usually packed together in VisitorPack. The function should return false when it does not need to check any hyphae anymore.
type Visitor func(Hypha) (keep bool)

// VisitorPack is a group of pending Visitor operations.
type VisitorPack struct {
	// HyphaIterator is a channel giving out many a Hypha in no particular order. The channel is read until the end on every iteration.
	HyphaIterator chan Hypha

	visitors []Visitor
}

func (pack *VisitorPack) AddVisitor(visitor Visitor) {
	pack.visitors = append(pack.visitors, visitor)
}

func (pack *VisitorPack) Iterate() {
	for hypha := range pack.HyphaIterator {
		for i, visitor := range pack.visitors {
			if keep := visitor(hypha); !keep {
				pack.removeVisitor(i)
			}
		}
	}
}

func (pack *VisitorPack) removeVisitor(i int) {
	size := len(pack.visitors)
	pack.visitors[i] = pack.visitors[size-1]
	pack.visitors = pack.visitors[:size-1]
}
