package opencollection

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// ---- Description ----

// Description can be unset, null, a plain string, or a typed {content, type} object.
type Description struct {
	IsSet    bool
	IsNull   bool
	Content  string
	MIMEType string // non-empty only when the typed {content, type} form is used
}

// IsZero reports whether the Description is unset (used by yaml.v3 omitempty).
func (d Description) IsZero() bool { return !d.IsSet }

// StringDescription returns a Description with a plain string value.
func StringDescription(s string) Description {
	return Description{IsSet: true, Content: s}
}

func (d Description) MarshalYAML() (any, error) {
	if !d.IsSet || d.IsNull {
		return nil, nil
	}
	if d.MIMEType != "" {
		return struct {
			Content string `yaml:"content"`
			Type    string `yaml:"type"`
		}{d.Content, d.MIMEType}, nil
	}
	return d.Content, nil
}

func (d *Description) UnmarshalYAML(value *yaml.Node) error {
	d.IsSet = true
	switch value.Kind {
	case yaml.ScalarNode:
		if value.Tag == "!!null" {
			d.IsNull = true
			return nil
		}
		d.Content = value.Value
		return nil
	case yaml.MappingNode:
		var obj struct {
			Content string `yaml:"content"`
			Type    string `yaml:"type"`
		}
		if err := value.Decode(&obj); err != nil {
			return err
		}
		d.Content = obj.Content
		d.MIMEType = obj.Type
		return nil
	}
	return fmt.Errorf("opencollection: unexpected node kind for Description")
}

// ---- Variables ----

// Variable is a named variable with an optional value.
type Variable struct {
	Name        string        `yaml:"name"`
	Value       VariableValue `yaml:"value,omitempty"`
	Description Description   `yaml:"description,omitempty"`
	Disabled    bool          `yaml:"disabled,omitempty"`
}

// SecretVariable is a secret named variable; its value is not stored inline.
type SecretVariable struct {
	Secret      bool        `yaml:"secret"`
	Name        string      `yaml:"name"`
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
	Type        string      `yaml:"type,omitempty"` // "string", "number", "boolean", "null", "object"
}

// VariableValue holds a variable value which can be a simple string, a typed
// object {type, data}, or an array of named variants.
type VariableValue struct {
	Simple   string
	Typed    *TypedVariableValue
	Variants []VariableValueVariant
}

// IsZero reports whether the value is unset (used by yaml.v3 omitempty).
func (v VariableValue) IsZero() bool {
	return v.Simple == "" && v.Typed == nil && len(v.Variants) == 0
}

func (v VariableValue) MarshalYAML() (any, error) {
	if len(v.Variants) > 0 {
		return v.Variants, nil
	}
	if v.Typed != nil {
		return v.Typed, nil
	}
	return v.Simple, nil
}

func (v *VariableValue) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		v.Simple = value.Value
		return nil
	case yaml.SequenceNode:
		return value.Decode(&v.Variants)
	case yaml.MappingNode:
		v.Typed = &TypedVariableValue{}
		return value.Decode(v.Typed)
	}
	return fmt.Errorf("opencollection: unexpected YAML node kind %v for VariableValue", value.Kind)
}

// TypedVariableValue holds a typed variable value.
type TypedVariableValue struct {
	Type string `yaml:"type"` // "string", "number", "boolean", "null", "object"
	Data string `yaml:"data"`
}

// VariableValueVariant is a named variant of a variable value.
type VariableValueVariant struct {
	Title    string        `yaml:"title"`
	Selected bool          `yaml:"selected,omitempty"`
	Value    VariableValue `yaml:"value"`
}

// ---- Scripts, assertions, actions ----

// Scripts is an ordered list of lifecycle scripts.
type Scripts []Script

// Script executes JavaScript at a specific lifecycle stage.
type Script struct {
	Type string `yaml:"type"` // "before-request", "after-response", "tests", "hooks"
	Code string `yaml:"code"`
}

