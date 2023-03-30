//go:generate mockgen -destination=mock/mock_s3_client.go -package=mock github.com/krobus00/storage-service/internal/model S3Client

package model

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error)
}
