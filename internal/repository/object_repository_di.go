package repository

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

func (r *objectRepository) InjectS3Client(client *s3.Client) error {
	if client == nil {
		return errors.New("invalid s3 client")
	}
	r.s3 = client
	return nil
}

func (r *objectRepository) InjectDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("invalid db")
	}
	r.db = db
	return nil
}
