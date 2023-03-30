//go:generate mockgen -destination=mock/mock_object_type_repository.go -package=mock github.com/krobus00/storage-service/internal/model ObjectTypeRepository

package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ErrObjectTypeNotFound  = errors.New("object type not found")
	ErrExtensionNotAllowed = errors.New("object extensions not allowed")
)

type ObjectType struct {
	ID   string
	Name string
}

func (ObjectType) TableName() string {
	return "object_types"
}

func NewObjectTypeCacheKeyByID(id string) string {
	return fmt.Sprintf("objects:type:typeID:%s", id)
}

func NewObjectTypeCacheKeyByName(name string) string {
	return fmt.Sprintf("objects:type:typeName:%s", name)
}

func GetObjectTypeCacheKeys(id string, name string) []string {
	return []string{
		NewObjectTypeCacheKeyByID(id),
		NewObjectTypeCacheKeyByName(name),
	}
}

type ObjectTypeRepository interface {
	Create(ctx context.Context, objectType *ObjectType) error
	FindByID(ctx context.Context, id string) (*ObjectType, error)
	FindByName(ctx context.Context, name string) (*ObjectType, error)
	DeleteByID(ctx context.Context, id string) error

	// DI
	InjectDB(db *gorm.DB) error
	InjectRedisClient(client *redis.Client) error
}
