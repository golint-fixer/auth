package main

import (
	"fmt"
	"gopkg.in/vinxi/apachelog.v0"
	"gopkg.in/vinxi/auth.v0"
	"gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
	// Create a new vinxi proxy
	vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

	// Attach the apachelog middleware
	vs.Use(apachelog.Default)

	// Attach the auth middleware
	vs.Use(auth.User("foo", "bar"))

	// Target server to forward
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
