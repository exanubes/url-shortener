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

func (q *queries) GetLink(ctx context.Context, id string) (*internal.LinkRow, error) {
	primary_key := internal.CreateLinkMetaPartitionKey(id)
	marshalled_key, err := attributevalue.MarshalMap(primary_key)

	if err != nil {
		return nil, err
	}

	result, err := q.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              &q.table_name,
		KeyConditionExpression: aws.String("#pk = :pk and #sk = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": marshalled_key["PK"],
			":sk": marshalled_key["SK"],
		},
		ExpressionAttributeNames: map[string]string{
			"#pk": "PK",
			"#sk": "SK",
		},
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
func (q *queries) LogLinkVisit(ctx context.Context, input internal.LogLinkVisitParams) error {
	visit_partition := internal.CreateLinkVisitPartitionKey(input.Shortcode, input.VisitedAt)
	row := internal.LogLinkVisitRow{
		PK:        visit_partition.PK,
		SK:        visit_partition.SK,
		VisitedAt: input.VisitedAt,
		IpAddress: input.IpAddress,
	}

	visit_row, err := attributevalue.MarshalMap(row)

	if err != nil {
		return err
	}

	bucket_updates, err := q.create_bucket_updates(input.Shortcode, input.VisitedAt)

	if err != nil {
		return err
	}

	visit_events := []types.TransactWriteItem{{
		Put: &types.Put{
			TableName: &q.table_name,
			Item:      visit_row,
		}},
	}

	_, err = q.db.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: append(visit_events, bucket_updates...),
	})

	return err

}

func (q *queries) CreateLink(ctx context.Context, input internal.LinkRow) error {
	primary_key := internal.CreateLinkMetaPartitionKey(input.Shortcode)
	input.PK = primary_key.PK
	input.SK = primary_key.SK

	item, err := attributevalue.MarshalMap(input)
	if err != nil {
		return err
	}

	if input.ConsumedAt.IsZero() {
		delete(item, "consumed_at")
	}

	_, err = q.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &q.table_name,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})

	return err
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

func (q *queries) create_bucket_updates(shortcode string, visited_at time.Time) ([]types.TransactWriteItem, error) {
	buckets := internal.CreateLinkVisitBucketPartitionKeys(shortcode, visited_at)

	bucket_updates := make([]types.TransactWriteItem, len(buckets))

	for index, bucket := range buckets {
		bucket_key, err := attributevalue.MarshalMap(bucket)

		if err != nil {
			return nil, err
		}
		bucket_updates[index] = types.TransactWriteItem{
			Update: &types.Update{
				TableName:        &q.table_name,
				Key:              bucket_key,
				UpdateExpression: aws.String("SET #c = if_not_exists(#c, :zero) + :inc"),
				ExpressionAttributeNames: map[string]string{
					"#c": "count",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":zero": &types.AttributeValueMemberN{Value: "0"},
					":inc":  &types.AttributeValueMemberN{Value: "1"},
				},
			},
		}
	}

	return bucket_updates, nil
}
