package util

import (
	"fmt"
	"path/filepath"
	"scraper-backend/config"
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

// TODO: map of tables
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
	commonName := GetEnvVariable("COMMON_NAME")
	env := GetEnvVariable("CLOUD_HOST")

	var AwsS3Client *awsS3.Client
	var AwsDynamodbClient *awsDynamodb.Client

	path, err := filepath.Abs("config/config.yml")
	if err != nil {
		return nil, err
	}
	configYml, err := config.ReadConfigFile(path)
	if err != nil {
		return nil, err
	}

	TablePictureProcessName := commonName + *configYml.Databases["tablePictureProcess"].Name
	TablePictureProcessPrimaryKey := *configYml.Databases["tablePictureProcess"].PrimaryKeyName
	TablePictureProcessSortKey := *configYml.Databases["tablePictureProcess"].SortKeyName
	TablePictureProcessPrimaryKeyType := *configYml.Databases["tablePictureProcess"].PrimaryKeyType
	TablePictureProcessSortKeyType := *configYml.Databases["tablePictureProcess"].SortKeyType
	TablePictureValidationName := commonName + *configYml.Databases["tablePictureValidation"].Name
	TablePictureValidationPrimaryKey := *configYml.Databases["tablePictureValidation"].PrimaryKeyName
	TablePictureValidationSortKey := *configYml.Databases["tablePictureValidation"].SortKeyName
	TablePictureValidationPrimaryKeyType := *configYml.Databases["tablePictureValidation"].PrimaryKeyType
	TablePictureValidationSortKeyType := *configYml.Databases["tablePictureValidation"].SortKeyType
	TablePictureProductionName := commonName + *configYml.Databases["tablePictureProduction"].Name
	TablePictureProductionPrimaryKey := *configYml.Databases["tablePictureProduction"].PrimaryKeyName
	TablePictureProductionSortKey := *configYml.Databases["tablePictureProduction"].SortKeyName
	TablePictureProductionPrimaryKeyType := *configYml.Databases["tablePictureProduction"].PrimaryKeyType
	TablePictureProductionSortKeyType := *configYml.Databases["tablePictureProduction"].SortKeyType
	TablePictureBlockedName := commonName + *configYml.Databases["tablePictureBlocked"].Name
	TablePictureBlockedPrimaryKey := *configYml.Databases["tablePictureBlocked"].PrimaryKeyName
	TablePictureBlockedSortKey := *configYml.Databases["tablePictureBlocked"].SortKeyName
	TablePictureBlockedPrimaryKeyType := *configYml.Databases["tablePictureBlocked"].PrimaryKeyType
	TablePictureBlockedSortKeyType := *configYml.Databases["tablePictureBlocked"].SortKeyType
	TableTagName := commonName + *configYml.Databases["tableTag"].Name
	TableTagPrimaryKey := *configYml.Databases["tableTag"].PrimaryKeyName
	TableTagSortKey := *configYml.Databases["tableTag"].SortKeyName
	TableTagPrimaryKeyType := *configYml.Databases["tableTag"].PrimaryKeyType
	TableTagSortKeyType := *configYml.Databases["tableTag"].SortKeyType
	TableUserName := commonName + *configYml.Databases["tableUser"].Name
	TableUserPrimaryKey := *configYml.Databases["tableUser"].PrimaryKeyName
	TableUserSortKey := *configYml.Databases["tableUser"].SortKeyName
	TableUserPrimaryKeyType := *configYml.Databases["tableUser"].PrimaryKeyType
	TableUserSortKeyType := *configYml.Databases["tableUser"].SortKeyType

	s3BucketNamePictures := commonName + *configYml.Buckets["env"].Name

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
			TablePictureProcessPrimaryKeyType,
			TablePictureProcessSortKey,
			TablePictureProcessSortKeyType,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureValidationName,
			TablePictureValidationPrimaryKey,
			TablePictureValidationPrimaryKeyType,
			TablePictureValidationSortKey,
			TablePictureValidationSortKeyType,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProductionName,
			TablePictureProductionPrimaryKey,
			TablePictureProductionPrimaryKeyType,
			TablePictureProductionSortKey,
			TablePictureProductionSortKeyType,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureBlockedName,
			TablePictureBlockedPrimaryKey,
			TablePictureBlockedPrimaryKeyType,
			TablePictureBlockedSortKey,
			TablePictureBlockedSortKeyType,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableTagName,
			TableTagPrimaryKey,
			TableTagPrimaryKeyType,
			TableTagSortKey,
			TableTagSortKeyType,
		)
		client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableUserName,
			TableUserPrimaryKey,
			TableUserPrimaryKeyType,
			TableUserSortKey,
			TableUserSortKeyType,
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
