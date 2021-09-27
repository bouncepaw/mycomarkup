package blocks

import "fmt"

// Quote is the block representing a quote.
type Quote struct {
	contents []Block
}

func (q Quote) isBlock() {}

// ID returns the quote's id which is quote- and a number.
func (q Quote) ID(counter *IDCounter) string {
	counter.quotes++
	return fmt.Sprintf("quote-%d", counter.quotes)
}

// Contents returns the quote's contents.
func (q *Quote) Contents() []Block {
	return q.contents
}

// AddBlock adds the block to the quote. V3
func (q *Quote) AddBlock(block Block) {
	q.contents = append(q.contents, block)
}
