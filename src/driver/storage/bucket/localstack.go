package bucket

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3CreateLocalstack(s3Client *s3.Client, bucketName string) error {
	if _, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: aws.String(bucketName)}); err != nil {
		return err
	}
	return nil
}
