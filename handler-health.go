package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svr *Service) Health(c *gin.Context) {
	successResponse(c, gin.H{
		"status":  http.StatusOK,
		"message": "Working",
		"error":   false,
	})
}
