package repository

import (
	"context"
	"errors"

	"github.com/krobus00/storage-service/internal/model"
	"gorm.io/gorm"
)

type objectWhitelistTypeRepository struct {
	db *gorm.DB
}

func NewObjectWhitelistTypeRepository() model.ObjectWhitelistTypeRepository {
	return new(objectWhitelistTypeRepository)
}

func (r *objectWhitelistTypeRepository) Create(ctx context.Context, objectWhitelistType *model.ObjectWhitelistType) error {
	err := r.db.WithContext(ctx).Create(objectWhitelistType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *objectWhitelistTypeRepository) FindByTypeIDAndExt(ctx context.Context, typeID string, ext string) (*model.ObjectWhitelistType, error) {
	objectType := new(model.ObjectWhitelistType)
	err := r.db.WithContext(ctx).First(objectType, "type_id = ? AND extension = ?", typeID, ext).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return objectType, nil
}
