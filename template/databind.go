package template

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON parses JSON data into the engine's data map.
// Nested keys can be accessed via dot notation: {{user.name}}.
func (e *TemplateEngine) LoadJSON(data []byte) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}
	e.SetDataMap(parsed)
	return nil
}

// LoadJSONFile reads a JSON file and loads it into the engine's data map.
func (e *TemplateEngine) LoadJSONFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read JSON file: %w", err)
	}
	return e.LoadJSON(data)
}
