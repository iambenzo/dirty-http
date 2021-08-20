package middleware

import (
	"context"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	Next http.Handler
}

func (tm TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}

	// Replace request context with a replica that also has a timeout
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
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
