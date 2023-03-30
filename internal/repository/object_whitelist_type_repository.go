package repository

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type objectWhitelistTypeRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewObjectWhitelistTypeRepository() model.ObjectWhitelistTypeRepository {
	return new(objectWhitelistTypeRepository)
}

func (r *objectWhitelistTypeRepository) Create(ctx context.Context, objectWhitelistType *model.ObjectWhitelistType) error {
	logger := logrus.WithFields(logrus.Fields{
		"typeID":    objectWhitelistType.TypeID,
		"extension": objectWhitelistType.Extension,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	err := db.WithContext(ctx).Create(objectWhitelistType).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectWhitelistTypCacheKeys(objectWhitelistType.TypeID))

	return nil
}

func (r *objectWhitelistTypeRepository) FindByTypeIDAndExt(ctx context.Context, typeID string, ext string) (*model.ObjectWhitelistType, error) {
	logger := logrus.WithFields(logrus.Fields{
		"typeID":    typeID,
		"extension": ext,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	objectWhitelistType := new(model.ObjectWhitelistType)
	cacheBucketKey := utils.NewBucketKey(model.NewObjectWhitelistTypeCacheKey(typeID), ext)

	cachedData, err := HGet(ctx, r.redisClient, cacheBucketKey, ext)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &objectWhitelistType)
	if err == nil {
		return objectWhitelistType, nil
	}

	objectWhitelistType = new(model.ObjectWhitelistType)

	err = db.WithContext(ctx).
		First(objectWhitelistType, "type_id = ? AND extension = ?", typeID, ext).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = HSetWithExpiry(ctx, r.redisClient, cacheBucketKey, ext, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		return nil, err
	}

	err = HSetWithExpiry(ctx, r.redisClient, cacheBucketKey, ext, objectWhitelistType)
	if err != nil {
		logger.Error(err.Error())
	}

	return objectWhitelistType, nil
}

func (r *objectWhitelistTypeRepository) DeleteByTypeIDAndExt(ctx context.Context, typeID string, ext string) error {
	logger := logrus.WithFields(logrus.Fields{
		"typeID":    typeID,
		"extension": ext,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	objectWhitelistType := new(model.ObjectWhitelistType)

	err := db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("type_id = ? AND extension = ?", typeID, ext).
		Delete(objectWhitelistType).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectWhitelistTypCacheKeys(objectWhitelistType.TypeID))

	return nil
}
