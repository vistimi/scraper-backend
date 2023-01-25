package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	controllerModel "scraper-backend/src/adapter/controller/model"
	dynamodbModel "scraper-backend/src/driver/database/dynamodb/model"

	"github.com/google/uuid"
)

type TablePicture struct {
	DynamoDbClient *awsDynamodb.Client
	TableName      string
	PrimaryKey     string // Origin
	SortKey        string // ID
}

func (table TablePicture) ReadPicture(ctx context.Context, primaryKey string, sortKey uuid.UUID) (*controllerModel.Picture, error) {
	input := &awsDynamodb.GetItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKey: &types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKey: &types.AttributeValueMemberB{
				Value: sortKey[:],
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

	return picture.DriverUnmarshal(), nil
}

func (table TablePicture) ReadPictures(ctx context.Context, projection *expression.ProjectionBuilder, filter *expression.ConditionBuilder) ([]controllerModel.Picture, error) {
	var err error
	var response *awsDynamodb.QueryOutput
	var pictures []dynamodbModel.Picture
	builder := expression.NewBuilder()

	// keyEx := expression.Key(table.PrimaryKey).Equal(expression.Value(primaryKey))
	// builder = builder.WithKeyCondition(keyEx)

	if projection != nil {
		builder = builder.WithProjection(*projection)
	}

	if filter != nil {
		builder = builder.WithFilter(*filter)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	response, err = table.DynamoDbClient.Query(ctx, &awsDynamodb.QueryInput{
		TableName:                 aws.String(table.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &pictures)
	if err != nil {
		return nil, err
	}

	var controllerPictures []controllerModel.Picture
	for _, picture := range pictures {
		pictureModel := picture.DriverUnmarshal()
		controllerPictures = append(controllerPictures, *pictureModel)
	}

	return controllerPictures, nil
}

func (table TablePicture) CreatePicture(ctx context.Context, id uuid.UUID, picture controllerModel.Picture) error {
	var driverPicture dynamodbModel.Picture
	driverPicture.DriverMarshal(picture)
	driverPicture.ID = id

	item, err := attributevalue.MarshalMap(picture)
	if err != nil {
		return err
	}
	_, err = table.DynamoDbClient.PutItem(ctx, &awsDynamodb.PutItemInput{
		TableName: aws.String(table.TableName),
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (table TablePicture) DeletePicture(ctx context.Context, primaryKey string, sortKey uuid.UUID) error {
	_, err := table.DynamoDbClient.DeleteItem(ctx, &awsDynamodb.DeleteItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			table.PrimaryKey: &types.AttributeValueMemberS{
				Value: primaryKey,
			},
			table.SortKey: &types.AttributeValueMemberB{
				Value: sortKey[:],
			},
		},
	})
	if err != nil {
		return err
	}
	return err
}

func (table TablePicture) DeletePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Delete(expression.Name("tags"), expression.Value(tagID.String()))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: &types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: &types.AttributeValueMemberB{
			Value: sortKey[:],
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &awsDynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) CreatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error {
	var driverTag dynamodbModel.PictureTag
	driverTag.DriverMarshal(tag)

	// Build the update expression
	tagMap := map[uuid.UUID]dynamodbModel.PictureTag{tagID: driverTag}
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Add(expression.Name("tags"), expression.Value(tagMap))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: &types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: &types.AttributeValueMemberB{
			Value: sortKey[:],
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &awsDynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) UpdatePictureTag(ctx context.Context, primaryKey string, sortKey uuid.UUID, tagID uuid.UUID, tag controllerModel.PictureTag) error {
	var driverTag dynamodbModel.PictureTag
	driverTag.DriverMarshal(tag)

	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Set(expression.Name("tags"), expression.Value(map[uuid.UUID]dynamodbModel.PictureTag{tagID: driverTag}))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: &types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: &types.AttributeValueMemberB{
			Value: sortKey[:],
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &awsDynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (table TablePicture) CreatePictureSize(ctx context.Context, primaryKey string, sortKey uuid.UUID, size controllerModel.PictureSize) error {
	var driverSize dynamodbModel.PictureSize
	driverSize.DriverMarshal(size)

	// Build the update expression
	sizeMap := map[uuid.UUID]dynamodbModel.PictureSize{uuid.New(): driverSize}
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Add(expression.Name("size"), expression.Value(sizeMap))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		table.PrimaryKey: &types.AttributeValueMemberS{
			Value: primaryKey,
		},
		table.SortKey: &types.AttributeValueMemberB{
			Value: sortKey[:],
		},
	}

	// Update the item in the table
	_, err = table.DynamoDbClient.UpdateItem(ctx, &awsDynamodb.UpdateItemInput{
		TableName:        aws.String(table.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}