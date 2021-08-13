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

func NewRestService(service *service.Service, authMiddleware *AuthMiddleware) *RestService {
	return &RestService{
		Service:    service,
		Middleware: authMiddleware,
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
		ctx.JSON(http.StatusForbidden, HTTPError{Code: http.StatusForbidden, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, claims)
}

// CreateUser godoc
// @Security ApiKeyAuth
// @Summary create user
// @ID createUser
// @Tags User
// @Description Create User
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "JSON body for create user"
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users [post]
func (s *RestService) CreateUser(ctx *gin.Context) {
	var json CreateUserRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.Service.CreateUser(ctx, json.Username, json.EmployeeID, s.Middleware.I.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusForbidden, HTTPError{Code: http.StatusForbidden, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, CreateUserResponse{ID: *id})
}

// FindUser godoc
// @Security ApiKeyAuth
// @Summary find user
// @Description Router for find user
// @ID findUser
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /users/{id} [get]
func (s *RestService) FindUser(ctx *gin.Context) {
	var req IDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	user, err := s.Service.FindUser(ctx, req.ID, s.Middleware.I.AccessToken)
	if err != nil {
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// SearchUsers godoc
// @Security ApiKeyAuth
// @Summary search users by filter
// @ID searchUsers
// @Tags User
// @Description Search users by `filter`. if the page and page size are empty, 0 and 10 will be considered respectively.
// @Accept json
// @Produce json
// @Param body query SearchUsersRequest true "JSON body for search users"
// @Success 200 {array} User
// @Failure 400 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Router /users [get]
func (s *RestService) SearchUsers(ctx *gin.Context) {
	var body SearchUsersRequest

	if err := ctx.ShouldBindQuery(&body); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	users, err := s.Service.SearchUsers(ctx, body.Filter.Username, body.Filter.Enabled, body.Filter.PageSize, body.Filter.Page, s.Middleware.I.AccessToken)
	if err != nil {
		ctx.JSON(
			http.StatusForbidden,
			HTTPError{
				Code:  http.StatusForbidden,
				Error: err.Error(),
			},
		)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// SetPassword godoc
// @Security ApiKeyAuth
// @Summary set user password
// @ID setPassword
// @Tags User
// @Description Set user password
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body SetPasswordRequest true "JSON body set user password"
// @Success 200 {object} HTTPResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users/{id}/password [post]
func (s *RestService) SetPassword(ctx *gin.Context) {
	var req IDRequest
	var json SetPasswordRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			HTTPError{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, HTTPError{Error: err.Error()})
		return
	}

	err := s.Service.SetPassword(ctx, req.ID, json.Password, json.Temporary, s.Middleware.I.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusForbidden, HTTPError{Code: http.StatusForbidden, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, HTTPResponse{Code: http.StatusOK, Message: "updated successfully"})
}
