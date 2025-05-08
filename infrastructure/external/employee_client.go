package external

import (
	"github.com/patricksferraz/auth-service/application/grpc/pb"
	"google.golang.org/grpc"
)

type EmployeeClient struct {
	C pb.EmployeeServiceClient
}

func NewEmployeeClient(cc *grpc.ClientConn) *EmployeeClient {
	return &EmployeeClient{
		C: pb.NewEmployeeServiceClient(cc),
	}
}
