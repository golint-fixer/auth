# auth [![Build Status](https://travis-ci.org/vinxi/auth.png)](https://travis-ci.org/vinxi/auth) [![GoDoc](https://godoc.org/github.com/vinxi/auth?status.svg)](https://godoc.org/github.com/vinxi/auth) [![Coverage Status](https://coveralls.io/repos/github/vinxi/auth/badge.svg?branch=master)](https://coveralls.io/github/vinxi/auth?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/auth)](https://goreportcard.com/report/github.com/vinxi/auth)

Simple HTTP authentication middleware, supporting Basic, Bearer, token and other authentication schemes.

## Installation

```bash
go get -u gopkg.in/vinxi/auth.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/auth) reference.

## Examples

#### Unique basic user auth

```go
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
  
  // Attach the auth middleware 
  vs.Use(auth.User("foo", "pas$w0rd"))
  
  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

#### Custom config allowing multiple credentials types

```go
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
```

## License

MIT
