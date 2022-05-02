package dirtyhttp

import (
	"fmt"
	"net/http"
	"strings"
)

type corsMiddleware struct {
	Options         *CorsConfig
	Logger          *logger
	HttpErrorWriter *httpErrorWriter
	Next            http.Handler
}

func (cm corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !cm.Options.Enabled {
		cm.Next.ServeHTTP(w, r)
		return
	}

	origin := r.Header.Get("Origin")
	allowedMethods := make([]string, 0)

	for k := range cm.Options.CorsAllowedMethods {
		allowedMethods = append(allowedMethods, k)
	}

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

		return
	}

	if !cm.testOrigin(origin) {
		cm.Logger.Error(fmt.Sprintf("A request from %s has been blocked", origin))
		cm.HttpErrorWriter.Unauthorised(w, "Requests from your domain aren't allowed")
		return
	}

	if !cm.testMethod(r.Method) {
		cm.Logger.Error(fmt.Sprintf("A %s request from %s has been blocked", r.Method, origin))
		cm.HttpErrorWriter.MethodNotAllowed(w, "Method not allowed")
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	cm.Next.ServeHTTP(w, r)
}

// Make sure that we're allowed to accept requests from the requestor
func (cm corsMiddleware) testOrigin(o string) bool {
	if _, ok := cm.Options.CorsAllowedOrigins["*"]; ok {
		return true
	}

	if _, ok := cm.Options.CorsAllowedOrigins[o]; ok {
		return true
	}

	return false
}

// make sure the request uses an acceptable method
func (cm corsMiddleware) testMethod(m string) bool {
	if _, ok := cm.Options.CorsAllowedMethods[m]; ok {
		return true
	}

	return false
}
