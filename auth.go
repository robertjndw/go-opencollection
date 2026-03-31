package opencollection

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Auth is a discriminated union of all supported authentication types.
// The zero value (IsSet == false) is omitted from YAML output.
type Auth struct {
	IsSet   bool
	Inherit bool
	AwsV4   *AuthAwsV4
	Basic   *AuthBasic
	Wsse    *AuthWsse
	Bearer  *AuthBearer
	Digest  *AuthDigest
	NTLM    *AuthNTLM
	APIKey  *AuthAPIKey
	OAuth2  *AuthOAuth2
}

// IsZero reports whether the auth is unset (used by yaml.v3 omitempty).
func (a Auth) IsZero() bool { return !a.IsSet }

func (a Auth) MarshalYAML() (any, error) {
	if a.Inherit {
		return "inherit", nil
	}
	if a.AwsV4 != nil {
		return a.AwsV4, nil
	}
	if a.Basic != nil {
		return a.Basic, nil
	}
	if a.Wsse != nil {
		return a.Wsse, nil
	}
	if a.Bearer != nil {
		return a.Bearer, nil
	}
	if a.Digest != nil {
		return a.Digest, nil
	}
	if a.NTLM != nil {
		return a.NTLM, nil
	}
	if a.APIKey != nil {
		return a.APIKey, nil
	}
	if a.OAuth2 != nil {
		return a.OAuth2, nil
	}
	return nil, errors.New("opencollection: Auth is set but has no concrete type")
}

func (a *Auth) UnmarshalYAML(value *yaml.Node) error {
	a.IsSet = true
	if value.Kind == yaml.ScalarNode && value.Value == "inherit" {
		a.Inherit = true
		return nil
	}
	var probe struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Type {
	case "awsv4":
		a.AwsV4 = &AuthAwsV4{}
		return value.Decode(a.AwsV4)
	case "basic":
		a.Basic = &AuthBasic{}
		return value.Decode(a.Basic)
	case "wsse":
		a.Wsse = &AuthWsse{}
		return value.Decode(a.Wsse)
	case "bearer":
		a.Bearer = &AuthBearer{}
		return value.Decode(a.Bearer)
	case "digest":
		a.Digest = &AuthDigest{}
		return value.Decode(a.Digest)
	case "ntlm":
		a.NTLM = &AuthNTLM{}
		return value.Decode(a.NTLM)
	case "apikey":
		a.APIKey = &AuthAPIKey{}
		return value.Decode(a.APIKey)
	case "oauth2":
		a.OAuth2 = &AuthOAuth2{}
		return value.Decode(a.OAuth2)
	default:
		return fmt.Errorf("opencollection: unknown auth type %q", probe.Type)
	}
}

// ---- Concrete auth types ----

// AuthAwsV4 holds AWS Signature Version 4 credentials.
type AuthAwsV4 struct {
	Type            string `yaml:"type"` // "awsv4"
	AccessKeyID     string `yaml:"accessKeyId,omitempty"`
	SecretAccessKey string `yaml:"secretAccessKey,omitempty"`
	SessionToken    string `yaml:"sessionToken,omitempty"`
	Service         string `yaml:"service,omitempty"`
	Region          string `yaml:"region,omitempty"`
	ProfileName     string `yaml:"profileName,omitempty"`
}

