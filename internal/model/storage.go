package model

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	authPB "github.com/krobus00/auth-service/pb/auth"
	pb "github.com/krobus00/storage-service/pb/storage"
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

type GetPresignedURLPayload struct {
	ObjectID string
}

type HTTPGetPresignedURLRequest struct {
	ObjectID string `query:"id"`
}

func (m *HTTPGetPresignedURLRequest) ToPayload() *GetPresignedURLPayload {
	return &GetPresignedURLPayload{
		ObjectID: m.ObjectID,
	}
}

type GetPresignedURLResponse struct {
	ID         string
	Filename   string
	URL        string
	ExpiredAt  time.Time
	IsPublic   bool
	UploadedBy string
	CreatedAt  time.Time
}

func (m *GetPresignedURLResponse) ToHTTPResponse() *HTTPGetPresignedURLResponse {
	expiredAt := m.ExpiredAt.UTC().Format(time.RFC3339Nano)
	createdAt := m.CreatedAt.UTC().Format(time.RFC3339Nano)
	return &HTTPGetPresignedURLResponse{
		ID:         m.ID,
		Filename:   m.Filename,
		URL:        m.URL,
		ExpiredAt:  expiredAt,
		IsPublic:   m.IsPublic,
		UploadedBy: m.UploadedBy,
		CreatedAt:  createdAt,
	}
}

func (m *GetPresignedURLResponse) ToGRPCResponse() *pb.Storage {
	expiredAt := m.ExpiredAt.UTC().Format(time.RFC3339Nano)
	createdAt := m.CreatedAt.UTC().Format(time.RFC3339Nano)
	return &pb.Storage{
		Id:         m.ID,
		FileName:   m.Filename,
		SignedUrl:  m.URL,
		ExpiredAt:  expiredAt,
		IsPublic:   m.IsPublic,
		UploadedBy: m.UploadedBy,
		CreatedAt:  createdAt,
	}
}

type HTTPGetPresignedURLResponse struct {
	ID         string
	Filename   string
	URL        string
	ExpiredAt  string
	IsPublic   bool
	UploadedBy string
	CreatedAt  string
}

type StorageRepository interface {
	Create(ctx context.Context, data *Storage) error
	FindByID(ctx context.Context, id string) (*Storage, error)
	GeneratePresignedURL(ctx context.Context, storage *Storage) (*GetPresignedURLResponse, error)

	// DI
	InjectS3Client(client *s3.Client) error
	InjectDB(db *gorm.DB) error
}

type StorageUsecase interface {
	Upload(ctx context.Context, payload *FileUploadPayload) (*Storage, error)
	GeneratePresignedURL(ctx context.Context, payload *GetPresignedURLPayload) (*GetPresignedURLResponse, error)

	// DI
	InjectStorageRepo(repo StorageRepository) error
	InjectAuthClient(client authPB.AuthServiceClient) error
}
