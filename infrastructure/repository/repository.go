package repository

import (
	"context"

	"github.com/Nerzal/gocloak/v8"
	"github.com/c-4u/auth-service/domain/entity"
	"github.com/c-4u/auth-service/infrastructure/external"
	"github.com/c-4u/auth-service/logger"
	"github.com/c-4u/auth-service/utils"
	"github.com/mitchellh/mapstructure"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type Repository struct {
	K *external.Keycloak
}

func NewRepository(keycloak *external.Keycloak) *Repository {
	return &Repository{
		K: keycloak,
	}
}

func (r *Repository) Login(ctx context.Context, auth *entity.Auth) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "Login", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := r.K.Client.Login(ctx, r.K.ClientID, r.K.ClientSecret, r.K.Realm, auth.Username, auth.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &entity.JWT{
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

func (r *Repository) RefreshToken(ctx context.Context, refreshToken string) (*entity.JWT, error) {
	span, ctx := apm.StartSpan(ctx, "RefreshToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, err := r.K.Client.RefreshToken(ctx, refreshToken, r.K.ClientID, r.K.ClientSecret, r.K.Realm)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	return &entity.JWT{
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

func (r *Repository) FindClaimsByToken(ctx context.Context, accessToken string) (*entity.Claims, error) {
	span, ctx := apm.StartSpan(ctx, "FindClaimsByToken", "repository")
	defer span.End()

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	jwt, _, err := r.K.Client.DecodeAccessToken(ctx, accessToken, r.K.Realm, r.K.Audience)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		return nil, err
	}

	Claims := new(entity.Claims)
	mapstructure.Decode(jwt.Claims, Claims)
	log.WithField("claims", Claims).Info("claims mapstructure")

	type ResourceAccess struct {
		ResourceAccess map[string]map[string][]string `mapstructure:"resource_access"`
	}

	ra := new(ResourceAccess)
	mapstructure.Decode(jwt.Claims, ra)

	roles := ra.ResourceAccess[r.K.ClientID]["roles"]
	Claims.Roles = roles
	log.WithField("claims", Claims).Info("claims with roles")

	return Claims, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entity.User, accessToken string) error {
	gUser := gocloak.User{
		Username:      &user.Username,
		FirstName:     &user.FirstName,
		LastName:      &user.LastName,
		Email:         &user.Email,
		Enabled:       &user.Enabled,
		EmailVerified: &user.EmailVerified,
	}
	gUser.Attributes = utils.StructToAttr(gUser)

	userID, err := r.K.Client.CreateUser(ctx, accessToken, r.K.Realm, gUser)
	if err != nil {
		return err
	}

	user.ID = userID
	return nil
}

func (r *Repository) SetPassword(ctx context.Context, pass *entity.PasswordInfo, accessToken string) error {
	err := r.K.Client.SetPassword(ctx, accessToken, pass.UserID, r.K.Realm, pass.Password, pass.Temporary)
	return err
}
