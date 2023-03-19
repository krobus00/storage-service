package repository

import (
	"context"
	"errors"

	"github.com/krobus00/storage-service/internal/model"
	"gorm.io/gorm"
)

type objectTypeRepository struct {
	db *gorm.DB
}

func NewObjectTypeRepository() model.ObjectTypeRepository {
	return new(objectTypeRepository)
}

func (r *objectTypeRepository) Create(ctx context.Context, objectType *model.ObjectType) error {
	err := r.db.WithContext(ctx).Create(objectType).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *objectTypeRepository) FindByID(ctx context.Context, id string) (*model.ObjectType, error) {
	objectType := new(model.ObjectType)
	err := r.db.WithContext(ctx).First(objectType, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return objectType, nil
}

func (r *objectTypeRepository) FindByName(ctx context.Context, name string) (*model.ObjectType, error) {
	objectType := new(model.ObjectType)
	err := r.db.WithContext(ctx).First(objectType, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return objectType, nil
}
