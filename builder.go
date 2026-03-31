package opencollection

// CollectionBuilder builds a Collection using a fluent API.
type CollectionBuilder struct {
	c Collection
}

// New creates a CollectionBuilder with the given collection name.
// The opencollection spec version defaults to "1".
func New(name string) *CollectionBuilder {
	return &CollectionBuilder{
		c: Collection{
			OpenCollection: "1",
			Info:           Info{Name: name},
		},
	}
}

// SpecVersion sets the opencollection spec version string (e.g. "1").
func (b *CollectionBuilder) SpecVersion(v string) *CollectionBuilder {
	b.c.OpenCollection = v
	return b
}

// Summary sets the collection summary.
func (b *CollectionBuilder) Summary(s string) *CollectionBuilder {
	b.c.Info.Summary = s
	return b
}

// CollectionVersion sets the collection's own version string (e.g. "1.0.0").
func (b *CollectionBuilder) CollectionVersion(v string) *CollectionBuilder {
	b.c.Info.Version = v
	return b
}

// Author adds an author entry to the collection.
func (b *CollectionBuilder) Author(name, email, url string) *CollectionBuilder {
	b.c.Info.Authors = append(b.c.Info.Authors, Author{
		Name:  name,
		Email: email,
		URL:   url,
	})
	return b
}

// Bundled marks the collection as a single bundled file.
func (b *CollectionBuilder) Bundled(bundled bool) *CollectionBuilder {
	b.c.Bundled = bundled
	return b
}

// Docs sets the collection-level documentation.
func (b *CollectionBuilder) Docs(content string) *CollectionBuilder {
	b.c.Docs = StringDescription(content)
	return b
}

// Extension adds or replaces a key in the extensions map.
func (b *CollectionBuilder) Extension(key string, value any) *CollectionBuilder {
	if b.c.Extensions == nil {
		b.c.Extensions = make(map[string]any)
	}
	b.c.Extensions[key] = value
	return b
}

// DefaultRequest sets collection-level default request settings.
func (b *CollectionBuilder) DefaultRequest(req RequestDefaults) *CollectionBuilder {
	b.c.Request = &req
	return b
}

// Environment adds one or more environments to the collection config.
func (b *CollectionBuilder) Environment(envs ...Environment) *CollectionBuilder {
	if b.c.Config == nil {
		b.c.Config = &CollectionConfig{}
	}
	b.c.Config.Environments = append(b.c.Config.Environments, envs...)
	return b
}

// Proxy sets the proxy configuration.
func (b *CollectionBuilder) Proxy(proxy Proxy) *CollectionBuilder {
	if b.c.Config == nil {
		b.c.Config = &CollectionConfig{}
	}
	b.c.Config.Proxy = &proxy
	return b
}

// AddItem appends one or more items (request, folder, or script) to the collection.
func (b *CollectionBuilder) AddItem(items ...Item) *CollectionBuilder {
	b.c.Items = append(b.c.Items, items...)
	return b
}

// AddHttpRequest appends one or more HTTP requests to the collection.
func (b *CollectionBuilder) AddHttpRequest(rs ...*HttpRequest) *CollectionBuilder {
	for _, r := range rs {
		b.AddItem(Item{HttpRequest: r})
	}
	return b
}

// AddGraphQLRequest appends one or more GraphQL requests to the collection.
func (b *CollectionBuilder) AddGraphQLRequest(rs ...*GraphQLRequest) *CollectionBuilder {
	for _, r := range rs {
		b.AddItem(Item{GraphQLRequest: r})
	}
	return b
}

// AddGrpcRequest appends one or more gRPC requests to the collection.
func (b *CollectionBuilder) AddGrpcRequest(rs ...*GrpcRequest) *CollectionBuilder {
	for _, r := range rs {
		b.AddItem(Item{GrpcRequest: r})
	}
	return b
}

