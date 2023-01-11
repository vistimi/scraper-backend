package adapter

import (
	"context"
	"io"
)

type DriverS3 interface {
	ItemCreate(ctx context.Context, buffer io.Reader, bucketName, path string) ( error)
	ItemRead(ctx context.Context, bucketName, path string) ([]byte, error)
	ItemCopy(ctx context.Context, bucketName, sourcePath, destinationPath string) ( error)
	ItemDelete(ctx context.Context, bucketName, destinationPath string) (error)
}
