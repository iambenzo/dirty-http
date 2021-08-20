package middleware

import (
	"log"
	"net/http"
)

type AuthMiddleware struct {
	User string
	Pass string
	Next http.Handler
}

func (am AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("\033[31m[ERROR] \033[0m")
	if am.Next == nil {
		am.Next = http.DefaultServeMux
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

	if username == am.User && password == am.Pass {
		am.Next.ServeHTTP(w, r)
	} else {
        w.WriteHeader(http.StatusUnauthorized)
        log.Println("Bad gateway credentials")
		return

	}

}
