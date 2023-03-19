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
	"gorm.io/gorm"
)

type objectRepository struct {
	s3 *s3.Client
	db *gorm.DB
}

func NewObjectRepository() model.ObjectRepository {
	return new(objectRepository)
}

func (r *objectRepository) uploadToS3(ctx context.Context, data *model.ObjectPayload) error {
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
		return err
	}
	return nil
}

func (r *objectRepository) Create(ctx context.Context, data *model.ObjectPayload) error {
	err := r.uploadToS3(ctx, data)
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Create(data.Object).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *objectRepository) FindByID(ctx context.Context, id string) (*model.Object, error) {
	object := new(model.Object)
	err := r.db.WithContext(ctx).First(object, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return object, nil
}

func (r *objectRepository) GeneratePresignedURL(ctx context.Context, object *model.Object) (*model.GetPresignedURLResponse, error) {
	expiration := time.Now().Add(config.GetS3SignDuration())
	bucketName := config.GetS3BucketName()
	getObjectArgs := s3.GetObjectInput{
		Bucket:          &bucketName,
		ResponseExpires: &expiration,
		Key:             &object.Key,
	}

	res, err := s3.NewPresignClient(r.s3).PresignGetObject(ctx, &getObjectArgs)
	if err != nil {
		return nil, err
	}
	return &model.GetPresignedURLResponse{
		ID:         object.ID,
		Filename:   object.FileName,
		URL:        res.URL,
		ExpiredAt:  expiration,
		IsPublic:   object.IsPublic,
		UploadedBy: object.UploadedBy,
		CreatedAt:  object.CreatedAt,
	}, nil

}
