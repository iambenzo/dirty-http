package dirtyhttp

import (
	"log"
	"net/http"
)

type authMiddleware struct {
	Options AuthConfig
	Next    http.Handler
}

func (am authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("\033[31m[ERROR] \033[0m")
	if am.Next == nil {
		am.Next = http.DefaultServeMux
	}

	if !am.Options.Enabled {
		am.Next.ServeHTTP(w, r)
		return
	}

	if r.URL.Path == "/health" {
		am.Next.ServeHTTP(w, r)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("No gateway credentials")
		return
	}

	if username == am.Options.ApiUser && password == am.Options.ApiPassword {
		am.Next.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Bad gateway credentials")
		return

	}

}
