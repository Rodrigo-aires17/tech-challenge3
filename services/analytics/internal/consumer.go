package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Consumer faz o long-polling da fila SQS e grava cada evento no Store.
// Autentica-se via IRSA (role associada ao Service Account `analytics-sa`).
type Consumer struct {
	client   *sqs.Client
	queueURL string
	store    Store
}

func NewConsumer(client *sqs.Client, queueURL string, store Store) *Consumer {
	return &Consumer{client: client, queueURL: queueURL, store: store}
}

// Run bloqueia processando mensagens até que o contexto seja cancelado.
func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("error receiving messages: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range out.Messages {
			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg types.Message) {
	var event Event
	if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &event); err != nil {
		log.Printf("discarding malformed message: %v", err)
		c.delete(ctx, msg)
		return
	}

	if err := c.store.Record(ctx, event); err != nil {
		log.Printf("error recording event, message will be retried: %v", err)
		return
	}

	c.delete(ctx, msg)
}

func (c *Consumer) delete(ctx context.Context, msg types.Message) {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Printf("error deleting message: %v", err)
	}
}
