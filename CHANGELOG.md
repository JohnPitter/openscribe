# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
