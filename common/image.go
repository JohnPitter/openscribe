package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ImageFormat supported image formats
type ImageFormat int

const (
	ImageFormatPNG ImageFormat = iota
	ImageFormatJPEG
	ImageFormatGIF
	ImageFormatBMP
	ImageFormatSVG
	ImageFormatTIFF
)

func (f ImageFormat) Extension() string {
	switch f {
	case ImageFormatPNG:
		return ".png"
	case ImageFormatJPEG:
		return ".jpeg"
	case ImageFormatGIF:
		return ".gif"
	case ImageFormatBMP:
		return ".bmp"
	case ImageFormatSVG:
		return ".svg"
	case ImageFormatTIFF:
		return ".tiff"
	default:
		return ".png"
	}
}

func (f ImageFormat) MimeType() string {
	switch f {
	case ImageFormatPNG:
		return "image/png"
	case ImageFormatJPEG:
		return "image/jpeg"
	case ImageFormatGIF:
		return "image/gif"
	case ImageFormatBMP:
		return "image/bmp"
	case ImageFormatSVG:
		return "image/svg+xml"
	case ImageFormatTIFF:
		return "image/tiff"
	default:
		return "image/png"
	}
}

// ImageData holds raw image data
type ImageData struct {
	Data   []byte
	Format ImageFormat
	Width  Measurement
	Height Measurement
}

// LoadImage loads an image from a file path
func LoadImage(path string) (*ImageData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var format ImageFormat
	switch ext {
	case ".png":
		format = ImageFormatPNG
	case ".jpg", ".jpeg":
		format = ImageFormatJPEG
	case ".gif":
		format = ImageFormatGIF
	case ".bmp":
		format = ImageFormatBMP
	case ".svg":
		format = ImageFormatSVG
	case ".tiff", ".tif":
		format = ImageFormatTIFF
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	return &ImageData{
		Data:   data,
		Format: format,
	}, nil
}
