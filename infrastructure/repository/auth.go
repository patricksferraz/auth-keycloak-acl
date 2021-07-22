package repository

import (
	"context"

	"github.com/c-4u/auth-service/domain/model"
	"github.com/c-4u/auth-service/infrastructure/external"
	"github.com/c-4u/auth-service/logger"
	"github.com/mitchellh/mapstructure"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthRepository struct {
	Service *external.Keycloak
}

func (a *AuthRepository) Login(ctx context.Context, auth *model.Auth) (*model.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "Login", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := a.Service.Client.Login(ctx, a.Service.ClientID, a.Service.ClientSecret, a.Service.Realm, auth.Username, auth.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &model.JWT{
		AccessToken:      jwt.AccessToken,
		IDToken:          jwt.IDToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  jwt.NotBeforePolicy,
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, nil
}

func (a *AuthRepository) RefreshToken(ctx context.Context, refreshToken string) (*model.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := a.Service.Client.RefreshToken(ctx, refreshToken, a.Service.ClientID, a.Service.ClientSecret, a.Service.Realm)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &model.JWT{
		AccessToken:      jwt.AccessToken,
		IDToken:          jwt.IDToken,
		ExpiresIn:        jwt.ExpiresIn,
		RefreshExpiresIn: jwt.RefreshExpiresIn,
		RefreshToken:     jwt.RefreshToken,
		TokenType:        jwt.TokenType,
		NotBeforePolicy:  jwt.NotBeforePolicy,
		SessionState:     jwt.SessionState,
		Scope:            jwt.Scope,
	}, nil
}

func (a *AuthRepository) FindClaimsByToken(ctx context.Context, accessToken string) (*model.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "FindClaimsByToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, _, err := a.Service.Client.DecodeAccessToken(ctx, accessToken, a.Service.Realm, a.Service.Audience)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	Claims := new(model.Claims)
	mapstructure.Decode(jwt.Claims, Claims)
	log.WithField("claims", Claims).Info("claims mapstructure")

	type ResourceAccess struct {
		ResourceAccess map[string]map[string][]string `mapstructure:"resource_access"`
	}

	ra := new(ResourceAccess)
	mapstructure.Decode(jwt.Claims, ra)

	roles := ra.ResourceAccess[a.Service.ClientID]["roles"]
	Claims.Roles = roles
	log.WithField("claims", Claims).Info("claims with roles")

	return Claims, nil
}

func NewAuthRepository(service *external.Keycloak) *AuthRepository {
	return &AuthRepository{
		Service: service,
	}
}
