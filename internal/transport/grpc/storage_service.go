package grpc

import (
	"github.com/krobus00/storage-service/internal/model"
	pb "github.com/krobus00/storage-service/pb/storage"
)

type Delivery struct {
	objectUC model.ObjectUsecase
	pb.UnsafeStorageServiceServer
}

func NewDelivery() *Delivery {
	return new(Delivery)
}