// AddWebSocketRequest appends one or more WebSocket requests to the collection.
func (b *CollectionBuilder) AddWebSocketRequest(rs ...*WebSocketRequest) *CollectionBuilder {
	for _, r := range rs {
		b.AddItem(Item{WebSocket: r})
	}
	return b
}

// AddFolder appends one or more folders to the collection.
func (b *CollectionBuilder) AddFolder(fs ...*Folder) *CollectionBuilder {
	for _, f := range fs {
		b.AddItem(Item{Folder: f})
	}
	return b
}

// Build returns the constructed Collection.
func (b *CollectionBuilder) Build() *Collection {
	c := b.c
	return &c
}

// ---- FolderBuilder ----

// FolderBuilder builds a Folder using a fluent API.
type FolderBuilder struct {
	f Folder
}

// NewFolder creates a FolderBuilder with the given folder name.
func NewFolder(name string) *FolderBuilder {
	return &FolderBuilder{
		f: Folder{
			Info: FolderInfo{Name: name, Type: "folder"},
		},
	}
}

// Description sets the folder description.
func (fb *FolderBuilder) Description(content string) *FolderBuilder {
	fb.f.Info.Description = StringDescription(content)
	return fb
}

// Docs sets the folder documentation.
func (fb *FolderBuilder) Docs(content string) *FolderBuilder {
	fb.f.Docs = StringDescription(content)
	return fb
}

// Seq sets the display order sequence number.
func (fb *FolderBuilder) Seq(seq float64) *FolderBuilder {
	fb.f.Info.Seq = &seq
	return fb
}

// Tag adds one or more tags to the folder.
func (fb *FolderBuilder) Tag(tags ...string) *FolderBuilder {
	fb.f.Info.Tags = append(fb.f.Info.Tags, tags...)
	return fb
}

// DefaultRequest sets folder-level default request settings.
func (fb *FolderBuilder) DefaultRequest(req RequestDefaults) *FolderBuilder {
	fb.f.Request = &req
	return fb
}

// AddItem appends one or more items to the folder.
func (fb *FolderBuilder) AddItem(items ...Item) *FolderBuilder {
	fb.f.Items = append(fb.f.Items, items...)
	return fb
}

// AddHttpRequest appends one or more HTTP requests to the folder.
func (fb *FolderBuilder) AddHttpRequest(rs ...*HttpRequest) *FolderBuilder {
	for _, r := range rs {
		fb.AddItem(Item{HttpRequest: r})
	}
	return fb
}

// AddGraphQLRequest appends one or more GraphQL requests to the folder.
func (fb *FolderBuilder) AddGraphQLRequest(rs ...*GraphQLRequest) *FolderBuilder {
	for _, r := range rs {
		fb.AddItem(Item{GraphQLRequest: r})
	}
	return fb
}

// AddGrpcRequest appends one or more gRPC requests to the folder.
func (fb *FolderBuilder) AddGrpcRequest(rs ...*GrpcRequest) *FolderBuilder {
	for _, r := range rs {
		fb.AddItem(Item{GrpcRequest: r})
	}
	return fb
}

// AddFolder adds one or more nested folders.
func (fb *FolderBuilder) AddFolder(nested ...*FolderBuilder) *FolderBuilder {
	for _, n := range nested {
		fb.AddItem(Item{Folder: n.Build()})
	}
	return fb
}

// Build returns the constructed Folder.
func (fb *FolderBuilder) Build() *Folder {
	f := fb.f
	return &f
}

// ---- EnvironmentBuilder ----

// EnvironmentBuilder builds an Environment using a fluent API.
type EnvironmentBuilder struct {
	e Environment
}

// NewEnvironment creates an EnvironmentBuilder with the given name.
func NewEnvironment(name string) *EnvironmentBuilder {
	return &EnvironmentBuilder{e: Environment{Name: name}}
}

// Color sets the environment display color.
func (eb *EnvironmentBuilder) Color(color string) *EnvironmentBuilder {
	eb.e.Color = color
	return eb
}

