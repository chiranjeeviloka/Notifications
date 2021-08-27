package main

import (
	"net/http"
	"notification-service/internal/notification"
	"notification-service/internal/store"

	"github.com/gin-gonic/gin"
)

const (
	notificationMessage = "User notifications"
)

// TODO: This is just an example handler.

func (svr *Service) ViewNotification(c *gin.Context) {

	// gorm.DB is embeded into our MySQLStore struct that implements Store interface.
	db := store.MySQLStore{DB: svr.db}
	response, err := notification.View(c, &db, svr.queue)
	if err != nil {
		handleError(c, err)
		return
	}

	successResponse(c, gin.H{
		"status":  http.StatusOK,
		"message": notificationMessage,
		"error":   false,
		"data":    response,
	})

}
