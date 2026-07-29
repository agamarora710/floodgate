package limiter

import (
	"fmt"
	"net"
	"net/http"
)

func RateLimitMiddleware(l Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		host, port, err := net.SplitHostPort(ip)

		if err != nil {
			fmt.Println("error :", err)
		}
		fmt.Println("Host :", host, "Port :", port)

		if !l.Allow(host) {
			http.Error(w, "Too Many Requests", 429)
		}

		next.ServeHTTP(w, r)
	})
}
