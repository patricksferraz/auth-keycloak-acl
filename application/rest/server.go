package rest

import (
	"fmt"
	"log"

	_ "github.com/c-4u/auth-service/application/rest/docs"
	_service "github.com/c-4u/auth-service/domain/service"
	"github.com/c-4u/auth-service/infrastructure/external"
	"github.com/c-4u/auth-service/infrastructure/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin"
)

// @title Auth Swagger API
// @version 1.0
// @description Swagger API for Golang Project Auth.
// @termsOfService http://swagger.io/terms/

// @contact.name Coding4u
// @contact.email contato@coding4u.com.br

// @BasePath /api/v1
func StartRestServer(keycloak *external.Keycloak, port int) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(apmgin.Middleware(r))

	repository := repository.NewRepository(keycloak)
	service := _service.NewService(repository)
	restService := NewRestService(service)

	v1 := r.Group("api/v1/auth")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.POST("/login", restService.Login)
		v1.POST("/refresh-token", restService.RefreshToken)
		v1.POST("/claims", restService.FindClaimsByToken)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("cannot start rest server", err)
	}

	log.Printf("rest server has been started on port %d", port)
}
