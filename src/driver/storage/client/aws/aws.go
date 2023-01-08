package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newConfig() (aws.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return config.LoadDefaultConfig(ctx)
}

func S3Client() (*s3.Client, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}
