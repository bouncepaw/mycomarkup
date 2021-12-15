package blocks

// ImgArrangement is an enumeration of possible Img internal arrangements.
//
// Keywords:
//     column   ImgArrangementColumn
//     grid     ImgArrangementGrid
type ImgArrangement int

const (
	// ImgArrangementColumn is the default arrangement. Images are in one column, text is below the images.
	ImgArrangementColumn ImgArrangement = iota
	// ImgArrangementGrid places the images in a grid. The grid size is unspecified.
	ImgArrangementGrid
)

// ImgPosition is an enumeration of possible Img external placements.
//
// Keywords:
//     stretch  ImgPositionStretch
//     start    ImgPositionStart
//     end      ImgPositionEnd
type ImgPosition int

const (
	// ImgPositionStretch is the default placement. The gallery is stretched horizontally.
	ImgPositionStretch ImgPosition = iota
	// ImgPositionStart places the gallery to the ‘start’. For LTR environments, that would be the left side. For RTL environments, that would be the right side.
	ImgPositionStart
	// ImgPositionEnd places the gallery to the ‘end’. For LTR environments, that would be the right side. For RTL environments, that would be the left side.
	ImgPositionEnd
)
