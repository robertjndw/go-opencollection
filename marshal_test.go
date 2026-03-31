package opencollection_test

import (
	"strings"
	"testing"

	oc "go-opencollection"
	"gopkg.in/yaml.v3"
)

// ---- Item union ----

func TestItem_UnmarshalHTTP(t *testing.T) {
	const src = `
info:
  name: List Users
  type: http
http:
  method: GET
  url: https://api.example.com/users
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.HttpRequest == nil {
		t.Fatal("expected HttpRequest, got nil")
	}
	if item.HttpRequest.Info.Name != "List Users" {
		t.Errorf("Name = %q", item.HttpRequest.Info.Name)
	}
	if item.HttpRequest.Http.Method != "GET" {
		t.Errorf("Method = %q", item.HttpRequest.Http.Method)
	}
}

func TestItem_UnmarshalFolder(t *testing.T) {
	const src = `
info:
  name: Users
  type: folder
items: []
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.Folder == nil {
		t.Fatal("expected Folder, got nil")
	}
	if item.Folder.Info.Name != "Users" {
		t.Errorf("Name = %q", item.Folder.Info.Name)
	}
}

func TestItem_UnmarshalGraphQL(t *testing.T) {
	const src = `
info:
  name: Get Posts
  type: graphql
graphql:
  url: https://api.example.com/graphql
  body:
    query: "{ posts { id } }"
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.GraphQLRequest == nil {
		t.Fatal("expected GraphQLRequest")
	}
}

func TestItem_UnmarshalGrpc(t *testing.T) {
	const src = `
info:
  name: Get User
  type: grpc
grpc:
  url: localhost:50051
  method: user.UserService/GetUser
  methodType: unary
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.GrpcRequest == nil {
		t.Fatal("expected GrpcRequest")
	}
	if item.GrpcRequest.Grpc.MethodType != "unary" {
		t.Errorf("MethodType = %q", item.GrpcRequest.Grpc.MethodType)
	}
}

func TestItem_UnmarshalWebSocket(t *testing.T) {
	const src = `
info:
  name: Live Feed
  type: websocket
websocket:
  url: wss://api.example.com/feed
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.WebSocket == nil {
		t.Fatal("expected WebSocketRequest")
	}
}

func TestItem_UnmarshalScript(t *testing.T) {
	const src = `
type: script
script: "console.log('hello')"
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if item.Script == nil {
		t.Fatal("expected Script")
	}
	if item.Script.Script != "console.log('hello')" {
		t.Errorf("Script = %q", item.Script.Script)
	}
}

func TestItem_UnmarshalUnknownType(t *testing.T) {
	const src = `
info:
  name: Unknown
  type: ftp
`
	var item oc.Item
	if err := yaml.Unmarshal([]byte(src), &item); err == nil {
		t.Error("expected error for unknown item type")
	}
}

func TestItem_MarshalHTTP(t *testing.T) {
	seq := 1.0
	item := oc.Item{
		HttpRequest: &oc.HttpRequest{
			Info: oc.HttpRequestInfo{Name: "Ping", Type: "http", Seq: &seq},
			Http: oc.HttpRequestDetails{Method: "GET", URL: "/ping"},
		},
	}
	data, err := yaml.Marshal(&item)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, "name: Ping") {
		t.Errorf("missing name in output:\n%s", s)
	}
	if !strings.Contains(s, "method: GET") {
		t.Errorf("missing method in output:\n%s", s)
	}
}

// ---- Auth union ----

func TestAuth_UnmarshalBearer(t *testing.T) {
	const src = `type: bearer
token: secret-token`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !auth.IsSet || auth.Bearer == nil {
		t.Fatal("expected Bearer auth")
	}
	if auth.Bearer.Token != "secret-token" {
		t.Errorf("Token = %q", auth.Bearer.Token)
	}
}

func TestAuth_UnmarshalBasic(t *testing.T) {
	const src = `type: basic
username: alice
password: p4ssw0rd`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if auth.Basic == nil {
		t.Fatal("expected Basic auth")
	}
	if auth.Basic.Username != "alice" {
		t.Errorf("Username = %q", auth.Basic.Username)
	}
}

