package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/goccy/go-json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-redis/redis/v8"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type objectRepository struct {
	s3          model.S3Client
	db          *gorm.DB
	redisClient *redis.Client
}

func NewObjectRepository() model.ObjectRepository {
	return new(objectRepository)
}

func (r *objectRepository) uploadToS3(ctx context.Context, data *model.ObjectPayload) error {
	logger := logrus.WithFields(logrus.Fields{
		"id":  data.Object.ID,
		"key": data.Object.Key,
	})

	buf := bytes.NewBuffer(data.Src)

	contentType := http.DetectContentType(data.Src)
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	data.Object.FileName = fmt.Sprintf("%s%s", data.Object.FileName, exts[0])
	data.Object.Key = fmt.Sprintf("%s%s", data.Object.Key, exts[0])

	bucketName := config.GetS3BucketName()
	_, err = r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &bucketName,
		Key:           &data.Object.Key,
		ACL:           types.ObjectCannedACLPrivate,
		ContentLength: int64(buf.Len()),
		Body:          buf,
		ContentType:   aws.String(contentType),
	})

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (r *objectRepository) Create(ctx context.Context, data *model.ObjectPayload) error {
	logger := logrus.WithFields(logrus.Fields{
		"id":  data.Object.ID,
		"key": data.Object.Key,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	err := r.uploadToS3(ctx, data)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = db.WithContext(ctx).Create(data.Object).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectCacheKeys(data.Object.ID))

	return nil
}

func (r *objectRepository) FindByID(ctx context.Context, id string) (*model.Object, error) {
	logger := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	object := new(model.Object)
	cacheKey := model.NewObjectCacheKey(id)

	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &object)
	if err == nil {
		return object, nil
	}

	object = new(model.Object)

	err = db.WithContext(ctx).First(object, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = SetWithExpiry(ctx, r.redisClient, cacheKey, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		logger.Error(err.Error())
		return nil, err
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, object)
	if err != nil {
		logger.Error(err.Error())
	}

	return object, nil
}

func (r *objectRepository) GeneratePresignedURL(ctx context.Context, object *model.Object) (*model.GetPresignedURLResponse, error) {
	logger := logrus.WithFields(logrus.Fields{
		"id":   object.ID,
		"key":  object.Key,
		"type": object.Type,
	})

	data := new(model.GetPresignedURLResponse)
	cacheKey := model.NewObjectPresignedURLCacheKey(object.ID)

	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &data)
	if err == nil {
		return data, nil
	}

	data = new(model.GetPresignedURLResponse)

	expiration := time.Now().Add(config.GetS3SignDuration())
	bucketName := config.GetS3BucketName()
	getObjectArgs := s3.GetObjectInput{
		Bucket:          &bucketName,
		ResponseExpires: &expiration,
		Key:             &object.Key,
	}

	res, err := r.s3.PresignGetObject(ctx, &getObjectArgs)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	data = &model.GetPresignedURLResponse{
		ID:         object.ID,
		Filename:   object.FileName,
		Type:       object.Type,
		URL:        res.URL,
		ExpiredAt:  expiration,
		IsPublic:   object.IsPublic,
		UploadedBy: object.UploadedBy,
		CreatedAt:  object.CreatedAt,
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, data)
	if err != nil {
		logger.Error(err.Error())
	}

	return data, nil
}

func (r *objectRepository) DeleteByID(ctx context.Context, id string) error {
	logger := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	object := new(model.Object)

	err := db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("id = ?", id).Delete(object).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectCacheKeys(id))

	return nil
}
