package blocks

// Quote is the block representing a quote.
type Quote struct {
	contents []Block
}

func (q Quote) IsBlock() {}

func (q *Quote) Contents() []Block {
	return q.contents
}

func (q *Quote) AddBlock(block Block) {
	q.contents = append(q.contents, block)
}
