package main

import (
	"log"
	"net/http"

	"github.com/agamarora710/ddos-mitigation-go/limiter"
)

func main() {
	l := limiter.NewTokenBucketLimitor(10, 2)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	wrapped := limiter.RateLimitMiddleware(l, mux)

	log.Fatal(http.ListenAndServe(":8080", wrapped))
}
