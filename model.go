// Package opencollection implements reading, writing, validating, and building
// OpenCollection YAML files — the open specification for API collections.
//
// # File structure
//
//   - collection.go — Collection, Info, Author
//   - config.go     — CollectionConfig, Environment, Proxy, Protobuf, certificates
//   - auth.go       — Auth union and all authentication types
//   - body.go       — HTTP, GraphQL, gRPC, and WebSocket body/message types
//   - item.go       — Item union and all request/folder types
//   - types.go      — Shared primitives: Variable, Description, Settings, Scripts…
//   - io.go         — Parse / Marshal / Open / Write
//   - validate.go   — JSON-schema validation
//   - dir.go        — Unbundled directory layout (bundled: false)
//   - builder.go    — Fluent builder API
package opencollection
