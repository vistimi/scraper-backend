package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func LocalS3() *s3.Client{
	
	awsEndpoint := GetEnvVariable("LOCALSTACK_URI")
	awsRegion := "us-east-1"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	// Create the resource client
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	bucketName := GetEnvVariable("IMAGES_BUCKET")

	_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	return client
}

func AwsS3() *s3.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Could not get config %v", err)
	}
	return s3.NewFromConfig(cfg)
}

func UploadItemS3(s3Client *s3.Client, buffer io.Reader, path string) (*manager.UploadOutput, error) {
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

func CopyItemS3(s3Client *s3.Client, sourcePath string, destinationPath string) (*string, error) {
	sourceUrl := filepath.Join(GetEnvVariable("IMAGES_BUCKET"), sourcePath)
    output, err := s3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket: aws.String(GetEnvVariable("IMAGES_BUCKET")),
        CopySource: aws.String(sourceUrl), 
		Key: aws.String(destinationPath),
	})
    if err != nil {
        return nil, fmt.Errorf("CopyObject has failed: %v", err)
    }
	return output.CopyObjectResult.ETag, nil
}