// Assertion validates an expression in the response.
type Assertion struct {
	Expression  string      `yaml:"expression"`
	Operator    string      `yaml:"operator"`
	Value       string      `yaml:"value,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
	Description Description `yaml:"description,omitempty"`
}

// Action is a runtime action. Currently only the set-variable type is defined.
type Action struct {
	SetVariable *ActionSetVariable
}

func (a Action) MarshalYAML() (any, error) {
	if a.SetVariable != nil {
		return a.SetVariable, nil
	}
	return nil, errors.New("opencollection: Action has no type set")
}

func (a *Action) UnmarshalYAML(value *yaml.Node) error {
	var probe struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Type {
	case "set-variable":
		a.SetVariable = &ActionSetVariable{}
		return value.Decode(a.SetVariable)
	default:
		return fmt.Errorf("opencollection: unknown action type %q", probe.Type)
	}
}

// ActionSetVariable sets a variable using a selector result.
type ActionSetVariable struct {
	Type        string               `yaml:"type"` // "set-variable"
	Description Description          `yaml:"description,omitempty"`
	Phase       string               `yaml:"phase,omitempty"` // "before-request" or "after-response"
	Selector    ActionSelector       `yaml:"selector"`
	Variable    ActionTargetVariable `yaml:"variable"`
	Disabled    bool                 `yaml:"disabled,omitempty"`
}

// ActionSelector selects a value from the response.
type ActionSelector struct {
	Expression string `yaml:"expression"`
	Method     string `yaml:"method"` // "jsonq"
}

// ActionTargetVariable identifies the variable to be set.
type ActionTargetVariable struct {
	Name  string `yaml:"name"`
	Scope string `yaml:"scope"` // "runtime", "request", "folder", "collection", "environment"
}

// ---- HTTP primitives ----

// HttpRequestHeader is an HTTP header with optional disabled state.
type HttpRequestHeader struct {
	Name        string      `yaml:"name"`
	Value       string      `yaml:"value"`
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
}

// HttpResponseHeader is an HTTP response header.
type HttpResponseHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// GrpcMetadata is a gRPC metadata entry.
type GrpcMetadata struct {
	Name        string      `yaml:"name"`
	Value       string      `yaml:"value"`
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
}

// HttpRequestParam is a query or path parameter.
type HttpRequestParam struct {
	Name        string      `yaml:"name"`
	Value       string      `yaml:"value"`
	Type        string      `yaml:"type"` // "query" or "path"
	Description Description `yaml:"description,omitempty"`
	Disabled    bool        `yaml:"disabled,omitempty"`
}

// ---- Request defaults and settings ----

// RequestDefaults holds default request configuration at the collection or folder level.
type RequestDefaults struct {
	Headers   []HttpRequestHeader `yaml:"headers,omitempty"`
	Metadata  []GrpcMetadata      `yaml:"metadata,omitempty"`
	Auth      Auth                `yaml:"auth,omitempty"`
	Variables []Variable          `yaml:"variables,omitempty"`
	Scripts   Scripts             `yaml:"scripts,omitempty"`
	Settings  *RequestSettings    `yaml:"settings,omitempty"`
}

// RequestSettings holds default settings scoped to HTTP and GraphQL requests.
type RequestSettings struct {
	HTTP    *HttpRequestSettings `yaml:"http,omitempty"`
	GraphQL *HttpRequestSettings `yaml:"graphql,omitempty"`
}

// HttpRequestSettings holds HTTP execution settings, each of which can be a
// concrete value or the string "inherit" to defer to the parent scope.
type HttpRequestSettings struct {
	EncodeURL       InheritableBool `yaml:"encodeUrl,omitempty"`
	Timeout         InheritableInt  `yaml:"timeout,omitempty"`
	FollowRedirects InheritableBool `yaml:"followRedirects,omitempty"`
	MaxRedirects    InheritableInt  `yaml:"maxRedirects,omitempty"`
}

// InheritableBool is a boolean that can also carry the value "inherit".
type InheritableBool struct {
	Set     bool
	Inherit bool
	Value   bool
}

// IsZero reports whether the field is unset (used by yaml.v3 omitempty).
func (b InheritableBool) IsZero() bool { return !b.Set }

func (b InheritableBool) MarshalYAML() (any, error) {
	if b.Inherit {
		return "inherit", nil
	}
	return b.Value, nil
}

func (b *InheritableBool) UnmarshalYAML(value *yaml.Node) error {
	b.Set = true
	if value.Tag == "!!str" && value.Value == "inherit" {
		b.Inherit = true
		return nil
	}
	return value.Decode(&b.Value)
}

// InheritableInt is an integer that can also carry the value "inherit".
type InheritableInt struct {
	Set     bool
	Inherit bool
	Value   int
}

// IsZero reports whether the field is unset (used by yaml.v3 omitempty).
func (i InheritableInt) IsZero() bool { return !i.Set }

func (i InheritableInt) MarshalYAML() (any, error) {
	if i.Inherit {
		return "inherit", nil
	}
	return i.Value, nil
}

func (i *InheritableInt) UnmarshalYAML(value *yaml.Node) error {
	i.Set = true
	if value.Tag == "!!str" && value.Value == "inherit" {
		i.Inherit = true
		return nil
	}
	return value.Decode(&i.Value)
}
