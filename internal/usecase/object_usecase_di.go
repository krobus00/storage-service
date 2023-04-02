package usecase

import (
	"errors"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/nats-io/nats.go"
)

func (uc *objectUsecase) InjectObjectRepo(repo model.ObjectRepository) error {
	if repo == nil {
		return errors.New("invalid object repository")
	}
	uc.objectRepo = repo
	return nil
}

func (uc *objectUsecase) InjectObjectTypeRepo(repo model.ObjectTypeRepository) error {
	if repo == nil {
		return errors.New("invalid object type repository")
	}
	uc.objectTypeRepo = repo
	return nil
}

func (uc *objectUsecase) InjectObjectWhitelistTypeRepo(repo model.ObjectWhitelistTypeRepository) error {
	if repo == nil {
		return errors.New("invalid object whitelist type repository")
	}
	uc.ObjectWhitelistTypeRepo = repo
	return nil
}

func (uc *objectUsecase) InjectAuthClient(client authPB.AuthServiceClient) error {
	if client == nil {
		return errors.New("invalid auth client")
	}
	uc.authClient = client
	return nil
}

func (uc *objectUsecase) InjectJetstreamClient(client nats.JetStreamContext) error {
	if client == nil {
		return errors.New("invalid jetstream client")
	}
	uc.jsClient = client
	return nil
}
