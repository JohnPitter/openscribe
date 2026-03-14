package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImageFormats(t *testing.T) {
	formats := []struct {
		fmt  ImageFormat
		ext  string
		mime string
	}{
		{ImageFormatPNG, ".png", "image/png"},
		{ImageFormatJPEG, ".jpeg", "image/jpeg"},
		{ImageFormatGIF, ".gif", "image/gif"},
		{ImageFormatBMP, ".bmp", "image/bmp"},
		{ImageFormatSVG, ".svg", "image/svg+xml"},
		{ImageFormatTIFF, ".tiff", "image/tiff"},
	}
	for _, f := range formats {
		if f.fmt.Extension() != f.ext {
			t.Errorf("expected ext %s, got %s", f.ext, f.fmt.Extension())
		}
		if f.fmt.MimeType() != f.mime {
			t.Errorf("expected mime %s, got %s", f.mime, f.fmt.MimeType())
		}
	}
}

func TestImageFormatDefault(t *testing.T) {
	var unknown ImageFormat = 99
	if unknown.Extension() != ".png" {
		t.Error("unknown format should default to .png")
	}
	if unknown.MimeType() != "image/png" {
		t.Error("unknown mime should default to image/png")
	}
}

func TestLoadImage(t *testing.T) {
	// Create temp PNG file
	tmpDir := t.TempDir()
	pngPath := filepath.Join(tmpDir, "test.png")
	// Minimal PNG data (just header for testing)
	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if err := os.WriteFile(pngPath, pngData, 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	img, err := LoadImage(pngPath)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if img.Format != ImageFormatPNG {
		t.Error("should be PNG format")
	}
	if len(img.Data) != len(pngData) {
		t.Error("data length mismatch")
	}
}

func TestLoadImageFormats(t *testing.T) {
	tmpDir := t.TempDir()
	data := []byte{0x00, 0x01, 0x02}

	exts := map[string]ImageFormat{
		".png":  ImageFormatPNG,
		".jpg":  ImageFormatJPEG,
		".jpeg": ImageFormatJPEG,
		".gif":  ImageFormatGIF,
		".bmp":  ImageFormatBMP,
		".svg":  ImageFormatSVG,
		".tiff": ImageFormatTIFF,
		".tif":  ImageFormatTIFF,
	}

	for ext, expectedFmt := range exts {
		path := filepath.Join(tmpDir, "test"+ext)
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("write error: %v", err)
		}
		img, err := LoadImage(path)
		if err != nil {
			t.Fatalf("load %s error: %v", ext, err)
		}
		if img.Format != expectedFmt {
			t.Errorf("ext %s: expected format %d, got %d", ext, expectedFmt, img.Format)
		}
	}
}

func TestLoadImageErrors(t *testing.T) {
	// Non-existent file
	_, err := LoadImage("/nonexistent/file.png")
	if err == nil {
		t.Error("should error on non-existent file")
	}

	// Unsupported format
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.xyz")
	os.WriteFile(path, []byte{0x00}, 0644)
	_, err = LoadImage(path)
	if err == nil {
		t.Error("should error on unsupported format")
	}
}

func TestColorString(t *testing.T) {
	c := NewColor(255, 128, 0)
	if c.String() != "#FF8000" {
		t.Errorf("expected #FF8000, got %s", c.String())
	}
}

func TestPixelConversion(t *testing.T) {
	m := Px(96)
	if !almostEqual(m.Inches(), 1.0, 0.01) {
		t.Errorf("96px should be ~1 inch, got %f", m.Inches())
	}

	m2 := In(1)
	if !almostEqual(m2.Pixels(), 96, 0.5) {
		t.Errorf("1 inch should be ~96px, got %f", m2.Pixels())
	}
}
