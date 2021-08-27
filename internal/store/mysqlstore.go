// Storage backend for sql databases (mysqlstore.go)
package store

import (
	"context"
	"errors"
	"notification-service/internal/model"
	"notification-service/internal/util"
	"time"

	"gorm.io/gorm"
)

// Store interface for common database operations for auth-service.
type Store interface {
	GetNotification(ctx context.Context, notifications *[]model.Notification, UserID string) error
	UpdateLastLogin(ctx context.Context, docID string, lastLoginTimestamp time.Time) error
}

// Store embeds gorm type to provide extra methods specific to auth-service.
type MySQLStore struct {
	*gorm.DB
}

// This is to make sure MySQLStore implements all of the Store interface functions.
var _ Store = (*MySQLStore)(nil)

var errNotificationNotFound = errors.New("This user does not have notifications")

// FindUser using either username or email.
func (ms *MySQLStore) GetNotification(ctx context.Context, notifications *[]model.Notification, UserID string) error {

	dbResult := ms.DB.Model(&model.Notification{}).WithContext(ctx).Where("user_id = ?", UserID).Find(notifications)
	// fmt.Println(dbResult)
	dbErr := dbResult.Error
	dbRowsAffected := dbResult.RowsAffected

	if dbErr == nil && dbRowsAffected == 0 {
		return &util.NotFound{ErrMessage: errNotificationNotFound.Error()}
	} else if dbErr != nil {
		return &util.InternalServer{ErrMessage: dbErr.Error()}
	}
	return nil
}

// UpdateLastLogin timestamp set on the users table.
func (ms *MySQLStore) UpdateLastLogin(ctx context.Context, docID string, lastLoginTimestamp time.Time) error {
	/*
		dbErr := ms.DB.WithContext(ctx).Model(&model.User{}).Where("document_id = ?", docID).Update("last_login_timestamp", lastLoginTimestamp).Error

		if dbErr != nil {
			return &util.InternalServer{ErrMessage: dbErr.Error()}
		}
	*/
	return nil
}
