package grpc

import (
	"context"

	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/utils"
	pb "github.com/krobus00/storage-service/pb/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *Delivery) GetObjectByID(ctx context.Context, req *pb.GetObjectByIDRequest) (*pb.Object, error) {
	ctx = setUserIDCtx(ctx, req.GetUserId())

	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

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

func (t *Delivery) DeleteObjectByID(ctx context.Context, req *pb.DeleteObjectByIDRequest) (*emptypb.Empty, error) {
	ctx = setUserIDCtx(ctx, req.GetUserId())

	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	err := t.objectUC.DeleteObject(ctx, req.GetObjectId())

	switch err {
	case nil:
	case model.ErrObjectNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &emptypb.Empty{}, nil
}
