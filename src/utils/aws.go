package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func ConnectS3() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}

func UploadItemS3(s3Client *s3.Client, buffer io.Reader, path string) (*manager.UploadOutput, error) {
	// upload in s3 the file
	uploader := manager.NewUploader(s3Client)
	return uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(DotEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
		Body:   buffer,
	})
}

func GetItemS3(s3Client *s3.Client, path string) ([]byte, error) {
	res, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(DotEnvVariable("IMAGES_BUCKET")),
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