// Description sets the environment description.
func (eb *EnvironmentBuilder) Description(content string) *EnvironmentBuilder {
	eb.e.Description = StringDescription(content)
	return eb
}

// Extends sets the name of an environment this one inherits from.
func (eb *EnvironmentBuilder) Extends(name string) *EnvironmentBuilder {
	eb.e.Extends = name
	return eb
}

// DotEnvFile sets the path to a .env file to load variables from.
func (eb *EnvironmentBuilder) DotEnvFile(path string) *EnvironmentBuilder {
	eb.e.DotEnvFilePath = path
	return eb
}

// Var adds a plain string variable to the environment.
func (eb *EnvironmentBuilder) Var(name, value string) *EnvironmentBuilder {
	eb.e.Variables = append(eb.e.Variables, EnvVariable{
		Variable: &Variable{
			Name:  name,
			Value: VariableValue{Simple: value},
		},
	})
	return eb
}

// Secret adds a secret (value-less) variable declaration.
func (eb *EnvironmentBuilder) Secret(name, typ string) *EnvironmentBuilder {
	eb.e.Variables = append(eb.e.Variables, EnvVariable{
		SecretVariable: &SecretVariable{
			Secret: true,
			Name:   name,
			Type:   typ,
		},
	})
	return eb
}

// Build returns the constructed Environment.
func (eb *EnvironmentBuilder) Build() Environment {
	return eb.e
}

// ---- HttpRequestBuilder ----

// HttpRequestBuilder builds an HttpRequest using a fluent API.
type HttpRequestBuilder struct {
	r HttpRequest
}

// BuildHttpRequest creates an HttpRequestBuilder for the given method and URL.
func BuildHttpRequest(name, method, url string) *HttpRequestBuilder {
	return &HttpRequestBuilder{
		r: HttpRequest{
			Info: HttpRequestInfo{Name: name, Type: "http"},
			Http: HttpRequestDetails{Method: method, URL: url},
		},
	}
}

// Description sets the request description.
func (rb *HttpRequestBuilder) Description(content string) *HttpRequestBuilder {
	rb.r.Info.Description = StringDescription(content)
	return rb
}

// Docs sets the request documentation.
func (rb *HttpRequestBuilder) Docs(docs string) *HttpRequestBuilder {
	rb.r.Docs = docs
	return rb
}

// Seq sets the display order sequence number.
func (rb *HttpRequestBuilder) Seq(seq float64) *HttpRequestBuilder {
	rb.r.Info.Seq = &seq
	return rb
}

// Tag adds one or more tags to the request.
func (rb *HttpRequestBuilder) Tag(tags ...string) *HttpRequestBuilder {
	rb.r.Info.Tags = append(rb.r.Info.Tags, tags...)
	return rb
}

// Header adds an HTTP request header.
func (rb *HttpRequestBuilder) Header(name, value string) *HttpRequestBuilder {
	rb.r.Http.Headers = append(rb.r.Http.Headers, HttpRequestHeader{Name: name, Value: value})
	return rb
}

// QueryParam adds a query parameter.
func (rb *HttpRequestBuilder) QueryParam(name, value string) *HttpRequestBuilder {
	rb.r.Http.Params = append(rb.r.Http.Params, HttpRequestParam{Name: name, Value: value, Type: "query"})
	return rb
}

// PathParam adds a path parameter.
func (rb *HttpRequestBuilder) PathParam(name, value string) *HttpRequestBuilder {
	rb.r.Http.Params = append(rb.r.Http.Params, HttpRequestParam{Name: name, Value: value, Type: "path"})
	return rb
}

// JSONBody sets a raw JSON body.
func (rb *HttpRequestBuilder) JSONBody(data string) *HttpRequestBuilder {
	rb.r.Http.Body = HttpBodyField{Body: &HttpBody{Raw: &RawBody{Type: "json", Data: data}}}
	return rb
}

