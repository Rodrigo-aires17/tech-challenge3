package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"togglemaster/targeting/internal"
)

func main() {
	var cache internal.Cache

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cache = internal.NewRedisCache(redisAddr, 5*time.Minute)
		log.Println("targeting service using Redis cache")
	} else {
		cache = internal.NewInMemoryCache()
		log.Println("targeting service using in-memory cache (REDIS_ADDR not set)")
	}

	server := internal.NewServer(cache)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("targeting service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}
