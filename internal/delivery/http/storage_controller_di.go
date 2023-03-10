package http

import (
	"errors"

	"github.com/krobus00/storage-service/internal/model"
)

func (d *StorageController) InjectStorageUsecase(uc model.StorageUsecase) error {
	if uc == nil {
		return errors.New("invalid storage usecase")
	}
	d.storageUC = uc
	return nil
}
