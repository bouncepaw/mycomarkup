package blocks

import (
	"fmt"
)

type LaunchPad struct {
	rockets []*RocketLink
}

func MakeLaunchPad(rockets ...*RocketLink) *LaunchPad {
	return &LaunchPad{rockets}
}

func (lp *LaunchPad) String() string {
	var s string
	for i, rocket := range lp.rockets {
		if i > 0 {
			s += "\n"
		}
		s += rocket.String()
	}
	return fmt.Sprintf(`LaunchPad() {
%s
};`, s)
}

func (lp *LaunchPad) IsNesting() bool {
	return false
}

func (lp *LaunchPad) Kind() BlockKind {
	return KindLaunchPad
}
