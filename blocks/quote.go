package blocks

import "github.com/bouncepaw/mycomarkup/util"

// Quote is the block representing a quote.
type Quote struct {
	contents string
}

func MakeQuote(line, hyphaName string) Quote {
	return Quote{
		ParagraphToHtml(hyphaName, util.Remover(">")(line)),
	}
}

func (q *Quote) Contents() string {
	return q.contents
}
