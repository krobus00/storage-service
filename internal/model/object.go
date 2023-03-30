//go:generate mockgen -destination=mock/mock_object_repository.go -package=mock github.com/krobus00/storage-service/internal/model ObjectRepository
//go:generate mockgen -destination=mock/mock_object_usecase.go -package=mock github.com/krobus00/storage-service/internal/model ObjectUsecase

package model

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	authPB "github.com/krobus00/auth-service/pb/auth"
	pb "github.com/krobus00/storage-service/pb/storage"
	"gorm.io/gorm"
)

const (
	DefaultPath = "DEFAULT-PATH"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type Object struct {
	ID         string
	FileName   string
	Key        string
	UploadedBy string
	IsPublic   bool
	TypeID     string
	Type       string `gorm:"-"`
	CreatedAt  time.Time
}

func (Object) TableName() string {
	return "objects"
}

func NewObjectCacheKey(id string) string {
	return fmt.Sprintf("objects:objectID:%s", id)
}

func NewObjectPresignedURLCacheKey(id string) string {
	return fmt.Sprintf("objects:objectID:%s:presignedURL", id)
}

func GetObjectCacheKeys(id string) []string {
	return []string{
		NewObjectCacheKey(id),
		NewObjectPresignedURLCacheKey(id),
	}
}

type ObjectPayload struct {
	Src    []byte
	Object *Object
}

func (m *ObjectPayload) SetObject(object *Object) *ObjectPayload {
	m.Object = object
	return m
}

func NewObject() *Object {
	return new(Object)
}

func (m *Object) SetID(id string) *Object {
	m.ID = id
	return m
}

func (m *Object) SetTypeID(id string) *Object {
	m.TypeID = id
	return m
}

func (m *Object) SetType(name string) *Object {
	m.Type = name
	return m
}

func (m *Object) SetKey(key string) *Object {
	ran := time.Now().UnixNano()
	if key == DefaultPath {
		key = m.UploadedBy
	}
	m.Key = fmt.Sprintf("%s/%d", key, ran)
	return m
}

func (m *Object) SetFileName(fileName string) *Object {
	re := regexp.MustCompile(`\.`)
	m.FileName = re.ReplaceAllString(fileName, "")
	return m
}

func (m *Object) SetUploadedBy(uploadedBy string) *Object {
	m.UploadedBy = uploadedBy
	return m
}

func (m *Object) SetIsPublic(isPublic bool) *Object {
	m.IsPublic = isPublic
	return m
}

type HTTPFileUploadRequest struct {
	Src      *multipart.FileHeader `form:"file"`
	Type     string                `form:"type"`
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
	Type       string
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
		Type:       m.Type,
		ExpiredAt:  expiredAt,
		IsPublic:   m.IsPublic,
		UploadedBy: m.UploadedBy,
		CreatedAt:  createdAt,
	}
}

func (m *GetPresignedURLResponse) ToGRPCResponse() *pb.Object {
	expiredAt := m.ExpiredAt.UTC().Format(time.RFC3339Nano)
	createdAt := m.CreatedAt.UTC().Format(time.RFC3339Nano)
	return &pb.Object{
		Id:         m.ID,
		FileName:   m.Filename,
		Type:       m.Type,
		SignedUrl:  m.URL,
		ExpiredAt:  expiredAt,
		IsPublic:   m.IsPublic,
		UploadedBy: m.UploadedBy,
		CreatedAt:  createdAt,
	}
}

type HTTPGetPresignedURLResponse struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	ExpiredAt  string `json:"expiredAt"`
	IsPublic   bool   `json:"isPublic"`
	UploadedBy string `json:"uploadedby"`
	CreatedAt  string `json:"createdAt"`
}

type HTTPUploadObjectResponse struct {
	ID         string    `json:"id"`
	FileName   string    `json:"filename"`
	Key        string    `json:"key"`
	UploadedBy string    `json:"uploadedBy"`
	IsPublic   bool      `json:"isPublic"`
	TypeID     string    `json:"typeID"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (m *Object) ToHTTPResponse() *HTTPUploadObjectResponse {
	return &HTTPUploadObjectResponse{
		ID:         m.ID,
		FileName:   m.FileName,
		Key:        m.Key,
		UploadedBy: m.UploadedBy,
		IsPublic:   m.IsPublic,
		TypeID:     m.TypeID,
		Type:       m.Type,
		CreatedAt:  m.CreatedAt,
	}
}

type ObjectRepository interface {
	Create(ctx context.Context, data *ObjectPayload) error
	FindByID(ctx context.Context, id string) (*Object, error)
	GeneratePresignedURL(ctx context.Context, object *Object) (*GetPresignedURLResponse, error)

	// DI
	InjectS3Client(client S3Client) error
	InjectDB(db *gorm.DB) error
	InjectRedisClient(client *redis.Client) error
}

type ObjectUsecase interface {
	Upload(ctx context.Context, payload *ObjectPayload) (*Object, error)
	GeneratePresignedURL(ctx context.Context, payload *GetPresignedURLPayload) (*GetPresignedURLResponse, error)

	// DI
	InjectObjectRepo(repo ObjectRepository) error
	InjectObjectTypeRepo(repo ObjectTypeRepository) error
	InjectObjectWhitelistTypeRepo(repo ObjectWhitelistTypeRepository) error
	InjectAuthClient(client authPB.AuthServiceClient) error
}
