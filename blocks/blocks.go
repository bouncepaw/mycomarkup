// Package blocks provides some of the Mycomarkup's blocks.
package blocks

// Block is a unit of Mycomarkup. It is somewhat analogous to HTML's tags.
type Block interface {
	// ID returns an id for the block. It should be unique when possible. The block should increment a value in the counter depending on its type.
	ID(counter *IDCounter) string
}

// IDCounter is a struct with counters of how many times some blocks have appeared. Block's ID depends on these counters. IDCounter is not a Block.
type IDCounter struct {
	// In some cases using the results of counting is not needed because the IDs are not needed themselves. This variable is false when this is the case. Do not modify it directly!
	shouldUseResults bool
	// Increment the fields below. You should not decrement them or set them to a specific value.
	codeblocks    uint
	hrs           uint
	imgs          uint
	launchpads    uint
	lists         uint
	paragraphs    uint
	quotes        uint
	tables        uint
	transclusions uint
}

// NewIDCounter returns a pointer to an IDCounter. ShouldUseResults is true for this counter.
func NewIDCounter() *IDCounter {
	return &IDCounter{
		shouldUseResults: true,
	}
}

// ShouldUseResults is true if you should use generated IDs.
func (c IDCounter) ShouldUseResults() bool {
	return c.shouldUseResults
}

// UnusableCopy returns a pointer to a copy of the counter with shouldUseResults set to false.
func (c IDCounter) UnusableCopy() *IDCounter {
	c.shouldUseResults = false
	return &c
}
