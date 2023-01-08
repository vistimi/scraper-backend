package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/google/uuid"
)

type Image struct {
	Origin       string // PK original werbsite
	OriginID     string // id from original website
	User         User
	Extension    string      // type of file
	Name         string      // SK name <originID>.<extension>
	Size         map[uuid.UUID]ImageSize // size cropping history
	Title        string
	Description  string // decription of image
	License      string // type of public license
	CreationDate *time.Time
	Tags         map[uuid.UUID]Tag
}

type ImageSize struct {
	CreationDate *time.Time
	Box          Box // absolut reference of the top left of new box based on the original sizes
}

type Box struct {
	Tlx    *int // top left x coordinate (pointer because 0 is a possible value)
	Tly    *int // top left y coordinate (pointer because 0 is a possible value)
	Width  *int // width (pointer because 0 is a possible value)
	Height *int // height (pointer because 0 is a possible value)
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKey     string
	SortKey        string
}

func (basics TableBasics) ReadOne(ctx context.Context, primaryKey, sortKey string) (*Image, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(basics.TableName),
		Key: map[string]types.AttributeValue{
			basics.PrimaryKey: types.AttributeValueMemberS{
				Value: primaryKey,
			},
			basics.SortKey: types.AttributeValueMemberS{
				Value: sortKey,
			},
		},
	}

	response, err := basics.DynamoDbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	var image Image
	err = attributevalue.UnmarshalMap(response.Item, &image)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (basics TableBasics) ReadManyA(ctx context.Context, primaryKey, attributes string) ([]Image, error) {
	var err error
	var response *dynamodb.QueryOutput
	var images []Image
	keyEx := expression.Key(basics.PrimaryKey).Equal(expression.Value(primaryKey))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	response, err = basics.DynamoDbClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(basics.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      aws.String(attributes),
	})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (basics TableBasics) CreateOne(ctx context.Context, image Image) error {
	item, err := attributevalue.MarshalMap(image)
	if err != nil {
		return err
	}
	_, err = basics.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName),
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (basics TableBasics) DeleteOneTag(ctx context.Context, primaryKey, sortKey string, uuid uuid.UUID) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Delete(expression.Name("tags"), expression.Value(uuid.String()))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		basics.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		basics.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = basics.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(basics.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (basics TableBasics) CreateOneTag(ctx context.Context, primaryKey, sortKey string, tag Tag) error {
	// Build the update expression
	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Add(expression.Name("tags"), expression.Value(map[uuid.UUID]Tag{
			uuid.New(): tag,
		}))).
		Build()
	if err != nil {
		return err
	}

	// Set the primary key values
	key := map[string]types.AttributeValue{
		basics.PrimaryKey: types.AttributeValueMemberS{
			Value: primaryKey,
		},
		basics.SortKey: types.AttributeValueMemberS{
			Value: sortKey,
		},
	}

	// Update the item in the table
	_, err = basics.DynamoDbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(basics.TableName),
		Key:              key,
		UpdateExpression: updateExpr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}