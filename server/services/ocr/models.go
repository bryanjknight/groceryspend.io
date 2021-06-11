package ocr

// Note: X goes from left to right, Y goes from top to bottom. Therefore
//       (0,0) is the top left, (N, 0) is the top right, (0,M) is the bottom left,
// 			 and (N, M) is bottom right

// TODO: we do multiple O(N) operations against the array of blocks
// an optimization could be to better structure the data for easier
// lookup; however, perf is not a driver at this point given the size of the
// receipts are small and take very little time (<500ms)

// we can use a tree to create an quick way to look up where the item is
// based on the block's top location. For example:
// t[0.0,1.0] = all blocks (where 1000 is the max pixel location)
// t[.5,1.0] == all blocks on the bottom half
// this could be combined with a similar tree but using the left/right position
// then doing a union on the resulting two sets to see what actually fits within
// those coordinates

// Point represents a location on a x-y coordinate plan
type Point struct {
	X float64
	Y float64
}

// Block is a quadrilateral encompassing a section of text
type Block struct {
	// unique ID for this string
	ID string
	// Top Left corner of block on a x,y plane
	TopLeft *Point
	// Top Right corner of block on a x,y plane
	TopRight *Point
	// Bottom Right corner of block on a x,y plane
	BottomRight *Point
	// Bottom Left corner of block on a x,y plane
	BottomLeft *Point
	// The text in the box
	Text string
	// Confidence from 0-100 that the text is accurage
	Confidence float64
}

// Image represents the location of text on an image from an OCR service
// TODO: optimize the image struct to make look ups faster
type Image struct {
	// The blocks represneting the text in the image. Note that the blocks try to be
	// in left to right, top to bottom but is not guaranteed
	Blocks []*Block

	// original order of block ids, useful when you need to order blocks based on
	// the original order
	BlockIDs []string

	// original response from ocr service
	OriginalResponse interface{}
}
