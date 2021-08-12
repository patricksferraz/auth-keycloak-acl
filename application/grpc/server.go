package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/c-4u/auth-service/application/grpc/pb"
	_service "github.com/c-4u/auth-service/domain/service"
	"github.com/c-4u/auth-service/infrastructure/external"
	"github.com/c-4u/auth-service/infrastructure/repository"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(keycloak *external.Keycloak, kafka *external.Kafka, port int) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery())),
	)
	reflection.Register(grpcServer)

	repository := repository.NewRepository(keycloak, kafka)
	service := _service.NewService(repository)
	grpcService := NewGrpcService(service)
	pb.RegisterAuthKeycloakAclServer(grpcServer, grpcService)

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
