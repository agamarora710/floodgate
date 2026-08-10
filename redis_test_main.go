package main

import (
	"fmt"
	"time"

	"github.com/agamarora710/ddos-mitigation-go/limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	l := limiter.NewRedisFixedWindowLimitor(5, 2*time.Second, client)

	for i := 1; i <= 7; i++ {
		allowed := l.Allow("1.2.3.4")
		fmt.Printf("Request %d: allowed=%v\n", i, allowed)
	}

	time.Sleep(2 * time.Second)

	for i := 1; i <= 5; i++ {
		allowed := l.Allow("1.2.3.4")
		fmt.Printf("Request %d: allowed=%v\n", i, allowed)
	}

}
