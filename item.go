package opencollection

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Item is a discriminated union of all collection item types.
// The concrete type is determined by the info.type field (or top-level type for scripts).
type Item struct {
	HttpRequest    *HttpRequest
	GraphQLRequest *GraphQLRequest
	GrpcRequest    *GrpcRequest
	WebSocket      *WebSocketRequest
	Folder         *Folder
	Script         *ScriptFile
}

func (i Item) MarshalYAML() (any, error) {
	if i.HttpRequest != nil {
		return i.HttpRequest, nil
	}
	if i.GraphQLRequest != nil {
		return i.GraphQLRequest, nil
	}
	if i.GrpcRequest != nil {
		return i.GrpcRequest, nil
	}
	if i.WebSocket != nil {
		return i.WebSocket, nil
	}
	if i.Folder != nil {
		return i.Folder, nil
	}
	if i.Script != nil {
		return i.Script, nil
	}
	return nil, errors.New("opencollection: Item has no type set")
}

func (i *Item) UnmarshalYAML(value *yaml.Node) error {
	// ScriptFile uses a top-level "type" field rather than info.type.
	var topLevel struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&topLevel); err != nil {
		return err
	}
	if topLevel.Type == "script" {
		i.Script = &ScriptFile{}
		return value.Decode(i.Script)
	}

	var probe struct {
		Info struct {
			Type string `yaml:"type"`
		} `yaml:"info"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Info.Type {
	case "http":
		i.HttpRequest = &HttpRequest{}
		return value.Decode(i.HttpRequest)
	case "graphql":
		i.GraphQLRequest = &GraphQLRequest{}
		return value.Decode(i.GraphQLRequest)
	case "grpc":
		i.GrpcRequest = &GrpcRequest{}
		return value.Decode(i.GrpcRequest)
	case "websocket":
		i.WebSocket = &WebSocketRequest{}
		return value.Decode(i.WebSocket)
	case "folder":
		i.Folder = &Folder{}
		return value.Decode(i.Folder)
	default:
		return fmt.Errorf("opencollection: unknown item type %q", probe.Info.Type)
	}
}

// ---- HTTP request ----

// HttpRequest represents a single HTTP request item.
type HttpRequest struct {
	Info     HttpRequestInfo      `yaml:"info"`
	Http     HttpRequestDetails   `yaml:"http"`
	Runtime  *HttpRequestRuntime  `yaml:"runtime,omitempty"`
	Settings *HttpRequestSettings `yaml:"settings,omitempty"`
	Examples []HttpRequestExample `yaml:"examples,omitempty"`
	Docs     string               `yaml:"docs,omitempty"`
}

// HttpRequestInfo holds HTTP request metadata.
type HttpRequestInfo struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "http"
	Description Description `yaml:"description,omitempty"`
	Seq         *float64    `yaml:"seq,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
}

// HttpRequestDetails holds the HTTP protocol fields.
type HttpRequestDetails struct {
	Method  string              `yaml:"method"`
	URL     string              `yaml:"url"`
	Headers []HttpRequestHeader `yaml:"headers,omitempty"`
	Params  []HttpRequestParam  `yaml:"params,omitempty"`
	Body    HttpBodyField       `yaml:"body,omitempty"`
}

// HttpRequestRuntime holds scripts, assertions, actions, and auth for a request.
type HttpRequestRuntime struct {
	Variables  []Variable  `yaml:"variables,omitempty"`
	Scripts    Scripts     `yaml:"scripts,omitempty"`
	Assertions []Assertion `yaml:"assertions,omitempty"`
	Actions    []Action    `yaml:"actions,omitempty"`
	Auth       Auth        `yaml:"auth,omitempty"`
}

// HttpRequestExample is a sample request/response pair.
type HttpRequestExample struct {
	Name        string           `yaml:"name,omitempty"`
	Description Description      `yaml:"description,omitempty"`
	Request     *ExampleRequest  `yaml:"request,omitempty"`
	Response    *ExampleResponse `yaml:"response,omitempty"`
}

// ExampleRequest is the request side of an example.
type ExampleRequest struct {
	URL     string              `yaml:"url,omitempty"`
	Method  string              `yaml:"method,omitempty"`
	Headers []HttpRequestHeader `yaml:"headers,omitempty"`
	Params  []HttpRequestParam  `yaml:"params,omitempty"`
	Body    HttpBodyField       `yaml:"body,omitempty"`
}

// ExampleResponse is the response side of an example.
type ExampleResponse struct {
	Status     int                  `yaml:"status,omitempty"`
	StatusText string               `yaml:"statusText,omitempty"`
	Headers    []HttpResponseHeader `yaml:"headers,omitempty"`
	Body       *ExampleResponseBody `yaml:"body,omitempty"`
}

// ExampleResponseBody holds the response body content.
type ExampleResponseBody struct {
	Type string `yaml:"type"` // "json", "text", "xml", "html", "binary"
	Data string `yaml:"data"`
}

// ---- GraphQL request ----

