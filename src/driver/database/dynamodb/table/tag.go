package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	controllerModel "scraper-backend/src/adapter/controller/model"
	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"
	"scraper-backend/src/driver/model"
)

const (
	TagPrimaryKeySearched = "searched"
	TagPrimaryKeyBlocked  = "blocked"
)

type TableTag struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKeyName string // Type
	PrimaryKeyType string
	SortKeyName    string // ID
	SortKeyType    string
}

func checkTablePK(primarykey string) error {
	switch primarykey {
	case TagPrimaryKeySearched, TagPrimaryKeyBlocked:
		return nil

	default:
		return fmt.Errorf("invalid primary key")
	}
}

func (table TableTag) CreateTag(ctx context.Context, tag controllerModel.Tag) error {
	var driverTag dynamodbModel.Tag
	driverTag.DriverMarshal(tag)

	item, err := attributevalue.MarshalMapWithOptions(driverTag)
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

func (table TableTag) DeleteTag(ctx context.Context, primaryKey string, sortKey model.UUID) error {
	if err := checkTablePK(primaryKey); err != nil {
		return err
	}
	_, err := table.DynamoDbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKeyName: &types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKeyName: &types.AttributeValueMemberB{
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
	if err := checkTablePK(primaryKey); err != nil {
		return nil, err
	}
	var err error
	var response *dynamodb.QueryOutput
	var tags []dynamodbModel.Tag
	keyEx := expression.Key(table.PrimaryKeyName).Equal(expression.Value(primaryKey))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	response, err = table.DynamoDbClient.Query(ctx, &dynamodb.QueryInput{
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
	var response *dynamodb.ScanOutput
	var tags []dynamodbModel.Tag

	response, err = table.DynamoDbClient.Scan(ctx, &dynamodb.ScanInput{
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
