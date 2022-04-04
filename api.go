package dirtyhttp

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"strings"
)

// Main entrypoint for using this module
type Api struct {
	Config          *Config
	Logger          *logger
	HttpErrorWriter *httpErrorWriter
	Upstream        *upstream
}

// Useful for registering extra handlers/contollers against a path.
func (api Api) RegisterHandler(path string, handler http.Handler) {
	http.Handle(path, handler)
}

// Attempt to obtain configuration from environment variables
// before registering a basic health check endpoint at "<server>:<port>/health".
func (api *Api) Init() {

	// Avoid initialising multiple times
	if api.Config != nil {
		api.Logger.Error("Attempt to initialise dirtyhttp multiple times")
		return
	}

	// First configure logger and error response helper
	api.Logger = &logger{}
	api.HttpErrorWriter = newHttpErrorWriter(api.Logger)
	api.Upstream = newUpstream()

	// Attempt to get config from environment variables
	getConfig()
	api.Config = config

	// Register a health check endpoint
	// (useful for k8s deployments and heartbeats)
	hh := newHealthHandler(api.Logger, api.Upstream.Db)
	api.RegisterHandler("/health", *hh)
}

// Initialise the service with a programmatically supplied configuration
func (api *Api) InitWithConfig(config *Config) {

	// Avoid initialising multiple times
	if api.Config != nil {
		api.Logger.Error("Attempt to initialise dirtyhttp multiple times")
		return
	}

	api.Logger = &logger{}
	api.HttpErrorWriter = newHttpErrorWriter(api.Logger)
	api.Upstream = newUpstream()
	api.Config = config

	// Register a health check endpoint
	// (useful for k8s deployments and heartbeats)
	hh := newHealthHandler(api.Logger, api.Upstream.Db)
	api.RegisterHandler("/health", *hh)
}

// Will start a HTTP Listener on port 8080, unless configured otherwise.
//
// Will make use of a default suite of middleware: Timeout, Gzip and Basic Authentication.
func (api Api) StartService() {

	// Don't start the service if the user just wants help
	if strings.Contains(strings.Join(os.Args, " "), "-h") {
		os.Exit(0)
	}

	if api.Config == nil {
		log := logger{}
		log.Fatal("dirtyhttp needs to be <Init()>ialised")
	}

	api.Logger.Info("Listening on http://localhost" + api.Config.ApiPort)
	http.ListenAndServe(api.Config.ApiPort,
		&timeoutMiddleware{
			Options: api.Config.Timeout,
			Next: &gzipMiddleware{
				Options: api.Config.Gzip,
				Next: &authMiddleware{
					Options: api.Config.Authentication,
				},
			},
		},
	)
}

// Deprecated: Use the new config setup. See docs for info
//
// Will start a HTTP Listener on port 8080, unless configured otherwise.
//
// Will make use of a default suite of middleware: Timeout and Gzip.
func (api *Api) StartServiceNoAuth() {

	if api.Config == nil {
		log := logger{}
		log.Fatal("dirtyhttp needs to be <Init()>ialised")
	} else {
		// We should probably quickly validate the custom config
		if api.Config.ApiPort == "" {
			// Default if empty
			api.Config.ApiPort = ":8080"
		} else if !strings.Contains(api.Config.ApiPort, ":") {
			// Ensuring that the port name has the correct formatting
			api.Config.ApiPort = ":" + api.Config.ApiPort
		}
	}

	api.Logger.Info("Listening on http://localhost" + api.Config.ApiPort)
	http.ListenAndServe(api.Config.ApiPort,
		&timeoutMiddleware{
			Options: api.Config.Timeout,
			Next: &gzipMiddleware{
				Options: api.Config.Gzip,
			},
		},
	)
}

// Generic function for marshalling structs into JSON output
func (api *Api) WriteResponseAsJSON(data interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

// Generic function for marshalling structs into XML output
func (api *Api) WriteResponseAsXML(data interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/xml")
	enc := xml.NewEncoder(w)
	enc.Encode(data)
}

// Generic HTTP response struct for passing messages back to users
type HttpMessageResponse struct {
	Message string `json:"message"`
}

// Deprecated: Use Api.WriteResponseAsJSON instead
//
// Generic function for marshalling structs into JSON output
func EncodeResponseAsJSON(data interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

// Deprecated: Use Api.WriteResponseAsXML instead
//
// Generic function for marshalling structs into XML output
func EncodeResponseAsXML(data interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/xml")
	enc := xml.NewEncoder(w)
	enc.Encode(data)
}
