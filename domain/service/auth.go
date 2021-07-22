package service

import (
	"context"

	"github.com/c-4u/auth-service/domain/model"
	"github.com/c-4u/auth-service/domain/repository"
	"github.com/c-4u/auth-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthService struct {
	AuthRepository repository.AuthRepositoryInterface
}

func (a *AuthService) Login(ctx context.Context, username, password string) (*model.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "Login", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	auth, err := model.NewAuth(username, password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	jwt, err := a.AuthRepository.Login(ctx, auth)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return jwt, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := a.AuthRepository.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return jwt, nil
}

func (a *AuthService) FindClaimsByToken(ctx context.Context, accessToken string) (*model.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "FindClaimsByToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	claims, err := a.AuthRepository.FindClaimsByToken(ctx, accessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("claims", claims).Info("claims response")

	return claims, nil
}

func NewAuthService(authRepository repository.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}
