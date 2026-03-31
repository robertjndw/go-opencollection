// build_and_write demonstrates building a collection with HTTP, GraphQL, and
// gRPC items using the fluent builder API, then writing it to disk.
package main

import (
	"fmt"
	"log"

	opencollection "github.com/robertjndw/go-opencollection"
)

func main() {
	// --- HTTP request ---
	listPets := opencollection.BuildHttpRequest("List Pets", "GET", "{{baseUrl}}/pets").
		Description("Returns all pets in the store").
		Tag("pets").
		Header("Accept", "application/json").
		QueryParam("limit", "20").
		QueryParam("status", "available").
		InheritAuth().
		Assert("response.status", "==", "200").
		Tests(`pm.test("has pets array", () => {
    const body = pm.response.json();
    pm.expect(body).to.be.an("array");
})`).
		Build()

	createPet := opencollection.BuildHttpRequest("Create Pet", "POST", "{{baseUrl}}/pets").
		Description("Adds a new pet to the store").
		Tag("pets").
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		JSONBody(`{"name":"{{petName}}","status":"available"}`).
		InheritAuth().
		Var("petName", "Fido").
		Assert("response.status", "==", "201").
		Build()

	// --- GraphQL request ---
	searchPets := opencollection.BuildGraphQLRequest("Search Pets", "{{graphqlUrl}}").
		Description("Full-text search via GraphQL").
		Tag("graphql").
		Header("Accept", "application/json").
		Query(
			`query SearchPets($query: String!, $limit: Int) {
  pets(search: $query, limit: $limit) { id name status }
}`,
			`{"query":"dog","limit":10}`,
		).
		BearerAuth("{{token}}").
		Assert("response.status", "==", "200").
		Build()

	// --- gRPC request ---
	getPet := opencollection.BuildGrpcRequest("Get Pet", "grpc.petstore.example.com:443", "PetService/GetPet").
		Description("Fetch a single pet by ID").
		MethodType("unary").
		ProtoFile("./proto/petstore.proto").
		Metadata("authorization", "Bearer {{token}}").
		Message(`{"id":"{{petId}}"}`).
		Build()

	// --- Folder: organise HTTP requests ---
	petsFolder := opencollection.NewFolder("Pets").
		Description("Pet management endpoints").
		Seq(1).
		Tag("pets").
		DefaultRequest(
			opencollection.NewRequestDefaults().
				Header("X-Namespace", "pets").
				Build(),
		).
		AddHttpRequest(listPets).
		AddHttpRequest(createPet).
		Build()

	// --- Environments ---
	prod := opencollection.NewEnvironment("Production").
		Color("#d32f2f").
		Var("baseUrl", "https://api.petstore.example.com").
		Var("graphqlUrl", "https://api.petstore.example.com/graphql").
		Secret("token", "string").
		Build()

	staging := opencollection.NewEnvironment("Staging").
		Color("#f57c00").
		Extends("Production").
		Var("baseUrl", "https://staging.petstore.example.com").
		Var("graphqlUrl", "https://staging.petstore.example.com/graphql").
		Build()

	// --- Collection ---
	c := opencollection.New("Petstore API").
		Summary("OpenAPI Petstore example").
		CollectionVersion("1.0.0").
		Author("Alice", "alice@example.com", "").
		Bundled(true).
		Docs("Full API docs at https://petstore.example.com/docs").
		Extension("x-team", "platform").
		Environment(prod).
		Environment(staging).
		DefaultRequest(
			opencollection.NewRequestDefaults().
				BearerAuth("{{token}}").
				Build(),
		).
		AddFolder(petsFolder).
		AddGraphQLRequest(searchPets).
		AddGrpcRequest(getPet).
		Build()

	// --- Validate before writing ---
	if err := opencollection.Validate(c); err != nil {
		log.Fatalf("validation failed: %v", err)
	}

	// --- Write to a single bundled YAML file ---
	if err := opencollection.WriteFile("petstore.yml", c); err != nil {
		log.Fatalf("write failed: %v", err)
	}

	fmt.Println("Written to petstore.yml")
}
