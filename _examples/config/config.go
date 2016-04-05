package main

import (
	"fmt"
	"gopkg.in/vinxi/auth.v0"
	"gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
	// Create a new vinxi proxy
	vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

	// Bind the auth middleware with custom config
	// Any of the following credentials will be authorized
	tokens := []auth.Token{
		{Type: "Basic", Value: "foo:s3cr3t"},
		{Type: "Bearer", Value: "s3cr3t"},
		{Value: "s3cr3t token"},
	}
	vs.Use(auth.New(&auth.Config{Tokens: tokens}))

	// Target server to forward
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
