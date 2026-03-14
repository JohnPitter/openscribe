# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-03-14

### Added

#### Charts
- **PDF charts**: Bar, line, pie, area, horizontal bar charts rendered with PDF primitives
- **XLSX charts**: OOXML DrawingML charts embedded in spreadsheets (bar, line, pie, area, column, scatter)
- **PPTX charts**: Chart elements on presentation slides

#### Document (DOCX)
- **Headers/Footers**: Left, center, right text with font customization
- **Table of Contents**: Auto-built from headings with configurable max level
- **Image embedding**: Inline drawing with blip fill (previously API-only, now functional)
- **Paragraph indentation**: Left, right, first-line indent support

#### Spreadsheet (XLSX)
- **Conditional formatting**: 15 condition types (greater than, less than, between, contains, etc.), color scales, data bars
- **Formula evaluation**: SUM, AVERAGE, MIN, MAX, COUNT, ABS, ROUND with cell range resolution
- **Column management**: Width, hidden, best fit, range-based column configuration
- **Chart support**: Bar, line, pie, area, column, scatter with OOXML serialization

#### Presentation (PPTX)
- **Transition serialization**: Fade, push, wipe, split, cover, cut, dissolve now written to slide XML
- **Image embedding**: Images on slides with relationship management
- **Slide masters**: 6 pre-built layouts (Title, Title+Content, Two Column, Blank, Section, Title Only)
- **AddSlideFromLayout**: Create slides from master layouts with placeholder population

#### PDF
- **Split & ExtractPages**: Split documents at any page, extract specific pages
- **Text extraction**: Extract text from elements and raw PDF content streams
- **HTML-to-PDF**: Convert HTML to PDF with headings, paragraphs, lists, formatting, entities
- **Image elements**: Add images to PDF pages
- **Chart rendering**: 5 chart types with axes, grid lines, legends, value labels

#### CLI & Testing
- **CLI tool**: `cmd/openscribe` with -list, -themes, -template, -output, -format, -level flags
- **Template Generate()**: All 32 templates now generate real documents with rich placeholder content
- **Benchmark tests**: Performance benchmarks for all formats
- **Common test coverage**: Increased from 50% to 95.3%

### Changed
- Renamed `MarshalXML()` to `toXML()` in document package (fixes go vet warnings)
- Renamed `WriteTo()` to `writeToWriter()` in packaging (fixes go vet interface mismatch)

### Fixed
- Go vet warnings eliminated (3 XML marshaling interface violations)
- CI workflow: aligned go.mod version, replaced incompatible golangci-lint with go vet + gofmt
- All code formatted with `gofmt -s`

## [0.1.0] - 2026-03-13

### Added

- **common package**: Color, Font, Measurement, Border, Image types with full conversion utilities
- **style package**: Design system with Theme, Palette, Typography; 6 pre-built themes across 4 design levels (Basic, Professional, Premium, Luxury)
- **internal/packaging**: ZIP-based OOXML packaging, relationships, and content types management
- **internal/xmlutil**: XML marshaling/unmarshaling helpers with proper declarations
- **document package**: Full DOCX support — create, open, edit, save, delete documents with paragraphs, headings, tables, images, sections, page breaks
- **spreadsheet package**: Full XLSX support — create, open, edit, save, delete workbooks with sheets, cells (string/number/boolean/formula), merged cells, formatting
- **presentation package**: Full PPTX support — create, open, edit, save, delete presentations with slides, text boxes, shapes (12 types), transitions, speaker notes
- **pdf package**: Pure Go PDF generation — create, save, delete PDFs with text, lines, rectangles, tables, watermarks, page backgrounds, document merging
- **template package**: Pre-built document templates at all design levels (Basic → Luxury) for all formats
- **E2E tests**: Full coverage for create, edit, delete operations across all formats and design levels
- **CI/CD**: GitHub Actions workflow for automated testing
- Project setup: README, LICENSE (MIT), Makefile, .gitignore
