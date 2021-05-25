// Package blocks provides some of the Mycomarkup's blocks.
package blocks

// Block is a unit of Mycomarkup. It is somewhat analogous to HTML's tags.
type Block interface {
	// IsBlock is the method that a block should implement. In the future, this method might be removed, but for now, we need it just to put something in the interface.
	IsBlock()
}
