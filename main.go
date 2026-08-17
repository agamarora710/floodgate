package main

import (
	"log"
	"net/http"
	"time"

	"github.com/agamarora710/ddos-mitigation-go/detector"
	"github.com/agamarora710/ddos-mitigation-go/limiter"
)

func main() {
	l := limiter.NewSlidingWindowLimitor(5, 1*time.Second)

	d := detector.NewDetector(10*time.Second, 3.0)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	wrapped := limiter.RateLimitMiddleware(l, d, mux)

	log.Fatal(http.ListenAndServe(":8080", wrapped))
}
