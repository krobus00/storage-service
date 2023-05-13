package usecase

import (
	"bytes"
	"context"
	"errors"
	"mime"
	"net/http"
	"sync"

	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/nats-io/nats.go"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type objectUsecase struct {
	objectRepo              model.ObjectRepository
	objectTypeRepo          model.ObjectTypeRepository
	ObjectWhitelistTypeRepo model.ObjectWhitelistTypeRepository
	authClient              authPB.AuthServiceClient
	jsClient                nats.JetStreamContext
}

func NewObjectUsecase() model.ObjectUsecase {
	return new(objectUsecase)
}

func (uc *objectUsecase) CreateStream() error {
	stream, _ := uc.jsClient.StreamInfo(model.ObjectStreamName)
	// stream not found, create it
	if stream == nil {
		logrus.Printf("Creating stream: %s\n", model.ObjectStreamName)

		_, err := uc.jsClient.AddStream(&nats.StreamConfig{
			Name:     model.ObjectStreamName,
			Subjects: []string{model.ObjectStreamSubjects},
			MaxAge:   config.JetstreamMaxAge(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *objectUsecase) Upload(ctx context.Context, payload *model.ObjectPayload) (*model.Object, error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"objectKey": payload.Object.Key,
		"fileName":  payload.Object.FileName,
		"isPublic":  payload.Object.IsPublic,
	})

	userID := getUserIDFromCtx(ctx)

	err := uc.hasAccess(ctx, constant.ActionCreate, nil)
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
		SetID(utils.GenerateUUID()).
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
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

	err = uc.hasAccess(ctx, constant.ActionRead, object)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	objectType, err := uc.objectTypeRepo.FindByID(ctx, object.TypeID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if objectType == nil {
		return nil, model.ErrObjectNotFound
	}
	object.SetType(objectType.Name)

	presignedObject, err := uc.objectRepo.GeneratePresignedURL(ctx, object)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return presignedObject, nil
}

func (uc *objectUsecase) DeleteObject(ctx context.Context, id string) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := logrus.WithFields(logrus.Fields{
		"objectID": id,
	})

	object, err := uc.objectRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if object == nil {
		return model.ErrObjectNotFound
	}

	err = uc.hasAccess(ctx, constant.ActionDelete, object)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = uc.objectRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	jsPayload := model.JSDeleteObjectPayload{
		ObjectID: object.ID,
	}

	wg := sync.WaitGroup{}
	for _, subject := range model.ObjectDeleteStreamSubjects {
		wg.Add(1)
		go func(subject string) {
			defer wg.Done()
			err = publishJS(ctx, uc.jsClient, subject, jsPayload)
			if err != nil {
				logger.Error(err.Error())
			}
		}(subject)
	}

	wg.Wait()

	return nil
}

func (uc *objectUsecase) hasAccess(ctx context.Context, action constant.ActionType, object *model.Object) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	permissions := []string{
		constant.PermissionObjectAll,
	}

	switch action {
	case constant.ActionCreate:
		permissions = append(permissions, constant.PermissionObjectCreate)
	case constant.ActionRead:
		if object.IsPublic {
			return nil
		}
		if !object.IsPublic && object.UploadedBy != userID {
			permissions = append(permissions, constant.PermissionObjectReadPrivate)
		} else {
			permissions = append(permissions, constant.PermissionObjectRead)
		}
	case constant.ActionDelete:
		if object.UploadedBy != userID {
			permissions = append(permissions, constant.PermissionObjectModifyOther)
		} else {
			permissions = append(permissions, constant.PermissionObjectDelete)
		}
	default:
		return errors.New("invaid action")
	}

	err := hasAccess(ctx, uc.authClient, permissions)
	if err != nil {
		return err
	}
	return nil
}

func (uc *objectUsecase) validationObjectType(ctx context.Context, data []byte, typeID string) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	buf := bytes.NewBuffer(data)

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
