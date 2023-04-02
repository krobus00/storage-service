package model

type PublisherUsecase interface {
	CreateStream() error
}

type JSDeleteObjectPayload struct {
	ObjectID string `json:"objectID"`
}
