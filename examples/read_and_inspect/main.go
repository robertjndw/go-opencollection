// read_and_inspect demonstrates opening a collection file (or directory) and
// iterating over its items to print a summary.
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

	fmt.Printf("Collection : %s\n", c.Info.Name)
	fmt.Printf("Version    : %s\n", c.Info.Version)
	fmt.Printf("Summary    : %s\n", c.Info.Summary)
	if c.Config != nil {
		fmt.Printf("Environments:\n")
		for _, env := range c.Config.Environments {
			fmt.Printf("  - %s\n", env.Name)
		}
	}

	fmt.Printf("\nItems:\n")
	printItems(c.Items, 1)
}

func printItems(items []opencollection.Item, depth int) {
	indent := ""
	for range depth {
		indent += "  "
	}

	for _, item := range items {
		switch {
		case item.HttpRequest != nil:
			r := item.HttpRequest
			fmt.Printf("%s[http] %s  %s %s\n", indent, r.Info.Name, r.Http.Method, r.Http.URL)

		case item.GraphQLRequest != nil:
			r := item.GraphQLRequest
			fmt.Printf("%s[graphql] %s  %s\n", indent, r.Info.Name, r.GraphQL.URL)

		case item.GrpcRequest != nil:
			r := item.GrpcRequest
			fmt.Printf("%s[grpc] %s  %s / %s\n", indent, r.Info.Name, r.Grpc.URL, r.Grpc.Method)

		case item.WebSocket != nil:
			r := item.WebSocket
			fmt.Printf("%s[ws] %s  %s\n", indent, r.Info.Name, r.WebSocket.URL)

		case item.Folder != nil:
			f := item.Folder
			fmt.Printf("%s[folder] %s\n", indent, f.Info.Name)
			printItems(f.Items, depth+1)

		case item.Script != nil:
			fmt.Printf("%s[script]\n", indent)
		}
	}
}
