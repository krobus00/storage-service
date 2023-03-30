package repository

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/krobus00/storage-service/internal/model"
	"gorm.io/gorm"
)

func (r *objectRepository) InjectS3Client(client model.S3Client) error {
	if client == nil {
		return errors.New("invalid s3 client")
	}
	r.s3 = client
	return nil
}

func (r *objectRepository) InjectRedisClient(client *redis.Client) error {
	if client == nil {
		return errors.New("invalid redis client")
	}
	r.redisClient = client
	return nil
}

func (r *objectRepository) InjectDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("invalid db")
	}
	r.db = db
	return nil
}
