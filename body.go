package opencollection

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ---- HTTP body ----

// HttpBody is a discriminated union of all HTTP request body formats.
type HttpBody struct {
	Raw            *RawBody
	FormURLEncoded *FormURLEncodedBody
	MultipartForm  *MultipartFormBody
	File           *FileBody
}

// IsZero reports whether no body is set (used by yaml.v3 omitempty).
func (b HttpBody) IsZero() bool {
	return b.Raw == nil && b.FormURLEncoded == nil && b.MultipartForm == nil && b.File == nil
}

func (b HttpBody) MarshalYAML() (any, error) {
	if b.Raw != nil {
		return b.Raw, nil
	}
	if b.FormURLEncoded != nil {
		return b.FormURLEncoded, nil
	}
	if b.MultipartForm != nil {
		return b.MultipartForm, nil
	}
	if b.File != nil {
		return b.File, nil
	}
	return nil, nil
}

func (b *HttpBody) UnmarshalYAML(value *yaml.Node) error {
	var probe struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Type {
	case "json", "text", "xml", "sparql":
		b.Raw = &RawBody{}
		return value.Decode(b.Raw)
	case "form-urlencoded":
		b.FormURLEncoded = &FormURLEncodedBody{}
		return value.Decode(b.FormURLEncoded)
	case "multipart-form":
		b.MultipartForm = &MultipartFormBody{}
		return value.Decode(b.MultipartForm)
	case "file":
		b.File = &FileBody{}
		return value.Decode(b.File)
	default:
		return fmt.Errorf("opencollection: unknown body type %q", probe.Type)
	}
}

// HttpBodyField holds either a single HttpBody or a list of named body variants.
type HttpBodyField struct {
	Body     *HttpBody
	Variants []HttpBodyVariant
}

// IsZero reports whether no body is set (used by yaml.v3 omitempty).
func (f HttpBodyField) IsZero() bool { return f.Body == nil && len(f.Variants) == 0 }

func (f HttpBodyField) MarshalYAML() (any, error) {
	if len(f.Variants) > 0 {
		return f.Variants, nil
	}
	if f.Body != nil {
		return f.Body, nil
	}
	return nil, nil
}

func (f *HttpBodyField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		return value.Decode(&f.Variants)
	}
	f.Body = &HttpBody{}
	return value.Decode(f.Body)
}

// HttpBodyVariant is a named, selectable body variant.
type HttpBodyVariant struct {
	Title    string   `yaml:"title"`
	Selected bool     `yaml:"selected,omitempty"`
	Body     HttpBody `yaml:"body"`
}

// RawBody is a raw text/JSON/XML/SPARQL body.
type RawBody struct {
	Type string `yaml:"type"` // "json", "text", "xml", "sparql"
	Data string `yaml:"data"`
}

// FormURLEncodedField is a single field in a form-urlencoded body.
type FormURLEncodedField struct {
	Name        string      `yaml:"name"`
	Value       string      `yaml:"value"`
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
}

// FormURLEncodedBody is an application/x-www-form-urlencoded body.
type FormURLEncodedBody struct {
	Type string                `yaml:"type"` // "form-urlencoded"
	Data []FormURLEncodedField `yaml:"data"`
}

// MultipartFormPart is a single part in a multipart/form-data body.
// Value is a string for text parts or []string for file parts.
type MultipartFormPart struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "text" or "file"
	Value       any         `yaml:"value"`
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
}

// MultipartFormBody is a multipart/form-data body.
type MultipartFormBody struct {
	Type string              `yaml:"type"` // "multipart-form"
	Data []MultipartFormPart `yaml:"data"`
}

// FileBodyVariant is a selectable file entry.
type FileBodyVariant struct {
	FilePath    string `yaml:"filePath"`
	ContentType string `yaml:"contentType"`
	Selected    bool   `yaml:"selected"`
}

// FileBody sends one or more files as the request body.
type FileBody struct {
	Type string            `yaml:"type"` // "file"
	Data []FileBodyVariant `yaml:"data"`
}

// ---- GraphQL body ----

// GraphQLBody holds a GraphQL query and its optional variables JSON string.
type GraphQLBody struct {
	Query     string `yaml:"query,omitempty"`
	Variables string `yaml:"variables,omitempty"`
}

// GraphQLBodyField holds either a single GraphQLBody or a list of named variants.
type GraphQLBodyField struct {
	Body     *GraphQLBody
	Variants []GraphQLBodyVariant
}

// IsZero reports whether no body is set (used by yaml.v3 omitempty).
func (f GraphQLBodyField) IsZero() bool { return f.Body == nil && len(f.Variants) == 0 }

func (f GraphQLBodyField) MarshalYAML() (any, error) {
	if len(f.Variants) > 0 {
		return f.Variants, nil
	}
	if f.Body != nil {
		return f.Body, nil
	}
	return nil, nil
}

func (f *GraphQLBodyField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		return value.Decode(&f.Variants)
	}
	f.Body = &GraphQLBody{}
	return value.Decode(f.Body)
}

// GraphQLBodyVariant is a named, selectable GraphQL body.
type GraphQLBodyVariant struct {
	Title    string      `yaml:"title"`
	Selected bool        `yaml:"selected,omitempty"`
	Body     GraphQLBody `yaml:"body"`
}

// ---- gRPC message ----

// GrpcMessageField holds either a raw message string or a list of named variants.
type GrpcMessageField struct {
	Message  string
	Variants []GrpcMessageVariant
}

// IsZero reports whether no message is set (used by yaml.v3 omitempty).
func (f GrpcMessageField) IsZero() bool { return f.Message == "" && len(f.Variants) == 0 }

func (f GrpcMessageField) MarshalYAML() (any, error) {
	if len(f.Variants) > 0 {
		return f.Variants, nil
	}
	return f.Message, nil
}

func (f *GrpcMessageField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		return value.Decode(&f.Variants)
	}
	f.Message = value.Value
	return nil
}

// GrpcMessageVariant is a named, selectable gRPC message.
type GrpcMessageVariant struct {
	Title    string `yaml:"title"`
	Selected bool   `yaml:"selected,omitempty"`
	Message  string `yaml:"message"`
}

// ---- WebSocket message ----

// WebSocketMessage is a single WebSocket message payload.
type WebSocketMessage struct {
	Type string `yaml:"type"` // "text", "json", "xml", "binary"
	Data string `yaml:"data"`
}

// WebSocketMessageField holds either a single WebSocketMessage or a list of named variants.
type WebSocketMessageField struct {
	Message  *WebSocketMessage
	Variants []WebSocketMessageVariant
}

// IsZero reports whether no message is set (used by yaml.v3 omitempty).
func (f WebSocketMessageField) IsZero() bool {
	return f.Message == nil && len(f.Variants) == 0
}

func (f WebSocketMessageField) MarshalYAML() (any, error) {
	if len(f.Variants) > 0 {
		return f.Variants, nil
	}
	if f.Message != nil {
		return f.Message, nil
	}
	return nil, nil
}

func (f *WebSocketMessageField) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		return value.Decode(&f.Variants)
	}
	f.Message = &WebSocketMessage{}
	return value.Decode(f.Message)
}

// WebSocketMessageVariant is a named, selectable WebSocket message.
type WebSocketMessageVariant struct {
	Title    string           `yaml:"title"`
	Selected bool             `yaml:"selected,omitempty"`
	Message  WebSocketMessage `yaml:"message"`
}
