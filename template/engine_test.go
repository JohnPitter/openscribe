package template

import (
	"testing"

	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/spreadsheet"
)

func TestEngineVariableSubstitution(t *testing.T) {
	engine := NewEngine()
	engine.SetData("name", "John")
	engine.SetData("company", "Acme")

	result := engine.replaceVars("Hello {{name}} from {{company}}")
	if result != "Hello John from Acme" {
		t.Errorf("expected 'Hello John from Acme', got %q", result)
	}
}

func TestEngineNestedDotNotation(t *testing.T) {
	engine := NewEngine()
	engine.SetData("user", map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
	})

	result := engine.replaceVars("Hello {{user.name}}, email: {{user.email}}")
	expected := "Hello Alice, email: alice@example.com"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEngineConditionalTruthy(t *testing.T) {
	engine := NewEngine()
	engine.SetData("showHeader", true)

	result := engine.processConditionals("{{#if showHeader}}HEADER{{/if}}")
	if result != "HEADER" {
		t.Errorf("expected 'HEADER', got %q", result)
	}
}

func TestEngineConditionalFalsy(t *testing.T) {
	engine := NewEngine()
	engine.SetData("showHeader", false)

	result := engine.processConditionals("{{#if showHeader}}HEADER{{/if}}")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestEngineConditionalMissing(t *testing.T) {
	engine := NewEngine()

	result := engine.processConditionals("{{#if missing}}CONTENT{{/if}}")
	if result != "" {
		t.Errorf("expected empty string for missing key, got %q", result)
	}
}

func TestEngineLoop(t *testing.T) {
	engine := NewEngine()
	engine.SetData("items", []interface{}{"A", "B", "C"})

	result := engine.processLoops("{{#each items}}[{{this}}]{{/each}}")
	expected := "[A][B][C]"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEngineLoopWithMaps(t *testing.T) {
	engine := NewEngine()
	engine.SetData("users", []interface{}{
		map[string]interface{}{"name": "Alice", "role": "Admin"},
		map[string]interface{}{"name": "Bob", "role": "User"},
	})

	result := engine.processLoops("{{#each users}}{{.name}}:{{.role}} {{/each}}")
	expected := "Alice:Admin Bob:User "
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEngineSetDataMap(t *testing.T) {
	engine := NewEngine()
	engine.SetDataMap(map[string]interface{}{
		"a": "1",
		"b": "2",
	})

	result := engine.replaceVars("{{a}}-{{b}}")
	if result != "1-2" {
		t.Errorf("expected '1-2', got %q", result)
	}
}

func TestEngineLoadJSON(t *testing.T) {
	engine := NewEngine()
	jsonData := []byte(`{"name": "Test", "count": 42, "nested": {"key": "value"}}`)
	err := engine.LoadJSON(jsonData)
	if err != nil {
		t.Fatalf("LoadJSON error: %v", err)
	}

	result := engine.replaceVars("{{name}} has {{nested.key}}")
	expected := "Test has value"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEngineLoadJSONInvalid(t *testing.T) {
	engine := NewEngine()
	err := engine.LoadJSON([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestEngineRenderDocx(t *testing.T) {
	doc := document.New()
	para := doc.AddParagraph()
	run := para.AddRun()
	run.SetText("Hello {{name}}")

	engine := NewEngine()
	engine.SetData("name", "World")

	result := engine.RenderDocx(doc)
	text := result.Paragraphs()[0].Text()
	if text != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", text)
	}
}

func TestEngineRenderPdf(t *testing.T) {
	doc := pdf.New()
	page := doc.AddPage()
	page.AddText("Invoice for {{client}}", 72, 72, pdf.DefaultHTMLOptions().DefaultFont)

	engine := NewEngine()
	engine.SetData("client", "Acme Corp")

	engine.RenderPdf(doc)

	// Verify the text was replaced
	for _, elem := range doc.Page(0).Elements() {
		if te, ok := elem.(*pdf.TextElement); ok {
			if te.Text() == "Invoice for Acme Corp" {
				return // success
			}
		}
	}
	t.Error("expected text element with replaced value")
}

func TestEngineRenderXlsx(t *testing.T) {
	wb := spreadsheet.New()
	sheet := wb.AddSheet("Sheet1")
	sheet.Cell(1, 1).SetString("{{title}}")
	sheet.Cell(1, 2).SetString("{{date}}")
	sheet.Cell(2, 1).SetNumber(42) // numbers should not be affected

	engine := NewEngine()
	engine.SetData("title", "Report")
	engine.SetData("date", "2026-03-14")

	engine.RenderXlsx(wb)

	if wb.Sheet(0).Cell(1, 1).String() != "Report" {
		t.Errorf("expected 'Report', got %q", wb.Sheet(0).Cell(1, 1).String())
	}
	if wb.Sheet(0).Cell(1, 2).String() != "2026-03-14" {
		t.Errorf("expected '2026-03-14', got %q", wb.Sheet(0).Cell(1, 2).String())
	}
}

func TestEngineCombinedDirectivesAndVars(t *testing.T) {
	engine := NewEngine()
	engine.SetData("showGreeting", true)
	engine.SetData("name", "Alice")

	input := "{{#if showGreeting}}Hello {{name}}!{{/if}}"
	result := engine.processDirectives(input)
	result = engine.replaceVars(result)
	expected := "Hello Alice!"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEngineIsTruthy(t *testing.T) {
	engine := NewEngine()
	engine.SetData("boolTrue", true)
	engine.SetData("boolFalse", false)
	engine.SetData("str", "hello")
	engine.SetData("emptyStr", "")
	engine.SetData("num", 42)
	engine.SetData("zero", 0)
	engine.SetData("list", []interface{}{1})
	engine.SetData("emptyList", []interface{}{})

	tests := []struct {
		key  string
		want bool
	}{
		{"boolTrue", true},
		{"boolFalse", false},
		{"str", true},
		{"emptyStr", false},
		{"num", true},
		{"zero", false},
		{"list", true},
		{"emptyList", false},
		{"missing", false},
	}

	for _, tt := range tests {
		if got := engine.isTruthy(tt.key); got != tt.want {
			t.Errorf("isTruthy(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}
