package opencollection_test

import (
	"testing"

	oc "go-opencollection"
)

// ---- CollectionBuilder ----

func TestNew_Defaults(t *testing.T) {
	c := oc.New("My API").Build()
	if c.Info.Name != "My API" {
		t.Errorf("Name = %q", c.Info.Name)
	}
	if c.OpenCollection != "1" {
		t.Errorf("OpenCollection = %q", c.OpenCollection)
	}
	if c.Bundled {
		t.Error("Bundled should default to false")
	}
}

func TestCollectionBuilder_SpecVersion(t *testing.T) {
	c := oc.New("API").SpecVersion("2").Build()
	if c.OpenCollection != "2" {
		t.Errorf("OpenCollection = %q", c.OpenCollection)
	}
}

func TestCollectionBuilder_Summary(t *testing.T) {
	c := oc.New("API").Summary("A test API").Build()
	if c.Info.Summary != "A test API" {
		t.Errorf("Summary = %q", c.Info.Summary)
	}
}

func TestCollectionBuilder_CollectionVersion(t *testing.T) {
	c := oc.New("API").CollectionVersion("3.1.0").Build()
	if c.Info.Version != "3.1.0" {
		t.Errorf("Version = %q", c.Info.Version)
	}
}

func TestCollectionBuilder_Author(t *testing.T) {
	c := oc.New("API").
		Author("Alice", "alice@example.com", "https://alice.dev").
		Author("Bob", "", "").
		Build()
	if len(c.Info.Authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(c.Info.Authors))
	}
	if c.Info.Authors[0].Name != "Alice" || c.Info.Authors[0].URL != "https://alice.dev" {
		t.Errorf("unexpected author[0]: %+v", c.Info.Authors[0])
	}
}

func TestCollectionBuilder_Bundled(t *testing.T) {
	c := oc.New("API").Bundled(true).Build()
	if !c.Bundled {
		t.Error("Bundled should be true")
	}
}

func TestCollectionBuilder_Extension(t *testing.T) {
	c := oc.New("API").
		Extension("x-owner", "team-alpha").
		Extension("x-version", 3).
		Build()
	if c.Extensions["x-owner"] != "team-alpha" {
		t.Errorf("extension x-owner = %v", c.Extensions["x-owner"])
	}
}

func TestCollectionBuilder_Environment(t *testing.T) {
	c := oc.New("API").
		Environment(oc.NewEnvironment("dev").Var("host", "localhost").Build()).
		Build()
	if c.Config == nil || len(c.Config.Environments) != 1 {
		t.Fatal("expected one environment")
	}
	env := c.Config.Environments[0]
	if env.Name != "dev" {
		t.Errorf("env.Name = %q", env.Name)
	}
	if len(env.Variables) != 1 || env.Variables[0].Variable.Value.Simple != "localhost" {
		t.Errorf("unexpected variables: %+v", env.Variables)
	}
}

func TestCollectionBuilder_DefaultRequest(t *testing.T) {
	c := oc.New("API").
		DefaultRequest(oc.NewRequestDefaults().Header("X-App", "test").Build()).
		Build()
	if c.Request == nil || len(c.Request.Headers) != 1 {
		t.Fatal("expected default request with header")
	}
	if c.Request.Headers[0].Name != "X-App" {
		t.Errorf("header.Name = %q", c.Request.Headers[0].Name)
	}
}

func TestCollectionBuilder_AddHttpRequest(t *testing.T) {
	c := oc.New("API").
		AddHttpRequest(oc.BuildHttpRequest("Ping", "GET", "/ping").Build()).
		Build()
	if len(c.Items) != 1 || c.Items[0].HttpRequest == nil {
		t.Fatal("expected one HTTP request item")
	}
}

func TestCollectionBuilder_AddFolder(t *testing.T) {
	c := oc.New("API").
		AddFolder(oc.NewFolder("Users").Build()).
		Build()
	if len(c.Items) != 1 || c.Items[0].Folder == nil {
		t.Fatal("expected one folder item")
	}
}

// ---- EnvironmentBuilder ----

func TestEnvironmentBuilder_Defaults(t *testing.T) {
	e := oc.NewEnvironment("staging").Build()
	if e.Name != "staging" {
		t.Errorf("Name = %q", e.Name)
	}
}

func TestEnvironmentBuilder_Color(t *testing.T) {
	e := oc.NewEnvironment("prod").Color("#ff0000").Build()
	if e.Color != "#ff0000" {
		t.Errorf("Color = %q", e.Color)
	}
}

func TestEnvironmentBuilder_Extends(t *testing.T) {
	e := oc.NewEnvironment("prod").Extends("dev").Build()
	if e.Extends != "dev" {
		t.Errorf("Extends = %q", e.Extends)
	}
}

func TestEnvironmentBuilder_Secret(t *testing.T) {
	e := oc.NewEnvironment("env").Secret("API_KEY", "string").Build()
	if len(e.Variables) != 1 || e.Variables[0].SecretVariable == nil {
		t.Fatal("expected SecretVariable")
	}
	if e.Variables[0].SecretVariable.Type != "string" {
		t.Errorf("Type = %q", e.Variables[0].SecretVariable.Type)
	}
}

