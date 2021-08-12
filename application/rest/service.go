package rest

import (
	http "net/http"

	"github.com/c-4u/auth-service/domain/service"
	"github.com/c-4u/auth-service/logger"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type RestService struct {
	Service    *service.Service
	Middleware *AuthMiddleware
}

func NewRestService(service *service.Service) *RestService {
	return &RestService{
		Service: service,
	}
}

// Login godoc
// @Summary log in
// @ID login
// @Tags Auth
// @Description System authentication
// @Accept json
// @Produce json
// @Param body body Auth true "JSON body for authentication"
// @Success 200 {object} JWT
// @Failure 401 {object} HTTPError
// @Router /auth/login [post]
func (s *RestService) Login(ctx *gin.Context) {
	var json Auth

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := s.Service.Login(ctx, json.Username, json.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

// RefreshToken godoc
// @Summary refresh token
// @ID refreshToken
// @Tags Auth
// @Description Refresh token route
// @Accept json
// @Produce json
// @Param body body RefreshToken true "JSON body for refresh token"
// @Success 200 {object} JWT
// @Failure 400 {object} HTTPError
// @Router /auth/refresh-token [post]
func (s *RestService) RefreshToken(ctx *gin.Context) {
	var json RefreshToken

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := s.Service.RefreshToken(ctx, json.RefreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

// FindClaimsByToken godoc
// @Summary get claims
// @ID findClaimsByToken
// @Tags Auth
// @Description Get Claims by access token
// @Accept json
// @Produce json
// @Param body body AccessToken true "JSON body for get claims"
// @Success 200 {object} Claims
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /auth/claims [post]
func (s *RestService) FindClaimsByToken(ctx *gin.Context) {
	var json AccessToken

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := s.Service.FindClaimsByToken(ctx, json.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, claims)
}

// CreateUser godoc
// @Summary create User
// @ID createUser
// @Tags User
// @Description Create User
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "JSON body for create user"
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /user [post]
func (s *RestService) CreateUser(ctx *gin.Context) {
	var json CreateUserRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.Service.CreateUser(ctx, json.Username, json.EmployeeID, s.Middleware.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, CreateUserResponse{ID: *id})
}

// SetPassword godoc
// @Summary create User
// @ID setPassword
// @Tags User
// @Description Set user password
// @Accept json
// @Produce json
// @Param body body SetPasswordRequest true "JSON body set user password"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /user/{id}/password [post]
func (s *RestService) SetPassword(ctx *gin.Context) {
	var req IDRequest
	var json SetPasswordRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, HTTPError{Error: err.Error()})
		return
	}

	err := s.Service.SetPassword(ctx, req.ID, json.Password, json.Temporary, s.Middleware.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HTTPError{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, HTTPResponse{Code: http.StatusOK, Message: "updated successfully"})
}
