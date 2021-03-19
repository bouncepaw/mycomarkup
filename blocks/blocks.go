package blocks

type BlockKind int

// From least complex to more complex:
const (
	KindNone BlockKind = iota
	KindHorizontalLine
	KindPreformatted
	KindFormatted
	KindHeading
	KindParagraph
	KindRocketLink
	KindLaunchpad
	KindList
	KindTable
)

type Block interface {
	String() string
	IsNesting() bool
	Kind() BlockKind
}
