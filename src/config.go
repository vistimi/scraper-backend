package main

import (
	"fmt"
	"scraper-backend/src/driver/cloud/client/aws"
	"scraper-backend/src/driver/cloud/client/localstack"
	"scraper-backend/src/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	AwsS3Client          *s3.Client
	S3BucketNamePictures string
}

func NewConfig() (*Config, error) {
	s3BucketNamePictures, err := util.GetEnvVariable("IMAGES_BUCKET")
	if err != nil {
		return nil, err
	}

	urlLocalstack, err := util.GetEnvVariable("LOCALSTACK_URI")
	if err != nil {
		return nil, err
	}

	env, err := util.GetEnvVariable("CLOUD_HOST")
	if err != nil {
		return nil, err
	}
	var AwsS3Client *s3.Client
	switch env {
	case "aws":
		if AwsS3Client, err = aws.S3Client(); err != nil {
			return nil, err
		}
	case "localstack":
		if AwsS3Client, err = localstack.S3Client(urlLocalstack); err != nil {
			return nil, err
		}
		if err = localstack.S3Create(AwsS3Client, s3BucketNamePictures); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("env variable not valid: %s", env)
	}

	// urlDatabase, err := util.GetEnvVariable("LOCALSTACK_URI")
	// if err != nil {
	// 	return nil, err
	// }

	return &Config{
		AwsS3Client:          AwsS3Client,
		S3BucketNamePictures: s3BucketNamePictures,
	}, nil
}
