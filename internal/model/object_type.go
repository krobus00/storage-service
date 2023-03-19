package model

import (
	"context"
	"errors"

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

type ObjectWhitelistType struct {
	ID        string
	TypeID    string
	Extension string
}

func (ObjectWhitelistType) TableName() string {
	return "object_whitelist_types"
}

type ObjectTypeRepository interface {
	Create(ctx context.Context, objectType *ObjectType) error
	FindByID(ctx context.Context, id string) (*ObjectType, error)
	FindByName(ctx context.Context, name string) (*ObjectType, error)

	// DI
	InjectDB(db *gorm.DB) error
}

type ObjectWhitelistTypeRepository interface {
	Create(ctx context.Context, objectWhitelistType *ObjectWhitelistType) error
	FindByTypeIDAndExt(ctx context.Context, typeID string, ext string) (*ObjectWhitelistType, error)

	// DI
	InjectDB(db *gorm.DB) error
}
