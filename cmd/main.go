package main

import (
	"fmt"
	"log"
	"net"
	"os"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/irrepressiblespirit/website-messages-service/pkg/core"
	"github.com/irrepressiblespirit/website-messages-service/pkg/grpcapi"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/authtoken"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/configservice"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/eventstorage"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/externaluserstorage"
	pb "github.com/irrepressiblespirit/website-messages-service/pkg/impl/grpcapi"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/storage"
	"github.com/irrepressiblespirit/website-messages-service/pkg/impl/tokenservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("-------------------Start application---------------")
	configService := configservice.NewConfigService(getEnv("config", "config.yaml"))
	config, err := configService.Load()
	if err != nil {
		panic(err)
	}
	dbstorage, err := storage.NewMongoStorage(config.Mongo)
	if err != nil {
		panic(err)
	}

	eventstorage := eventstorage.NewGrpcEventStorage(config.Centrifugo.Grpc)
	if err := eventstorage.Start(); err != nil {
		panic(err)
	}
	defer eventstorage.Close()

	core := core.Core{
		Storage:      dbstorage,
		StorageUsers: externaluserstorage.NewStorageUserDecorator(config.ExternalUser, dbstorage),
		TokenService: tokenservice.NewTokenService(config.Centrifugo.Secret),
		EventStorage: eventstorage,
		AuthToken:    authtoken.New(config.Token),
		Config:       config,
	}
	startGRPC(&core)
}

func startGRPC(core *core.Core) {
	listenerGrpc, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	authInterceptor := pb.UnaryAuthInterceptor{
		Auth:    core.AuthToken,
		PrintMD: false,
	}
	serverGRPC := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.ChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor, pb.TimeInterceptor, authInterceptor.Interceptor),
	)
	grpcapi.RegisterServiceMessagesServer(serverGRPC, &pb.GRPCServer{Core: *core})
	reflection.Register(serverGRPC)
	log.Printf("gRPC server started at :8081")
	if err = serverGRPC.Serve(listenerGrpc); err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
