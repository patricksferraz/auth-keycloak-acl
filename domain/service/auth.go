package service

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/model"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/repository"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/logger"
	"github.com/sirupsen/logrus"
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
	log.WithFields(
		logrus.Fields{
			"username": username,
			"password": password,
		},
	).Info("Login attributes")

	auth, err := model.NewAuth(username, password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("auth", auth).Info("auth request")

	jwt, err := a.AuthRepository.Login(ctx, auth)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("jwt", jwt).Info("jwt response")

	return jwt, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("refreshToken", refreshToken).Info("RefreshToken attributes")

	jwt, err := a.AuthRepository.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("jwt", jwt).Info("jwt response")

	return jwt, nil
}

func (a *AuthService) FindEmployeeClaimsByToken(ctx context.Context, accessToken string) (*model.EmployeeClaims, error) {
	span, ctx := apm.StartSpan(ctx, "FindEmployeeClaimsByToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("accessToken", accessToken).Info("FindEmployeeClaimsByToken attributes")

	employee, err := a.AuthRepository.FindEmployeeClaimsByToken(ctx, accessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("employee", employee).Info("employee response")

	return employee, nil
}

func NewAuthService(authRepository repository.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}
