package blocks

import "fmt"

// Quote is the block representing a quote.
type Quote struct {
	contents []Block
}

// NewQuote returns Quote with the given contents.
func NewQuote(contents []Block) Quote {
	return Quote{contents: contents}
}

// ID returns the quote's id which is quote- and a number.
func (q Quote) ID(counter *IDCounter) string {
	counter.quotes++
	return fmt.Sprintf("quote-%d", counter.quotes)
}

// Contents returns the quote's contents.
func (q *Quote) Contents() []Block {
	return q.contents
}
