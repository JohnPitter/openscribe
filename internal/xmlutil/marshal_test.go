package xmlutil

import (
	"encoding/xml"
	"strings"
	"testing"
)

type testElement struct {
	XMLName xml.Name `xml:"test"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

func TestMarshalXML(t *testing.T) {
	elem := testElement{Name: "hello", Value: "world"}
	data, err := MarshalXML(elem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(data)
	if !strings.Contains(str, "<?xml version=") {
		t.Error("should contain XML declaration")
	}
	if !strings.Contains(str, `name="hello"`) {
		t.Error("should contain attribute")
	}
	if !strings.Contains(str, "world") {
		t.Error("should contain value")
	}
}

func TestMarshalXMLFragment(t *testing.T) {
	elem := testElement{Name: "hello", Value: "world"}
	data, err := MarshalXMLFragment(elem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(data)
	if strings.Contains(str, "<?xml version=") {
		t.Error("fragment should not contain XML declaration")
	}
}

func TestUnmarshalXML(t *testing.T) {
	xmlData := `<test name="hello">world</test>`
	var elem testElement
	err := UnmarshalXML([]byte(xmlData), &elem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elem.Name != "hello" || elem.Value != "world" {
		t.Errorf("unexpected result: %+v", elem)
	}
}
