package grpc

import (
	"fmt"
	"log"
	"net"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/infrastructure/external"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/infrastructure/repository"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(_service *external.Keycloak, port int) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery())),
	)
	reflection.Register(grpcServer)

	authRepository := repository.NewAuthRepository(_service)
	authService := service.NewAuthService(authRepository)
	authGrpcService := NewAuthGrpcService(authService)
	pb.RegisterAuthServiceServer(grpcServer, authGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}

	log.Printf("gRPC server has been started on port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}
}
