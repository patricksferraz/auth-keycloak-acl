package repository

import (
	"context"

	"github.com/c-4u/auth-service/domain/entity"
)

type RepositoryInterface interface {
	Login(ctx context.Context, auth *entity.Auth) (*entity.JWT, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.JWT, error)
	FindClaimsByToken(ctx context.Context, accessToken string) (*entity.Claims, error)

	CreateUser(ctx context.Context, user *entity.User, accessToken string) error
	SetPassword(ctx context.Context, pass *entity.PasswordInfo, accessToken string) error
}
