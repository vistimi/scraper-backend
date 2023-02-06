package dynamodb

import (
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	table "scraper-backend/src/driver/database/dynamodb/table"
	interfaceDatabase "scraper-backend/src/driver/interface/database"
)

func ConstructorPicture(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKeyName string,
	PrimaryKeyType string,
	SortKeyName string,
	SortKeyType string,
) interfaceDatabase.DriverDynamodbPicture {
	return &table.TablePicture{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKeyName: PrimaryKeyName,
		PrimaryKeyType: PrimaryKeyType,
		SortKeyName:    SortKeyName,
		SortKeyType:    SortKeyType,
	}
}

func ConstructorTag(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKeyName string,
	PrimaryKeyType string,
	SortKeyName string,
	SortKeyType string,
) interfaceDatabase.DriverDynamodbTag {
	return &table.TableTag{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKeyName: PrimaryKeyName,
		PrimaryKeyType: PrimaryKeyType,
		SortKeyName:    SortKeyName,
		SortKeyType:    SortKeyType,
	}
}

func ConstructorUser(
	client *awsDynamodb.Client,
	TableName string,
	PrimaryKeyName string,
	PrimaryKeyType string,
	SortKeyName string,
	SortKeyType string,
) interfaceDatabase.DriverDynamodbUser {
	return &table.TableUser{
		DynamoDbClient: client,
		TableName:      TableName,
		PrimaryKeyName: PrimaryKeyName,
		PrimaryKeyType: PrimaryKeyType,
		SortKeyName:    SortKeyName,
		SortKeyType:    SortKeyType,
	}
}
