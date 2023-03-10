package http

import (
	"net/http"

	"github.com/krobus00/storage-service/internal/model"
	"github.com/labstack/echo/v4"
)

type StorageController struct {
	storageUC model.StorageUsecase
}

func NewStorageController() *StorageController {
	return new(StorageController)
}

func (d *StorageController) Upload(eCtx echo.Context) (err error) {
	var (
		ctx = buildContext(eCtx)
		res = new(model.Response)
		req = new(model.HTTPFileUploadRequest)
	)

	err = eCtx.Bind(req)
	if err != nil {
		res = model.NewResponse().WithMessage("bad request")
		return eCtx.JSON(http.StatusBadRequest, res)
	}
	req.Src, err = eCtx.FormFile("file")
	if err != nil {
		return err
	}

	storage, err := d.storageUC.Upload(ctx, &model.FileUploadPayload{
		Src:      req.Src,
		Filename: req.Filename,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		return err
	}

	res = model.NewResponse().WithData(storage)
	return eCtx.JSON(http.StatusCreated, res)
}

func (d *StorageController) GetPresignURL(eCtx echo.Context) (err error) {
	var (
		ctx = buildContext(eCtx)
		res = new(model.Response)
		req = new(model.HTTPGetPresignURLRequest)
	)

	err = eCtx.Bind(req)
	if err != nil {
		res = model.NewResponse().WithMessage("bad request")
		return eCtx.JSON(http.StatusBadRequest, res)
	}

	storage, err := d.storageUC.GeneratePresignURL(ctx, req.ToPayload())
	if err != nil {
		return err
	}

	res = model.NewResponse().WithData(storage.ToHTTPResponse())
	return eCtx.JSON(http.StatusCreated, res)
}
