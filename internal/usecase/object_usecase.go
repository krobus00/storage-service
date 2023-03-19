package usecase

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/util"
	log "github.com/sirupsen/logrus"
)

type objectUsecase struct {
	objectRepo              model.ObjectRepository
	objectTypeRepo          model.ObjectTypeRepository
	ObjectWhitelistTypeRepo model.ObjectWhitelistTypeRepository
	authClient              authPB.AuthServiceClient
}

func NewObjectUsecase() model.ObjectUsecase {
	return new(objectUsecase)
}

func (uc *objectUsecase) Upload(ctx context.Context, payload *model.ObjectPayload) (*model.Object, error) {
	logger := log.WithFields(log.Fields{
		"objectKey": payload.Object.Key,
		"fileName":  payload.Object.FileName,
		"isPublic":  payload.Object.IsPublic,
	})
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	objectType, err := uc.objectTypeRepo.FindByName(ctx, payload.Object.Type)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if objectType == nil {
		return nil, model.ErrObjectTypeNotFound
	}

	err = uc.validationObjectType(ctx, payload.Src, objectType.ID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	newObject := model.NewObject().
		SetID(util.GenerateUUID()).
		SetTypeID(objectType.ID).
		SetType(objectType.Name).
		SetUploadedBy(userID).
		SetFileName(payload.Object.FileName).
		SetKey(model.DefaultPath).
		SetIsPublic(payload.Object.IsPublic)

	payload.SetObject(newObject)

	err = uc.objectRepo.Create(ctx, payload)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return payload.Object, nil
}

func (uc *objectUsecase) GeneratePresignedURL(ctx context.Context, payload *model.GetPresignedURLPayload) (*model.GetPresignedURLResponse, error) {
	logger := log.WithFields(log.Fields{
		"objectID": payload.ObjectID,
	})
	object, err := uc.objectRepo.FindByID(ctx, payload.ObjectID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if object == nil {
		return nil, model.ErrObjectNotFound
	}

	err = uc.hasAccess(ctx, object)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	presignObject, err := uc.objectRepo.GeneratePresignedURL(ctx, object)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return presignObject, nil
}

func (uc *objectUsecase) hasAccess(ctx context.Context, object *model.Object) error {
	if object.IsPublic {
		return nil
	}
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return model.ErrUnauthorizedObjectAccess
	}
	if object.UploadedBy == userID {
		return nil
	}
	if !object.IsPublic {
		allowAccess, err := hasAccess(ctx, uc.authClient, []string{constant.FULL_ACCESS})
		if err != nil {
			return model.ErrUnauthorizedObjectAccess
		}
		if allowAccess {
			return nil
		}
	}
	return model.ErrUnauthorizedObjectAccess
}

func (uc *objectUsecase) validationObjectType(ctx context.Context, data *multipart.FileHeader, typeID string) error {
	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return err
	}

	contentType := http.DetectContentType(buf.Bytes())
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return err
	}

	whiteList, err := uc.ObjectWhitelistTypeRepo.FindByTypeIDAndExt(ctx, typeID, exts[0])
	if err != nil {
		return err
	}
	if whiteList == nil {
		return model.ErrExtensionNotAllowed
	}
	return nil
}
