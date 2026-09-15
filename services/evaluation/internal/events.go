package internal

import (
	"context"
	"encoding/json"
	"time"
)

// Event representa o resultado de uma avaliação de flag, publicado para o
// serviço `analytics` processar de forma assíncrona.
type Event struct {
	FlagID    string    `json:"flag_id"`
	UserID    string    `json:"user_id"`
	Matched   bool      `json:"matched"`
	Timestamp time.Time `json:"timestamp"`
}

// Publisher abstrai o envio de eventos. Em produção, SQSPublisher envia para a
// fila SQS provisionada via Terraform; em testes usamos InMemoryPublisher.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type InMemoryPublisher struct {
	Events []Event
}

func (p *InMemoryPublisher) Publish(_ context.Context, event Event) error {
	p.Events = append(p.Events, event)
	return nil
}

func (e Event) MarshalToJSON() ([]byte, error) {
	return json.Marshal(e)
}
