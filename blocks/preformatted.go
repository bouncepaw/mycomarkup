package blocks

import (
	"fmt"
)

type Preformatted struct {
	altText string
	content string
}

func MakePreformatted(altText, content string) *Preformatted {
	return &Preformatted{altText, content}
}

func (pre *Preformatted) String() string {
	return fmt.Sprintf(`Preformatted("%s") {
%s
};`, pre.altText, pre.content)
}

func (pre *Preformatted) IsNesting() bool {
	return false
}

func (pre *Preformatted) Kind() BlockKind {
	return KindPreformatted
}
