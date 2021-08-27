package main

import (
	"net/http"
	"notification-service/internal/util"

	"github.com/gin-gonic/gin"
)

// TODO: These are optional to use. Create more here if you need extra utilities for working with Gin.

// successResponse send successful response back to client.
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *util.InternalServer:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"status":  http.StatusInternalServerError,
			"message": e.ErrMessage,
		})
		return
	case *util.BadRequest:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"status":  http.StatusBadRequest,
			"message": e.ErrMessage,
		})
		return
	case *util.NotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"status":  http.StatusNotFound,
			"message": e.ErrMessage,
		})
		return
	case *util.UnAuthorized:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"status":  http.StatusUnauthorized,
			"message": e.ErrMessage,
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"status":  http.StatusInternalServerError,
			"message": e.Error(),
		})
		return
	}
}
