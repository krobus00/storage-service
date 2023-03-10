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
	ctx, err := setUserInfoContext(ctx, uc.authClient)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	user, err := getUserInfoFromContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	newStorage := model.NewStorage().
		SetID(util.GenerateUUID()).
		SetSrc(payload.Src).
		SetUploadedBy(user.GetId()).
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

func (uc *storageUsecase) GeneratePresignURL(ctx context.Context, payload *model.GetPresignURLPayload) (*model.GetPresignURLResponse, error) {
	logger := log.WithFields(log.Fields{
		"objectKey": payload.ObjectKey,
	})
	ctx, _ = setUserInfoContext(ctx, uc.authClient)
	storage, err := uc.storageRepo.FindByObjectKey(ctx, payload.ObjectKey)
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

	presignObject, err := uc.storageRepo.GeneratePresignURL(ctx, storage)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return presignObject, nil
}

func (uc *storageUsecase) hasAccess(ctx context.Context, storage *model.Storage) error {
	user, _ := getUserInfoFromContext(ctx)
	if storage.IsPublic {
		return nil
	}
	if !storage.IsPublic {
		if user == nil {
			return errors.New("user not found")
		}

		allowAccess, err := HasAccess(ctx, uc.authClient, []string{constant.FULL_ACCESS})
		if err != nil {
			return err
		}
		if allowAccess {
			return nil
		}
		if storage.UploadedBy != user.GetId() {
			return errors.New("unauthorized access")
		}
	}
	return nil
}
