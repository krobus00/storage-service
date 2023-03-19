package usecase

import (
	"context"
	"fmt"

	authPB "github.com/krobus00/auth-service/pb/auth"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
)

func getUserIDFromCtx(ctx context.Context) (string, error) {
	ctxUserID := ctx.Value(constant.KeyUserIDCtx)

	userID := fmt.Sprintf("%v", ctxUserID)
	if userID == "" {
		return "", model.ErrUserNotFound
	}
	return userID, nil
}

func hasAccess(ctx context.Context, authClient authPB.AuthServiceClient, accessList []string) (bool, error) {
	hasAccess := false
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return hasAccess, err
	}
	res, err := authClient.HasAccess(ctx, &authPB.HasAccessRequest{
		UserId:      userID,
		AccessNames: accessList,
	})

	if err != nil {
		return hasAccess, err
	}
	if res == nil {
		return hasAccess, model.ErrUnauthorizedObjectAccess
	}
	hasAccess = res.Value
	return hasAccess, nil
}
