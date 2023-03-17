package usecase

import (
	"context"
	"errors"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/util"
	log "github.com/sirupsen/logrus"
)

type storageUsecase struct {
	storageRepo model.StorageRepository
	authClient  authPB.AuthServiceClient
}

func NewStorageUsecase() model.StorageUsecase {
	return new(storageUsecase)
}

func (uc *storageUsecase) Upload(ctx context.Context, payload *model.FileUploadPayload) (*model.Storage, error) {
	logger := log.WithFields(log.Fields{
		"path":     payload.Path,
		"fileName": payload.Filename,
		"isPublic": payload.IsPublic,
	})
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	newStorage := model.NewStorage().
		SetID(util.GenerateUUID()).
		SetSrc(payload.Src).
		SetUploadedBy(userID).
		SetFileName(payload.Filename).
		SetObjectKey(model.DefaultPath).
		SetIsPublic(payload.IsPublic)

	err = uc.storageRepo.Create(ctx, newStorage)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return newStorage, nil
}

func (uc *storageUsecase) GeneratePresignedURL(ctx context.Context, payload *model.GetPresignedURLPayload) (*model.GetPresignedURLResponse, error) {
	logger := log.WithFields(log.Fields{
		"objectID": payload.ObjectID,
	})
	storage, err := uc.storageRepo.FindByID(ctx, payload.ObjectID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if storage == nil {
		return nil, errors.New("object not found")
	}

	err = uc.hasAccess(ctx, storage)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	presignObject, err := uc.storageRepo.GeneratePresignedURL(ctx, storage)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return presignObject, nil
}

func (uc *storageUsecase) hasAccess(ctx context.Context, storage *model.Storage) error {
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	if storage.IsPublic {
		return nil
	}
	if !storage.IsPublic {
		allowAccess, err := hasAccess(ctx, uc.authClient, []string{constant.FULL_ACCESS})
		if err != nil {
			return err
		}
		if allowAccess {
			return nil
		}
		if storage.UploadedBy != userID {
			return errors.New("unauthorized access")
		}
	}
	return nil
}
