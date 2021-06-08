package grpc

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/application/grpc/pb"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthGrpcService struct {
	pb.UnimplementedAuthServiceServer
	AuthService *service.AuthService
}

func (a *AuthGrpcService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.JWT, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling Login request")

	jwt, err := a.AuthService.Login(ctx, in.Username, in.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.JWT{}, err
	}
	log.WithField("jwt", jwt).Info("JWT response")

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

func (a *AuthGrpcService) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.JWT, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling RefreshToken request")

	jwt, err := a.AuthService.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.JWT{}, err
	}
	log.WithField("jwt", jwt).Info("JWT response")

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

func (a *AuthGrpcService) FindEmployeeClaimsByToken(ctx context.Context, in *pb.FindEmployeeClaimsByTokenRequest) (*pb.EmployeeClaims, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling FindEmployeeClaimsByToken request")

	employee, err := a.AuthService.FindEmployeeClaimsByToken(ctx, in.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return &pb.EmployeeClaims{}, err
	}
	log.WithField("employeeClaims", employee).Info("employeeClaims response")

	return &pb.EmployeeClaims{
		Id:    employee.ID,
		Roles: employee.Roles,
	}, nil
}

func NewAuthGrpcService(service *service.AuthService) *AuthGrpcService {
	return &AuthGrpcService{
		AuthService: service,
	}
}
