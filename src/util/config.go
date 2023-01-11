package util

import (
	"fmt"
	"scraper-backend/src/driver/client"
	dynamodbTable "scraper-backend/src/driver/database/dynamodb/table"
	"scraper-backend/src/driver/database/dynamodb"
	"scraper-backend/src/driver/storage/bucket"

	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	AwsS3Client          *awsS3.Client
	S3BucketNamePictures string
	TablePicture         dynamodbTable.TablePicture
	TableTag             dynamodbTable.TableTag
	TableUser            dynamodbTable.TableTag
}

func NewConfig() (*Config, error) {
	s3BucketNamePictures := GetEnvVariable("IMAGES_BUCKET")
	env := GetEnvVariable("CLOUD_HOST")

	var AwsS3Client *awsS3.Client
	var AwsDynamodbClient *awsDynamodb.Client
	TablePicture := dynamodbTable.TablePicture{
		DynamoDbClient: nil,
		TableName:      GetEnvVariable("TABLE_PICTURE_NAME"),
		PrimaryKey:     GetEnvVariable("TABLE_PICTURE_PK"),
		SortKey:        GetEnvVariable("TABLE_PICTURE_SK"),
	}
	TableTag := dynamodbTable.TableTag{
		DynamoDbClient: nil,
		TableName:      GetEnvVariable("TABLE_TAG_NAME"),
		PrimaryKey:     GetEnvVariable("TABLE_TAG_PK"),
		SortKey:        GetEnvVariable("TABLE_TAG_SK"),
	}
	TableUser := dynamodbTable.TableTag{
		DynamoDbClient: nil,
		TableName:      GetEnvVariable("TABLE_USER_NAME"),
		PrimaryKey:     GetEnvVariable("TABLE_USER_PK"),
		SortKey:        GetEnvVariable("TABLE_USER_SK"),
	}

	switch env {
	case "aws":
		awsConfig, err := client.NewConfigAws()
		if err != nil {
			return nil, err
		}

		AwsS3Client = bucket.S3Client(awsConfig)
		AwsDynamodbClient = dynamodb.DynamodbClient(awsConfig)
	case "localstack":
		urlLocalstack := GetEnvVariable("LOCALSTACK_URI")

		awsConfig, err := client.NewConfigLocalstack(urlLocalstack)
		if err != nil {
			return nil, err
		}

		AwsS3Client = bucket.S3ClientPathStyle(awsConfig)
		AwsDynamodbClient = dynamodb.DynamodbClient(awsConfig)

		if err = bucket.S3CreateLocalstack(AwsS3Client, s3BucketNamePictures); err != nil {
			return nil, err
		}

		client.DynamodbCreateTableStandardPkSk(AwsDynamodbClient, TablePicture.TableName, TablePicture.PrimaryKey, TablePicture.SortKey)
		client.DynamodbCreateTableStandardPkSk(AwsDynamodbClient, TableTag.TableName, TableTag.PrimaryKey, TableTag.SortKey)
		client.DynamodbCreateTableStandardPkSk(AwsDynamodbClient, TableUser.TableName, TableUser.PrimaryKey, TableUser.SortKey)
	default:
		return nil, fmt.Errorf("env variable not valid: %s", env)
	}

	TablePicture.DynamoDbClient = AwsDynamodbClient
	TableTag.DynamoDbClient = AwsDynamodbClient
	TableUser.DynamoDbClient = AwsDynamodbClient

	return &Config{
		AwsS3Client:          AwsS3Client,
		S3BucketNamePictures: s3BucketNamePictures,
		TablePicture:         TablePicture,
		TableTag:             TableTag,
		TableUser:            TableUser,
	}, nil
}
