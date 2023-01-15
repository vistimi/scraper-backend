package bucket

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	interfaceStorage "scraper-backend/src/driver/interface/storage"
)

type S3 struct {
	Client *s3.Client
}

func Constructor(client *s3.Client) interfaceStorage.DriverS3 {
	return &S3{
		Client: client,
	}
}

func (s *S3) ItemCreate(ctx context.Context, buffer io.Reader, bucketName, path string) error {
	uploader := manager.NewUploader(s.Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
		Body:   buffer,
	})
	if err != nil {
		return nil
	}
	return nil
}

func (s *S3) ItemRead(ctx context.Context, bucketName, path string) ([]byte, error) {
	res, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *S3) ItemCopy(ctx context.Context, bucketName, sourcePath, destinationPath string) error {
	sourceUrl := filepath.Join(bucketName, sourcePath)
	_, err := s.Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(sourceUrl),
		Key:        aws.String(destinationPath),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3) ItemDelete(ctx context.Context, bucketName, destinationPath string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationPath),
	})
	if err != nil {
		return err
	}
	return nil
}
