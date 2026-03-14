package document

import "github.com/JohnPitter/openscribe/common"

// ImageRef represents an image reference in the document
type ImageRef struct {
	id     string
	data   *common.ImageData
	width  common.Measurement
	height common.Measurement
	relID  string
}

// ID returns the image identifier
func (img *ImageRef) ID() string { return img.id }

// Width returns the image width
func (img *ImageRef) Width() common.Measurement { return img.width }

// Height returns the image height
func (img *ImageRef) Height() common.Measurement { return img.height }

// SetSize sets the image dimensions
func (img *ImageRef) SetSize(width, height common.Measurement) {
	img.width = width
	img.height = height
}
