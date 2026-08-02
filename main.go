package main

import (
	"log"
	"net/http"
	"time"

	"github.com/agamarora710/ddos-mitigation-go/limiter"
)

func main() {
	l := limiter.NewSlidingWindowLimitor(5, 1*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	wrapped := limiter.RateLimitMiddleware(l, mux)

	log.Fatal(http.ListenAndServe(":8080", wrapped))
}
