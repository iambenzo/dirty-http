# Dirty HTTP

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