func TestAuth_UnmarshalAPIKey(t *testing.T) {
	const src = `type: apikey
key: X-API-Key
value: my-api-key
placement: header`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if auth.APIKey == nil {
		t.Fatal("expected APIKey auth")
	}
	if auth.APIKey.Placement != "header" {
		t.Errorf("Placement = %q", auth.APIKey.Placement)
	}
}

func TestAuth_UnmarshalInherit(t *testing.T) {
	const src = `inherit`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !auth.Inherit {
		t.Error("expected Inherit = true")
	}
}

func TestAuth_UnmarshalOAuth2ClientCredentials(t *testing.T) {
	const src = `
type: oauth2
flow: client_credentials
accessTokenUrl: https://auth.example.com/token
credentials:
  clientId: my-client
  clientSecret: my-secret
`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if auth.OAuth2 == nil || auth.OAuth2.ClientCredentials == nil {
		t.Fatal("expected OAuth2 client_credentials flow")
	}
	if auth.OAuth2.ClientCredentials.Credentials.ClientID != "my-client" {
		t.Errorf("ClientID = %q", auth.OAuth2.ClientCredentials.Credentials.ClientID)
	}
}

func TestAuth_UnmarshalOAuth2AuthorizationCode(t *testing.T) {
	const src = `
type: oauth2
flow: authorization_code
authorizationUrl: https://auth.example.com/authorize
accessTokenUrl: https://auth.example.com/token
`
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(src), &auth); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if auth.OAuth2 == nil || auth.OAuth2.AuthorizationCode == nil {
		t.Fatal("expected OAuth2 authorization_code flow")
	}
}

func TestAuth_UnmarshalUnknown(t *testing.T) {
	var auth oc.Auth
	if err := yaml.Unmarshal([]byte(`type: saml`), &auth); err == nil {
		t.Error("expected error for unknown auth type")
	}
}

func TestAuth_IsZero(t *testing.T) {
	var auth oc.Auth
	if !auth.IsZero() {
		t.Error("zero Auth should report IsZero() = true")
	}
	auth.IsSet = true
	if auth.IsZero() {
		t.Error("set Auth should report IsZero() = false")
	}
}

func TestAuth_MarshalBearer(t *testing.T) {
	auth := oc.BearerAuth("tok")
	data, err := yaml.Marshal(&auth)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, "type: bearer") {
		t.Errorf("missing type in output:\n%s", s)
	}
	if !strings.Contains(s, "token: tok") {
		t.Errorf("missing token in output:\n%s", s)
	}
}

func TestAuth_MarshalInherit(t *testing.T) {
	auth := oc.InheritAuth()
	data, err := yaml.Marshal(&auth)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.TrimSpace(string(data)) != "inherit" {
		t.Errorf("expected 'inherit', got %q", string(data))
	}
}

// ---- Body types ----

