// Package blocks provides some of the Mycomarkup's blocks.
package blocks

// Block is a unit of Mycomarkup. It is somewhat analogous to HTML's tags.
type Block interface {
	// String returns a debug string representation of the block.
	String() string

	// ID returns an id for the block which may be utilised in markup languages. It may not be unique.
	ID() string
}

type NestingBlock struct{}

func (ns *NestingBlock) IsNesting() bool {
	return true
}

type TerminalBlock struct{}

func (tb *TerminalBlock) IsNesting() bool {
	return false
}
