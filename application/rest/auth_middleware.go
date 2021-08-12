package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/c-4u/auth-service/domain/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Service     *service.Service
	AccessToken string
}

func (a *AuthMiddleware) Require() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.Request.Header.Get("Authorization")
		if accessToken == "" {
			err := errors.New("authorization token is not provided")
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		_, err := a.Service.FindClaimsByToken(ctx, accessToken)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("access token is invalid: %v", err)})
			ctx.Abort()
			return
		}

		a.AccessToken = accessToken

		// TODO: adds retricted permissions
		// for _, role := range claims.Roles {
		// 	if role == method {
		// 		return nil
		// 	}
		// }

		// return status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}
}

func NewAuthMiddleware(service *service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		Service: service,
	}
}
