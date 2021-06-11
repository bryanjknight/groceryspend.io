package ocr

// Service allows a caller to extract text from an image
type Service interface {
	// DetectTextInImage requests an OCR operation on the file
	DetectTextInImage(filePath string) (*Image, error)
}

// SortBlocksByLogicalOrder - sort all blocks by their top left point. the goal should be that all blocks
// are ordered from left to right and top to bottom
func SortBlocksByLogicalOrder(resp *Image, blocks []*Block) []*Block {

	// v1 - just re-apply the ordering the response had already

	// assuming only one relationship
	childrenIds := resp.BlockIDs

	newList := []*Block{}

	blockIDToBlock := make(map[string]*Block)
	for _, block := range blocks {
		blockIDToBlock[block.ID] = block
	}

	for _, childID := range childrenIds {

		if val, ok := blockIDToBlock[childID]; ok {
			newList = append(newList, val)
		}
	}

	return newList

}
