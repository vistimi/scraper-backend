package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func newConfig() (aws.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return config.LoadDefaultConfig(ctx)
}

func DynamodbClient() (*dynamodb.Client, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}
