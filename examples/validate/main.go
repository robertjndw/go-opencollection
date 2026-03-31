// validate parses a collection file and validates it against the embedded
// OpenCollection JSON schema, printing any validation errors.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
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
		var ve *jsonschema.ValidationError
		if errors.As(err, &ve) {
			fmt.Println("Validation failed:")
			for _, cause := range ve.Causes {
				fmt.Printf("  - %s: %s\n", cause.InstanceLocation, cause.Message)
			}
		} else {
			fmt.Printf("Validation error: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("OK: %q is a valid OpenCollection document\n", path)
}
