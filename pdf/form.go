package pdf

// FormFieldType represents the type of a form field
type FormFieldType int

const (
	FieldText FormFieldType = iota
	FieldCheckbox
	FieldRadio
	FieldDropdown
	FieldSignature
)

// formField holds common form field properties
type formField struct {
	name     string
	x, y     float64
	width    float64
	height   float64
	value    string
	readOnly bool
	required bool
}

func (f *formField) pdfElement() {}

// SetValue sets the field value
func (f *formField) SetValue(v string) { f.value = v }

// SetReadOnly sets read-only state
func (f *formField) SetReadOnly(b bool) { f.readOnly = b }

// SetRequired sets required state
func (f *formField) SetRequired(b bool) { f.required = b }

// Name returns the field name
func (f *formField) Name() string { return f.name }

// Value returns the field value
func (f *formField) Value() string { return f.value }

// IsReadOnly returns read-only state
func (f *formField) IsReadOnly() bool { return f.readOnly }

// IsRequired returns required state
func (f *formField) IsRequired() bool { return f.required }

// TextField represents a text input field
type TextField struct {
	formField
	maxLength int
	multiline bool
}

// SetMaxLength sets the maximum number of characters
func (t *TextField) SetMaxLength(n int) { t.maxLength = n }

// MaxLength returns the max length
func (t *TextField) MaxLength() int { return t.maxLength }

// SetMultiline sets whether the field supports multiline input
func (t *TextField) SetMultiline(b bool) { t.multiline = b }

// IsMultiline returns whether the field is multiline
func (t *TextField) IsMultiline() bool { return t.multiline }

// FieldType returns FieldText
func (t *TextField) FieldType() FormFieldType { return FieldText }

// Checkbox represents a checkbox field
type Checkbox struct {
	formField
	checked bool
}

// SetChecked sets the checked state
func (c *Checkbox) SetChecked(b bool) {
	c.checked = b
	if b {
		c.value = "Yes"
	} else {
		c.value = "Off"
	}
}

// IsChecked returns the checked state
func (c *Checkbox) IsChecked() bool { return c.checked }

// FieldType returns FieldCheckbox
func (c *Checkbox) FieldType() FormFieldType { return FieldCheckbox }

// Dropdown represents a dropdown/select field
type Dropdown struct {
	formField
	options  []string
	selected int
}

// SetSelected sets the selected option index
func (d *Dropdown) SetSelected(index int) {
	if index >= 0 && index < len(d.options) {
		d.selected = index
		d.value = d.options[index]
	}
}

// Selected returns the selected index
func (d *Dropdown) Selected() int { return d.selected }

// Options returns the dropdown options
func (d *Dropdown) Options() []string { return d.options }

// FieldType returns FieldDropdown
func (d *Dropdown) FieldType() FormFieldType { return FieldDropdown }

// AddTextField adds a text input field to the page
func (p *Page) AddTextField(name string, x, y, width, height float64) *TextField {
	t := &TextField{
		formField: formField{
			name:   name,
			x:      x,
			y:      y,
			width:  width,
			height: height,
		},
	}
	p.elements = append(p.elements, t)
	return t
}

// AddCheckbox adds a checkbox field to the page
func (p *Page) AddCheckbox(name string, x, y, size float64) *Checkbox {
	c := &Checkbox{
		formField: formField{
			name:   name,
			x:      x,
			y:      y,
			width:  size,
			height: size,
		},
	}
	c.value = "Off"
	p.elements = append(p.elements, c)
	return c
}

// AddDropdown adds a dropdown field to the page
func (p *Page) AddDropdown(name string, x, y, width, height float64, options []string) *Dropdown {
	d := &Dropdown{
		formField: formField{
			name:   name,
			x:      x,
			y:      y,
			width:  width,
			height: height,
		},
		options:  options,
		selected: -1,
	}
	p.elements = append(p.elements, d)
	return d
}