// AuthBasic holds HTTP Basic Auth credentials.
type AuthBasic struct {
	Type     string `yaml:"type"` // "basic"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// AuthWsse holds WSSE credentials.
type AuthWsse struct {
	Type     string `yaml:"type"` // "wsse"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// AuthBearer holds a bearer token.
type AuthBearer struct {
	Type  string `yaml:"type"` // "bearer"
	Token string `yaml:"token,omitempty"`
}

// AuthDigest holds HTTP Digest Auth credentials.
type AuthDigest struct {
	Type     string `yaml:"type"` // "digest"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// AuthNTLM holds NTLM credentials.
type AuthNTLM struct {
	Type     string `yaml:"type"` // "ntlm"
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Domain   string `yaml:"domain,omitempty"`
}

// AuthAPIKey holds an API key credential.
type AuthAPIKey struct {
	Type      string `yaml:"type"` // "apikey"
	Key       string `yaml:"key,omitempty"`
	Value     string `yaml:"value,omitempty"`
	Placement string `yaml:"placement,omitempty"` // "header" or "query"
}

// ---- OAuth 2.0 ----

// AuthOAuth2 is a discriminated union of OAuth 2.0 flows, keyed by the flow field.
type AuthOAuth2 struct {
	ClientCredentials                *OAuth2ClientCredentialsFlow
	ResourceOwnerPasswordCredentials *OAuth2ResourceOwnerPasswordFlow
	AuthorizationCode                *OAuth2AuthorizationCodeFlow
	Implicit                         *OAuth2ImplicitFlow
}

func (a AuthOAuth2) MarshalYAML() (any, error) {
	if a.ClientCredentials != nil {
		return a.ClientCredentials, nil
	}
	if a.ResourceOwnerPasswordCredentials != nil {
		return a.ResourceOwnerPasswordCredentials, nil
	}
	if a.AuthorizationCode != nil {
		return a.AuthorizationCode, nil
	}
	if a.Implicit != nil {
		return a.Implicit, nil
	}
	return nil, errors.New("opencollection: AuthOAuth2 has no flow set")
}

func (a *AuthOAuth2) UnmarshalYAML(value *yaml.Node) error {
	var probe struct {
		Flow string `yaml:"flow"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Flow {
	case "client_credentials":
		a.ClientCredentials = &OAuth2ClientCredentialsFlow{}
		return value.Decode(a.ClientCredentials)
	case "resource_owner_password_credentials":
		a.ResourceOwnerPasswordCredentials = &OAuth2ResourceOwnerPasswordFlow{}
		return value.Decode(a.ResourceOwnerPasswordCredentials)
	case "authorization_code":
		a.AuthorizationCode = &OAuth2AuthorizationCodeFlow{}
		return value.Decode(a.AuthorizationCode)
	case "implicit":
		a.Implicit = &OAuth2ImplicitFlow{}
		return value.Decode(a.Implicit)
	default:
		return fmt.Errorf("opencollection: unknown OAuth2 flow %q", probe.Flow)
	}
}

// OAuth2ClientCredentials holds client ID and secret for OAuth 2.0 requests.
type OAuth2ClientCredentials struct {
	ClientID     string `yaml:"clientId,omitempty"`
	ClientSecret string `yaml:"clientSecret,omitempty"`
	Placement    string `yaml:"placement,omitempty"` // "basic_auth_header" or "body"
}

// OAuth2ResourceOwner holds resource owner credentials.
type OAuth2ResourceOwner struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// OAuth2PKCE configures Proof Key for Code Exchange.
type OAuth2PKCE struct {
	Enabled bool   `yaml:"enabled,omitempty"`
	Method  string `yaml:"method,omitempty"` // "S256" or "plain"
}

// OAuth2AdditionalParameter is an extra parameter sent alongside an OAuth 2.0 request.
type OAuth2AdditionalParameter struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	Placement string `yaml:"placement,omitempty"` // "header", "query", "body"
}

// OAuth2TokenConfig configures how a token is stored and transported.
type OAuth2TokenConfig struct {
	ID        string                `yaml:"id,omitempty"`
	Placement *OAuth2TokenPlacement `yaml:"placement,omitempty"`
}

// OAuth2TokenPlacement specifies where the token is injected (header or query param).
type OAuth2TokenPlacement struct {
	Header string `yaml:"header,omitempty"`
	Query  string `yaml:"query,omitempty"`
}

// OAuth2Settings controls automatic token fetch and refresh behaviour.
type OAuth2Settings struct {
	AutoFetchToken   bool `yaml:"autoFetchToken,omitempty"`
	AutoRefreshToken bool `yaml:"autoRefreshToken,omitempty"`
}

// OAuth2AdditionalParameters groups extra parameters by request phase.
type OAuth2AdditionalParameters struct {
	AuthorizationRequest []OAuth2AdditionalParameter `yaml:"authorizationRequest,omitempty"`
	AccessTokenRequest   []OAuth2AdditionalParameter `yaml:"accessTokenRequest,omitempty"`
	RefreshTokenRequest  []OAuth2AdditionalParameter `yaml:"refreshTokenRequest,omitempty"`
}

// OAuth2ImplicitAdditionalParameters groups extra parameters for the implicit flow.
type OAuth2ImplicitAdditionalParameters struct {
	AuthorizationRequest []OAuth2AdditionalParameter `yaml:"authorizationRequest,omitempty"`
}

// OAuth2ImplicitCredentials holds only the client ID (implicit flow only needs clientId).
type OAuth2ImplicitCredentials struct {
	ClientID string `yaml:"clientId,omitempty"`
}

// OAuth2ClientCredentialsFlow is the OAuth 2.0 Client Credentials flow.
type OAuth2ClientCredentialsFlow struct {
	Type                 string                      `yaml:"type"` // "oauth2"
	Flow                 string                      `yaml:"flow"` // "client_credentials"
	AccessTokenURL       string                      `yaml:"accessTokenUrl,omitempty"`
	RefreshTokenURL      string                      `yaml:"refreshTokenUrl,omitempty"`
	Credentials          *OAuth2ClientCredentials    `yaml:"credentials,omitempty"`
	Scope                string                      `yaml:"scope,omitempty"`
	AdditionalParameters *OAuth2AdditionalParameters `yaml:"additionalParameters,omitempty"`
	TokenConfig          *OAuth2TokenConfig          `yaml:"tokenConfig,omitempty"`
	Settings             *OAuth2Settings             `yaml:"settings,omitempty"`
}

// OAuth2ResourceOwnerPasswordFlow is the OAuth 2.0 Resource Owner Password flow.
type OAuth2ResourceOwnerPasswordFlow struct {
	Type                 string                      `yaml:"type"` // "oauth2"
	Flow                 string                      `yaml:"flow"` // "resource_owner_password_credentials"
	AccessTokenURL       string                      `yaml:"accessTokenUrl,omitempty"`
	RefreshTokenURL      string                      `yaml:"refreshTokenUrl,omitempty"`
	Credentials          *OAuth2ClientCredentials    `yaml:"credentials,omitempty"`
	ResourceOwner        *OAuth2ResourceOwner        `yaml:"resourceOwner,omitempty"`
	Scope                string                      `yaml:"scope,omitempty"`
	AdditionalParameters *OAuth2AdditionalParameters `yaml:"additionalParameters,omitempty"`
	TokenConfig          *OAuth2TokenConfig          `yaml:"tokenConfig,omitempty"`
	Settings             *OAuth2Settings             `yaml:"settings,omitempty"`
}

// OAuth2AuthorizationCodeFlow is the OAuth 2.0 Authorization Code flow.
type OAuth2AuthorizationCodeFlow struct {
	Type                 string                      `yaml:"type"` // "oauth2"
	Flow                 string                      `yaml:"flow"` // "authorization_code"
	AuthorizationURL     string                      `yaml:"authorizationUrl,omitempty"`
	AccessTokenURL       string                      `yaml:"accessTokenUrl,omitempty"`
	RefreshTokenURL      string                      `yaml:"refreshTokenUrl,omitempty"`
	CallbackURL          string                      `yaml:"callbackUrl,omitempty"`
	Credentials          *OAuth2ClientCredentials    `yaml:"credentials,omitempty"`
	Scope                string                      `yaml:"scope,omitempty"`
	State                string                      `yaml:"state,omitempty"`
	PKCE                 *OAuth2PKCE                 `yaml:"pkce,omitempty"`
	AdditionalParameters *OAuth2AdditionalParameters `yaml:"additionalParameters,omitempty"`
	TokenConfig          *OAuth2TokenConfig          `yaml:"tokenConfig,omitempty"`
	Settings             *OAuth2Settings             `yaml:"settings,omitempty"`
}

// OAuth2ImplicitFlow is the OAuth 2.0 Implicit flow.
type OAuth2ImplicitFlow struct {
	Type                 string                              `yaml:"type"` // "oauth2"
	Flow                 string                              `yaml:"flow"` // "implicit"
	AuthorizationURL     string                              `yaml:"authorizationUrl,omitempty"`
	CallbackURL          string                              `yaml:"callbackUrl,omitempty"`
	Credentials          *OAuth2ImplicitCredentials          `yaml:"credentials,omitempty"`
	Scope                string                              `yaml:"scope,omitempty"`
	State                string                              `yaml:"state,omitempty"`
	AdditionalParameters *OAuth2ImplicitAdditionalParameters `yaml:"additionalParameters,omitempty"`
	TokenConfig          *OAuth2TokenConfig                  `yaml:"tokenConfig,omitempty"`
	Settings             *OAuth2Settings                     `yaml:"settings,omitempty"`
}
