package main

import (
	"github.com/FreedomCentral/central/secret"
	"github.com/gin-gonic/gin"
)

const jwtKey = "jwtKey"

func setupRouter(svr *Service, sec secret.Secret) *gin.Engine {

	/*key, err := sec.Get(jwtKey)
	if err != nil {
		logger.Fatalf("Failed to get jwtkey from secrets %q: %v", jwtKey, err)
	}
	*/
	// Create service instance. Main instance shared by all http handlers.
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	/*router.Use(auth.ValidateToken([]byte(key)))*/

	router.GET("/health", svr.Health)
	//router.GET("/view/:user_id/*type", svr.ViewNotification)

	// TODO: add other gin middlewares here.

	// IMPORTANT handlers are methods on the instance of the Service struct.
	// that is how databases and other common parameters get passed in into the handlers.
	// Prefix URLs with service name and API version number.
	router.Group("/notifications")
	{
		v1Group := router.Group("/notifications/v1")
		{
			v1Group.GET("/view", svr.ViewNotification)
		}
	}
	return router
}
