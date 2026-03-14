package pdf

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestAddTextField(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("username", 72, 100, 200, 24)

	if tf.Name() != "username" {
		t.Errorf("expected name 'username', got %q", tf.Name())
	}
	if tf.FieldType() != FieldText {
		t.Error("expected FieldText type")
	}
	if p.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", p.ElementCount())
	}
}

func TestTextFieldValue(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("email", 72, 100, 200, 24)

	tf.SetValue("test@example.com")
	if tf.Value() != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %q", tf.Value())
	}
}

func TestTextFieldMaxLength(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("code", 72, 100, 100, 24)

	tf.SetMaxLength(6)
	if tf.MaxLength() != 6 {
		t.Errorf("expected maxLength 6, got %d", tf.MaxLength())
	}
}

func TestTextFieldMultiline(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("comments", 72, 100, 300, 100)

	tf.SetMultiline(true)
	if !tf.IsMultiline() {
		t.Error("expected multiline to be true")
	}
}

func TestTextFieldReadOnly(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("readonly_field", 72, 100, 200, 24)

	tf.SetReadOnly(true)
	if !tf.IsReadOnly() {
		t.Error("expected readOnly to be true")
	}
}

func TestTextFieldRequired(t *testing.T) {
	d := New()
	p := d.AddPage()
	tf := p.AddTextField("required_field", 72, 100, 200, 24)

	tf.SetRequired(true)
	if !tf.IsRequired() {
		t.Error("expected required to be true")
	}
}

func TestAddCheckbox(t *testing.T) {
	d := New()
	p := d.AddPage()
	cb := p.AddCheckbox("agree", 72, 200, 14)

	if cb.Name() != "agree" {
		t.Errorf("expected name 'agree', got %q", cb.Name())
	}
	if cb.FieldType() != FieldCheckbox {
		t.Error("expected FieldCheckbox type")
	}
	if cb.IsChecked() {
		t.Error("checkbox should not be checked by default")
	}
	if cb.Value() != "Off" {
		t.Errorf("expected value 'Off', got %q", cb.Value())
	}
}

func TestCheckboxSetChecked(t *testing.T) {
	d := New()
	p := d.AddPage()
	cb := p.AddCheckbox("terms", 72, 200, 14)

	cb.SetChecked(true)
	if !cb.IsChecked() {
		t.Error("expected checkbox to be checked")
	}
	if cb.Value() != "Yes" {
		t.Errorf("expected value 'Yes', got %q", cb.Value())
	}

	cb.SetChecked(false)
	if cb.IsChecked() {
		t.Error("expected checkbox to be unchecked")
	}
	if cb.Value() != "Off" {
		t.Errorf("expected value 'Off', got %q", cb.Value())
	}
}

func TestCheckboxReadOnlyRequired(t *testing.T) {
	d := New()
	p := d.AddPage()
	cb := p.AddCheckbox("locked", 72, 200, 14)

	cb.SetReadOnly(true)
	cb.SetRequired(true)
	if !cb.IsReadOnly() {
		t.Error("expected readOnly")
	}
	if !cb.IsRequired() {
		t.Error("expected required")
	}
}

func TestAddDropdown(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}
	d := New()
	p := d.AddPage()
	dd := p.AddDropdown("country", 72, 300, 200, 24, options)

	if dd.Name() != "country" {
		t.Errorf("expected name 'country', got %q", dd.Name())
	}
	if dd.FieldType() != FieldDropdown {
		t.Error("expected FieldDropdown type")
	}
	if len(dd.Options()) != 3 {
		t.Errorf("expected 3 options, got %d", len(dd.Options()))
	}
	if dd.Selected() != -1 {
		t.Errorf("expected no selection (-1), got %d", dd.Selected())
	}
}

func TestDropdownSetSelected(t *testing.T) {
	options := []string{"Red", "Green", "Blue"}
	d := New()
	p := d.AddPage()
	dd := p.AddDropdown("color", 72, 300, 200, 24, options)

	dd.SetSelected(1)
	if dd.Selected() != 1 {
		t.Errorf("expected selected index 1, got %d", dd.Selected())
	}
	if dd.Value() != "Green" {
		t.Errorf("expected value 'Green', got %q", dd.Value())
	}

	// Out of range should not change
	dd.SetSelected(10)
	if dd.Selected() != 1 {
		t.Error("out-of-range index should not change selection")
	}
}

func TestDropdownReadOnlyRequired(t *testing.T) {
	d := New()
	p := d.AddPage()
	dd := p.AddDropdown("locked", 72, 300, 200, 24, []string{"A"})

	dd.SetReadOnly(true)
	dd.SetRequired(true)
	if !dd.IsReadOnly() {
		t.Error("expected readOnly")
	}
	if !dd.IsRequired() {
		t.Error("expected required")
	}
}

func TestFormFieldsBuildPDF(t *testing.T) {
	d := New()
	p := d.AddPage()

	tf := p.AddTextField("name", 72, 100, 200, 24)
	tf.SetValue("John")
	tf.SetRequired(true)

	cb := p.AddCheckbox("agree", 72, 140, 14)
	cb.SetChecked(true)

	dd := p.AddDropdown("role", 72, 170, 200, 24, []string{"Admin", "User", "Guest"})
	dd.SetSelected(0)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	content := string(data)

	// Verify PDF structure
	if !strings.Contains(content, "%PDF-") {
		t.Error("should start with PDF header")
	}
	if !strings.Contains(content, "/AcroForm") {
		t.Error("should contain AcroForm in catalog")
	}
	if !strings.Contains(content, "/FT /Tx") {
		t.Error("should contain text field type")
	}
	if !strings.Contains(content, "/FT /Btn") {
		t.Error("should contain button/checkbox field type")
	}
	if !strings.Contains(content, "/FT /Ch") {
		t.Error("should contain choice field type")
	}
	if !strings.Contains(content, "/T (name)") {
		t.Error("should contain field name 'name'")
	}
	if !strings.Contains(content, "/V (John)") {
		t.Error("should contain field value 'John'")
	}
	if !strings.Contains(content, "/Annots") {
		t.Error("page should have /Annots array")
	}
}

func TestFormFieldsWithoutFormFieldsNoAcroForm(t *testing.T) {
	d := New()
	p := d.AddPage()
	p.AddText("No forms here", 72, 72, common.NewFont("Helvetica", 12))

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "/AcroForm") {
		t.Error("should NOT contain AcroForm when no form fields exist")
	}
}
