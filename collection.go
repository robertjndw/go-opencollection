package opencollection

// Collection is the root OpenCollection document.
type Collection struct {
	OpenCollection string            `yaml:"opencollection"`
	Info           Info              `yaml:"info"`
	Config         *CollectionConfig `yaml:"config,omitempty"`
	Items          []Item            `yaml:"items,omitempty"`
	Bundled        bool              `yaml:"bundled,omitempty"`
	Request        *RequestDefaults  `yaml:"request,omitempty"`
	Docs           Description       `yaml:"docs,omitempty"`
	Extensions     map[string]any    `yaml:"extensions,omitempty"`
}

// Info holds collection metadata.
type Info struct {
	Name    string   `yaml:"name"`
	Summary string   `yaml:"summary,omitempty"`
	Version string   `yaml:"version,omitempty"`
	Authors []Author `yaml:"authors,omitempty"`
}

// Author describes a collection author.
type Author struct {
	Name  string `yaml:"name,omitempty"`
	Email string `yaml:"email,omitempty"`
	URL   string `yaml:"url,omitempty"`
}
