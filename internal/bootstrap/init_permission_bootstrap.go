package bootstrap

import (
	"context"
	"fmt"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func StartInitPermission() {
	authConn, err := grpc.Dial(config.AuthGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.ContinueOrFatal(err)
	authClient := authPB.NewAuthServiceClient(authConn)

	ctx := context.Background()

	for _, permission := range constant.SeedPermissions {
		currentPermission, _ := authClient.FindPermissionByName(ctx, &authPB.FindPermissionByNameRequest{
			SessionUserId: constant.SystemID,
			Name:          permission,
		})
		if currentPermission == nil {
			_, err = authClient.CreatePermission(ctx, &authPB.CreatePermissionRequest{
				SessionUserId: constant.SystemID,
				Name:          permission,
			})
			utils.ContinueOrFatal(err)
			logrus.Info(fmt.Sprintf("permission %s created", permission))
		} else {
			logrus.Info(fmt.Sprintf("permission %s already exist", permission))
		}
	}

	for userGroup, permissions := range constant.SeedGroupPermissios {
		group, err := authClient.FindGroupByName(ctx, &authPB.FindGroupByNameRequest{
			SessionUserId: constant.SystemID,
			Name:          userGroup,
		})
		utils.ContinueOrFatal(err)
		for _, permission := range permissions {
			currentPermission, _ := authClient.FindPermissionByName(ctx, &authPB.FindPermissionByNameRequest{
				SessionUserId: constant.SystemID,
				Name:          permission,
			})

			if currentPermission != nil {
				_, err = authClient.CreateGroupPermission(ctx, &authPB.CreateGroupPermissionRequest{
					SessionUserId: constant.SystemID,
					GroupId:       group.GetId(),
					PermissionId:  currentPermission.GetId(),
				})
				e, ok := status.FromError(err)
				if !ok {
					utils.ContinueOrFatal(err)
				}
				switch e.Code() {
				case codes.AlreadyExists:
					logrus.Info(fmt.Sprintf("group permission %s:%s already exist", userGroup, permission))
				case codes.OK:
					logrus.Info(fmt.Sprintf("group permission %s:%s created", userGroup, permission))
				default:
					logrus.Error(fmt.Sprintf("group permission %s:%s failed: %s", userGroup, permission, e.Message()))
				}
			}
		}
	}
}
