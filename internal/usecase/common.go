package usecase

import (
	"context"
	"fmt"

	authPB "github.com/krobus00/auth-service/pb/auth"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
)

func getUserIDFromCtx(ctx context.Context) string {
	ctxUserID := ctx.Value(constant.KeyUserIDCtx)

	userID := fmt.Sprintf("%v", ctxUserID)
	if userID == "" {
		return constant.GuestID
	}
	return userID
}

func hasAccess(ctx context.Context, authClient authPB.AuthServiceClient, permissions []string) error {
	userID := getUserIDFromCtx(ctx)

	res, err := authClient.HasAccess(ctx, &authPB.HasAccessRequest{
		UserId:      userID,
		Permissions: permissions,
	})

	if err != nil {
		return model.ErrUnauthorizeAccess
	}
	if res == nil {
		return model.ErrUnauthorizeAccess
	}

	if res.Value {
		return nil
	}
	return model.ErrUnauthorizeAccess
}
