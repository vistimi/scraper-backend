package adapter

import (
	"context"
	"io"
)

type DriverS3 interface {
	ItemUpload(ctx context.Context, buffer io.Reader, bucketName, path string) (*string, error)
	ItemGet(ctx context.Context, bucketName, path string) ([]byte, error)
	ItemCopy(ctx context.Context, bucketName, sourcePath, destinationPath string) (*string, error)
}
