package opencollection

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schema.json
var schemaData []byte

var resolvedSchema *jsonschema.Resolved

func init() {
	var s jsonschema.Schema
	if err := json.Unmarshal(schemaData, &s); err != nil {
		panic(fmt.Sprintf("opencollection: failed to load embedded schema: %v", err))
	}
	var err error
	resolvedSchema, err = s.Resolve(nil)
	if err != nil {
		panic(fmt.Sprintf("opencollection: failed to resolve embedded schema: %v", err))
	}
}

// Validate validates a Collection against the OpenCollection JSON schema.
func Validate(c *Collection) error {
	// Marshal to YAML so custom marshalers produce the canonical representation.
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("opencollection: validate: %w", err)
	}

	// Unmarshal YAML into a generic document.
	var doc any
	if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
		return fmt.Errorf("opencollection: validate: %w", err)
	}

	// Round-trip through JSON to normalise types (e.g. YAML int → JSON float64)
	// so the schema validator receives standard JSON-decoded values.
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("opencollection: validate: %w", err)
	}
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return fmt.Errorf("opencollection: validate: %w", err)
	}

	if err := resolvedSchema.Validate(doc); err != nil {
		return fmt.Errorf("opencollection: validate: %w", err)
	}
	return nil
}
