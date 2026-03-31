# go-opencollection

Go library for reading, writing, validating, and building [OpenCollection](https://opencollection.dev) YAML files — the open specification for API collections.

## Installation

```sh
go get go-opencollection
```

## Quick start

```go
// Read a collection from a file or directory
c, err := opencollection.Open("my-api.yml")

// Build a collection programmatically
c := opencollection.New("My API").
    Summary("Internal REST API").
    CollectionVersion("1.0.0").
    AddHttpRequest(
        opencollection.BuildHttpRequest("Get Users", "GET", "{{baseUrl}}/users").
            Header("Accept", "application/json").
            BearerAuth("{{token}}").
            Build(),
    ).
    Build()

// Validate against the OpenCollection schema
if err := opencollection.Validate(c); err != nil {
    log.Fatal(err)
}

// Write to disk
err = opencollection.Write("my-api.yml", c)
```

## I/O

| Function | Description |
|---|---|
| `Open(path)` | Read a collection from a file **or** an unbundled directory |
| `ParseFile(path)` | Read a bundled YAML file |
| `Parse(data)` | Unmarshal a YAML byte slice |
| `Write(path, c)` | Write to a file or directory depending on `c.Bundled` |
| `WriteFile(path, c)` | Always write a single bundled YAML file |
| `Marshal(c)` | Serialize to YAML bytes |
| `ReadDir(dir)` | Read an unbundled directory layout |
| `WriteDir(dir, c)` | Write an unbundled directory layout |

### Bundled vs. unbundled

Set `c.Bundled = true` to write everything into a single YAML file.  
Set `c.Bundled = false` (the default) to use a directory layout:

```
my-api/
  opencollection.yml       # collection root
  environments/
    production.yml
    staging.yml
  get-users.yml            # request items
  users/
    folder.yml             # folder metadata
    create-user.yml
    get-user-by-id.yml
```

## Validation

```go
err := opencollection.Validate(c)
```

`Validate` marshals the collection to its canonical YAML form, round-trips it through JSON, and validates it against the embedded OpenCollection JSON schema. Errors are returned as `*jsonschema.ValidationError` with full detail.

## Builders

All builders return `*Builder` for chaining and end with `.Build()`.

### Collection

```go
c := opencollection.New("Petstore API").
    Summary("OpenAPI Petstore example").
    CollectionVersion("2.0.0").
    Author("Alice", "alice@example.com", "").
    Bundled(true).
    Docs("Full API documentation at https://petstore.example.com/docs").
    Extension("x-team", "platform").
    Environment(
        opencollection.NewEnvironment("Production").
            Color("#d32f2f").
            Var("baseUrl", "https://api.petstore.example.com").
            Secret("apiKey", "string").
            Build(),
    ).
    DefaultRequest(
        opencollection.NewRequestDefaults().
            Header("X-Client-Version", "1.0").
            BearerAuth("{{apiKey}}").
            Build(),
    ).
    AddHttpRequest(req).
    AddFolder(folder).
    Build()
```

### HTTP request

```go
req := opencollection.BuildHttpRequest("Create Pet", "POST", "{{baseUrl}}/pets").
    Description("Creates a new pet in the store").
    Tag("pets").
    Header("Content-Type", "application/json").
    Header("Accept", "application/json").
    QueryParam("dryRun", "false").
    JSONBody(`{"name":"{{petName}}","status":"available"}`).
    BearerAuth("{{token}}").
    Var("petName", "Fido").
    BeforeRequest(`pm.request.headers.add({ key: "X-Request-ID", value: uuid() })`).
    Tests(`pm.test("status is 201", () => pm.response.to.have.status(201))`).
    Assert("response.status", "==", "201").
    Build()
```

Available body helpers: `JSONBody`, `TextBody`, `XMLBody`, `FormBody`.  
Available auth helpers: `BearerAuth`, `BasicAuth`, `InheritAuth`, `Auth(auth)`.  
Available script helpers: `BeforeRequest`, `AfterResponse`, `Tests`, `Script(type, code)`.

### GraphQL request

```go
gql := opencollection.BuildGraphQLRequest("List Pets", "https://api.example.com/graphql").
    Header("Accept", "application/json").
    Query(`query ListPets($limit: Int) { pets(limit: $limit) { id name } }`, `{"limit":10}`).
    BearerAuth("{{token}}").
    Assert("response.status", "==", "200").
    Build()
```

### gRPC request

```go
grpc := opencollection.BuildGrpcRequest("Get Pet", "grpc.petstore.example.com:443", "PetService/GetPet").
    MethodType("unary").
    ProtoFile("./proto/petstore.proto").
    Metadata("authorization", "Bearer {{token}}").
    Message(`{"id":"{{petId}}"}`).
    Build()
```

Method types: `"unary"`, `"client-streaming"`, `"server-streaming"`, `"bidi-streaming"`.

### Folder

```go
folder := opencollection.NewFolder("Pets").
    Description("Pet management endpoints").
    Seq(1).
    Tag("pets").
    DefaultRequest(
        opencollection.NewRequestDefaults().
            Header("X-Namespace", "pets").
            Build(),
    ).
    AddHttpRequest(req).
    AddFolder(nestedFolder).
    Build()
```

### Environment

```go
env := opencollection.NewEnvironment("Staging").
    Color("#f57c00").
    Extends("Production").
    DotEnvFile(".env.staging").
    Var("baseUrl", "https://staging.petstore.example.com").
    Secret("apiKey", "string").
    Build()
```

### Request defaults

```go
defaults := opencollection.NewRequestDefaults().
    Header("X-Client", "go-sdk").
    BearerAuth("{{token}}").
    Var("retries", "3").
    Script("before-request", `console.log("sending", pm.request.url)`).
    Build()
```

## Auth

Standalone auth constructors can be passed anywhere an `Auth` is accepted:

```go
opencollection.BearerAuth("{{token}}")
opencollection.BasicAuth("user", "pass")
opencollection.APIKeyAuth("X-API-Key", "{{apiKey}}", "header") // placement: "header" or "query"
opencollection.InheritAuth()  // inherit from parent scope
```

OAuth 2.0 is configured directly on the `Auth` struct:

```go
auth := opencollection.Auth{
    IsSet: true,
    OAuth2: &opencollection.AuthOAuth2{
        ClientCredentials: &opencollection.OAuth2ClientCredentialsFlow{
            Type:           "oauth2",
            Flow:           "client_credentials",
            AccessTokenURL: "https://auth.example.com/token",
            Credentials: &opencollection.OAuth2ClientCredentials{
                ClientID:     "{{clientId}}",
                ClientSecret: "{{clientSecret}}",
            },
            Scope: "read:pets write:pets",
        },
    },
}
```

## Data types

### `Description`

Can be a plain string, a typed content+MIME object, or null:

```go
opencollection.StringDescription("plain text")

opencollection.Description{
    IsSet:    true,
    Content:  "# Heading\nmarkdown body",
    MIMEType: "text/markdown",
}

opencollection.Description{IsSet: true, IsNull: true}
```

### `InheritableBool` / `InheritableInt`

Fields on `HttpRequestSettings` (e.g. `FollowRedirects`, `Timeout`) can hold a concrete value or `"inherit"`:

```go
settings := opencollection.HttpRequestSettings{
    FollowRedirects: opencollection.InheritableBool{Set: true, Value: true},
    Timeout:         opencollection.InheritableInt{Set: true, Inherit: true},
}
```

### `VariableValue`

A variable value can be a simple string, a typed `{type, data}` object, or a list of named variants:

```go
// simple
opencollection.VariableValue{Simple: "hello"}

// typed
opencollection.VariableValue{Typed: &opencollection.TypedVariableValue{Type: "number", Data: "42"}}

// variants
opencollection.VariableValue{Variants: []opencollection.VariableValueVariant{
    {Title: "US East", Value: opencollection.VariableValue{Simple: "us-east-1"}},
    {Title: "EU West", Selected: true, Value: opencollection.VariableValue{Simple: "eu-west-1"}},
}}
```

## Examples

See the [examples/](examples/) directory for complete, runnable programs:

| Example | Description |
|---|---|
| [build_and_write](examples/build_and_write/main.go) | Build a collection with HTTP, GraphQL, gRPC items and write it to disk |
| [read_and_inspect](examples/read_and_inspect/main.go) | Open a collection file and iterate over its items |
| [validate](examples/validate/main.go) | Parse a collection file and validate it against the schema |
| [environments](examples/environments/main.go) | Add multiple environments with plain and secret variables |
| [unbundled_dir](examples/unbundled_dir/main.go) | Write and read back an unbundled directory layout |
