# OpenScribe Roadmap

This document outlines planned improvements and new features for OpenScribe.

## Current State (v0.2.0)

OpenScribe is a fully functional, pure Go office document library with:
- 97 Go source files, 449 tests, 85%+ coverage
- 4 document formats: DOCX, XLSX, PPTX, PDF
- 6 design themes across 4 quality levels
- 32 templates with Generate() methods
- CLI tool for template-based generation
- Charts, conditional formatting, formula evaluation
- HTML-to-PDF conversion, text extraction
- Zero external dependencies

---

## Phase 1: Core Improvements (v0.3.0)

### Document (DOCX)
- [ ] **Footnotes & endnotes** — Add footnote/endnote support with auto-numbering
- [ ] **Numbered & bulleted lists** — Native list support with nesting levels (currently workaround via paragraphs)
- [ ] **Hyperlinks** — Clickable links within document text
- [ ] **Table styling presets** — Pre-built table styles (striped rows, header bands, etc.)
- [ ] **Comments** — Add review comments to document elements
- [ ] **Track changes** — Basic tracked changes support for collaborative editing
- [ ] **Custom styles** — User-defined paragraph and character styles beyond headings

### Spreadsheet (XLSX)
- [ ] **Data validation** — Dropdown lists, number ranges, custom validation rules per cell
- [ ] **Cell comments** — Add/read comments on cells
- [ ] **Sheet protection** — Password-protect sheets and lock specific cells
- [ ] **Auto-filter** — Enable column filters on data ranges
- [ ] **Freeze panes** — Freeze rows/columns for scrolling
- [ ] **Print settings** — Print area, page breaks, headers/footers, fit-to-page
- [ ] **Named ranges** — Define and use named ranges in formulas
- [ ] **More chart types** — Donut, radar, stock, combo charts

### Presentation (PPTX)
- [ ] **Animations** — Element entrance/exit/emphasis animations with timing
- [ ] **Embedded video/audio** — Media embedding with playback controls
- [ ] **Smart shapes** — Connector lines, callouts, flowchart shapes
- [ ] **Speaker notes formatting** — Rich text speaker notes (currently plain text)
- [ ] **Slide numbers** — Automatic slide numbering
- [ ] **Tables in slides** — Native table support within slides

### PDF
- [ ] **Font embedding** — Embed TrueType/OpenType fonts for consistent rendering
- [ ] **PDF/A compliance** — Generate archival-standard PDF/A documents
- [ ] **Digital signatures** — Sign PDFs with X.509 certificates
- [ ] **Form fields** — Create interactive PDF forms (text fields, checkboxes, dropdowns)
- [ ] **Annotations** — Add highlights, sticky notes, and markup annotations
- [ ] **Advanced text layout** — Word wrapping, text columns, justified text with hyphenation
- [ ] **Image compression** — JPEG compression for embedded images to reduce file size

---

## Phase 2: Advanced Features (v0.4.0)

### Cross-Format Conversion
- [ ] **DOCX to PDF** — Convert Word documents to PDF with layout preservation
- [ ] **PPTX to PDF** — Convert presentations to PDF (one page per slide)
- [ ] **XLSX to PDF** — Print spreadsheet to PDF with pagination
- [ ] **Markdown to DOCX** — Convert Markdown files to styled Word documents
- [ ] **Markdown to PDF** — Convert Markdown directly to PDF

### Template Engine
- [ ] **Variable substitution** — `{{variable}}` placeholders in templates
- [ ] **Conditional sections** — Show/hide document sections based on data
- [ ] **Loop/repeat sections** — Generate repeated content from data arrays
- [ ] **Custom template loading** — Load templates from DOCX/XLSX/PPTX files
- [ ] **Template inheritance** — Base templates with overridable sections
- [ ] **JSON/YAML data binding** — Populate templates from structured data files

### Design System Expansion
- [ ] **10 additional themes** — Industry-specific themes (Healthcare, Finance, Education, Tech, Legal, etc.)
- [ ] **Custom color palette builder** — API to create palettes from brand colors
- [ ] **Theme export/import** — Save and load custom themes as JSON
- [ ] **Dark mode themes** — First-class dark theme support
- [ ] **Responsive layouts** — Adaptive layouts for different page sizes

---

## Phase 3: Enterprise & Integration (v0.5.0)

### Performance & Scale
- [ ] **Streaming writer** — Write large documents without loading entire file in memory
- [ ] **Concurrent sheet building** — Build multiple XLSX sheets in parallel
- [ ] **Memory-mapped reading** — Efficient reading of large files
- [ ] **Compression options** — Configurable ZIP compression levels

### Security
- [ ] **Document encryption** — Password-protect DOCX/XLSX/PPTX files
- [ ] **PDF encryption** — AES-256 encryption for PDFs
- [ ] **Redaction** — Permanently remove sensitive content from documents
- [ ] **Metadata cleaning** — Strip author, revision, and tracking metadata

### Integration
- [ ] **S3/GCS/Azure storage** — Direct save/load from cloud storage
- [ ] **HTTP handler** — `http.Handler` for serving generated documents
- [ ] **gRPC service** — Document generation as a microservice
- [ ] **WASM support** — Compile to WebAssembly for browser-side generation
- [ ] **Plugin system** — Custom element types and renderers

### Accessibility
- [ ] **PDF/UA compliance** — Tagged PDFs for screen readers
- [ ] **Alt text on images** — Accessibility descriptions for embedded images
- [ ] **Document language** — Proper language tagging for assistive technology
- [ ] **High-contrast themes** — WCAG AA compliant color schemes

---

## Phase 4: Ecosystem (v1.0.0)

### Developer Experience
- [ ] **Interactive playground** — Web-based document builder with live preview
- [ ] **VSCode extension** — Template previewer and code snippets
- [ ] **Documentation site** — Full API docs with interactive examples
- [ ] **Example gallery** — Showcase of generated documents at each design level

### Community
- [ ] **Template marketplace** — Community-contributed templates
- [ ] **Theme gallery** — Shareable custom themes
- [ ] **Plugin registry** — Third-party extensions and integrations

### Compatibility
- [ ] **Microsoft 365 validation** — Automated testing with Microsoft 365 compatibility
- [ ] **LibreOffice validation** — Automated testing with LibreOffice
- [ ] **Google Docs import** — Verify documents open correctly in Google Workspace
- [ ] **Apple Pages/Numbers/Keynote** — Basic compatibility testing

---

## Contributing

Want to help? Pick any unchecked item above and submit a PR! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

### Priority Items (Good First Issues)
1. Hyperlinks in DOCX — Relatively straightforward OOXML addition
2. Numbered lists in DOCX — Add `w:numPr` to paragraph properties
3. Freeze panes in XLSX — Simple XML addition to worksheet
4. Slide numbers in PPTX — Add to slide master footer
5. More chart types — Extend existing chart infrastructure

### Architecture Decisions
- **No external dependencies** — All features must use only the Go standard library
- **Format fidelity** — Generated files must open correctly in Microsoft Office, LibreOffice, and Google Workspace
- **API consistency** — All packages follow the same patterns (New, Open, Save, SaveToBytes, Delete)
- **Theme-first design** — Every visual element should respect the active theme
