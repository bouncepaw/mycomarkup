package blocks

import (
	"fmt"
)

type Paragraph struct {
	content *Formatted
}

func MakeParagraph(content *Formatted) *Paragraph {
	return &Paragraph{content}
}

func (p *Paragraph) String() string {
	return fmt.Sprintf(`Paragraph() {
%s
};`, p.content)
}

func (p *Paragraph) IsNesting() bool {
	return false
}

func (p *Paragraph) Kind() BlockKind {
	return KindParagraph
}
