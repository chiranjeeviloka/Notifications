package main

import (
	"time"

	"github.com/FreedomCentral/central/secret"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func setupRouter(svr *Service, sec secret.Secret) *gin.Engine {

	// Create service instance. Main instance shared by all http handlers.
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, HEAD, OPTIONS, PATCH",
		RequestHeaders:  "Origin, Authorization, Content-Type, Content-Length",
		ExposedHeaders:  "",
		MaxAge:          12 * time.Hour,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.GET("/health", svr.Health)
	// router.Use(middlewares.VerifyTokenMiddleware()
	router.Group("/notifications")
	{
		v1Group := router.Group("/notifications/v1")
		{
			v1Group.GET("/view", svr.ViewNotification)
		}
	}
	return router
}