// TextBody sets a plain-text body.
func (rb *HttpRequestBuilder) TextBody(data string) *HttpRequestBuilder {
	rb.r.Http.Body = HttpBodyField{Body: &HttpBody{Raw: &RawBody{Type: "text", Data: data}}}
	return rb
}

// XMLBody sets an XML body.
func (rb *HttpRequestBuilder) XMLBody(data string) *HttpRequestBuilder {
	rb.r.Http.Body = HttpBodyField{Body: &HttpBody{Raw: &RawBody{Type: "xml", Data: data}}}
	return rb
}

// FormBody sets a form-urlencoded body.
func (rb *HttpRequestBuilder) FormBody(fields ...FormURLEncodedField) *HttpRequestBuilder {
	rb.r.Http.Body = HttpBodyField{Body: &HttpBody{FormURLEncoded: &FormURLEncodedBody{
		Type: "form-urlencoded",
		Data: fields,
	}}}
	return rb
}

// Auth sets the request-level auth, overriding any inherited auth.
func (rb *HttpRequestBuilder) Auth(auth Auth) *HttpRequestBuilder {
	rb.runtime().Auth = auth
	return rb
}

// BearerAuth sets bearer token authentication.
func (rb *HttpRequestBuilder) BearerAuth(token string) *HttpRequestBuilder {
	return rb.Auth(Auth{IsSet: true, Bearer: &AuthBearer{Type: "bearer", Token: token}})
}

// BasicAuth sets HTTP Basic authentication.
func (rb *HttpRequestBuilder) BasicAuth(username, password string) *HttpRequestBuilder {
	return rb.Auth(Auth{IsSet: true, Basic: &AuthBasic{Type: "basic", Username: username, Password: password}})
}

// InheritAuth marks this request as inheriting auth from its parent.
func (rb *HttpRequestBuilder) InheritAuth() *HttpRequestBuilder {
	return rb.Auth(Auth{IsSet: true, Inherit: true})
}

// Var adds a runtime variable to the request.
func (rb *HttpRequestBuilder) Var(name, value string) *HttpRequestBuilder {
	rb.runtime().Variables = append(rb.runtime().Variables, Variable{
		Name:  name,
		Value: VariableValue{Simple: value},
	})
	return rb
}

// Script adds a lifecycle script to the request.
func (rb *HttpRequestBuilder) Script(typ, code string) *HttpRequestBuilder {
	rb.runtime().Scripts = append(rb.runtime().Scripts, Script{Type: typ, Code: code})
	return rb
}

// BeforeRequest adds a before-request script.
func (rb *HttpRequestBuilder) BeforeRequest(code string) *HttpRequestBuilder {
	return rb.Script("before-request", code)
}

// AfterResponse adds an after-response script.
func (rb *HttpRequestBuilder) AfterResponse(code string) *HttpRequestBuilder {
	return rb.Script("after-response", code)
}

// Tests adds a test script.
func (rb *HttpRequestBuilder) Tests(code string) *HttpRequestBuilder {
	return rb.Script("tests", code)
}

// Assert adds a response assertion.
func (rb *HttpRequestBuilder) Assert(expression, operator, value string) *HttpRequestBuilder {
	rb.runtime().Assertions = append(rb.runtime().Assertions, Assertion{
		Expression: expression,
		Operator:   operator,
		Value:      value,
	})
	return rb
}

// Settings configures HTTP execution settings for this request.
func (rb *HttpRequestBuilder) Settings(s HttpRequestSettings) *HttpRequestBuilder {
	rb.r.Settings = &s
	return rb
}

// Build returns the constructed HttpRequest.
func (rb *HttpRequestBuilder) Build() *HttpRequest {
	r := rb.r
	return &r
}

func (rb *HttpRequestBuilder) runtime() *HttpRequestRuntime {
	if rb.r.Runtime == nil {
		rb.r.Runtime = &HttpRequestRuntime{}
	}
	return rb.r.Runtime
}

// ---- GraphQLRequestBuilder ----

