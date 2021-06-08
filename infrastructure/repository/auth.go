package repository

import (
	"context"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/model"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/infrastructure/external"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/logger"
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
	log.WithField("auth", auth).Info("Auth attributes")

	jwt, err := a.Service.Client.Login(ctx, a.Service.ClientID, a.Service.ClientSecret, a.Service.Realm, auth.Username, auth.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("jwt", jwt).Info("jwt response")

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
	log.WithField("refreshToken", refreshToken).Info("refreshToken attributes")

	jwt, err := a.Service.Client.RefreshToken(ctx, refreshToken, a.Service.ClientID, a.Service.ClientSecret, a.Service.Realm)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("jwt", jwt).Info("jwt response")

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

func (a *AuthRepository) FindEmployeeClaimsByToken(ctx context.Context, accessToken string) (*model.EmployeeClaims, error) {
	span, ctx := apm.StartSpan(ctx, "FindEmployeeClaimsByToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))
	log.WithField("accessToken", accessToken).Info("accessToken attributes")

	jwt, _, err := a.Service.Client.DecodeAccessToken(ctx, accessToken, a.Service.Realm, a.Service.Audience)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}
	log.WithField("jwt", jwt).Info("jwt decoded")

	employeeClaims := new(model.EmployeeClaims)
	mapstructure.Decode(jwt.Claims, employeeClaims)
	log.WithField("employeeClaims", employeeClaims).Info("employeeClaims mapstructure")

	type ResourceAccess struct {
		ResourceAccess map[string]map[string][]string `mapstructure:"resource_access"`
	}

	ra := new(ResourceAccess)
	mapstructure.Decode(jwt.Claims, ra)

	roles := ra.ResourceAccess[a.Service.ClientID]["roles"]
	employeeClaims.Roles = roles
	log.WithField("employeeClaims", employeeClaims).Info("employeeClaims with roles")

	return employeeClaims, nil
}

func NewAuthRepository(service *external.Keycloak) *AuthRepository {
	return &AuthRepository{
		Service: service,
	}
}
