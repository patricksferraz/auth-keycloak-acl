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
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func StartRestServer(keycloak *external.Keycloak, kafka *external.Kafka, port int) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))
	r.Use(apmgin.Middleware(r))

	repository := repository.NewRepository(keycloak, kafka)
	service := _service.NewService(repository)
	authMiddlerare := NewAuthMiddleware(service)
	restService := NewRestService(service, authMiddlerare)

	v1 := r.Group("api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		auth := v1.Group("/auth")
		{
			auth.POST("/login", restService.Login)
			auth.POST("/refresh-token", restService.RefreshToken)
			auth.POST("/claims", restService.FindClaimsByToken)
		}

		user := v1.Group("/users", authMiddlerare.Require())
		{
			user.POST("", restService.CreateUser)
			user.POST("/:id/password", restService.SetPassword)

			user.GET("/", restService.SearchUsers)
			user.GET("/:id", restService.FindUser)
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("cannot start rest server", err)
	}

	log.Printf("rest server has been started on port %d", port)
}
