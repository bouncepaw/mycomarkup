// Package blocks provides some of the Mycomarkup's blocks.
package blocks

// Block is a unit of Mycomarkup. It is somewhat analogous to HTML's tags.
type Block interface {
	// ID returns an id for the block. It should be unique when possible. The block should increment a value in the counter depending on its type.
	ID(counter *IDCounter) string
}

// IDCounter is a struct with counters of how many times some blocks have appeared. Block's ID depends on these counters.
type IDCounter struct {
	// In some cases using the results of counting is not needed because the IDs are not needed themselves. This variable is true when this is the case.
	ShouldUseResults bool
	codeblocks       uint
	hrs              uint
	imgs             uint
	launchpads       uint
	lists            uint
	paragraphs       uint
	quotes           uint
	tables           uint
	transclusions    uint
}

// UnusableCopy returns a copy of the counter with ShouldUseResults set to false.
func (c IDCounter) UnusableCopy() *IDCounter {
	c.ShouldUseResults = false
	return &c
}
