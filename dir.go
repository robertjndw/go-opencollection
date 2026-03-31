package opencollection

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// ReadDir reads an unbundled OpenCollection (bundled: false) from a directory.
//
// Expected layout:
//
//	<dir>/
//	  opencollection.yml   — collection root (bundled: false, no items)
//	  environments/
//	    <name>.yml         — one Environment per file
//	  <name>.yml           — a request item (http, graphql, grpc, websocket)
//	  <name>/
//	    folder.yml         — folder metadata (info.type: "folder")
//	    ...                — nested items, same rules recursively
//
// Items are ordered by their seq field when present, otherwise by filename.
func ReadDir(dir string) (*Collection, error) {
	rootPath := filepath.Join(dir, "opencollection.yml")
	c, err := ParseFile(rootPath)
	if err != nil {
		return nil, fmt.Errorf("opencollection: read dir: root: %w", err)
	}
	if c.Bundled {
		return nil, fmt.Errorf("opencollection: read dir: collection has bundled: true — use ParseFile instead")
	}

	// Environments directory
	envs, err := readEnvironmentsDir(filepath.Join(dir, "environments"))
	if err != nil {
		return nil, fmt.Errorf("opencollection: read dir: environments: %w", err)
	}
	if len(envs) > 0 {
		if c.Config == nil {
			c.Config = &CollectionConfig{}
		}
		c.Config.Environments = append(c.Config.Environments, envs...)
	}

	// Items
	items, err := scanDir(dir, map[string]bool{"opencollection.yml": true}, map[string]bool{"environments": true})
	if err != nil {
		return nil, fmt.Errorf("opencollection: read dir: %w", err)
	}
	c.Items = items

	return c, nil
}

// WriteDir writes a Collection as an unbundled directory layout (bundled: false).
//
// Environments are written to <dir>/environments/<slug>.yml.
// Folders are written as subdirectories containing a folder.yml and their items.
// All other items are written as <dir>/<slug>.yml.
//
// Any existing files in dir are overwritten; extra files are not removed.
func WriteDir(dir string, c *Collection) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("opencollection: write dir: %w", err)
	}

	// Build a root-file copy: bundled: false, no items, environments stripped out.
	root := *c
	root.Bundled = false
	root.Items = nil

	var envs []Environment
	if c.Config != nil {
		cfg := *c.Config
		envs = cfg.Environments
		cfg.Environments = nil
		if isEmptyCollectionConfig(&cfg) {
			root.Config = nil
		} else {
			root.Config = &cfg
		}
	}

	if err := marshalYAMLFile(filepath.Join(dir, "opencollection.yml"), &root); err != nil {
		return err
	}

	// Write environments
	if len(envs) > 0 {
		envDir := filepath.Join(dir, "environments")
		if err := os.MkdirAll(envDir, 0o755); err != nil {
			return fmt.Errorf("opencollection: write dir: environments: %w", err)
		}
		for i := range envs {
			name := slugify(envs[i].Name) + ".yml"
			if err := marshalYAMLFile(filepath.Join(envDir, name), &envs[i]); err != nil {
				return err
			}
		}
	}

	return writeItems(dir, c.Items)
}

// ---- internal helpers ----

func readEnvironmentsDir(dir string) ([]Environment, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var envs []Environment
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var env Environment
		if err := yaml.Unmarshal(data, &env); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		envs = append(envs, env)
	}
	return envs, nil
}

// scanDir reads items from a directory, skipping excluded files/dirs.
// Results are sorted by seq (if present) then by original directory order.
func scanDir(dir string, skipFiles, skipDirs map[string]bool) ([]Item, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			if skipDirs[name] {
				continue
			}
			folder, err := readFolderDir(filepath.Join(dir, name))
			if err != nil {
				return nil, fmt.Errorf("folder %q: %w", name, err)
			}
			items = append(items, Item{Folder: folder})
		} else if strings.HasSuffix(name, ".yml") && !skipFiles[name] {
			item, err := readItemFile(filepath.Join(dir, name))
			if err != nil {
				return nil, fmt.Errorf("item %q: %w", name, err)
			}
			items = append(items, *item)
		}
	}

	sortBySeq(items)
	return items, nil
}

