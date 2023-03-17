package grpc

import (
	"errors"

	"github.com/krobus00/storage-service/internal/model"
)

func (t *Server) InjectStorageUsecase(uc model.StorageUsecase) error {
	if uc == nil {
		return errors.New("invalid storage usecase")
	}
	t.storageUC = uc
	return nil
}
