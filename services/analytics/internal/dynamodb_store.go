package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBStore persiste os eventos agregados na tabela ToggleMasterAnalytics
// (chave de partição flag_id, chave de ordenação event_timestamp).
type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStore(client *dynamodb.Client, tableName string) *DynamoDBStore {
	return &DynamoDBStore{client: client, tableName: tableName}
}

func (s *DynamoDBStore) Record(ctx context.Context, event Event) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item: map[string]types.AttributeValue{
			"flag_id":         &types.AttributeValueMemberS{Value: event.FlagID},
			"event_timestamp": &types.AttributeValueMemberS{Value: event.Timestamp.Format("2006-01-02T15:04:05.000000000Z07:00")},
			"user_id":         &types.AttributeValueMemberS{Value: event.UserID},
			"matched":         &types.AttributeValueMemberBOOL{Value: event.Matched},
		},
	})
	return err
}

func (s *DynamoDBStore) Counts(ctx context.Context, flagID string) (int, int, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		KeyConditionExpression: aws.String("flag_id = :flag_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":flag_id": &types.AttributeValueMemberS{Value: flagID},
		},
	})
	if err != nil {
		return 0, 0, err
	}

	matched, total := 0, 0
	for _, item := range out.Items {
		total++
		if b, ok := item["matched"].(*types.AttributeValueMemberBOOL); ok && b.Value {
			matched++
		}
	}
	return matched, total, nil
}
