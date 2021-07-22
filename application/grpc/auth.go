package grpc

import (
	"context"

	"github.com/c-4u/auth-service/application/grpc/pb"
	"github.com/c-4u/auth-service/domain/service"
	"github.com/c-4u/auth-service/logger"
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

func (a *AuthGrpcService) FindClaimsByToken(ctx context.Context, in *pb.FindClaimsByTokenRequest) (*pb.Claims, error) {
	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("in", in).Info("handling FindClaimsByToken request")

	claims, err := a.AuthService.FindClaimsByToken(ctx, in.AccessToken)
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

func NewAuthGrpcService(service *service.AuthService) *AuthGrpcService {
	return &AuthGrpcService{
		AuthService: service,
	}
}
