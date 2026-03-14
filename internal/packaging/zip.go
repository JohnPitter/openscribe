// Package packaging provides ZIP and OOXML packaging utilities.
package packaging

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// Package represents a ZIP-based Office document package
type Package struct {
	Files map[string][]byte
}

// NewPackage creates a new empty package
func NewPackage() *Package {
	return &Package{
		Files: make(map[string][]byte),
	}
}

// AddFile adds a file to the package
func (p *Package) AddFile(path string, data []byte) {
	p.Files[path] = data
}

// GetFile retrieves a file from the package
func (p *Package) GetFile(path string) ([]byte, bool) {
	data, ok := p.Files[path]
	return data, ok
}

// RemoveFile removes a file from the package
func (p *Package) RemoveFile(path string) {
	delete(p.Files, path)
}

// HasFile checks if a file exists in the package
func (p *Package) HasFile(path string) bool {
	_, ok := p.Files[path]
	return ok
}

// Save writes the package to a file
func (p *Package) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	return p.WriteTo(f)
}

// WriteTo writes the package to a writer
func (p *Package) WriteTo(w io.Writer) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	for name, data := range p.Files {
		fw, err := zw.Create(name)
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", name, err)
		}
		if _, err := fw.Write(data); err != nil {
			return fmt.Errorf("write zip entry %s: %w", name, err)
		}
	}

	return nil
}

// OpenPackage reads a ZIP package from a file
func OpenPackage(path string) (*Package, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return OpenPackageFromBytes(data)
}

// OpenPackageFromBytes reads a ZIP package from bytes
func OpenPackageFromBytes(data []byte) (*Package, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	pkg := NewPackage()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}
		fileData, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read zip entry %s: %w", f.Name, err)
		}
		pkg.Files[f.Name] = fileData
	}

	return pkg, nil
}

// ToBytes returns the package as a byte slice
func (p *Package) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := p.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
