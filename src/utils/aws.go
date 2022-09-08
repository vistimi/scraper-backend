package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elgohr/go-localstack"
)

func LocalS3() (aws.Config, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	l, err := localstack.NewInstance()
	if err != nil {
		log.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.StartWithContext(ctx); err != nil {
		log.Fatalf("Could not start localstack %v", err)
	}

	url := l.EndpointV2(localstack.S3)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               url,
				SigningRegion:     "us-east-1",
				HostnameImmutable: false,
			}, nil
		})),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		log.Fatalf("Could not get config %v", err)
	}
	return cfg, url
}

func AwsS3() aws.Config {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Could not get config %v", err)
	}
	return cfg
}

func ConnectS3(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func UploadItemS3(s3Client *s3.Client, buffer io.Reader, path string) (*manager.UploadOutput, error) {
	// upload in s3 the file
	uploader := manager.NewUploader(s3Client)
	return uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(GetEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
		Body:   buffer,
	})
}

func GetItemS3(s3Client *s3.Client, path string) ([]byte, error) {
	res, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(GetEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("s3Client.GetObject has failed: %v", err)
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll has failed: %v", err)
	}
	return buffer, nil
}
