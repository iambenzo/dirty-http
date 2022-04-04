# Dirty HTTP

<p>
    <a href="https://pkg.go.dev/github.com/iambenzo/dirtyhttp"><img src="https://pkg.go.dev/badge/github.com/iambenzo/dirtyhttp.svg" alt="Go Reference"></a>
</p>

A quick and dirty utility module for getting a HTTP application up and running:

```go
package main

import (
	"net/http"

	"github.com/iambenzo/dirtyhttp"
)

// Handler/Controller struct
type helloHandler struct{}

// Implement http.Handler
//
// Your logic goes here
func (hey helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:

        var queryParameters = r.URL.Query()

        // Make use of a pre-defined "Message" response
        response := dirtyhttp.HttpMessageResponse{}

        if queryParameters.Get("name") != "" {
            response.Message = "Hello, " + queryParameters.Get("name")
        } else {
            response.Message = "Hello, world!"
        }

        // Encode an object as JSON and send the response
        dirtyhttp.EncodeResponseAsJSON(response, w)

    default:
        // Write a timestamped log entry
        api.Logger.Error("A non-implemented method was attempted")

        // Return a pre-defined error with a custom message
        api.HttpErrorWriter.MethodNotAllowed(w, "Naughty, naughty.")
        return
    }
}

var api dirtyhttp.Api = dirtyhttp.Api{}

func main() {
    // Initialisation
    api.Init()

    // Register a handler
    hello := &helloHandler{}
    api.RegisterHandler("/hello", *hello)

    // Go, baby, go!
    api.StartService()
}

```

By default, the Service will look for values in `API_USER` and `API_PASSWORD` to use as Basic Authentication credentials.

If you want to start a service without the default requirement for Basic Auth, you can swap out the above main function for one that looks a little like this:

```go
var api dirtyhttp.Api = dirtyhttp.Api{}

func main() {
    // Initialisation with custom config
    config := dirtyhttp.EnvConfig{}
    api.InitWithConfig(&config)

    // Register a handler
    hello := &helloHandler{}
    api.RegisterHandler("/hello", *hello)

    // Go, baby, go!
    api.StartServiceNoAuth()
}
```
