package rest

import (
	http "net/http"

	"dev.azure.com/c4ut/TimeClock/_git/auth-service/domain/service"
	"dev.azure.com/c4ut/TimeClock/_git/auth-service/logger"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

type AuthRestService struct {
	AuthService *service.AuthService
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
func (a *AuthRestService) Login(ctx *gin.Context) {
	var json Auth

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("auth", json).Info("handling Auth request")

	jwt, err := a.AuthService.Login(ctx, json.Username, json.Password)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	log.WithField("jwt", jwt).Info("JWT response")

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
// @Router /auth/refreshToken [post]
func (a *AuthRestService) RefreshToken(ctx *gin.Context) {
	var json RefreshToken

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("refreshToken", json).Info("handling RefreshToken request")

	jwt, err := a.AuthService.RefreshToken(ctx, json.RefreshToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("jwt", jwt).Info("JWT response")

	ctx.JSON(http.StatusOK, jwt)
}

// FindEmployeeClaimsByToken godoc
// @Summary get employee claims
// @ID findEmployeeClaimsByToken
// @Tags Auth
// @Description Get Employee Claims by access token
// @Accept json
// @Produce json
// @Param body body AccessToken true "JSON body for get employee claims"
// @Success 200 {object} EmployeeClaims
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /auth/employeeClaims [post]
func (a *AuthRestService) FindEmployeeClaimsByToken(ctx *gin.Context) {
	var json AccessToken

	log := logger.Log.WithFields(apmlogrus.TraceContext(ctx))

	if err := ctx.ShouldBindJSON(&json); err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.WithField("accessToken", json).Info("handling AccessToken request")

	employee, err := a.AuthService.FindEmployeeClaimsByToken(ctx, json.AccessToken)
	if err != nil {
		log.WithError(err)
		apm.CaptureError(ctx, err).Send()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.WithField("employeeClaims", employee).Info("employeeClaims response")

	ctx.JSON(http.StatusOK, employee)
}

func NewAuthRestService(service *service.AuthService) *AuthRestService {
	return &AuthRestService{
		AuthService: service,
	}
}
