package limiter

import (
	"fmt"
	"net"
	"net/http"

	"github.com/agamarora710/ddos-mitigation-go/detector"
)

func RateLimitMiddleware(l Limiter, d *detector.Detector, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		host, port, err := net.SplitHostPort(ip)

		if err != nil {
			fmt.Println("error :", err)
		}
		fmt.Println("Host :", host, "Port :", port)

		d.Record(host)

		if d.IsAnamalous(host) {
			fmt.Println("ANOMALY DETECTED FOR IP :", host)
		}
		if !l.Allow(host) {
			http.Error(w, "Too Many Requests", 429)
			return
		}

		next.ServeHTTP(w, r)
	})
}
