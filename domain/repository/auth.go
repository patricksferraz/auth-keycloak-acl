package repository

import (
	"context"

	"github.com/c-4u/auth-service/domain/model"
)

type AuthRepositoryInterface interface {
	Login(ctx context.Context, auth *model.Auth) (*model.JWT, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.JWT, error)
	FindClaimsByToken(ctx context.Context, accessToken string) (*model.Claims, error)
}
