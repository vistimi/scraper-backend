package dynamodb

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type TableUser struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
	PrimaryKey     string
	SortKey        string
}