// ---- HttpRequestBuilder ----

func TestBuildHttpRequest_Defaults(t *testing.T) {
	r := oc.BuildHttpRequest("List", "GET", "/items").Build()
	if r.Info.Name != "List" {
		t.Errorf("Name = %q", r.Info.Name)
	}
	if r.Info.Type != "http" {
		t.Errorf("Type = %q", r.Info.Type)
	}
	if r.Http.Method != "GET" || r.Http.URL != "/items" {
		t.Errorf("method/url mismatch: %q %q", r.Http.Method, r.Http.URL)
	}
	if r.Runtime != nil {
		t.Error("Runtime should be nil when nothing is set")
	}
}

func TestBuildHttpRequest_Description(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Description("Returns items").Build()
	if !r.Info.Description.IsSet || r.Info.Description.Content != "Returns items" {
		t.Errorf("unexpected description: %+v", r.Info.Description)
	}
}

func TestBuildHttpRequest_Seq(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Seq(3).Build()
	if r.Info.Seq == nil || *r.Info.Seq != 3 {
		t.Errorf("Seq = %v", r.Info.Seq)
	}
}

func TestBuildHttpRequest_Tags(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Tag("smoke").Tag("regression").Build()
	if len(r.Info.Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", r.Info.Tags)
	}
}

func TestBuildHttpRequest_Headers(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Header("Accept", "application/json").Build()
	if len(r.Http.Headers) != 1 || r.Http.Headers[0].Value != "application/json" {
		t.Errorf("unexpected headers: %+v", r.Http.Headers)
	}
}

func TestBuildHttpRequest_QueryParam(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").QueryParam("limit", "10").Build()
	if len(r.Http.Params) != 1 || r.Http.Params[0].Type != "query" {
		t.Errorf("unexpected params: %+v", r.Http.Params)
	}
}

func TestBuildHttpRequest_PathParam(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/users/:id").PathParam("id", "42").Build()
	if len(r.Http.Params) != 1 || r.Http.Params[0].Type != "path" {
		t.Errorf("unexpected params: %+v", r.Http.Params)
	}
}

func TestBuildHttpRequest_JSONBody(t *testing.T) {
	r := oc.BuildHttpRequest("R", "POST", "/").JSONBody(`{"k":"v"}`).Build()
	if r.Http.Body.Body == nil || r.Http.Body.Body.Raw == nil {
		t.Fatal("expected raw body")
	}
	if r.Http.Body.Body.Raw.Type != "json" {
		t.Errorf("body.Type = %q", r.Http.Body.Body.Raw.Type)
	}
}

func TestBuildHttpRequest_BearerAuth(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").BearerAuth("t0k3n").Build()
	if r.Runtime == nil || !r.Runtime.Auth.IsSet || r.Runtime.Auth.Bearer == nil {
		t.Fatal("expected bearer auth in runtime")
	}
	if r.Runtime.Auth.Bearer.Token != "t0k3n" {
		t.Errorf("Token = %q", r.Runtime.Auth.Bearer.Token)
	}
}

func TestBuildHttpRequest_BasicAuth(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").BasicAuth("user", "pass").Build()
	if r.Runtime.Auth.Basic == nil {
		t.Fatal("expected basic auth")
	}
}

func TestBuildHttpRequest_InheritAuth(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").InheritAuth().Build()
	if !r.Runtime.Auth.Inherit {
		t.Error("expected Inherit = true")
	}
}

func TestBuildHttpRequest_Assertion(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").
		Assert("res.status", "eq", "200").
		Build()
	if r.Runtime == nil || len(r.Runtime.Assertions) != 1 {
		t.Fatal("expected one assertion")
	}
	a := r.Runtime.Assertions[0]
	if a.Expression != "res.status" || a.Operator != "eq" || a.Value != "200" {
		t.Errorf("unexpected assertion: %+v", a)
	}
}

func TestBuildHttpRequest_Scripts(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").
		BeforeRequest("// setup").
		AfterResponse("// cleanup").
		Tests("expect(res.status).to.equal(200)").
		Build()
	if r.Runtime == nil || len(r.Runtime.Scripts) != 3 {
		t.Fatalf("expected 3 scripts, got %d", len(r.Runtime.Scripts))
	}
	types := []string{r.Runtime.Scripts[0].Type, r.Runtime.Scripts[1].Type, r.Runtime.Scripts[2].Type}
	expected := []string{"before-request", "after-response", "tests"}
	for i, tt := range expected {
		if types[i] != tt {
			t.Errorf("script[%d].Type = %q, want %q", i, types[i], tt)
		}
	}
}

func TestBuildHttpRequest_Var(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Var("token", "abc").Build()
	if r.Runtime == nil || len(r.Runtime.Variables) != 1 {
		t.Fatal("expected one variable")
	}
	if r.Runtime.Variables[0].Name != "token" {
		t.Errorf("Name = %q", r.Runtime.Variables[0].Name)
	}
}

func TestBuildHttpRequest_Docs(t *testing.T) {
	r := oc.BuildHttpRequest("R", "GET", "/").Docs("# Docs").Build()
	if r.Docs != "# Docs" {
		t.Errorf("Docs = %q", r.Docs)
	}
}

