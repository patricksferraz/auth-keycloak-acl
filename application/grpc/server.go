package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/patricksferraz/auth-service/application/grpc/pb"
	_service "github.com/patricksferraz/auth-service/domain/service"
	"github.com/patricksferraz/auth-service/infrastructure/external"
	"github.com/patricksferraz/auth-service/infrastructure/repository"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer(keycloak *external.Keycloak, kafka *external.Kafka, employeeServiceAddr string, port int) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery())),
	)
	reflection.Register(grpcServer)

	authInterceptor := external.NewAuthInterceptor()
	employeeConn, err := external.GrpcClient(employeeServiceAddr, authInterceptor.TransportOpts()...)
	if err != nil {
		log.Fatal(err)
	}
	defer employeeConn.Close()

	employeeClient := external.NewEmployeeClient(employeeConn)
	repository := repository.NewRepository(keycloak, kafka, employeeClient)
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
