package usecase

import (
	"context"
	"errors"
	"fmt"

	authPB "github.com/krobus00/auth-service/pb/auth"

	"github.com/krobus00/storage-service/internal/constant"
)

func setUserInfoContext(ctx context.Context, authClient authPB.AuthServiceClient) (context.Context, error) {
	accessToken := fmt.Sprintf("%s", ctx.Value(constant.KeyTokenCtx))
	user, err := authClient.GetUserInfo(ctx, &authPB.AuthRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return ctx, err
	}
	if user == nil {
		return ctx, errors.New("user not found")
	}
	ctx = context.WithValue(ctx, constant.KeyUserInfoCtx, user)

	return ctx, nil
}

func getUserInfoFromContext(ctx context.Context) (*authPB.User, error) {
	ctxUser := ctx.Value(constant.KeyUserInfoCtx)

	val, ok := ctxUser.(*authPB.User)
	if !ok {
		return nil, errors.New("user not found")
	}
	return val, nil
}

func HasAccess(ctx context.Context, authClient authPB.AuthServiceClient, accessList []string) (bool, error) {
	hasAccess := false
	user, _ := getUserInfoFromContext(ctx)
	if user == nil {
		return hasAccess, errors.New("user not found")
	}
	res, err := authClient.HasAccess(ctx, &authPB.HasAccessRequest{
		UserId:      user.GetId(),
		AccessNames: accessList,
	})

	if err != nil {
		return hasAccess, err
	}
	if res == nil {
		return hasAccess, errors.New("not allowed")
	}
	hasAccess = res.Value
	return hasAccess, nil
}
