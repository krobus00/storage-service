package grpc

import (
	"context"

	"github.com/krobus00/storage-service/internal/constant"
)

func setUserIDCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, constant.KeyUserIDCtx, userID)
}
