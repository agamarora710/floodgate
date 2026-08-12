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

	l := limiter.NewRedisSlidingWindow(client, 5, 1*time.Second)

	for i := 1; i <= 5; i++ {
		allowed := l.Allow("1.2.3.4")
		fmt.Printf("Request %d: allowed=%v\n", i, allowed)
	}

	time.Sleep(1100 * time.Millisecond)

	for i := 1; i <= 5; i++ {
		allowed := l.Allow("1.2.3.4")
		fmt.Printf("Request %d: allowed=%v\n", i, allowed)
	}
}
