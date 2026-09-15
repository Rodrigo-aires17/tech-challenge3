package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"togglemaster/auth/internal"
)

func main() {
	secret := os.Getenv("AUTH_JWT_SECRET")
	tokens, err := internal.NewTokenIssuer(secret, time.Hour)
	if err != nil {
		log.Fatalf("failed to initialize token issuer: %v", err)
	}

	users := internal.NewUserStore()
	// Usuário de exemplo para ambientes de desenvolvimento/demo.
	if seedUser := os.Getenv("AUTH_SEED_USERNAME"); seedUser != "" {
		if err := users.Register(seedUser, os.Getenv("AUTH_SEED_PASSWORD")); err != nil {
			log.Fatalf("failed to seed user: %v", err)
		}
	}

	server := internal.NewServer(users, tokens)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("auth service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}
