package blocks

import "fmt"

// Quote is the block representing a quote.
type Quote struct {
	contents []Block
}

func (q Quote) IsBlock() {}

func (q Quote) ID(counter *IDCounter) string {
	counter.quotes++
	return fmt.Sprintf("quote-%d")
}

func (q *Quote) Contents() []Block {
	return q.contents
}

func (q *Quote) AddBlock(block Block) {
	q.contents = append(q.contents, block)
}
