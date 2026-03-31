// environments demonstrates adding multiple environments with plain variables,
// secret declarations, environment inheritance, and .env file references.
package main

import (
	"fmt"
	"log"

	opencollection "github.com/robertjndw/go-opencollection"
)

func main() {
	// Base environment — defines the variable schema.
	base := opencollection.NewEnvironment("Base").
		Color("#607d8b").
		Var("apiVersion", "v1").
		Var("timeout", "30000").
		Secret("token", "string").
		Secret("dbPassword", "string").
		Build()

	// Production extends Base, overrides baseUrl, keeps secrets.
	prod := opencollection.NewEnvironment("Production").
		Color("#d32f2f").
		Extends("Base").
		Var("baseUrl", "https://api.example.com").
		Var("logLevel", "error").
		Build()

	// Staging extends Base, uses a .env file for the rest.
	staging := opencollection.NewEnvironment("Staging").
		Color("#f57c00").
		Extends("Base").
		Var("baseUrl", "https://staging.api.example.com").
		Var("logLevel", "info").
		DotEnvFile(".env.staging").
		Build()

	// Local development — all values inlined, no secrets needed.
	local := opencollection.NewEnvironment("Local").
		Color("#388e3c").
		Extends("Base").
		Var("baseUrl", "http://localhost:8080").
		Var("logLevel", "debug").
		Var("token", "dev-token-insecure").
		Build()

	c := opencollection.New("My API").
		Summary("Example with multiple environments").
		Environment(base).
		Environment(prod).
		Environment(staging).
		Environment(local).
		AddHttpRequest(
			opencollection.BuildHttpRequest("Health Check", "GET", "{{baseUrl}}/{{apiVersion}}/health").
				Header("Accept", "application/json").
				InheritAuth().
				Build(),
		).
		Build()

	if err := opencollection.Validate(c); err != nil {
		log.Fatalf("validation failed: %v", err)
	}

	if err := opencollection.WriteFile("environments.yml", c); err != nil {
		log.Fatalf("write failed: %v", err)
	}

	fmt.Println("Written to environments.yml")

	// Read back and list environments.
	loaded, err := opencollection.ParseFile("environments.yml")
	if err != nil {
		log.Fatalf("read back: %v", err)
	}
	fmt.Println("\nEnvironments in collection:")
	for _, env := range loaded.Config.Environments {
		fmt.Printf("  %s", env.Name)
		if env.Extends != "" {
			fmt.Printf(" (extends %s)", env.Extends)
		}
		fmt.Printf("  [%d variables]\n", len(env.Variables))
	}
}
