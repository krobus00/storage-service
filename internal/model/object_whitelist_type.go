//go:generate mockgen -destination=mock/mock_object_whitelisty_type_repository.go -package=mock github.com/krobus00/storage-service/internal/model ObjectWhitelistTypeRepository

package model

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ObjectWhitelistType struct {
	TypeID    string
	Extension string
}

func (ObjectWhitelistType) TableName() string {
	return "object_whitelist_types"
}

func NewObjectWhitelistTypeCacheKey(typeID string) string {
	return fmt.Sprintf("object-whitelist-types:typeID:%s:extension", typeID)
}

func GetObjectWhitelistTypCacheKeys(typeID string) []string {
	return []string{
		NewObjectWhitelistTypeCacheKey(typeID),
		"object-whitelist-types:typeID:*",
	}
}

type ObjectWhitelistTypeRepository interface {
	Create(ctx context.Context, objectWhitelistType *ObjectWhitelistType) error
	FindByTypeIDAndExt(ctx context.Context, typeID string, ext string) (*ObjectWhitelistType, error)
	DeleteByTypeIDAndExt(ctx context.Context, typeID string, ext string) error

	// DI
	InjectDB(db *gorm.DB) error
	InjectRedisClient(client *redis.Client) error
}
