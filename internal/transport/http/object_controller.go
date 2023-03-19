package http

import (
	"net/http"

	"github.com/krobus00/storage-service/internal/model"
	"github.com/labstack/echo/v4"
)

type ObjectController struct {
	objectUC model.ObjectUsecase
}

func NewObjectController() *ObjectController {
	return new(ObjectController)
}

func (t *ObjectController) Upload(eCtx echo.Context) (err error) {
	var (
		ctx = buildContext(eCtx)
		res = model.NewResponse()
		req = new(model.HTTPFileUploadRequest)
	)

	err = eCtx.Bind(req)
	if err != nil {
		res = model.WithBadRequestResponse(nil)
		return eCtx.JSON(http.StatusBadRequest, res)
	}

	req.Src, err = eCtx.FormFile("file")
	if err != nil {
		return err
	}

	object, err := t.objectUC.Upload(ctx, &model.ObjectPayload{
		Src: req.Src,
		Object: &model.Object{
			FileName: req.Filename,
			Type:     req.Type,
			IsPublic: req.IsPublic,
		},
	})
	switch err {
	case nil:
	case model.ErrExtensionNotAllowed:
		return eCtx.JSON(http.StatusBadRequest, res.WithMessage(err.Error()))
	case model.ErrObjectTypeNotFound:
		return eCtx.JSON(http.StatusBadRequest, res.WithMessage(err.Error()))
	case model.ErrUserNotFound:
		return eCtx.JSON(http.StatusBadRequest, res.WithMessage(err.Error()))
	default:
		return eCtx.JSON(http.StatusInternalServerError, res.WithMessage("internal server error"))
	}

	res.WithData(object)
	return eCtx.JSON(http.StatusCreated, res)
}

func (t *ObjectController) GetPresignURL(eCtx echo.Context) (err error) {
	var (
		ctx = buildContext(eCtx)
		res = model.NewResponse()
		req = new(model.HTTPGetPresignedURLRequest)
	)

	err = eCtx.Bind(req)
	if err != nil {
		res = model.WithBadRequestResponse(nil)
		return eCtx.JSON(http.StatusBadRequest, res)
	}

	presignedObject, err := t.objectUC.GeneratePresignedURL(ctx, req.ToPayload())
	switch err {
	case nil:
	case model.ErrObjectNotFound:
		return eCtx.JSON(http.StatusBadRequest, res.WithMessage(err.Error()))
	case model.ErrUserNotFound:
		return eCtx.JSON(http.StatusBadRequest, res.WithMessage(err.Error()))
	default:
		return eCtx.JSON(http.StatusInternalServerError, res.WithMessage("internal server error"))
	}

	res.WithData(presignedObject.ToHTTPResponse())
	return eCtx.JSON(http.StatusCreated, res)
}