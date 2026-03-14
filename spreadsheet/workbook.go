// Package spreadsheet provides XLSX spreadsheet creation, reading, and editing.
package spreadsheet

import (
	"fmt"
	"os"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/style"
)

// Workbook represents an XLSX workbook
type Workbook struct {
	pkg           *packaging.Package
	sheets        []*Sheet
	sharedStrings []string
	stringIndex   map[string]int
	theme         *style.Theme
	namedRanges   []*NamedRange
	security      *common.SecurityOptions
}

// New creates a new empty workbook
func New() *Workbook {
	theme := style.BasicClean()
	wb := &Workbook{
		pkg:         packaging.NewPackage(),
		stringIndex: make(map[string]int),
		theme:       &theme,
	}
	return wb
}

// NewWithTheme creates a workbook with a specific theme
func NewWithTheme(theme style.Theme) *Workbook {
	wb := New()
	wb.theme = &theme
	return wb
}

// Theme returns the current theme
func (wb *Workbook) Theme() *style.Theme {
	return wb.theme
}

// SetTheme sets the workbook theme
func (wb *Workbook) SetTheme(theme style.Theme) {
	wb.theme = &theme
}

// AddSheet adds a new worksheet
func (wb *Workbook) AddSheet(name string) *Sheet {
	s := newSheet(wb, name, len(wb.sheets)+1)
	wb.sheets = append(wb.sheets, s)
	return s
}

// Sheet returns a sheet by index (0-based)
func (wb *Workbook) Sheet(index int) *Sheet {
	if index < 0 || index >= len(wb.sheets) {
		return nil
	}
	return wb.sheets[index]
}

// SheetByName returns a sheet by name
func (wb *Workbook) SheetByName(name string) *Sheet {
	for _, s := range wb.sheets {
		if s.name == name {
			return s
		}
	}
	return nil
}

// SheetCount returns the number of sheets
func (wb *Workbook) SheetCount() int {
	return len(wb.sheets)
}

// RemoveSheet removes a sheet by index
func (wb *Workbook) RemoveSheet(index int) error {
	if index < 0 || index >= len(wb.sheets) {
		return fmt.Errorf("sheet index %d out of range", index)
	}
	wb.sheets = append(wb.sheets[:index], wb.sheets[index+1:]...)
	return nil
}

// addSharedString adds a string to the shared string table and returns its index
func (wb *Workbook) addSharedString(s string) int {
	if idx, ok := wb.stringIndex[s]; ok {
		return idx
	}
	idx := len(wb.sharedStrings)
	wb.sharedStrings = append(wb.sharedStrings, s)
	wb.stringIndex[s] = idx
	return idx
}

// Save writes the workbook to a file
func (wb *Workbook) Save(path string) error {
	if err := wb.build(); err != nil {
		return fmt.Errorf("build workbook: %w", err)
	}
	return wb.pkg.Save(path)
}

// SaveToBytes returns the workbook as bytes
func (wb *Workbook) SaveToBytes() ([]byte, error) {
	if err := wb.build(); err != nil {
		return nil, fmt.Errorf("build workbook: %w", err)
	}
	return wb.pkg.ToBytes()
}

// Open reads an XLSX file
func Open(path string) (*Workbook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return OpenFromBytes(data)
}

// OpenFromBytes reads an XLSX from bytes
func OpenFromBytes(data []byte) (*Workbook, error) {
	pkg, err := packaging.OpenPackageFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}

	wb := &Workbook{
		pkg:         pkg,
		stringIndex: make(map[string]int),
	}

	// Parse shared strings
	if ssData, ok := pkg.GetFile("xl/sharedStrings.xml"); ok {
		wb.parseSharedStrings(ssData)
	}

	// Parse workbook to get sheet names
	if wbData, ok := pkg.GetFile("xl/workbook.xml"); ok {
		wb.parseWorkbook(wbData)
	}

	return wb, nil
}

// Delete removes an XLSX file from disk
func Delete(path string) error {
	return os.Remove(path)
}