// ---- FolderBuilder ----

func TestNewFolder_Defaults(t *testing.T) {
	f := oc.NewFolder("Products").Build()
	if f.Info.Name != "Products" {
		t.Errorf("Name = %q", f.Info.Name)
	}
	if f.Info.Type != "folder" {
		t.Errorf("Type = %q", f.Info.Type)
	}
}

func TestFolderBuilder_AddFolder(t *testing.T) {
	outer := oc.NewFolder("Outer").
		AddFolder(oc.NewFolder("Inner")).
		Build()
	if len(outer.Items) != 1 || outer.Items[0].Folder == nil {
		t.Fatal("expected nested folder")
	}
	if outer.Items[0].Folder.Info.Name != "Inner" {
		t.Errorf("inner folder name = %q", outer.Items[0].Folder.Info.Name)
	}
}

func TestFolderBuilder_DefaultRequest(t *testing.T) {
	f := oc.NewFolder("F").
		DefaultRequest(oc.NewRequestDefaults().BearerAuth("tok").Build()).
		Build()
	if f.Request == nil || !f.Request.Auth.IsSet {
		t.Error("expected default request with auth")
	}
}

// ---- GraphQLRequestBuilder ----

func TestBuildGraphQLRequest(t *testing.T) {
	r := oc.BuildGraphQLRequest("Get Posts", "https://api.example.com/graphql").
		Query("{ posts { id title } }", "").
		Header("X-Custom", "value").
		BearerAuth("tok").
		Build()

	if r.Info.Type != "graphql" {
		t.Errorf("Type = %q", r.Info.Type)
	}
	if r.GraphQL.Body.Body == nil || r.GraphQL.Body.Body.Query != "{ posts { id title } }" {
		t.Errorf("unexpected body: %+v", r.GraphQL.Body)
	}
	if len(r.GraphQL.Headers) != 1 {
		t.Errorf("expected 1 header, got %d", len(r.GraphQL.Headers))
	}
	if r.Runtime == nil || r.Runtime.Auth.Bearer == nil {
		t.Error("expected bearer auth")
	}
}

// ---- GrpcRequestBuilder ----

func TestBuildGrpcRequest(t *testing.T) {
	r := oc.BuildGrpcRequest("Get User", "localhost:50051", "user.UserService/GetUser").
		MethodType("unary").
		ProtoFile("./user.proto").
		Metadata("x-trace", "123").
		Message(`{"id": "42"}`).
		Build()

	if r.Info.Type != "grpc" {
		t.Errorf("Type = %q", r.Info.Type)
	}
	if r.Grpc.MethodType != "unary" {
		t.Errorf("MethodType = %q", r.Grpc.MethodType)
	}
	if r.Grpc.ProtoFilePath != "./user.proto" {
		t.Errorf("ProtoFilePath = %q", r.Grpc.ProtoFilePath)
	}
	if len(r.Grpc.Metadata) != 1 {
		t.Errorf("expected 1 metadata, got %d", len(r.Grpc.Metadata))
	}
	if r.Grpc.Message.Message != `{"id": "42"}` {
		t.Errorf("Message = %q", r.Grpc.Message.Message)
	}
}

// ---- Auth constructors ----

func TestBearerAuth(t *testing.T) {
	a := oc.BearerAuth("tok")
	if !a.IsSet || a.Bearer == nil || a.Bearer.Token != "tok" {
		t.Errorf("unexpected: %+v", a)
	}
}

func TestBasicAuth(t *testing.T) {
	a := oc.BasicAuth("u", "p")
	if a.Basic == nil || a.Basic.Username != "u" {
		t.Errorf("unexpected: %+v", a)
	}
}

func TestAPIKeyAuth(t *testing.T) {
	a := oc.APIKeyAuth("X-Key", "v", "header")
	if a.APIKey == nil || a.APIKey.Placement != "header" {
		t.Errorf("unexpected: %+v", a)
	}
}

func TestInheritAuthConstructor(t *testing.T) {
	a := oc.InheritAuth()
	if !a.IsSet || !a.Inherit {
		t.Errorf("unexpected: %+v", a)
	}
}

// ---- RequestDefaultsBuilder ----

func TestRequestDefaultsBuilder(t *testing.T) {
	rd := oc.NewRequestDefaults().
		BearerAuth("tok").
		Header("X-App", "myapp").
		Var("env", "dev").
		Script("before-request", "// init").
		Build()

	if !rd.Auth.IsSet || rd.Auth.Bearer == nil {
		t.Error("expected bearer auth")
	}
	if len(rd.Headers) != 1 || rd.Headers[0].Name != "X-App" {
		t.Errorf("unexpected headers: %+v", rd.Headers)
	}
	if len(rd.Variables) != 1 || rd.Variables[0].Name != "env" {
		t.Errorf("unexpected variables: %+v", rd.Variables)
	}
	if len(rd.Scripts) != 1 || rd.Scripts[0].Type != "before-request" {
		t.Errorf("unexpected scripts: %+v", rd.Scripts)
	}
}
