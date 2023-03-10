package model

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	authPB "github.com/krobus00/auth-service/pb/auth"
	"gorm.io/gorm"
)

const (
	DefaultPath = "DEFAULT-PATH"
)

type Storage struct {
	ID         string
	Src        *multipart.FileHeader `gorm:"-" json:"-"`
	ObjectKey  string
	FileName   string
	UploadedBy string
	IsPublic   bool
	CreatedAt  time.Time
}

func NewStorage() *Storage {
	return new(Storage)
}

func (m *Storage) SetID(id string) *Storage {
	m.ID = id
	return m
}

func (m *Storage) SetSrc(src *multipart.FileHeader) *Storage {
	m.Src = src
	return m
}

func (m *Storage) SetObjectKey(key string) *Storage {
	ran := time.Now().UnixNano()
	if key == DefaultPath {
		key = m.UploadedBy
	}
	m.ObjectKey = fmt.Sprintf("%s/%d", key, ran)
	return m
}

func (m *Storage) SetFileName(fileName string) *Storage {
	re := regexp.MustCompile(`\.`)
	m.FileName = re.ReplaceAllString(fileName, "")
	return m
}

func (m *Storage) SetUploadedBy(uploadedBy string) *Storage {
	m.UploadedBy = uploadedBy
	return m
}

func (m *Storage) SetIsPublic(isPublic bool) *Storage {
	m.IsPublic = isPublic
	return m
}

type FileUploadPayload struct {
	Src      *multipart.FileHeader
	Filename string
	Path     string
	IsPublic bool
}

type HTTPFileUploadRequest struct {
	Src      *multipart.FileHeader `form:"file"`
	Filename string                `form:"fileName"`
	IsPublic bool                  `form:"isPublic"`
}

type GetPresignURLPayload struct {
	ObjectKey string
}

type HTTPGetPresignURLRequest struct {
	ObjectKey string `query:"objecyKey"`
}

func (m *HTTPGetPresignURLRequest) ToPayload() *GetPresignURLPayload {
	return &GetPresignURLPayload{
		ObjectKey: m.ObjectKey,
	}
}

type GetPresignURLResponse struct {
	URL       string
	ExpiredAt time.Time
}

func (m *GetPresignURLResponse) ToHTTPResponse() *HTTPGetPresignURLResponse {
	expiredAt := m.ExpiredAt.UTC().Format(time.RFC3339Nano)
	return &HTTPGetPresignURLResponse{
		URL:       m.URL,
		ExpiredAt: expiredAt,
	}
}

type HTTPGetPresignURLResponse struct {
	URL       string `json:"url"`
	ExpiredAt string `json:"expiredAt"`
}

type StorageRepository interface {
	Create(ctx context.Context, data *Storage) error
	FindByID(ctx context.Context, id string) (*Storage, error)
	FindByObjectKey(ctx context.Context, objectKey string) (*Storage, error)
	GeneratePresignURL(ctx context.Context, storage *Storage) (*GetPresignURLResponse, error)

	// DI
	InjectS3Client(client *s3.Client) error
	InjectDB(db *gorm.DB) error
}

type StorageUsecase interface {
	Upload(ctx context.Context, payload *FileUploadPayload) (*Storage, error)
	GeneratePresignURL(ctx context.Context, payload *GetPresignURLPayload) (*GetPresignURLResponse, error)

	// DI
	InjectStorageRepo(repo StorageRepository) error
	InjectAuthClient(client authPB.AuthServiceClient) error
}
