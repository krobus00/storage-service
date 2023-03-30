package grpc

import (
	"errors"

	"github.com/krobus00/storage-service/internal/model"
)

func (t *Delivery) InjectObjectUsecase(uc model.ObjectUsecase) error {
	if uc == nil {
		return errors.New("invalid object usecase")
	}
	t.objectUC = uc
	return nil
}
