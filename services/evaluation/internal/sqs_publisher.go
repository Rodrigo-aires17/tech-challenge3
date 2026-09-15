package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSPublisher publica eventos de avaliação na fila SQS provisionada via
// Terraform. Autentica-se via IRSA (role associada ao Service Account
// `evaluation-sa`), sem access keys estáticas no container.
type SQSPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSPublisher(ctx context.Context, queueURL string) (*SQSPublisher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &SQSPublisher{client: sqs.NewFromConfig(cfg), queueURL: queueURL}, nil
}

func (p *SQSPublisher) Publish(ctx context.Context, event Event) error {
	body, err := event.MarshalToJSON()
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}
