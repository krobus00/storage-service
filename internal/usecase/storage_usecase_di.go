package usecase

import (
	"errors"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/model"
)

func (uc *storageUsecase) InjectStorageRepo(repo model.StorageRepository) error {
	if repo == nil {
		return errors.New("invalid storage repository")
	}
	uc.storageRepo = repo
	return nil
}

func (uc *storageUsecase) InjectAuthClient(client authPB.AuthServiceClient) error {
	if client == nil {
		return errors.New("invalid auth client")
	}
	uc.authClient = client
	return nil
}
