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
	SortKey    *string
}

type Config struct {
	AwsS3Client                       *awsS3.Client
	S3BucketNamePictures              string
	AwsDynamodbClient                 *awsDynamodb.Client
	AwsDynamodbTablePictureProcess    AwsDynamodbTable
	AwsDynamodbTablePictureValidation AwsDynamodbTable
	AwsDynamodbTablePictureProduction AwsDynamodbTable
	AwsDynamodbTablePictureBlocked    AwsDynamodbTable
	AwsDynamodbTableTag               AwsDynamodbTable
	AwsDynamodbTableUser              AwsDynamodbTable
}

func NewConfig() (*Config, error) {
	s3BucketNamePictures := GetEnvVariable("IMAGES_BUCKET")
	env := GetEnvVariable("CLOUD_HOST")

	var AwsS3Client *awsS3.Client
	var AwsDynamodbClient *awsDynamodb.Client

	TablePictureProcessName := GetEnvVariable("TablePictureProcessName")
	TablePictureProcessPrimaryKey := GetEnvVariable("TablePicturePK")
	TablePictureProcessSortKey := GetEnvVariable("TablePictureSK")

	TablePictureValidationName := GetEnvVariable("TablePictureValidationName")
	TablePictureValidationPrimaryKey := GetEnvVariable("TablePicturePK")
	TablePictureValidationSortKey := GetEnvVariable("TablePictureSK")

	TablePictureProductionName := GetEnvVariable("TablePictureProductionName")
	TablePictureProductionPrimaryKey := GetEnvVariable("TablePicturePK")
	TablePictureProductionSortKey := GetEnvVariable("TablePictureSK")

	TablePictureBlockedName := GetEnvVariable("TablePictureBlockedName")
	TablePictureBlockedPrimaryKey := GetEnvVariable("TablePicturePK")
	TablePictureBlockedSortKey := GetEnvVariable("TablePictureSK")

	TableTagName := GetEnvVariable("TableTagName")
	TableTagPrimaryKey := GetEnvVariable("TableTagPK")
	TableTagSortKey := GetEnvVariable("TableTagSK")

	TableUserName := GetEnvVariable("TableUserName")
	TableUserPrimaryKey := GetEnvVariable("TableUserPK")
	TableUserSortKey := GetEnvVariable("TableUserSK")

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
			TablePictureProcessPrimaryKey,
			TablePictureProcessSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureValidationName,
			TablePictureProcessPrimaryKey,
			TablePictureProcessSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProductionName,
			TablePictureProcessPrimaryKey,
			TablePictureProcessSortKey,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureBlockedName,
			TablePictureProcessPrimaryKey,
			TablePictureProcessSortKey,
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
		AwsS3Client:          AwsS3Client,
		S3BucketNamePictures: s3BucketNamePictures,
		AwsDynamodbClient:    AwsDynamodbClient,
		AwsDynamodbTablePictureProcess: AwsDynamodbTable{
			TableName:  TablePictureProcessName,
			PrimaryKey: TablePictureProcessPrimaryKey,
			SortKey:    &TablePictureProcessSortKey,
		},
		AwsDynamodbTablePictureValidation: AwsDynamodbTable{
			TableName:  TablePictureValidationName,
			PrimaryKey: TablePictureValidationPrimaryKey,
			SortKey:    &TablePictureValidationSortKey,
		},
		AwsDynamodbTablePictureProduction: AwsDynamodbTable{
			TableName:  TablePictureProductionName,
			PrimaryKey: TablePictureProductionPrimaryKey,
			SortKey:    &TablePictureProductionSortKey,
		},
		AwsDynamodbTablePictureBlocked: AwsDynamodbTable{
			TableName:  TablePictureBlockedName,
			PrimaryKey: TablePictureBlockedPrimaryKey,
			SortKey:    &TablePictureBlockedSortKey,
		},
		AwsDynamodbTableTag: AwsDynamodbTable{
			TableName:  TableTagName,
			PrimaryKey: TableTagPrimaryKey,
			SortKey:    &TableTagSortKey,
		},
		AwsDynamodbTableUser: AwsDynamodbTable{
			TableName:  TableUserName,
			PrimaryKey: TableUserPrimaryKey,
			SortKey:    &TableUserSortKey,
		},
	}, nil
}
