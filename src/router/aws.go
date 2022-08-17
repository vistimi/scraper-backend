package router

import (
	"bytes"
	"context"
	"fmt"
	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func UploadS3(s3Client *s3.Client, link string, path string) (*manager.UploadOutput, error) {
	// get buffer of image
	buffer, err := GetFile(link)
	if err != nil {
		return nil, fmt.Errorf("GetFile has failed: %v", err)
	}

	// upload in s3 the file
	uploader := manager.NewUploader(s3Client)
	return uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(utils.DotEnvVariable("IMAGES_BUCKET")),
		Key:    aws.String(path),
		Body:   bytes.NewReader(buffer),
	})
}
