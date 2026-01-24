package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb/internal"
)

type queries struct {
	db         *dynamodb.Client
	table_name string
}

func new_queries(db *dynamodb.Client) *queries {
	return &queries{db: db, table_name: "url_shortener"}
}

var get_link_condition = "pk = :pk and sk = :sk"

func (q *queries) GetLink(ctx context.Context, id string) (*internal.LinkRow, error) {
	primary_key := internal.CreateLinkMetaPartitionKey(id)
	marshalled_key, err := attributevalue.MarshalMap(primary_key)
	if err != nil {
		return nil, err
	}

	result, err := q.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &q.table_name,
		KeyConditionExpression:    &get_link_condition,
		ExpressionAttributeValues: marshalled_key,
	})

	if err != nil {
		return nil, err
	}

	if result.Count == 0 {
		return nil, errors.New("Not found")
	}

	var output internal.LinkRow

	if err := attributevalue.UnmarshalMap(result.Items[0], &output); err != nil {
		return nil, err
	}

	return &output, err
}

// TODO:
func (q *queries) LogLinkVisit(ctx context.Context) {}

func (q *queries) CreateLink(ctx context.Context, input internal.LinkRow) error {
	primary_key := internal.CreateLinkMetaPartitionKey(input.Shortcode)
	input.PK = primary_key.PK
	input.SK = primary_key.SK

	item, err := attributevalue.MarshalMap(input)
	if err != nil {
		return err
	}

	q.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &q.table_name,
		Item:      item,
	})

	return nil
}
func (q *queries) ConsumeSingleUseLink(ctx context.Context, input internal.ConsumeSingleUseLinkParams) error {
	primary_key := internal.CreateLinkMetaPartitionKey(input.Shortcode)
	consumed_at, err := attributevalue.Marshal(input.ConsumedAt)
	if err != nil {
		return err
	}

	marshalled_key, err := attributevalue.MarshalMap(primary_key)
	if err != nil {
		return err
	}

	_, err = q.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &q.table_name,
		UpdateExpression: aws.String("set consumed_at = :consumed_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":consumed_at": consumed_at,
		},
		ConditionExpression: aws.String("attribute_not_exists(consumed_at)"),
		Key:                 marshalled_key,
	})

	return err
}
