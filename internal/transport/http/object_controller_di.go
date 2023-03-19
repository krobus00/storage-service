package http

import (
	"errors"

	"github.com/krobus00/storage-service/internal/model"
)

func (t *ObjectController) InjectObjectUsecase(uc model.ObjectUsecase) error {
	if uc == nil {
		return errors.New("invalid object usecase")
	}
	t.objectUC = uc
	return nil
}
