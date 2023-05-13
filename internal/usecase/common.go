package usecase

import (
	"context"

	"fmt"

	"github.com/goccy/go-json"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/nats-io/nats.go"

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

	permissions = append(permissions, constant.PermissionFullAccess)
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

func publishJS(ctx context.Context, jsClient nats.JetStreamContext, subjectName string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = jsClient.Publish(
		subjectName,
		payload,
		nats.Context(ctx),
	)

	if err != nil {
		return err
	}
	return nil
}
