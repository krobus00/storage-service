package infrastructure

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/krobus00/storage-service/internal/config"
)

func NewS3Client() (*s3.Client, error) {
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
	return client, nil
}