// GraphQLRequestBuilder builds a GraphQLRequest using a fluent API.
type GraphQLRequestBuilder struct {
	r GraphQLRequest
}

// BuildGraphQLRequest creates a GraphQLRequestBuilder.
func BuildGraphQLRequest(name, url string) *GraphQLRequestBuilder {
	return &GraphQLRequestBuilder{
		r: GraphQLRequest{
			Info:    GraphQLRequestInfo{Name: name, Type: "graphql"},
			GraphQL: GraphQLRequestDetails{Method: "POST", URL: url},
		},
	}
}

// Description sets the request description.
func (rb *GraphQLRequestBuilder) Description(content string) *GraphQLRequestBuilder {
	rb.r.Info.Description = StringDescription(content)
	return rb
}

// Tag adds one or more tags to the request.
func (rb *GraphQLRequestBuilder) Tag(tags ...string) *GraphQLRequestBuilder {
	rb.r.Info.Tags = append(rb.r.Info.Tags, tags...)
	return rb
}

// Header adds a request header.
func (rb *GraphQLRequestBuilder) Header(name, value string) *GraphQLRequestBuilder {
	rb.r.GraphQL.Headers = append(rb.r.GraphQL.Headers, HttpRequestHeader{Name: name, Value: value})
	return rb
}

// Query sets the GraphQL query and optional variables JSON string.
func (rb *GraphQLRequestBuilder) Query(query, variables string) *GraphQLRequestBuilder {
	rb.r.GraphQL.Body = GraphQLBodyField{Body: &GraphQLBody{Query: query, Variables: variables}}
	return rb
}

// Auth sets the request-level auth.
func (rb *GraphQLRequestBuilder) Auth(auth Auth) *GraphQLRequestBuilder {
	rb.runtime().Auth = auth
	return rb
}

// BearerAuth sets bearer token authentication.
func (rb *GraphQLRequestBuilder) BearerAuth(token string) *GraphQLRequestBuilder {
	return rb.Auth(Auth{IsSet: true, Bearer: &AuthBearer{Type: "bearer", Token: token}})
}

// Assert adds a response assertion.
func (rb *GraphQLRequestBuilder) Assert(expression, operator, value string) *GraphQLRequestBuilder {
	rb.runtime().Assertions = append(rb.runtime().Assertions, Assertion{
		Expression: expression,
		Operator:   operator,
		Value:      value,
	})
	return rb
}

// Build returns the constructed GraphQLRequest.
func (rb *GraphQLRequestBuilder) Build() *GraphQLRequest {
	r := rb.r
	return &r
}

func (rb *GraphQLRequestBuilder) runtime() *GraphQLRequestRuntime {
	if rb.r.Runtime == nil {
		rb.r.Runtime = &GraphQLRequestRuntime{}
	}
	return rb.r.Runtime
}

// ---- GrpcRequestBuilder ----

// GrpcRequestBuilder builds a GrpcRequest using a fluent API.
type GrpcRequestBuilder struct {
	r GrpcRequest
}

// BuildGrpcRequest creates a GrpcRequestBuilder for the given service method.
func BuildGrpcRequest(name, url, method string) *GrpcRequestBuilder {
	return &GrpcRequestBuilder{
		r: GrpcRequest{
			Info: GrpcRequestInfo{Name: name, Type: "grpc"},
			Grpc: GrpcRequestDetails{URL: url, Method: method, MethodType: "unary"},
		},
	}
}

// Description sets the request description.
func (rb *GrpcRequestBuilder) Description(content string) *GrpcRequestBuilder {
	rb.r.Info.Description = StringDescription(content)
	return rb
}

// MethodType sets the streaming type: "unary", "client-streaming", "server-streaming", "bidi-streaming".
func (rb *GrpcRequestBuilder) MethodType(t string) *GrpcRequestBuilder {
	rb.r.Grpc.MethodType = t
	return rb
}

