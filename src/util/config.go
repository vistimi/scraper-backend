package util

import (
	"fmt"
	"path/filepath"
	"scraper-backend/config"
	"scraper-backend/src/driver/client"
	"scraper-backend/src/driver/database/dynamodb"
	"scraper-backend/src/driver/storage/bucket"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsDynamodbTable struct {
	TableName      string
	PrimaryKeyName string
	PrimaryKeyType string
	SortKeyName    *string
	SortKeyType    *string
}

type Config struct {
	Port                              int
	HealthCheckPath                   string
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
	cloudHost := GetEnvVariable("CLOUD_HOST")
	awsRegion := GetEnvVariable("AWS_REGION")
	accessKeyID := GetEnvVariable("AWS_ACCESS_KEY")
	secretAccessKey := GetEnvVariable("AWS_SECRET_KEY")

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

	port := *configYml.Port
	healthCheckPath := *configYml.HealthCheckPath

	TablePictureProcessName := commonName + "-" + *configYml.Databases["tablePictureProcess"].Name
	TablePictureProcessPrimaryKeyName := *configYml.Databases["tablePictureProcess"].PrimaryKeyName
	TablePictureProcessSortKeyName := *configYml.Databases["tablePictureProcess"].SortKeyName
	TablePictureProcessPrimaryKeyType := *configYml.Databases["tablePictureProcess"].PrimaryKeyType
	TablePictureProcessSortKeyType := *configYml.Databases["tablePictureProcess"].SortKeyType
	TablePictureValidationName := commonName + "-" + *configYml.Databases["tablePictureValidation"].Name
	TablePictureValidationPrimaryKeyName := *configYml.Databases["tablePictureValidation"].PrimaryKeyName
	TablePictureValidationSortKeyName := *configYml.Databases["tablePictureValidation"].SortKeyName
	TablePictureValidationPrimaryKeyType := *configYml.Databases["tablePictureValidation"].PrimaryKeyType
	TablePictureValidationSortKeyType := *configYml.Databases["tablePictureValidation"].SortKeyType
	TablePictureProductionName := commonName + "-" + *configYml.Databases["tablePictureProduction"].Name
	TablePictureProductionPrimaryKeyName := *configYml.Databases["tablePictureProduction"].PrimaryKeyName
	TablePictureProductionSortKeyName := *configYml.Databases["tablePictureProduction"].SortKeyName
	TablePictureProductionPrimaryKeyType := *configYml.Databases["tablePictureProduction"].PrimaryKeyType
	TablePictureProductionSortKeyType := *configYml.Databases["tablePictureProduction"].SortKeyType
	TablePictureBlockedName := commonName + "-" + *configYml.Databases["tablePictureBlocked"].Name
	TablePictureBlockedPrimaryKeyName := *configYml.Databases["tablePictureBlocked"].PrimaryKeyName
	TablePictureBlockedSortKeyName := *configYml.Databases["tablePictureBlocked"].SortKeyName
	TablePictureBlockedPrimaryKeyType := *configYml.Databases["tablePictureBlocked"].PrimaryKeyType
	TablePictureBlockedSortKeyType := *configYml.Databases["tablePictureBlocked"].SortKeyType
	TableTagName := commonName + "-" + *configYml.Databases["tableTag"].Name
	TableTagPrimaryKeyName := *configYml.Databases["tableTag"].PrimaryKeyName
	TableTagSortKeyName := *configYml.Databases["tableTag"].SortKeyName
	TableTagPrimaryKeyType := *configYml.Databases["tableTag"].PrimaryKeyType
	TableTagSortKeyType := *configYml.Databases["tableTag"].SortKeyType
	TableUserName := commonName + "-" + *configYml.Databases["tableUser"].Name
	TableUserPrimaryKeyName := *configYml.Databases["tableUser"].PrimaryKeyName
	TableUserSortKeyName := *configYml.Databases["tableUser"].SortKeyName
	TableUserPrimaryKeyType := *configYml.Databases["tableUser"].PrimaryKeyType
	TableUserSortKeyType := *configYml.Databases["tableUser"].SortKeyType

	s3BucketNamePictures := commonName + "-" + *configYml.Buckets["picture"].Name

	switch cloudHost {
	case "aws":
		optFnsRegion := awsConfig.WithRegion(awsRegion)
		optFnsCredentials := awsConfig.WithCredentialsProvider(awsCredentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))

		// sess, err := session.NewSession(&aws.Config{
		// 	Region:      aws.String("us-west-2"),
		// 	Credentials: credentials.NewStaticCredentials(conf.AWS_ACCESS_KEY_ID, conf.AWS_SECRET_ACCESS_KEY, ""),
		// })

		awsConfig, err := client.NewConfigAws(optFnsRegion, optFnsCredentials)
		if err != nil {
			return nil, err
		}

		AwsS3Client = bucket.S3Client(*awsConfig)
		AwsDynamodbClient = dynamodb.DynamodbClient(*awsConfig)
	case "localstack":
		urlLocalstack := GetEnvVariable("LOCALSTACK_URI")
		optFnsRegion := awsConfig.WithRegion(awsRegion)
		optFnsCredentials := awsConfig.WithCredentialsProvider(awsCredentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKey, SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		})
		awsConfig, err := client.NewConfigLocalstack(urlLocalstack, optFnsRegion, optFnsCredentials)
		if err != nil {
			return nil, err
		}

		AwsS3Client = bucket.S3ClientPathStyle(awsConfig)
		AwsDynamodbClient = dynamodb.DynamodbClient(awsConfig)

		if err = bucket.S3CreateLocalstack(AwsS3Client, s3BucketNamePictures); err != nil {
			return nil, err
		}

		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProcessName,
			TablePictureProcessPrimaryKeyName,
			TablePictureProcessPrimaryKeyType,
			TablePictureProcessSortKeyName,
			TablePictureProcessSortKeyType,
		); err != nil {
			return nil, err
		}
		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureValidationName,
			TablePictureValidationPrimaryKeyName,
			TablePictureValidationPrimaryKeyType,
			TablePictureValidationSortKeyName,
			TablePictureValidationSortKeyType,
		); err != nil {
			return nil, err
		}
		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureProductionName,
			TablePictureProductionPrimaryKeyName,
			TablePictureProductionPrimaryKeyType,
			TablePictureProductionSortKeyName,
			TablePictureProductionSortKeyType,
		); err != nil {
			return nil, err
		}
		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TablePictureBlockedName,
			TablePictureBlockedPrimaryKeyName,
			TablePictureBlockedPrimaryKeyType,
			TablePictureBlockedSortKeyName,
			TablePictureBlockedSortKeyType,
		); err != nil {
			return nil, err
		}
		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableTagName,
			TableTagPrimaryKeyName,
			TableTagPrimaryKeyType,
			TableTagSortKeyName,
			TableTagSortKeyType,
		); err != nil {
			return nil, err
		}
		if err := client.DynamodbCreateTableStandardPkSk(
			AwsDynamodbClient,
			TableUserName,
			TableUserPrimaryKeyName,
			TableUserPrimaryKeyType,
			TableUserSortKeyName,
			TableUserSortKeyType,
		); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("cloud host variable not valid: %s", cloudHost)
	}

	config := Config{
		Port:                 port,
		HealthCheckPath:      healthCheckPath,
		AwsS3Client:          AwsS3Client,
		S3BucketNamePictures: s3BucketNamePictures,
		AwsDynamodbClient:    AwsDynamodbClient,
		AwsDynamodbTablePictureProcess: AwsDynamodbTable{
			TableName:      TablePictureProcessName,
			PrimaryKeyName: TablePictureProcessPrimaryKeyName,
			PrimaryKeyType: TablePictureProcessPrimaryKeyType,
			SortKeyName:    &TablePictureProcessSortKeyName,
			SortKeyType:    &TablePictureProcessSortKeyType,
		},
		AwsDynamodbTablePictureValidation: AwsDynamodbTable{
			TableName:      TablePictureValidationName,
			PrimaryKeyName: TablePictureValidationPrimaryKeyName,
			PrimaryKeyType: TablePictureValidationPrimaryKeyType,
			SortKeyName:    &TablePictureValidationSortKeyName,
			SortKeyType:    &TablePictureValidationSortKeyType,
		},
		AwsDynamodbTablePictureProduction: AwsDynamodbTable{
			TableName:      TablePictureProductionName,
			PrimaryKeyName: TablePictureProductionPrimaryKeyName,
			PrimaryKeyType: TablePictureProductionPrimaryKeyType,
			SortKeyName:    &TablePictureProductionSortKeyName,
			SortKeyType:    &TablePictureProductionSortKeyType,
		},
		AwsDynamodbTablePictureBlocked: AwsDynamodbTable{
			TableName:      TablePictureBlockedName,
			PrimaryKeyName: TablePictureBlockedPrimaryKeyName,
			PrimaryKeyType: TablePictureBlockedPrimaryKeyType,
			SortKeyName:    &TablePictureBlockedSortKeyName,
			SortKeyType:    &TablePictureBlockedSortKeyType,
		},
		AwsDynamodbTableTag: AwsDynamodbTable{
			TableName:      TableTagName,
			PrimaryKeyName: TableTagPrimaryKeyName,
			PrimaryKeyType: TableTagPrimaryKeyType,
			SortKeyName:    &TableTagSortKeyName,
			SortKeyType:    &TableTagSortKeyType,
		},
		AwsDynamodbTableUser: AwsDynamodbTable{
			TableName:      TableUserName,
			PrimaryKeyName: TableUserPrimaryKeyName,
			PrimaryKeyType: TableUserPrimaryKeyType,
			SortKeyName:    &TableUserSortKeyName,
			SortKeyType:    &TableUserSortKeyType,
		},
	}

	return &config, nil
}
