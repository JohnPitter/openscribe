// Package xmlutil provides XML marshaling/unmarshaling utilities for Office Open XML.
package xmlutil

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// MarshalXML marshals a value to XML with proper declaration
func MarshalXML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("xml marshal: %w", err)
	}
	if err := enc.Flush(); err != nil {
		return nil, fmt.Errorf("xml flush: %w", err)
	}

	return buf.Bytes(), nil
}

// MarshalXMLFragment marshals without XML declaration
func MarshalXMLFragment(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("xml marshal: %w", err)
	}
	if err := enc.Flush(); err != nil {
		return nil, fmt.Errorf("xml flush: %w", err)
	}
	return buf.Bytes(), nil
}
