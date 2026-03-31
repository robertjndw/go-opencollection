package opencollection

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse unmarshals an OpenCollection YAML document from a byte slice.
func Parse(data []byte) (*Collection, error) {
	var c Collection
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("opencollection: parse: %w", err)
	}
	return &c, nil
}

// ParseFile reads and parses an OpenCollection YAML file.
func ParseFile(path string) (*Collection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("opencollection: read file: %w", err)
	}
	return Parse(data)
}

// Marshal serializes a Collection to YAML bytes.
func Marshal(c *Collection) ([]byte, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("opencollection: marshal: %w", err)
	}
	return data, nil
}

// WriteFile serializes a Collection and writes it to a file.
func WriteFile(path string, c *Collection) error {
	data, err := Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("opencollection: write file: %w", err)
	}
	return nil
}

// Open reads a Collection from either a YAML file (bundled: true) or a
// directory (bundled: false). When path points to a directory Open delegates
// to ReadDir; otherwise it delegates to ParseFile.
func Open(path string) (*Collection, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("opencollection: open: %w", err)
	}
	if info.IsDir() {
		return ReadDir(path)
	}
	return ParseFile(path)
}

// Write serializes a Collection to disk. When c.Bundled is true the collection
// is written as a single YAML file at path; otherwise it uses the unbundled
// directory layout via WriteDir.
func Write(path string, c *Collection) error {
	if c.Bundled {
		return WriteFile(path, c)
	}
	return WriteDir(path, c)
}
