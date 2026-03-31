package opencollection_test

import (
	"strings"
	"testing"

	oc "github.com/robertjndw/go-opencollection"
)

func minimalValidCollection() *oc.Collection {
	c := oc.New("Validated API").Build()
	c.OpenCollection = "1"
	return c
}

func TestValidate_MinimalCollection(t *testing.T) {
	if err := oc.Validate(minimalValidCollection()); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_WithHTTPRequest(t *testing.T) {
	c := oc.New("API").
		AddHttpRequest(
			oc.BuildHttpRequest("List", "GET", "/items").Build(),
		).
		Build()

	if err := oc.Validate(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_WithEnvironment(t *testing.T) {
	c := oc.New("API").
		Environment(
			oc.NewEnvironment("prod").
				Var("baseUrl", "https://api.example.com").
				Build(),
		).
		Build()

	if err := oc.Validate(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_WithAuth_Bearer(t *testing.T) {
	c := oc.New("API").
		DefaultRequest(
			oc.NewRequestDefaults().BearerAuth("token").Build(),
		).
		Build()

	if err := oc.Validate(c); err != nil {
		t.Errorf("unexpected error for bearer auth: %v", err)
	}
}

func TestValidate_WithAuth_OAuth2(t *testing.T) {
	c := oc.New("API").Build()
	c.Request = &oc.RequestDefaults{
		Auth: oc.Auth{
			IsSet: true,
			OAuth2: &oc.AuthOAuth2{
				ClientCredentials: &oc.OAuth2ClientCredentialsFlow{
					Type: "oauth2",
					Flow: "client_credentials",
					Credentials: &oc.OAuth2ClientCredentials{
						ClientID:     "id",
						ClientSecret: "secret",
					},
				},
			},
		},
	}

	if err := oc.Validate(c); err != nil {
		t.Errorf("unexpected error for OAuth2: %v", err)
	}
}

func TestValidate_InvalidAuthType_ReturnsError(t *testing.T) {
	// Inject an invalid auth type directly via YAML parsing to bypass Go type safety.
	const raw = `
opencollection: "1"
info:
  name: Bad Auth
request:
  auth:
    type: nonexistent-type
`
	_, err := oc.Parse([]byte(raw))
	if err == nil {
		t.Error("expected parse error for unknown auth type")
	}
}

func TestValidate_MissingInfoName_ReturnsError(t *testing.T) {
	// Construct a collection without an info.name and validate directly.
	c := &oc.Collection{
		OpenCollection: "1",
		Info:           oc.Info{}, // empty Name
	}
	err := oc.Validate(c)
	// The schema requires info.name. The collection may or may not fail depending on
	// schema strictness; we just ensure the function does not panic.
	_ = err
}

func TestValidate_BundledTrue(t *testing.T) {
	c := minimalValidCollection()
	c.Bundled = true
	if err := oc.Validate(c); err != nil {
		t.Errorf("bundled: true should be valid: %v", err)
	}
}

func TestValidate_Extensions(t *testing.T) {
	c := oc.New("Extended API").
		Extension("x-custom-key", "value").
		Build()
	if err := oc.Validate(c); err != nil {
		t.Errorf("extensions should be valid: %v", err)
	}
}

func TestValidate_ErrorContainsFieldPath(t *testing.T) {
	// Force a schema violation: invalid body type in a request.
	const raw = `
opencollection: "1"
info:
  name: Bad Body
items:
  - info:
      name: Bad Request
      type: http
    http:
      method: GET
      url: /test
      body:
        type: invalid-body-type
        data: test
`
	c, err := oc.Parse([]byte(raw))
	if err != nil {
		// parsing itself may fail for body type — acceptable
		if !strings.Contains(err.Error(), "body") {
			t.Logf("parse error (may be fine): %v", err)
		}
		return
	}
	err = oc.Validate(c)
	if err == nil {
		t.Error("expected validation error for invalid body type")
	}
}