func readFolderDir(dir string) (*Folder, error) {
	var folder Folder

	folderYML := filepath.Join(dir, "folder.yml")
	if data, err := os.ReadFile(folderYML); err == nil {
		if err := yaml.Unmarshal(data, &folder); err != nil {
			return nil, fmt.Errorf("parse folder.yml: %w", err)
		}
	} else if os.IsNotExist(err) {
		// No folder.yml — synthesise minimal info from directory name.
		folder.Info = FolderInfo{Name: filepath.Base(dir), Type: "folder"}
	} else {
		return nil, err
	}

	nested, err := scanDir(dir, map[string]bool{"folder.yml": true}, nil)
	if err != nil {
		return nil, err
	}
	folder.Items = nested
	return &folder, nil
}

func readItemFile(path string) (*Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var item Item
	if err := yaml.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("parse %s: %w", filepath.Base(path), err)
	}
	return &item, nil
}

func writeItems(dir string, items []Item) error {
	scriptIdx := 0
	for _, item := range items {
		switch {
		case item.HttpRequest != nil:
			slug := slugify(item.HttpRequest.Info.Name)
			if err := marshalYAMLFile(filepath.Join(dir, slug+".yml"), item.HttpRequest); err != nil {
				return err
			}
		case item.GraphQLRequest != nil:
			slug := slugify(item.GraphQLRequest.Info.Name)
			if err := marshalYAMLFile(filepath.Join(dir, slug+".yml"), item.GraphQLRequest); err != nil {
				return err
			}
		case item.GrpcRequest != nil:
			slug := slugify(item.GrpcRequest.Info.Name)
			if err := marshalYAMLFile(filepath.Join(dir, slug+".yml"), item.GrpcRequest); err != nil {
				return err
			}
		case item.WebSocket != nil:
			slug := slugify(item.WebSocket.Info.Name)
			if err := marshalYAMLFile(filepath.Join(dir, slug+".yml"), item.WebSocket); err != nil {
				return err
			}
		case item.Folder != nil:
			subdir := filepath.Join(dir, slugify(item.Folder.Info.Name))
			if err := os.MkdirAll(subdir, 0o755); err != nil {
				return err
			}
			// Write folder.yml without nested items.
			meta := *item.Folder
			nestedItems := meta.Items
			meta.Items = nil
			if err := marshalYAMLFile(filepath.Join(subdir, "folder.yml"), &meta); err != nil {
				return err
			}
			if err := writeItems(subdir, nestedItems); err != nil {
				return err
			}
		case item.Script != nil:
			scriptIdx++
			name := fmt.Sprintf("script-%d.yml", scriptIdx)
			if err := marshalYAMLFile(filepath.Join(dir, name), item.Script); err != nil {
				return err
			}
		}
	}
	return nil
}

func marshalYAMLFile(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("opencollection: marshal %s: %w", path, err)
	}
	return os.WriteFile(path, data, 0o644)
}

func isEmptyCollectionConfig(c *CollectionConfig) bool {
	return len(c.Environments) == 0 &&
		c.Protobuf == nil &&
		c.Proxy == nil &&
		len(c.ClientCertificates) == 0
}

// sortBySeq sorts items stably by their seq field.
// Items without a seq value keep their original relative order after
// items that do have one.
func sortBySeq(items []Item) {
	seqs := make([]*float64, len(items))
	for i, item := range items {
		seqs[i] = itemSeq(item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		si, sj := seqs[i], seqs[j]
		if si == nil || sj == nil {
			return false
		}
		return *si < *sj
	})
}

func itemSeq(item Item) *float64 {
	switch {
	case item.HttpRequest != nil:
		return item.HttpRequest.Info.Seq
	case item.GraphQLRequest != nil:
		return item.GraphQLRequest.Info.Seq
	case item.GrpcRequest != nil:
		return item.GrpcRequest.Info.Seq
	case item.WebSocket != nil:
		return item.WebSocket.Info.Seq
	case item.Folder != nil:
		return item.Folder.Info.Seq
	}
	return nil
}

// slugify converts a human-readable name to a lowercase, hyphen-separated
// filename slug suitable for use on any filesystem.
//
//	"Get User by ID" → "get-user-by-id"
//	"POST /api/v1"   → "post-api-v1"
func slugify(name string) string {
	var buf []rune
	prevDash := true // start true to suppress leading dash
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf = append(buf, unicode.ToLower(r))
			prevDash = false
		} else if !prevDash {
			buf = append(buf, '-')
			prevDash = true
		}
	}
	// Trim trailing dash.
	for len(buf) > 0 && buf[len(buf)-1] == '-' {
		buf = buf[:len(buf)-1]
	}
	if len(buf) == 0 {
		return "item"
	}
	return string(buf)
}