// ProtoFile sets the path to the .proto file.
func (rb *GrpcRequestBuilder) ProtoFile(path string) *GrpcRequestBuilder {
	rb.r.Grpc.ProtoFilePath = path
	return rb
}

// Metadata adds a gRPC metadata entry.
func (rb *GrpcRequestBuilder) Metadata(name, value string) *GrpcRequestBuilder {
	rb.r.Grpc.Metadata = append(rb.r.Grpc.Metadata, GrpcMetadata{Name: name, Value: value})
	return rb
}

// Message sets the request message JSON string.
func (rb *GrpcRequestBuilder) Message(msg string) *GrpcRequestBuilder {
	rb.r.Grpc.Message = GrpcMessageField{Message: msg}
	return rb
}

// Auth sets the request-level auth.
func (rb *GrpcRequestBuilder) Auth(auth Auth) *GrpcRequestBuilder {
	rb.runtime().Auth = auth
	return rb
}

// Build returns the constructed GrpcRequest.
func (rb *GrpcRequestBuilder) Build() *GrpcRequest {
	r := rb.r
	return &r
}

func (rb *GrpcRequestBuilder) runtime() *GrpcRequestRuntime {
	if rb.r.Runtime == nil {
		rb.r.Runtime = &GrpcRequestRuntime{}
	}
	return rb.r.Runtime
}

// ---- Auth constructors ----

// BearerAuth returns an Auth configured with a bearer token.
func BearerAuth(token string) Auth {
	return Auth{IsSet: true, Bearer: &AuthBearer{Type: "bearer", Token: token}}
}

// BasicAuth returns an Auth configured with basic credentials.
func BasicAuth(username, password string) Auth {
	return Auth{IsSet: true, Basic: &AuthBasic{Type: "basic", Username: username, Password: password}}
}

// APIKeyAuth returns an Auth configured with an API key.
// placement is "header" or "query".
func APIKeyAuth(key, value, placement string) Auth {
	return Auth{IsSet: true, APIKey: &AuthAPIKey{Type: "apikey", Key: key, Value: value, Placement: placement}}
}

// InheritAuth returns an Auth that inherits from the parent scope.
func InheritAuth() Auth {
	return Auth{IsSet: true, Inherit: true}
}

// ---- RequestDefaults builder ----

// RequestDefaultsBuilder builds a RequestDefaults using a fluent API.
type RequestDefaultsBuilder struct {
	r RequestDefaults
}

// NewRequestDefaults creates a RequestDefaultsBuilder.
func NewRequestDefaults() *RequestDefaultsBuilder {
	return &RequestDefaultsBuilder{}
}

// Auth sets the default auth.
func (rb *RequestDefaultsBuilder) Auth(auth Auth) *RequestDefaultsBuilder {
	rb.r.Auth = auth
	return rb
}

// BearerAuth sets bearer token as the default auth.
func (rb *RequestDefaultsBuilder) BearerAuth(token string) *RequestDefaultsBuilder {
	return rb.Auth(BearerAuth(token))
}

// Header adds a default header.
func (rb *RequestDefaultsBuilder) Header(name, value string) *RequestDefaultsBuilder {
	rb.r.Headers = append(rb.r.Headers, HttpRequestHeader{Name: name, Value: value})
	return rb
}

// Var adds a default variable.
func (rb *RequestDefaultsBuilder) Var(name, value string) *RequestDefaultsBuilder {
	rb.r.Variables = append(rb.r.Variables, Variable{
		Name:  name,
		Value: VariableValue{Simple: value},
	})
	return rb
}

// Script adds a lifecycle script to the defaults.
func (rb *RequestDefaultsBuilder) Script(typ, code string) *RequestDefaultsBuilder {
	rb.r.Scripts = append(rb.r.Scripts, Script{Type: typ, Code: code})
	return rb
}

// Build returns the constructed RequestDefaults.
func (rb *RequestDefaultsBuilder) Build() RequestDefaults {
	return rb.r
}
