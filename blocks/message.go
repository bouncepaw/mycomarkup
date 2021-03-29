package blocks

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/links"
)

type Message struct {
	author *links.Link
	body   []Block
}

func MakeMessage(author *links.Link, body []Block) *Message {
	return &Message{author, body}
}

func (msg *Message) String() string {
	var s string
	for _, child := range msg.body {
		s += child.String() + "\n"
	}
	return fmt.Sprintf(`Message(%s) {
%s};`, msg.author, s)
}

func (msg *Message) IsNesting() bool {
	return true
}

func (msg *Message) Kind() BlockKind {
	return KindMessage
}
