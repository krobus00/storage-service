package grpc

import (
	"context"

	"github.com/krobus00/storage-service/internal/model"
	pb "github.com/krobus00/storage-service/pb/storage"
)

// Server :nodoc:
type Server struct {
	storageUC model.StorageUsecase
	pb.UnsafeStorageServiceServer
}

// NewGRPCServer :nodoc:
func NewGRPCServer() *Server {
	return new(Server)
}

func (t *Server) GetObjectByID(ctx context.Context, req *pb.GetObjectByIDRequest) (*pb.Storage, error) {
	ctx = setUserIDCtx(ctx, req.GetUserId())

	storage, err := t.storageUC.GeneratePresignedURL(ctx, &model.GetPresignedURLPayload{
		ObjectID: req.GetObjectId(),
	})
	if err != nil {
		return nil, err
	}

	return storage.ToGRPCResponse(), nil
}
