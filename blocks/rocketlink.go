package blocks

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/links"
)

type RocketLink struct {
	link *links.Link
}

func MakeRocketLink(link *links.Link) *RocketLink {
	return &RocketLink{link}
}

func (rl *RocketLink) String() string {
	return fmt.Sprintf(`Rocket%s;`, rl.link.String())
}

func (rl *RocketLink) IsNesting() bool {
	return false
}

func (rl *RocketLink) Kind() BlockKind {
	return KindRocketLink
}