func TestHttpBody_UnmarshalJSON(t *testing.T) {
	const src = `type: json
data: '{"key":"value"}'`
	var body oc.HttpBody
	if err := yaml.Unmarshal([]byte(src), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Raw == nil {
		t.Fatal("expected RawBody")
	}
	if body.Raw.Type != "json" {
		t.Errorf("Type = %q", body.Raw.Type)
	}
}

func TestHttpBody_UnmarshalFormURLEncoded(t *testing.T) {
	const src = `
type: form-urlencoded
data:
  - name: username
    value: alice
`
	var body oc.HttpBody
	if err := yaml.Unmarshal([]byte(src), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.FormURLEncoded == nil {
		t.Fatal("expected FormURLEncodedBody")
	}
	if len(body.FormURLEncoded.Data) != 1 || body.FormURLEncoded.Data[0].Name != "username" {
		t.Errorf("unexpected data: %+v", body.FormURLEncoded.Data)
	}
}

func TestHttpBody_UnmarshalMultipart(t *testing.T) {
	const src = `
type: multipart-form
data:
  - name: file
    type: file
    value: /tmp/upload.txt
`
	var body oc.HttpBody
	if err := yaml.Unmarshal([]byte(src), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.MultipartForm == nil {
		t.Fatal("expected MultipartFormBody")
	}
}

func TestHttpBody_UnmarshalUnknown(t *testing.T) {
	var body oc.HttpBody
	if err := yaml.Unmarshal([]byte(`type: binary`), &body); err == nil {
		t.Error("expected error for unknown body type")
	}
}

func TestHttpBodyField_UnmarshalVariants(t *testing.T) {
	const src = `
- title: With pagination
  selected: true
  body:
    type: json
    data: '{"page":1}'
- title: Without pagination
  body:
    type: json
    data: '{}'
`
	var field oc.HttpBodyField
	if err := yaml.Unmarshal([]byte(src), &field); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(field.Variants) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(field.Variants))
	}
	if field.Variants[0].Title != "With pagination" {
		t.Errorf("variant[0].Title = %q", field.Variants[0].Title)
	}
}

// ---- Description ----

func TestDescription_UnmarshalString(t *testing.T) {
	const src = `description: plain string`
	var s struct {
		Description oc.Description `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !s.Description.IsSet || s.Description.Content != "plain string" {
		t.Errorf("unexpected Description: %+v", s.Description)
	}
}

func TestDescription_UnmarshalTyped(t *testing.T) {
	const src = `
description:
  content: "# Hello"
  type: text/markdown
`
	var s struct {
		Description oc.Description `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Description.MIMEType != "text/markdown" {
		t.Errorf("MIMEType = %q", s.Description.MIMEType)
	}
}

func TestDescription_MarshalOmitWhenUnset(t *testing.T) {
	type wrapper struct {
		Desc oc.Description `yaml:"desc,omitempty"`
	}
	data, err := yaml.Marshal(wrapper{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "desc") {
		t.Errorf("unset Description should be omitted:\n%s", data)
	}
}

// ---- VariableValue ----

func TestVariableValue_UnmarshalSimple(t *testing.T) {
	const src = `value: hello`
	var s struct {
		Value oc.VariableValue `yaml:"value"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Value.Simple != "hello" {
		t.Errorf("Simple = %q", s.Value.Simple)
	}
}

func TestVariableValue_UnmarshalTyped(t *testing.T) {
	const src = `
value:
  type: number
  data: "42"
`
	var s struct {
		Value oc.VariableValue `yaml:"value"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Value.Typed == nil || s.Value.Typed.Type != "number" {
		t.Errorf("unexpected Typed: %+v", s.Value.Typed)
	}
}

func TestVariableValue_UnmarshalVariants(t *testing.T) {
	const src = `
value:
  - title: v1
    selected: true
    value: "one"
  - title: v2
    value: "two"
`
	var s struct {
		Value oc.VariableValue `yaml:"value"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(s.Value.Variants) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(s.Value.Variants))
	}
}

// ---- InheritableBool / InheritableInt ----

func TestInheritableBool_UnmarshalTrue(t *testing.T) {
	const src = `flag: true`
	var s struct {
		Flag oc.InheritableBool `yaml:"flag"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !s.Flag.Set || !s.Flag.Value {
		t.Errorf("unexpected: %+v", s.Flag)
	}
}

func TestInheritableBool_UnmarshalInherit(t *testing.T) {
	const src = `flag: inherit`
	var s struct {
		Flag oc.InheritableBool `yaml:"flag"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !s.Flag.Inherit {
		t.Error("expected Inherit = true")
	}
}

func TestInheritableInt_UnmarshalInherit(t *testing.T) {
	const src = `timeout: inherit`
	var s struct {
		Timeout oc.InheritableInt `yaml:"timeout"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !s.Timeout.Inherit {
		t.Error("expected Inherit = true")
	}
}

func TestInheritableInt_UnmarshalValue(t *testing.T) {
	const src = `timeout: 5000`
	var s struct {
		Timeout oc.InheritableInt `yaml:"timeout"`
	}
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.Timeout.Value != 5000 {
		t.Errorf("Value = %d", s.Timeout.Value)
	}
}
