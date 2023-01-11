package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"
	controllerModel "scraper-backend/src/adapter/controller/model"

	"github.com/google/uuid"
)

type TablePicture struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKey     string
	SortKey        string
}

func (table TablePicture) ReadPicture(ctx context.Context, primaryKey, sortKey string) (*dynamodbModel.Picture, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKey: types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKey: types.AttributeValueMemberS{
				Value: sortKey,
			},
		},
	}

	response, err := table.DynamoDbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	var picture dynamodbModel.Picture
	err = attributevalue.UnmarshalMap(response.Item, &picture)
	if err != nil {
		return nil, err
	}

	return &picture, nil
}

func (table TablePicture) ReadPictures(ctx context.Context, primaryKey string) ([]dynamodbModel.Picture, error) {
	var err error
	var response *dynamodb.QueryOutput
	var pictures []dynamodbModel.Picture
	keyEx := expression.Key(table.PrimaryKey).Equal(expression.Value(primaryKey))
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

	err = attributevalue.UnmarshalListOfMaps(response.Items, &pictures)
	if err != nil {
		return nil, err
	}

	return pictures, nil
}

// attributes is all the desired attributes, e.g "att1 att2"
func (table TablePicture) ReadPicturesA(ctx context.Context, primaryKey, attributes string) ([]dynamodbModel.Picture, error) {
	var err error
	var response *dynamodb.QueryOutput
	var pictures []dynamodbModel.Picture
	keyEx := expression.Key(table.PrimaryKey).Equal(expression.Value(primaryKey))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	response, err = table.DynamoDbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(table.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      aws.String(attributes),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &pictures)
	if err != nil {
		return nil, err
	}

	return pictures, nil
}

func (table TablePicture) CreatePicture(ctx context.Context, picture controllerModel.Picture) error {
	var driverPicture  dynamodbModel.Picture
	if err := driverPicture.ToDriverModel(picture); err != nil{
		return err
	}

	item, err := attributevalue.MarshalMap(picture)
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

func (table TablePicture) DeletePicture(ctx context.Context, primaryKey, sortKey string) error {
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

func (table TablePicture) DeletePictureTag(ctx context.Context, primaryKey, sortKey string, tagID uuid.UUID) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Delete(expression.Name("tags"), expression.Value(tagID.String()))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) CreatePictureTag(ctx context.Context, primaryKey, sortKey string, tag dynamodbModel.PictureTag) error {
	// Build the update expression
	tagMap := map[uuid.UUID]dynamodbModel.PictureTag{uuid.New(): tag}
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Add(expression.Name("tags"), expression.Value(tagMap))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) UpdatePictureTag(ctx context.Context, primaryKey, sortKey string, tag map[uuid.UUID]dynamodbModel.PictureTag) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Set(expression.Name("tags"), expression.Value(tag))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) UpdatePictureSize(ctx context.Context, primaryKey, sortKey string, sizeMap map[uuid.UUID]dynamodbModel.PictureSize) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Set(expression.Name("size"), expression.Value(sizeMap))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}
