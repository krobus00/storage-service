package infrastructure

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/model"
)

type s3Client struct {
	client *s3.Client
}

func NewS3Client() (model.S3Client, error) {
	client := s3.NewFromConfig(aws.Config{
		Credentials: config.GetS3Credential(),
		Region:      config.GetS3Region(),
	}, s3.WithEndpointResolver(
		s3.EndpointResolverFromURL(config.GetS3Endpoint()),
	), func(opts *s3.Options) {
		opts.UsePathStyle = true
	})

	if client == nil {
		return nil, errors.New("creating an S3 SDK client failed")
	}

	return &s3Client{
		client: client,
	}, nil
}

func (i *s3Client) PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return i.client.PutObject(ctx, params)
}

func (i *s3Client) PresignGetObject(ctx context.Context, params *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return s3.NewPresignClient(i.client).PresignGetObject(ctx, params)
}
