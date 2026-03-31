// unbundled_dir demonstrates writing a collection in the unbundled directory
// layout (one file per item) and reading it back.
package main

import (
	"fmt"
	"log"
	"os"

	opencollection "go-opencollection"
)

func main() {
	dir := "petstore-dir"

	// Build a small collection.
	c := opencollection.New("Petstore").
		Summary("Unbundled directory example").
		CollectionVersion("1.0.0").
		// Bundled must be false (the default) for WriteDir.
		Bundled(false).
		Environment(
			opencollection.NewEnvironment("Production").
				Var("baseUrl", "https://api.petstore.example.com").
				Secret("token", "string").
				Build(),
		).
		Environment(
			opencollection.NewEnvironment("Local").
				Var("baseUrl", "http://localhost:8080").
				Var("token", "dev-token").
				Build(),
		).
		AddFolder(
			opencollection.NewFolder("Pets").
				Seq(1).
				AddHttpRequest(
					opencollection.BuildHttpRequest("List Pets", "GET", "{{baseUrl}}/pets").
						Header("Accept", "application/json").
						InheritAuth().
						Build(),
				).
				AddHttpRequest(
					opencollection.BuildHttpRequest("Get Pet", "GET", "{{baseUrl}}/pets/{{petId}}").
						PathParam("petId", "1").
						InheritAuth().
						Build(),
				).
				AddHttpRequest(
					opencollection.BuildHttpRequest("Create Pet", "POST", "{{baseUrl}}/pets").
						Header("Content-Type", "application/json").
						JSONBody(`{"name":"{{name}}","status":"available"}`).
						InheritAuth().
						Build(),
				).
				Build(),
		).
		AddFolder(
			opencollection.NewFolder("Orders").
				Seq(2).
				AddHttpRequest(
					opencollection.BuildHttpRequest("Place Order", "POST", "{{baseUrl}}/store/order").
						Header("Content-Type", "application/json").
						JSONBody(`{"petId":{{petId}},"quantity":1}`).
						InheritAuth().
						Build(),
				).
				Build(),
		).
		Build()

	// Write to the unbundled directory layout.
	if err := opencollection.WriteDir(dir, c); err != nil {
		log.Fatalf("write dir: %v", err)
	}
	fmt.Printf("Written to %s/\n", dir)

	// Show what was created.
	printDir(dir, "")

	// Read back using Open (detects directory automatically).
	loaded, err := opencollection.Open(dir)
	if err != nil {
		log.Fatalf("open dir: %v", err)
	}
	fmt.Printf("\nRead back: %q with %d top-level items\n", loaded.Info.Name, len(loaded.Items))

	// Clean up.
	if err := os.RemoveAll(dir); err != nil {
		log.Fatalf("cleanup: %v", err)
	}
}

func printDir(dir, prefix string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for i, e := range entries {
		connector := "├── "
		if i == len(entries)-1 {
			connector = "└── "
		}
		fmt.Printf("%s%s%s\n", prefix, connector, e.Name())
		if e.IsDir() {
			extension := "│   "
			if i == len(entries)-1 {
				extension = "    "
			}
			printDir(dir+"/"+e.Name(), prefix+extension)
		}
	}
}
