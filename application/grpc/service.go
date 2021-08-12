package grpc

import (
	"context"

	"github.com/c-4u/auth-service/application/grpc/pb"
	"github.com/c-4u/auth-service/domain/service"
	"github.com/c-4u/auth-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcService struct {
	pb.UnimplementedAuthKeycloakAclServer
	Service *service.Service
}

func NewGrpcService(service *service.Service) *GrpcService {
	return &GrpcService{
		Service: service,
	}
}

func (s *GrpcService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.JWT, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling Login request")

	jwt, err := s.Service.Login(ctx, in.Username, in.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.JWT{}, err
	}

	return &pb.JWT{
		AccessToken:      jwt.AccessToken,
		IdToken:          jwt.IDToken,
		ExpiresIn:        int64(jwt.ExpiresIn),
		RefreshExpiresIn: int64(jwt.RefreshExpiresIn),
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  int64(jwt.NotBeforePolicy),
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, err
}

func (s *GrpcService) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.JWT, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RefreshToken request")

	jwt, err := s.Service.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.JWT{}, err
	}

	return &pb.JWT{
		AccessToken:      jwt.AccessToken,
		IdToken:          jwt.IDToken,
		ExpiresIn:        int64(jwt.ExpiresIn),
		RefreshExpiresIn: int64(jwt.RefreshExpiresIn),
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  int64(jwt.NotBeforePolicy),
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, err
}

func (s *GrpcService) FindClaimsByToken(ctx context.Context, in *pb.FindClaimsByTokenRequest) (*pb.Claims, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling FindClaimsByToken request")

	claims, err := s.Service.FindClaimsByToken(ctx, in.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.Claims{}, err
	}

	return &pb.Claims{
		EmployeeId: claims.EmployeeID,
		Roles:      claims.Roles,
	}, nil
}

func (s *GrpcService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := s.Service.CreateUser(ctx, in.User.Username, in.User.EmployeeId, in.AccessToken)
	if err != nil {
		return &pb.CreateUserResponse{}, err
	}

	return &pb.CreateUserResponse{
		Id: *userID,
	}, nil
}

func (s *GrpcService) SetPassword(ctx context.Context, in *pb.SetPasswordRequest) (*pb.StatusResponse, error) {
	err := s.Service.SetPassword(ctx, in.UserId, in.Password, in.Temporary, in.AccessToken)
	if err != nil {
		return &pb.StatusResponse{
			Code:  uint32(status.Code(err)),
			Error: err.Error(),
		}, err
	}
	return &pb.StatusResponse{
		Code:    uint32(codes.OK),
		Message: "password updated successfully",
	}, nil
}
