package bootstrap

import (
	"context"
	"fmt"
	"time"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/delivery/http"
	"github.com/krobus00/storage-service/internal/infrastructure"
	"github.com/krobus00/storage-service/internal/repository"
	"github.com/krobus00/storage-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	infrastructure.InitializeDBConn()

	// init infra
	db, err := infrastructure.DB.DB()
	continueOrFatal(err)

	redisClient, err := infrastructure.NewRedisClient()
	continueOrFatal(err)

	s3Client, err := infrastructure.NewS3Client()
	continueOrFatal(err)

	echo := infrastructure.NewEcho()

	// init grpc client
	authConn, err := grpc.Dial(config.AuthGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	continueOrFatal(err)
	authClient := authPB.NewAuthServiceClient(authConn)

	// init repository
	storageRepo := repository.NewStorageRepository()
	err = storageRepo.InjectDB(infrastructure.DB)
	continueOrFatal(err)
	err = storageRepo.InjectS3Client(s3Client)
	continueOrFatal(err)

	// init usecase
	storageUsecase := usecase.NewStorageUsecase()
	err = storageUsecase.InjectStorageRepo(storageRepo)
	continueOrFatal(err)
	err = storageUsecase.InjectAuthClient(authClient)
	continueOrFatal(err)

	// init delivery layer
	// ini http
	storageCtrl := http.NewStorageController()
	err = storageCtrl.InjectStorageUsecase(storageUsecase)
	continueOrFatal(err)

	httpDelivery := http.NewHTTPDelivery()
	err = httpDelivery.InjectEcho(echo)
	continueOrFatal(err)
	err = httpDelivery.InjectStorageController(storageCtrl)
	continueOrFatal(err)
	httpDelivery.InitRoutes()

	go func() {
		_ = echo.Start(":" + config.HTTPPort())
	}()
	log.Info(fmt.Sprintf("http server started on :%s", config.HTTPPort()))

	wait := gracefulShutdown(context.Background(), 30*time.Second, map[string]operation{
		"redis connection": func(ctx context.Context) error {
			return redisClient.Close()
		},
		"database connection": func(ctx context.Context) error {
			infrastructure.StopTickerCh <- true
			return db.Close()
		},
		"http": func(ctx context.Context) error {
			return echo.Shutdown(ctx)
		},
	})

	<-wait
}
