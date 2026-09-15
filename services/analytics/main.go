package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"togglemaster/analytics/internal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var store internal.Store

	tableName := os.Getenv("DYNAMODB_TABLE")
	queueURL := os.Getenv("SQS_QUEUE_URL")

	if tableName != "" && queueURL != "" {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("failed to load AWS config: %v", err)
		}

		store = internal.NewDynamoDBStore(dynamodb.NewFromConfig(cfg), tableName)

		consumer := internal.NewConsumer(sqs.NewFromConfig(cfg), queueURL, store)
		go func() {
			log.Println("analytics consumer started, polling SQS")
			if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
				log.Printf("consumer stopped: %v", err)
			}
		}()
	} else {
		store = internal.NewInMemoryStore()
		log.Println("analytics service using in-memory store (DYNAMODB_TABLE/SQS_QUEUE_URL not set)")
	}

	server := internal.NewServer(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{Addr: ":" + port, Handler: server}
	go func() {
		log.Printf("analytics service listening on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownCtx)
}
