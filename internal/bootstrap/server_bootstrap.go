package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/http"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/storage-service/internal/config"
	"github.com/krobus00/storage-service/internal/infrastructure"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/krobus00/storage-service/internal/repository"
	grpcServer "github.com/krobus00/storage-service/internal/transport/grpc"
	httpServer "github.com/krobus00/storage-service/internal/transport/http"
	"github.com/krobus00/storage-service/internal/usecase"
	pb "github.com/krobus00/storage-service/pb/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/sirupsen/logrus"
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

	nc, js, err := infrastructure.NewJetstreamClient()
	continueOrFatal(err)

	tp, err := infrastructure.JaegerTraceProvider()
	continueOrFatal(err)

	echo := infrastructure.NewEcho()

	// init grpc client
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	authConn, err := grpc.Dial(config.AuthGRPCHost(), grpcOpts...)
	continueOrFatal(err)
	authClient := authPB.NewAuthServiceClient(authConn)

	// init repository
	objectRepo := repository.NewObjectRepository()
	err = objectRepo.InjectDB(infrastructure.DB)
	continueOrFatal(err)
	err = objectRepo.InjectS3Client(s3Client)
	continueOrFatal(err)
	err = objectRepo.InjectRedisClient(redisClient)
	continueOrFatal(err)

	objectTypeRepo := repository.NewObjectTypeRepository()
	err = objectTypeRepo.InjectDB(infrastructure.DB)
	continueOrFatal(err)
	err = objectTypeRepo.InjectRedisClient(redisClient)
	continueOrFatal(err)

	objectWhitelistTypeRepo := repository.NewObjectWhitelistTypeRepository()
	err = objectWhitelistTypeRepo.InjectDB(infrastructure.DB)
	continueOrFatal(err)
	err = objectWhitelistTypeRepo.InjectRedisClient(redisClient)
	continueOrFatal(err)

	// init usecase
	objectUsecase := usecase.NewObjectUsecase()
	err = objectUsecase.InjectObjectRepo(objectRepo)
	continueOrFatal(err)
	err = objectUsecase.InjectObjectTypeRepo(objectTypeRepo)
	continueOrFatal(err)
	err = objectUsecase.InjectObjectWhitelistTypeRepo(objectWhitelistTypeRepo)
	continueOrFatal(err)
	err = objectUsecase.InjectAuthClient(authClient)
	continueOrFatal(err)
	err = objectUsecase.InjectJetstreamClient(js)
	continueOrFatal(err)

	// init stream
	publisherUsecase := []model.PublisherUsecase{
		objectUsecase,
	}

	for _, uc := range publisherUsecase {
		err = uc.CreateStream()
		continueOrFatal(err)
	}

	// init delivery layer
	// init http
	objectCtrl := httpServer.NewObjectController()
	err = objectCtrl.InjectObjectUsecase(objectUsecase)
	continueOrFatal(err)

	httpDelivery := httpServer.NewDelivery()
	err = httpDelivery.InjectEcho(echo)
	continueOrFatal(err)
	err = httpDelivery.InjectObjectController(objectCtrl)
	continueOrFatal(err)
	httpDelivery.InitRoutes()

	// init grpc
	grpcDelivery := grpcServer.NewDelivery()
	err = grpcDelivery.InjectObjectUsecase(objectUsecase)
	continueOrFatal(err)

	storageGrpcServer := grpc.NewServer()

	pb.RegisterStorageServiceServer(storageGrpcServer, grpcDelivery)
	if config.Env() == "development" {
		reflection.Register(storageGrpcServer)
	}
	lis, _ := net.Listen("tcp", ":"+config.PortGRPC())
	go func() {
		_ = storageGrpcServer.Serve(lis)
	}()
	logrus.Info(fmt.Sprintf("grpc server started on :%s", config.PortGRPC()))

	go func() {
		_ = echo.Start(":" + config.PortHTTP())
	}()
	logrus.Info(fmt.Sprintf("http server started on :%s", config.PortHTTP()))

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%s", config.PortMetrics()), nil)
	}()
	logrus.Info(fmt.Sprintf("metrics server started on :%s", config.PortMetrics()))

	wait := gracefulShutdown(context.Background(), config.GracefulShutdownTimeOut(), map[string]operation{
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
		"grpc": func(ctx context.Context) error {
			return lis.Close()
		},
		"nats connection": func(ctx context.Context) error {
			return nc.Drain()
		},
		"trace provider": func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})

	<-wait
}
