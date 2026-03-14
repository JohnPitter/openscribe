package document

// Hyperlink represents a clickable hyperlink in a paragraph
type Hyperlink struct {
	text  string
	url   string
	relID string // relationship ID assigned during build
}

// NewHyperlink creates a new hyperlink
func NewHyperlink(text, url string) *Hyperlink {
	return &Hyperlink{
		text: text,
		url:  url,
	}
}

// Text returns the display text
func (h *Hyperlink) Text() string {
	return h.text
}

// URL returns the target URL
func (h *Hyperlink) URL() string {
	return h.url
}

// RelID returns the relationship ID (set during build)
func (h *Hyperlink) RelID() string {
	return h.relID
}
