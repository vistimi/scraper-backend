package aws

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	interfaceCloud "scraper-backend/src/adapter/interface/cloud"
)

type S3 struct {
	Client *s3.Client
}

func (s *S3) Constructor(client *s3.Client) interfaceCloud.S3{
return &S3{
		Client: client,
	}
}

func (s *S3) ItemUpload(ctx context.Context, buffer io.Reader, bucketName, path string) (*string, error) {
	uploader := manager.NewUploader(s.Client)
	output, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
		Body:   buffer,
	})
	if err != nil {
		return nil, fmt.Errorf("UploadItem has failed: %v", err)
	}
	return output.ETag, nil
}

func (s *S3) ItemGet(ctx context.Context, bucketName, path string) ([]byte, error) {
	res, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
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

func (s *S3) ItemCopy(ctx context.Context, bucketName, sourcePath, destinationPath string) (*string, error) {
	sourceUrl := filepath.Join(bucketName, sourcePath)
	output, err := s.Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(sourceUrl),
		Key:        aws.String(destinationPath),
	})
	if err != nil {
		return nil, fmt.Errorf("CopyObject has failed: %v", err)
	}
	return output.CopyObjectResult.ETag, nil
}
