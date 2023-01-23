package dynamodb

import (
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	interfaceDatabase "scraper-backend/src/driver/interface/database"
	dynamodbTable "scraper-backend/src/driver/database/dynamodb/table"
)

func ConstructorPicture(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKey string,
	SortKey string,
) interfaceDatabase.DriverDynamodbPicture {
	return &dynamodbTable.TablePicture{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKey:     PrimaryKey,
		SortKey:        SortKey,
	}
}

func ConstructorTag(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKey string,
	SortKey string,
) interfaceDatabase.DriverDynamodbTag {
	return &dynamodbTable.TableTag{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKey:     PrimaryKey,
		SortKey:        SortKey,
	}
}

func ConstructorUser(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKey string,
	SortKey string,
) interfaceDatabase.DriverDynamodbUser {
	return &dynamodbTable.TableUser{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKey:     PrimaryKey,
		SortKey:        SortKey,
	}
}
