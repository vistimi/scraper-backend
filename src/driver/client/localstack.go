package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func NewConfigLocalstack(url string) (aws.Config, error) {
	awsRegion := "us-east-1"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if url != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
}

func convertKeyType(keyType string) types.ScalarAttributeType {
	switch keyType {
	case "S":
		return types.ScalarAttributeTypeS
	case "N":
		return types.ScalarAttributeTypeN
	case "B":
		return types.ScalarAttributeTypeB
	}
	return ""
}

// not global table with primary and secondary keys
func DynamodbCreateTableStandardPkSk(client *dynamodb.Client, tableName, primaryKeyName, primaryKeyType, sortKeyName, sortKeyType string) error {

	primaryKeyAttributeType := convertKeyType(primaryKeyType)
	sortKeyAttributeType := convertKeyType(sortKeyType)
	if primaryKeyAttributeType == "" || sortKeyAttributeType == "" {
		return fmt.Errorf("invalid key type: %s, %s", primaryKeyType, sortKeyType)
	}

	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err == nil {
		return nil
	}

	if _, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(primaryKeyName),
				AttributeType: primaryKeyAttributeType,
			},
			{
				AttributeName: aws.String(sortKeyName),
				AttributeType: sortKeyAttributeType,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(primaryKeyName),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(sortKeyName),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}); err != nil {
		return err
	}

	return nil
}
