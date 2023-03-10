package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/util"
	"gorm.io/gorm"
)

type storageRepository struct {
	s3 *s3.Client
	db *gorm.DB
}

func NewStorageRepository() model.StorageRepository {
	return new(storageRepository)
}

func (r *storageRepository) uploadToS3(ctx context.Context, data *model.Storage) error {
	src, err := data.Src.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return err
	}

	contentType := http.DetectContentType(buf.Bytes())
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return err
	}

	if !util.Contains(config.FileTypeWhitelist(), exts[0]) {
		return errors.New("file extention not allowed")
	}
	data.FileName = fmt.Sprintf("%s%s", data.FileName, exts[0])

	data.ObjectKey = fmt.Sprintf("%s%s", data.ObjectKey, exts[0])

	bucketName := config.GetS3BucketName()
	_, err = r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &bucketName,
		Key:           &data.ObjectKey,
		ACL:           types.ObjectCannedACLPrivate,
		ContentLength: int64(buf.Len()),
		Body:          buf,
		ContentType:   aws.String(contentType),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *storageRepository) Create(ctx context.Context, data *model.Storage) error {
	err := r.uploadToS3(ctx, data)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *storageRepository) FindByID(ctx context.Context, id string) (*model.Storage, error) {
	storage := new(model.Storage)
	err := r.db.WithContext(ctx).First(storage, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return storage, nil
}
func (r *storageRepository) FindByObjectKey(ctx context.Context, objectKey string) (*model.Storage, error) {
	storage := new(model.Storage)
	err := r.db.WithContext(ctx).First(storage, "object_key = ?", objectKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return storage, nil
}

func (r *storageRepository) GeneratePresignURL(ctx context.Context, storage *model.Storage) (*model.GetPresignURLResponse, error) {
	expiration := time.Now().Add(config.GetS3SignDuration())
	bucketName := config.GetS3BucketName()
	getObjectArgs := s3.GetObjectInput{
		Bucket:          &bucketName,
		ResponseExpires: &expiration,
		Key:             &storage.ObjectKey,
	}

	res, err := s3.NewPresignClient(r.s3).PresignGetObject(ctx, &getObjectArgs)
	if err != nil {
		return nil, err
	}
	return &model.GetPresignURLResponse{
		URL:       res.URL,
		ExpiredAt: expiration,
	}, nil

}
