package dynamodb

import (
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	table "scraper-backend/src/driver/database/dynamodb/table"
	interfaceDatabase "scraper-backend/src/driver/interface/database"
)

func ConstructorPicture(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKey string,
	SortKey string,
) interfaceDatabase.DriverDynamodbPicture {
	return &table.TablePicture{
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
	return &table.TableTag{
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
	return &table.TableUser{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKey:     PrimaryKey,
		SortKey:        SortKey,
	}
}
