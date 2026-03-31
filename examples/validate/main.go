// validate parses a collection file and validates it against the embedded
// OpenCollection JSON schema, printing any validation errors.
package main

import (
	"fmt"
	"log"
	"os"

	opencollection "go-opencollection"
)

func main() {
	path := "petstore.yml"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	c, err := opencollection.Open(path)
	if err != nil {
		log.Fatalf("open: %v", err)
	}

	if err := opencollection.Validate(c); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OK: %q is a valid OpenCollection document\n", path)
}
