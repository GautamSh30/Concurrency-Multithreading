package main

import (
	"fmt"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("Handler started")

	select {
	case <-time.After(10 * time.Second):
		fmt.Fprintln(w, "Response from slow handler")
	case <-ctx.Done():
		fmt.Println("Handler cancelled:", ctx.Err())
		return
	}
}

func main() {
	mux := http.NewServeMux()

	// Wrap handler with a 2-second timeout
	timeoutHandler := http.TimeoutHandler(
		http.HandlerFunc(slowHandler),
		2*time.Second,
		"Service Unavailable\n",
	)
	mux.Handle("/slow", timeoutHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Server starting on :8080")
	fmt.Println("Try: curl http://localhost:8080/slow")
	fmt.Println("(will timeout after 2 seconds)")
	srv.ListenAndServe()
}
