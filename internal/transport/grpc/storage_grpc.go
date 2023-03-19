package grpc

import (
	"context"

	"github.com/krobus00/storage-service/internal/model"
	pb "github.com/krobus00/storage-service/pb/storage"
)

// Server :nodoc:
type Server struct {
	objectUC model.ObjectUsecase
	pb.UnsafeStorageServiceServer
}

// NewGRPCServer :nodoc:
func NewGRPCServer() *Server {
	return new(Server)
}

func (t *Server) GetObjectByID(ctx context.Context, req *pb.GetObjectByIDRequest) (*pb.Object, error) {
	ctx = setUserIDCtx(ctx, req.GetUserId())

	presignedObject, err := t.objectUC.GeneratePresignedURL(ctx, &model.GetPresignedURLPayload{
		ObjectID: req.GetObjectId(),
	})
	if err != nil {
		return nil, err
	}

	return presignedObject.ToGRPCResponse(), nil
}
