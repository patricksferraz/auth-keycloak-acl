package repository

import (
	"context"

	"github.com/patricksferraz/auth-service/domain/entity"
)

type RepositoryInterface interface {
	Login(ctx context.Context, auth *entity.Auth) (*entity.JWT, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*entity.JWT, error)
	FindClaimsByToken(ctx context.Context, accessToken string) (*entity.Claims, error)

	CreateUser(ctx context.Context, user *entity.User, accessToken string) error
	FindUser(ctx context.Context, id string, accessToken string) (*entity.User, error)
	SearchUsers(ctx context.Context, filter *entity.Filter, accessToken string) ([]*entity.User, error)
	SetPassword(ctx context.Context, pass *entity.PasswordInfo, accessToken string) error

	PublishEvent(ctx context.Context, msg, topic, key string) error

	FindEmployee(ctx context.Context, employeeID string) error
}
