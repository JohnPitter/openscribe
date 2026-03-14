# OpenScribe Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a pure Go open-source library for creating, editing, and manipulating office documents (DOCX, XLSX, PPTX, PDF) with support for basic to premium design levels.

**Architecture:** Modular package design with `common` for shared types (colors, fonts, measurements, styles), format-specific packages (`document`, `spreadsheet`, `presentation`, `pdf`), a `template` package for pre-built designs ranging from basic to premium quality, and a `style` package for design system themes. Each format package exposes a high-level API for common operations and raw access for advanced manipulation.

**Tech Stack:** Go 1.22+, encoding/xml for OOXML, archive/zip for packaging, go-pdf for PDF generation, testify for assertions

---

## Architecture Overview

```
openscribe/
├── common/          # Shared types: Color, Font, Measurement, Border, Alignment
├── style/           # Design system: Theme, Palette, Typography, Spacing
├── template/        # Pre-built templates: Basic, Professional, Premium, Luxury
├── document/        # DOCX: Create, Read, Edit, Delete, Styles, Tables, Images
├── spreadsheet/     # XLSX: Create, Read, Edit, Delete, Formulas, Charts, Styles
├── presentation/    # PPTX: Create, Read, Edit, Delete, Slides, Animations
├── pdf/             # PDF: Create, Read, Edit, Merge, Split, Watermark
├── internal/        # Internal utilities: XML helpers, zip packaging
├── e2e/             # End-to-end tests
├── testdata/        # Test fixtures
└── cmd/openscribe/  # CLI tool
```

---

## Task 1: Common Package — Shared Types & Design System

**Files:**
- Create: `common/color.go`
- Create: `common/font.go`
- Create: `common/measurement.go`
- Create: `common/border.go`
- Create: `common/alignment.go`
- Create: `common/image.go`
- Create: `common/common.go`
- Test: `common/color_test.go`
- Test: `common/font_test.go`
- Test: `common/measurement_test.go`

## Task 2: Style Package — Design Themes & Palettes

**Files:**
- Create: `style/theme.go`
- Create: `style/palette.go`
- Create: `style/typography.go`
- Create: `style/presets.go` (Basic, Professional, Premium, Luxury, Behance, Freepik, Slidesgo)
- Test: `style/theme_test.go`
- Test: `style/presets_test.go`

## Task 3: Internal Package — XML & Zip Utilities

**Files:**
- Create: `internal/xmlutil/marshal.go`
- Create: `internal/xmlutil/unmarshal.go`
- Create: `internal/packaging/zip.go`
- Create: `internal/packaging/relationships.go`
- Create: `internal/packaging/content_types.go`
- Test: `internal/xmlutil/marshal_test.go`
- Test: `internal/packaging/zip_test.go`

## Task 4: Document Package — DOCX Support

**Files:**
- Create: `document/document.go` (New, Open, Save, Close)
- Create: `document/paragraph.go` (Add, Style, Align)
- Create: `document/run.go` (Text, Bold, Italic, Font, Color)
- Create: `document/table.go` (Create, Rows, Cells, Style)
- Create: `document/image.go` (Insert, Resize, Position)
- Create: `document/header_footer.go`
- Create: `document/section.go` (PageSize, Margins, Orientation)
- Create: `document/styles.go` (Heading1-6, Body, Quote, Code)
- Test: `document/document_test.go`
- Test: `document/paragraph_test.go`
- Test: `document/table_test.go`

## Task 5: Spreadsheet Package — XLSX Support

**Files:**
- Create: `spreadsheet/workbook.go` (New, Open, Save, Close)
- Create: `spreadsheet/sheet.go` (Add, Remove, Rename)
- Create: `spreadsheet/cell.go` (Value, Formula, Style)
- Create: `spreadsheet/row.go` (Add, Height, Style)
- Create: `spreadsheet/column.go` (Width, Style, AutoFit)
- Create: `spreadsheet/styles.go` (NumberFormat, Fill, Border, Font)
- Create: `spreadsheet/chart.go` (Bar, Line, Pie, Area)
- Create: `spreadsheet/merge.go` (MergeCells, UnmergeCells)
- Test: `spreadsheet/workbook_test.go`
- Test: `spreadsheet/cell_test.go`
- Test: `spreadsheet/chart_test.go`

## Task 6: Presentation Package — PPTX Support

**Files:**
- Create: `presentation/presentation.go` (New, Open, Save, Close)
- Create: `presentation/slide.go` (Add, Remove, Reorder, Layout)
- Create: `presentation/textbox.go` (Create, Style, Position)
- Create: `presentation/shape.go` (Rectangle, Circle, Arrow, Custom)
- Create: `presentation/image.go` (Insert, Crop, Position)
- Create: `presentation/transition.go` (Fade, Slide, Zoom)
- Create: `presentation/master.go` (SlideMaster, Layout templates)
- Test: `presentation/presentation_test.go`
- Test: `presentation/slide_test.go`
- Test: `presentation/shape_test.go`

## Task 7: PDF Package — PDF Support

**Files:**
- Create: `pdf/pdf.go` (New, Open, Save, Close)
- Create: `pdf/page.go` (Add, Remove, Size, Margins)
- Create: `pdf/text.go` (Write, Style, Position)
- Create: `pdf/table.go` (Create, Style, Cells)
- Create: `pdf/image.go` (Insert, Position, Scale)
- Create: `pdf/watermark.go` (Text, Image, Opacity)
- Create: `pdf/merge.go` (Merge, Split, Extract)
- Create: `pdf/header_footer.go`
- Test: `pdf/pdf_test.go`
- Test: `pdf/page_test.go`
- Test: `pdf/table_test.go`

## Task 8: Template Package — Pre-built Designs

**Files:**
- Create: `template/template.go` (Interface, Registry)
- Create: `template/basic.go` (Clean minimal designs)
- Create: `template/professional.go` (Business-grade designs)
- Create: `template/premium.go` (Behance/Freepik quality)
- Create: `template/luxury.go` (Slidesgo/high-end quality)
- Create: `template/catalog.go` (List, Get, Preview templates)
- Test: `template/template_test.go`
- Test: `template/catalog_test.go`

## Task 9: E2E Tests — Full Coverage

**Files:**
- Create: `e2e/document_e2e_test.go` (Create, Edit, Delete DOCX)
- Create: `e2e/spreadsheet_e2e_test.go` (Create, Edit, Delete XLSX)
- Create: `e2e/presentation_e2e_test.go` (Create, Edit, Delete PPTX)
- Create: `e2e/pdf_e2e_test.go` (Create, Edit, Delete, Merge, Split PDF)
- Create: `e2e/template_e2e_test.go` (All template levels for all formats)
- Create: `e2e/design_levels_e2e_test.go` (Basic → Luxury design quality)

## Task 10: README, CI, and GitHub Setup

**Files:**
- Create: `README.md`
- Create: `.github/workflows/ci.yml`
- Create: `.gitignore`
- Create: `LICENSE`
- Create: `Makefile`
- Create: `CHANGELOG.md`
- Create: `CONTRIBUTING.md`