// GraphQLRequest represents a GraphQL request item.
type GraphQLRequest struct {
	Info     GraphQLRequestInfo    `yaml:"info"`
	GraphQL  GraphQLRequestDetails `yaml:"graphql"`
	Runtime  *GraphQLRequestRuntime `yaml:"runtime,omitempty"`
	Settings *HttpRequestSettings   `yaml:"settings,omitempty"`
	Docs     string                 `yaml:"docs,omitempty"`
}

// GraphQLRequestInfo holds GraphQL request metadata.
type GraphQLRequestInfo struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "graphql"
	Description Description `yaml:"description,omitempty"`
	Seq         *float64    `yaml:"seq,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
}

// GraphQLRequestDetails holds the GraphQL protocol fields.
type GraphQLRequestDetails struct {
	Method  string              `yaml:"method,omitempty"`
	URL     string              `yaml:"url,omitempty"`
	Headers []HttpRequestHeader `yaml:"headers,omitempty"`
	Params  []HttpRequestParam  `yaml:"params,omitempty"`
	Body    GraphQLBodyField    `yaml:"body,omitempty"`
}

// GraphQLRequestRuntime holds scripts, assertions, and auth for a GraphQL request.
type GraphQLRequestRuntime struct {
	Variables  []Variable  `yaml:"variables,omitempty"`
	Scripts    Scripts     `yaml:"scripts,omitempty"`
	Assertions []Assertion `yaml:"assertions,omitempty"`
	Actions    []Action    `yaml:"actions,omitempty"`
	Auth       Auth        `yaml:"auth,omitempty"`
}

// ---- gRPC request ----

// GrpcRequest represents a gRPC request item.
type GrpcRequest struct {
	Info    GrpcRequestInfo     `yaml:"info"`
	Grpc    GrpcRequestDetails  `yaml:"grpc"`
	Runtime *GrpcRequestRuntime `yaml:"runtime,omitempty"`
	Docs    string              `yaml:"docs,omitempty"`
}

// GrpcRequestInfo holds gRPC request metadata.
type GrpcRequestInfo struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "grpc"
	Description Description `yaml:"description,omitempty"`
	Seq         *float64    `yaml:"seq,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
}

// GrpcRequestDetails holds the gRPC protocol fields.
type GrpcRequestDetails struct {
	URL           string           `yaml:"url,omitempty"`
	Method        string           `yaml:"method,omitempty"`
	MethodType    string           `yaml:"methodType,omitempty"` // "unary", "client-streaming", "server-streaming", "bidi-streaming"
	ProtoFilePath string           `yaml:"protoFilePath,omitempty"`
	Metadata      []GrpcMetadata   `yaml:"metadata,omitempty"`
	Message       GrpcMessageField `yaml:"message,omitempty"`
}

// GrpcRequestRuntime holds scripts, assertions, and auth for a gRPC request.
type GrpcRequestRuntime struct {
	Variables  []Variable  `yaml:"variables,omitempty"`
	Scripts    Scripts     `yaml:"scripts,omitempty"`
	Assertions []Assertion `yaml:"assertions,omitempty"`
	Auth       Auth        `yaml:"auth,omitempty"`
}

// ---- WebSocket request ----

// WebSocketRequest represents a WebSocket request item.
type WebSocketRequest struct {
	Info      WebSocketRequestInfo     `yaml:"info"`
	WebSocket WebSocketRequestDetails  `yaml:"websocket"`
	Runtime   *WebSocketRequestRuntime `yaml:"runtime,omitempty"`
	Docs      string                   `yaml:"docs,omitempty"`
}

// WebSocketRequestInfo holds WebSocket request metadata.
type WebSocketRequestInfo struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "websocket"
	Description Description `yaml:"description,omitempty"`
	Seq         *float64    `yaml:"seq,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
}

// WebSocketRequestDetails holds the WebSocket protocol fields.
type WebSocketRequestDetails struct {
	URL     string                `yaml:"url,omitempty"`
	Headers []HttpRequestHeader   `yaml:"headers,omitempty"`
	Message WebSocketMessageField `yaml:"message,omitempty"`
}

// WebSocketRequestRuntime holds scripts and auth for a WebSocket request.
type WebSocketRequestRuntime struct {
	Variables []Variable `yaml:"variables,omitempty"`
	Scripts   Scripts    `yaml:"scripts,omitempty"`
	Auth      Auth       `yaml:"auth,omitempty"`
}

// ---- Folder ----

// Folder groups related items and can carry default request settings.
type Folder struct {
	Info    FolderInfo       `yaml:"info"`
	Items   []Item           `yaml:"items,omitempty"`
	Request *RequestDefaults `yaml:"request,omitempty"`
	Docs    Description      `yaml:"docs,omitempty"`
}

// FolderInfo holds folder metadata.
type FolderInfo struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // "folder"
	Description Description `yaml:"description,omitempty"`
	Seq         *float64    `yaml:"seq,omitempty"`
	Tags        []string    `yaml:"tags,omitempty"`
}

// ---- Script ----

// ScriptFile is a shared JavaScript module item.
type ScriptFile struct {
	Type   string `yaml:"type"` // "script"
	Script string `yaml:"script"`
}
