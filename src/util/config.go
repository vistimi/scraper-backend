package util

import (
	"fmt"
	"scraper-backend/src/driver/client"
	"scraper-backend/src/driver/database/dynamodb"
	"scraper-backend/src/driver/storage/bucket"

	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsDynamodbTable struct {
	TableName  string
	PrimaryKey string
	SortKey    string
}

type Config struct {
	AwsS3Client                *awsS3.Client
	S3BucketNamePictures       string
	AwsDynamoDbClient          *awsDynamodb.Client
	TablePictureProcessName    string
	TablePictureValidationName string
	TablePictureProductionName string
	TablePictureBlockedName    string
	TablePicturePrimaryKey     string
	TablePictureSortKey        string
	TableTagName               string
	TableTagPrimaryKey         string
	TableTagSortKey            string
	TableUserName              string
	TableUserPrimaryKey        string
	TableUserSortKey           string
}

func NewConfig() (*Config, error) {
	s3BucketNamePictures := GetEnvVariable("IMAGES_BUCKET")
	env := GetEnvVariable("CLOUD_HOST")

	var AwsS3Client *awsS3.Client
	var AwsDynamodbClient *awsDynamodb.Client

	TablePictureProcessName := GetEnvVariable("TABLE_PICTURE_PROCESS_NAME")
	TablePictureValidationName := GetEnvVariable("TABLE_PICTURE_VALIDATION_NAME")
	TablePictureProductionName := GetEnvVariable("TABLE_PICTURE_PRODUCTION_NAME")
	TablePictureBlockedName := GetEnvVariable("TABLE_PICTURE_BLOCKED_NAME")
	TablePicturePrimaryKey := GetEnvVariable("TABLE_PICTURE_PK")
	TablePictureSortKey := GetEnvVariable("TABLE_PICTURE_SK")

	TableTagName := GetEnvVariable("TABLE_TAG_NAME")
	TableTagPrimaryKey := GetEnvVariable("TABLE_TAG_PK")
	TableTagSortKey := GetEnvVariable("TABLE_TAG_SK")

	TableUserName := GetEnvVariable("TABLE_USER_NAME")
	TableUserPrimaryKey := GetEnvVariable("TABLE_USER_PK")
	TableUserSortKey := GetEnvVariable("TABLE_USER_SK")

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

		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProcessName,
			TablePicturePrimaryKey,
			TablePictureSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureValidationName,
			TablePicturePrimaryKey,
			TablePictureSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProductionName,
			TablePicturePrimaryKey,
			TablePictureSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureBlockedName,
			TablePicturePrimaryKey,
			TablePictureSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableTagName,
			TableTagPrimaryKey,
			TableTagSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableUserName,
			TableUserPrimaryKey,
			TableUserSortKey,
		)
	default:
		return nil, fmt.Errorf("env variable not valid: %s", env)
	}

	return &Config{
		AwsS3Client:                AwsS3Client,
		S3BucketNamePictures:       s3BucketNamePictures,
		AwsDynamoDbClient:          AwsDynamodbClient,
		TablePictureProcessName:    TablePictureProcessName,
		TablePictureValidationName: TablePictureValidationName,
		TablePictureProductionName: TablePictureProductionName,
		TablePictureBlockedName:    TablePictureBlockedName,
		TablePicturePrimaryKey:     TablePicturePrimaryKey,
		TablePictureSortKey:        TablePictureSortKey,
		TableTagName:               TableTagName,
		TableTagPrimaryKey:         TableTagPrimaryKey,
		TableTagSortKey:            TableTagSortKey,
		TableUserName:              TableUserName,
		TableUserPrimaryKey:        TableUserPrimaryKey,
		TableUserSortKey:           TableUserSortKey,
	}, nil
}
