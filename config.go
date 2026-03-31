package opencollection

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// CollectionConfig holds optional collection-wide configuration.
type CollectionConfig struct {
	Environments       []Environment       `yaml:"environments,omitempty"`
	Protobuf           *Protobuf           `yaml:"protobuf,omitempty"`
	Proxy              *Proxy              `yaml:"proxy,omitempty"`
	ClientCertificates []ClientCertificate `yaml:"clientCertificates,omitempty"`
}

// ---- Environments ----

// Environment describes a named environment with variables and optional certificates.
type Environment struct {
	Name               string              `yaml:"name"`
	Color              string              `yaml:"color,omitempty"`
	Description        Description         `yaml:"description,omitempty"`
	Variables          []EnvVariable       `yaml:"variables,omitempty"`
	ClientCertificates []ClientCertificate `yaml:"clientCertificates,omitempty"`
	Extends            string              `yaml:"extends,omitempty"`
	DotEnvFilePath     string              `yaml:"dotEnvFilePath,omitempty"`
}

// EnvVariable is a discriminated union of Variable and SecretVariable.
type EnvVariable struct {
	Variable       *Variable
	SecretVariable *SecretVariable
}

func (v EnvVariable) MarshalYAML() (any, error) {
	if v.SecretVariable != nil {
		return v.SecretVariable, nil
	}
	if v.Variable != nil {
		return v.Variable, nil
	}
	return nil, errors.New("opencollection: EnvVariable has neither Variable nor SecretVariable set")
}

func (v *EnvVariable) UnmarshalYAML(value *yaml.Node) error {
	var probe struct {
		Secret bool `yaml:"secret"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	if probe.Secret {
		v.SecretVariable = &SecretVariable{}
		return value.Decode(v.SecretVariable)
	}
	v.Variable = &Variable{}
	return value.Decode(v.Variable)
}

// ---- Protobuf ----

// Protobuf holds protobuf file and import-path configuration.
type Protobuf struct {
	ProtoFiles  []ProtoFile       `yaml:"protoFiles,omitempty"`
	ImportPaths []ProtoImportPath `yaml:"importPaths,omitempty"`
}

// ProtoFile references a .proto file.
type ProtoFile struct {
	Type string `yaml:"type"` // "file"
	Path string `yaml:"path"`
}

// ProtoImportPath is an import path used when resolving proto imports.
type ProtoImportPath struct {
	Path     string `yaml:"path"`
	Disabled bool   `yaml:"disabled,omitempty"`
}

// ---- Proxy ----

// Proxy holds proxy configuration for the collection.
type Proxy struct {
	Enabled bool                   `yaml:"enabled,omitempty"`
	Inherit bool                   `yaml:"inherit,omitempty"`
	Config  *ProxyConnectionConfig `yaml:"config,omitempty"`
}

// ProxyConnectionConfig holds proxy connection details.
type ProxyConnectionConfig struct {
	Protocol    string         `yaml:"protocol,omitempty"`
	Hostname    string         `yaml:"hostname,omitempty"`
	Port        int            `yaml:"port,omitempty"`
	Auth        ProxyAuthField `yaml:"auth,omitempty"`
	BypassProxy string         `yaml:"bypassProxy,omitempty"`
}

// ProxyAuthField is either disabled (false) or a ProxyAuth credentials object.
type ProxyAuthField struct {
	Disabled bool
	Auth     *ProxyAuth
}

// IsZero reports whether no auth is configured (used by yaml.v3 omitempty).
func (p ProxyAuthField) IsZero() bool { return !p.Disabled && p.Auth == nil }

func (p ProxyAuthField) MarshalYAML() (any, error) {
	if p.Disabled {
		return false, nil
	}
	if p.Auth != nil {
		return p.Auth, nil
	}
	return nil, nil
}

func (p *ProxyAuthField) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!bool" && value.Value == "false" {
		p.Disabled = true
		return nil
	}
	p.Auth = &ProxyAuth{}
	return value.Decode(p.Auth)
}

// ProxyAuth holds proxy authentication credentials.
type ProxyAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ---- Client certificates ----

// ClientCertificate is a discriminated union of PemCertificate and Pkcs12Certificate.
type ClientCertificate struct {
	PEM    *PemCertificate
	PKCS12 *Pkcs12Certificate
}

func (c ClientCertificate) MarshalYAML() (any, error) {
	if c.PEM != nil {
		return c.PEM, nil
	}
	if c.PKCS12 != nil {
		return c.PKCS12, nil
	}
	return nil, errors.New("opencollection: ClientCertificate has neither PEM nor PKCS12 set")
}

func (c *ClientCertificate) UnmarshalYAML(value *yaml.Node) error {
	var probe struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&probe); err != nil {
		return err
	}
	switch probe.Type {
	case "pem":
		c.PEM = &PemCertificate{}
		return value.Decode(c.PEM)
	case "pkcs12":
		c.PKCS12 = &Pkcs12Certificate{}
		return value.Decode(c.PKCS12)
	default:
		return fmt.Errorf("opencollection: unknown client certificate type %q", probe.Type)
	}
}

// PemCertificate uses separate PEM-encoded certificate and key files.
type PemCertificate struct {
	Domain              string `yaml:"domain"`
	Type                string `yaml:"type"` // "pem"
	CertificateFilePath string `yaml:"certificateFilePath"`
	PrivateKeyFilePath  string `yaml:"privateKeyFilePath"`
	Passphrase          string `yaml:"passphrase,omitempty"`
}

// Pkcs12Certificate uses a PKCS#12/PFX file.
type Pkcs12Certificate struct {
	Domain         string `yaml:"domain"`
	Type           string `yaml:"type"` // "pkcs12"
	PKCS12FilePath string `yaml:"pkcs12FilePath"`
	Passphrase     string `yaml:"passphrase,omitempty"`
}
