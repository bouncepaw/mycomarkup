package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/lexer"
)

type Paragraph struct {
	tokens []lexer.Token
}

func (p *Paragraph) String() string {
	return fmt.Sprintf(`Paragraph() {
%s
};`, p.tokens)
}

func (p *Paragraph) IsNesting() bool {
	return false
}

func (p *Paragraph) Kind() BlockKind {
	return KindParagraph
}

func (p *Paragraph) AsHTML() string {
	for _, _ = range p.tokens {

	}
	return ""
}
