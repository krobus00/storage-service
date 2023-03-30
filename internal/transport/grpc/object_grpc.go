package grpc

import (
	"context"

	"github.com/krobus00/storage-service/internal/model"
	pb "github.com/krobus00/storage-service/pb/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *Delivery) GetObjectByID(ctx context.Context, req *pb.GetObjectByIDRequest) (*pb.Object, error) {
	ctx = setUserIDCtx(ctx, req.GetUserId())

	presignedObject, err := t.objectUC.GeneratePresignedURL(ctx, &model.GetPresignedURLPayload{
		ObjectID: req.GetObjectId(),
	})

	switch err {
	case nil:
	case model.ErrObjectNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return presignedObject.ToGRPCResponse(), nil
}
