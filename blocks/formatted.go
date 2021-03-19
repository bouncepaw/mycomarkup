package blocks

import (
	"fmt"
)

// TODO:
type Formatted struct {
	src string
}

func MakeFormatted(src string) *Formatted {
	return &Formatted{src}
}

func (f *Formatted) String() string {
	return fmt.Sprintf(`Formatted("%s");`, f.src)
}

func (f *Formatted) IsNesting() bool {
	return false
}

func (f *Formatted) Kind() BlockKind {
	return KindFormatted
}
