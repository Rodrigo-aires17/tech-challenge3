package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"togglemaster/evaluation/internal"
)

func main() {
	var publisher internal.Publisher

	if queueURL := os.Getenv("SQS_QUEUE_URL"); queueURL != "" {
		sqsPublisher, err := internal.NewSQSPublisher(context.Background(), queueURL)
		if err != nil {
			log.Fatalf("failed to initialize SQS publisher: %v", err)
		}
		publisher = sqsPublisher
		log.Println("evaluation service publishing to SQS")
	} else {
		publisher = &internal.InMemoryPublisher{}
		log.Println("evaluation service using in-memory publisher (SQS_QUEUE_URL not set)")
	}

	server := internal.NewServer(publisher)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("evaluation service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}
