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

type objectTypeRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewObjectTypeRepository() model.ObjectTypeRepository {
	return new(objectTypeRepository)
}

func (r *objectTypeRepository) Create(ctx context.Context, objectType *model.ObjectType) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"id":   objectType.ID,
		"name": objectType.Name,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	err := db.WithContext(ctx).Create(objectType).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectTypeCacheKeys(objectType.ID, objectType.Name))

	return nil
}

func (r *objectTypeRepository) FindByID(ctx context.Context, id string) (*model.ObjectType, error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	objectType := new(model.ObjectType)
	cacheKey := model.NewObjectTypeCacheKeyByID(id)

	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &objectType)
	if err == nil {
		return objectType, nil
	}

	objectType = new(model.ObjectType)

	err = db.WithContext(ctx).First(objectType, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = SetWithExpiry(ctx, r.redisClient, cacheKey, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		return nil, err
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, objectType)
	if err != nil {
		logger.Error(err.Error())
	}

	return objectType, nil
}

func (r *objectTypeRepository) FindByName(ctx context.Context, name string) (*model.ObjectType, error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"name": name,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	objectType := new(model.ObjectType)
	cacheKey := model.NewObjectTypeCacheKeyByName(name)

	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &objectType)
	if err == nil {
		return objectType, nil
	}

	objectType = new(model.ObjectType)

	err = db.WithContext(ctx).First(objectType, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = SetWithExpiry(ctx, r.redisClient, cacheKey, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		return nil, err
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, objectType)
	if err != nil {
		logger.Error(err.Error())
	}

	return objectType, nil
}

func (r *objectTypeRepository) DeleteByID(ctx context.Context, id string) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	objectType := new(model.ObjectType)

	err := db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("id = ?", id).
		Delete(objectType).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetObjectTypeCacheKeys(objectType.ID, objectType.Name))

	return nil
}
