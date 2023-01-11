package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"
)

type TableTag struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKey     string
	SortKey        string
}

func (table TableTag) CreateTag(ctx context.Context, tag dynamodbModel.Tag) error {
	item, err := attributevalue.MarshalMap(tag)
	if err != nil {
		return err
	}
	_, err = table.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(table.TableName),
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (table TableTag) DeleteTag(ctx context.Context, primaryKey, sortKey string) error {
	_, err := table.DynamoDbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKey: types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKey: types.AttributeValueMemberS{
				Value: sortKey,
			},
		},
	})
	if err != nil {
		return err
	}
	return err
}

// TODO: read many tags
