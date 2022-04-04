package dirtyhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type timeoutMiddleware struct {
	Options TimeoutConfig
	Next    http.Handler
}

func (tm timeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}

	if !tm.Options.Enabled {
		tm.Next.ServeHTTP(w, r)
		return
	}

	// Replace request context with a replica that also has a timeout
	ctx := r.Context()
	duration, err := time.ParseDuration(fmt.Sprintf("%d%s", tm.Options.Length, "s"))
	if err != nil {
		fmt.Printf("Couldn't parse timeout config \n %v \n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(ctx, duration)

	defer cancel()
	r.WithContext(ctx)

	// Create an empty channel to receive a done signal
	// from either the upstream handler, or the context's timeout
	ch := make(chan struct{})

	// Call upstream handler and send done signal upon completion
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()

	// Return if upstream handler completes first,
	// send timeout reqponse if context timout completes first
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}

}
