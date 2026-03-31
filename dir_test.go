package opencollection_test

import (
	"os"
	"path/filepath"
	"testing"

	oc "github.com/robertjndw/go-opencollection"
)

// buildTestCollection builds a collection with nested folders and multiple
// request types for use in WriteDir / ReadDir tests.
func buildTestCollection() *oc.Collection {
	c := oc.New("Dir Test API").
		CollectionVersion("0.1.0").
		Author("Alice", "alice@example.com", "").
		Environment(
			oc.NewEnvironment("dev").
				Var("baseUrl", "http://localhost:8080").
				Build(),
		).
		Environment(
			oc.NewEnvironment("prod").
				Var("baseUrl", "https://api.example.com").
				Build(),
		).
		AddFolder(
			oc.NewFolder("Users").
				AddHttpRequest(
					oc.BuildHttpRequest("List Users", "GET", "{{baseUrl}}/users").
						QueryParam("page", "1").
						Build(),
				).
				AddHttpRequest(
					oc.BuildHttpRequest("Create User", "POST", "{{baseUrl}}/users").
						JSONBody(`{"name":"Alice"}`).
						BearerAuth("{{token}}").
						Build(),
				).
				Build(),
		).
		AddHttpRequest(
			oc.BuildHttpRequest("Health Check", "GET", "{{baseUrl}}/health").Build(),
		).
		Build()

	// Fix info.type (builder sets it; ensure it's correct for test clarity)
	return c
}

func TestWriteDir_CreatesExpectedFiles(t *testing.T) {
	dir := t.TempDir()
	c := buildTestCollection()

	if err := oc.WriteDir(dir, c); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	expected := []string{
		"opencollection.yml",
		filepath.Join("environments", "dev.yml"),
		filepath.Join("environments", "prod.yml"),
		filepath.Join("users", "folder.yml"),
		filepath.Join("users", "list-users.yml"),
		filepath.Join("users", "create-user.yml"),
		"health-check.yml",
	}
	for _, rel := range expected {
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s: %v", rel, err)
		}
	}
}

func TestWriteDir_RootFileHasNoBundled(t *testing.T) {
	dir := t.TempDir()
	c := buildTestCollection()
	c.Bundled = true // should be forced to false in the root file

	if err := oc.WriteDir(dir, c); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	root, err := oc.ParseFile(filepath.Join(dir, "opencollection.yml"))
	if err != nil {
		t.Fatalf("ParseFile root: %v", err)
	}
	if root.Bundled {
		t.Error("root opencollection.yml should have bundled: false (or omitted)")
	}
	if len(root.Items) > 0 {
		t.Error("root opencollection.yml should not contain items")
	}
}

func TestWriteDir_EnvironmentsStrippedFromRoot(t *testing.T) {
	dir := t.TempDir()
	c := buildTestCollection()

	if err := oc.WriteDir(dir, c); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	root, err := oc.ParseFile(filepath.Join(dir, "opencollection.yml"))
	if err != nil {
		t.Fatalf("ParseFile root: %v", err)
	}
	if root.Config != nil && len(root.Config.Environments) > 0 {
		t.Error("environments should be written to environments/ dir, not embedded in root")
	}
}

func TestReadDir_ReadsAllItems(t *testing.T) {
	dir := t.TempDir()
	if err := oc.WriteDir(dir, buildTestCollection()); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	c, err := oc.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	if c.Info.Name != "Dir Test API" {
		t.Errorf("Name = %q", c.Info.Name)
	}
	if c.Config == nil || len(c.Config.Environments) != 2 {
		t.Errorf("expected 2 environments, got %+v", c.Config)
	}
	// Expect: users/ folder + health-check.yml = 2 top-level items
	if len(c.Items) != 2 {
		t.Fatalf("expected 2 top-level items, got %d", len(c.Items))
	}
}

func TestReadDir_FolderContainsNestedRequests(t *testing.T) {
	dir := t.TempDir()
	if err := oc.WriteDir(dir, buildTestCollection()); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	c, err := oc.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	var folder *oc.Folder
	for _, item := range c.Items {
		if item.Folder != nil {
			folder = item.Folder
			break
		}
	}
	if folder == nil {
		t.Fatal("expected a folder in items")
	}
	if len(folder.Items) != 2 {
		t.Errorf("expected 2 nested items, got %d", len(folder.Items))
	}
}

func TestWriteReadDir_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	original := buildTestCollection()

	if err := oc.WriteDir(dir, original); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	restored, err := oc.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	// Collection metadata
	if restored.Info.Name != original.Info.Name {
		t.Errorf("Name: got %q, want %q", restored.Info.Name, original.Info.Name)
	}
	if restored.Info.Version != original.Info.Version {
		t.Errorf("Version: got %q, want %q", restored.Info.Version, original.Info.Version)
	}

	// Top-level item count
	if len(restored.Items) != len(original.Items) {
		t.Errorf("top-level items: got %d, want %d", len(restored.Items), len(original.Items))
	}
}

func TestReadDir_BundledError(t *testing.T) {
	dir := t.TempDir()
	// Write a root file with bundled: true
	c := oc.New("Bundled").Build()
	c.Bundled = true
	if err := oc.WriteFile(filepath.Join(dir, "opencollection.yml"), c); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := oc.ReadDir(dir)
	if err == nil {
		t.Error("expected error when ReadDir finds bundled: true")
	}
}

func TestReadDir_FolderWithoutFolderYML(t *testing.T) {
	dir := t.TempDir()
	// Write a valid root
	root := oc.New("Test").Build()
	if err := oc.WriteFile(filepath.Join(dir, "opencollection.yml"), root); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	// Create a subdirectory without folder.yml
	if err := os.MkdirAll(filepath.Join(dir, "misc"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// Write a request item YAML directly (not a full collection document).
	const pingYAML = `info:
  name: Ping
  type: http
http:
  method: GET
  url: /ping
`
	if err := os.WriteFile(filepath.Join(dir, "misc", "ping.yml"), []byte(pingYAML), 0o644); err != nil {
		t.Fatalf("WriteFile ping.yml: %v", err)
	}

	c, err := oc.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	// misc/ should be treated as a folder with a synthesised name
	if len(c.Items) == 0 {
		t.Error("expected at least one item (the misc folder)")
	}
	if c.Items[0].Folder == nil {
		t.Error("expected a Folder item for misc/")
	}
	if c.Items[0].Folder.Info.Name != "misc" {
		t.Errorf("Folder name = %q, want %q", c.Items[0].Folder.Info.Name, "misc")
	}
}

func TestWriteDir_NestedFolders(t *testing.T) {
	dir := t.TempDir()

	c := oc.New("Nested").
		AddFolder(
			oc.NewFolder("API").
				AddFolder(
					oc.NewFolder("V1").
						AddHttpRequest(
							oc.BuildHttpRequest("Ping", "GET", "/ping").Build(),
						),
				).
				Build(),
		).
		Build()
	if err := oc.WriteDir(dir, c); err != nil {
		t.Fatalf("WriteDir: %v", err)
	}

	pingPath := filepath.Join(dir, "api", "v1", "ping.yml")
	if _, err := os.Stat(pingPath); err != nil {
		t.Errorf("expected %s: %v", pingPath, err)
	}

	restored, err := oc.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(restored.Items) == 0 || restored.Items[0].Folder == nil {
		t.Fatal("expected top-level folder")
	}
	nested := restored.Items[0].Folder
	if len(nested.Items) == 0 || nested.Items[0].Folder == nil {
		t.Fatal("expected nested folder")
	}
}
