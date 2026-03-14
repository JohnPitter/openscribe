package xmlutil

import (
	"encoding/xml"
	"fmt"
)

// UnmarshalXML unmarshals XML data into a value
func UnmarshalXML(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("xml unmarshal: %w", err)
	}
	return nil
}
