package service

import (
	"context"

	"github.com/patricksferraz/auth-service/domain/entity"
	"github.com/patricksferraz/auth-service/domain/repository"
	"github.com/patricksferraz/auth-service/infrastructure/external/topic"
	"github.com/patricksferraz/auth-service/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type Service struct {
	Repository repository.RepositoryInterface
}

func NewService(repository repository.RepositoryInterface) *Service {
	return &Service{
		Repository: repository,
	}
}

func (a *Service) Login(ctx context.Context, username, password string) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "Login", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	auth, err := entity.NewAuth(username, password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	jwt, err := a.Repository.Login(ctx, auth)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return jwt, nil
}

func (a *Service) RefreshToken(ctx context.Context, refreshToken string) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := a.Repository.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return jwt, nil
}

func (a *Service) FindClaimsByToken(ctx context.Context, accessToken string) (*entity.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "FindClaimsByToken", "service")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	claims, err := a.Repository.FindClaimsByToken(ctx, accessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("claims", claims).Info("claims response")

	return claims, nil
}

func (s *Service) CreateUser(ctx context.Context, username, employeeID, accessToken string) (*string, error) {
	user, err := entity.NewUser(username, employeeID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.FindEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.CreateUser(ctx, user, accessToken)
	if err != nil {
		return nil, err
	}

	event, err := entity.NewUserEvent(user)
	if err != nil {
		return nil, err
	}

	msg, err := event.ToJson()
	if err != nil {
		return nil, err
	}

	err = s.Repository.PublishEvent(ctx, string(msg), topic.NEW_USER, user.ID)
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (s *Service) FindUser(ctx context.Context, id, accessToken string) (*entity.User, error) {
	user, err := s.Repository.FindUser(ctx, id, accessToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) SearchUsers(ctx context.Context, username *string, enabled *bool, pageSize, page *int, accessToken string) ([]*entity.User, error) {
	filter, err := entity.NewFilter(username, enabled, pageSize, page)
	if err != nil {
		return nil, err
	}

	users, err := s.Repository.SearchUsers(ctx, filter, accessToken)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) SetPassword(ctx context.Context, userID string, password string, temporary bool, accessToken string) error {
	pass := entity.NewPasswordInfo(userID, password, temporary)
	err := s.Repository.SetPassword(ctx, pass, accessToken)
	return err
}

func (a *Service) Logout(ctx context.Context, refreshToken string) error {
	err := a.Repository.Logout(ctx, refreshToken)
	return err
}
