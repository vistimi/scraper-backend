package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	controllerModel "scraper-backend/src/adapter/controller/model"
	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"
)

type TableTag struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKey     string	// Type
	SortKey        string	// ID
}

func (table TableTag) CreateTag(ctx context.Context, tag controllerModel.Tag) error {
	var driverTag dynamodbModel.Tag
	driverTag.DriverMarshal(tag)

	item, err := attributevalue.MarshalMap(driverTag)
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

func (table TableTag) DeleteTag(ctx context.Context, primaryKey string, sortKey uuid.UUID) error {
	_, err := table.DynamoDbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKey: types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKey: types.AttributeValueMemberB{
				Value: sortKey[:],
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (table TableTag) ReadTags(ctx context.Context, primaryKey string) ([]controllerModel.Tag, error) {
	var err error
	var response *awsDynamodb.QueryOutput
	var tags []dynamodbModel.Tag
	keyEx := expression.Key(table.PrimaryKey).Equal(expression.Value(primaryKey))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	response, err = table.DynamoDbClient.Query(ctx, &awsDynamodb.QueryInput{
		TableName:                 aws.String(table.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &tags)
	if err != nil {
		return nil, err
	}

	var tagsModel []controllerModel.Tag
	for _, tag := range tags {
		tagsModel = append(tagsModel, tag.DriverUnmarshal())
	}

	return tagsModel, nil
}

func (table TableTag) ScanTags(ctx context.Context) ([]controllerModel.Tag, error) {
	var err error
	var response *awsDynamodb.ScanOutput
	var tags []dynamodbModel.Tag

	response, err = table.DynamoDbClient.Scan(ctx, &awsDynamodb.ScanInput{
		TableName: aws.String(table.TableName),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &tags)
	if err != nil {
		return nil, err
	}

	var controllerTags []controllerModel.Tag
	for _, tag := range tags {
		controllerTags = append(controllerTags, tag.DriverUnmarshal())
	}

	return controllerTags, nil
}
