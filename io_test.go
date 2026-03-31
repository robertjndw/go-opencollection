package opencollection_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	oc "github.com/robertjndw/go-opencollection"
)

// ---- helpers ----

func mustMarshal(t *testing.T, c *oc.Collection) []byte {
	t.Helper()
	b, err := oc.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	return b
}

func mustParse(t *testing.T, yaml string) *oc.Collection {
	t.Helper()
	c, err := oc.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return c
}

// ---- Parse ----

func TestParse_MinimalCollection(t *testing.T) {
	const yaml = `
opencollection: "1"
info:
  name: My API
`
	c := mustParse(t, yaml)
	if c.Info.Name != "My API" {
		t.Errorf("Info.Name = %q, want %q", c.Info.Name, "My API")
	}
	if c.OpenCollection != "1" {
		t.Errorf("OpenCollection = %q, want %q", c.OpenCollection, "1")
	}
}

func TestParse_FullCollection(t *testing.T) {
	const yaml = `
opencollection: "1"
info:
  name: Full API
  version: 2.0.0
  summary: A complete collection
  authors:
    - name: Alice
      email: alice@example.com
config:
  environments:
    - name: dev
      color: green
      variables:
        - name: baseUrl
          value: http://localhost:8080
bundled: true
`
	c := mustParse(t, yaml)

	if c.Info.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", c.Info.Version, "2.0.0")
	}
	if len(c.Info.Authors) != 1 || c.Info.Authors[0].Email != "alice@example.com" {
		t.Errorf("unexpected authors: %+v", c.Info.Authors)
	}
	if !c.Bundled {
		t.Error("Bundled should be true")
	}
	if c.Config == nil || len(c.Config.Environments) != 1 {
		t.Fatal("expected one environment")
	}
	env := c.Config.Environments[0]
	if env.Name != "dev" || env.Color != "green" {
		t.Errorf("unexpected env: %+v", env)
	}
	if len(env.Variables) != 1 || env.Variables[0].Variable.Value.Simple != "http://localhost:8080" {
		t.Errorf("unexpected env variable: %+v", env.Variables)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := oc.Parse([]byte(":\tinvalid: [yaml"))
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

// ---- Marshal ----

func TestMarshal_RoundTrip(t *testing.T) {
	original := oc.New("Round-trip API").
		CollectionVersion("1.2.3").
		Author("Bob", "bob@example.com", "").
		Environment(oc.NewEnvironment("staging").Var("host", "https://staging.example.com").Build()).
		AddHttpRequest(
			oc.BuildHttpRequest("Health", "GET", "{{host}}/health").
				Tag("smoke").
				Build(),
		).
		Build()

	data := mustMarshal(t, original)
	restored, err := oc.Parse(data)
	if err != nil {
		t.Fatalf("Parse after Marshal: %v", err)
	}

	if restored.Info.Name != original.Info.Name {
		t.Errorf("Name mismatch: got %q", restored.Info.Name)
	}
	if restored.Info.Version != "1.2.3" {
		t.Errorf("Version mismatch: got %q", restored.Info.Version)
	}
	if len(restored.Info.Authors) != 1 || restored.Info.Authors[0].Name != "Bob" {
		t.Errorf("Authors mismatch: %+v", restored.Info.Authors)
	}
	if len(restored.Items) != 1 || restored.Items[0].HttpRequest == nil {
		t.Fatal("expected one HTTP request item")
	}
	req := restored.Items[0].HttpRequest
	if req.Info.Name != "Health" {
		t.Errorf("request name = %q", req.Info.Name)
	}
	if len(req.Info.Tags) != 1 || req.Info.Tags[0] != "smoke" {
		t.Errorf("tags = %v", req.Info.Tags)
	}
}

func TestMarshal_OmitsZeroFields(t *testing.T) {
	c := &oc.Collection{
		OpenCollection: "1",
		Info:           oc.Info{Name: "Minimal"},
	}
	data := mustMarshal(t, c)
	s := string(data)

	for _, absent := range []string{"config", "items", "bundled", "request", "docs", "extensions"} {
		if strings.Contains(s, absent+":") {
			t.Errorf("output should not contain %q:\n%s", absent, s)
		}
	}
}

// ---- ParseFile / WriteFile ----

func TestParseFile_WriteFile_RoundTrip(t *testing.T) {
	c := oc.New("File API").Build()
	c.Bundled = true

	dir := t.TempDir()
	path := filepath.Join(dir, "opencollection.yml")

	if err := oc.WriteFile(path, c); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	restored, err := oc.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if restored.Info.Name != "File API" {
		t.Errorf("Name = %q", restored.Info.Name)
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := oc.ParseFile("/nonexistent/path/file.yml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// ---- Open / Write ----

func TestOpen_File(t *testing.T) {
	c := oc.New("Open File Test").Build()
	c.Bundled = true

	dir := t.TempDir()
	path := filepath.Join(dir, "opencollection.yml")
	if err := oc.WriteFile(path, c); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	restored, err := oc.Open(path)
	if err != nil {
		t.Fatalf("Open(file): %v", err)
	}
	if restored.Info.Name != "Open File Test" {
		t.Errorf("Name = %q", restored.Info.Name)
	}
}

func TestOpen_Directory(t *testing.T) {
	c := oc.New("Open Dir Test").Build()

	dir := t.TempDir()
	if err := oc.WriteDir(dir, c); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	restored, err := oc.Open(dir)
	if err != nil {
		t.Fatalf("Open(dir): %v", err)
	}
	if restored.Info.Name != "Open Dir Test" {
		t.Errorf("Name = %q", restored.Info.Name)
	}
}

func TestWrite_Bundled_WritesFile(t *testing.T) {
	c := oc.New("Bundled Write").Build()
	c.Bundled = true

	dir := t.TempDir()
	path := filepath.Join(dir, "opencollection.yml")

	if err := oc.Write(path, c); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s: %v", path, err)
	}
}

func TestWrite_Unbundled_WritesDirectory(t *testing.T) {
	c := oc.New("Unbundled Write").Build()

	dir := t.TempDir()
	if err := oc.Write(dir, c); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "opencollection.yml")); err != nil {
		t.Errorf("expected opencollection.yml: %v", err)
	}
}
