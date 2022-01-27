package blocks

// ImgLayout represents the layout of the Img element
type ImgLayout int

const (
	// ImgLayoutNormal is a main-column-wide one-column stack of images.
	ImgLayoutNormal ImgLayout = iota
	// ImgLayoutGrid is main-column-wide two-column stack of images.
	ImgLayoutGrid
	// ImgLayoutSide is a thin right-floating stack of images.
	ImgLayoutSide
)

func (l ImgLayout) String() string {
	switch l {
	case ImgLayoutSide:
		return "side"
	case ImgLayoutGrid:
		return "grid"
	default:
		return "normal"
	}
}
