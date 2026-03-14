package presentation

import "github.com/JohnPitter/openscribe/common"

// SlideImage represents an image on a slide
type SlideImage struct {
	data   *common.ImageData
	x, y   common.Measurement
	width  common.Measurement
	height common.Measurement
	relID  string
}

func (img *SlideImage) elementType() string { return "image" }

// SetPosition sets the image position
func (img *SlideImage) SetPosition(x, y common.Measurement) { img.x = x; img.y = y }

// SetSize sets the image size
func (img *SlideImage) SetSize(width, height common.Measurement) { img.width = width; img.height = height }
