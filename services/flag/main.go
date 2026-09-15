package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"togglemaster/flag/internal"
)

func main() {
	var repo internal.Repository

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		pgRepo, err := internal.NewPostgresRepository(context.Background(), internal.DatabaseURLFromEnv())
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}
		repo = pgRepo
		log.Println("flag service using PostgreSQL repository")
	} else {
		repo = internal.NewInMemoryRepository()
		log.Println("flag service using in-memory repository (DB_HOST not set)")
	}

	server := internal.NewServer(repo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("flag service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}
