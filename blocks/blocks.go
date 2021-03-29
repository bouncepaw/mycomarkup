package blocks

type BlockKind int

const (
	KindNone BlockKind = iota
	KindHorizontalLine
	KindPreformatted
	KindFormatted
	KindHeading
	KindParagraph

	KindRocketLink
	KindLaunchPad

	KindMessage

	KindImage
	KindImageGallery

	KindTableCell
	KindTableRow
	KindTable

	KindListElement
	KindListIndent
	KindList

	KindBlockQuote
)

type Block interface {
	String() string
	IsNesting() bool
	Kind() BlockKind
}
