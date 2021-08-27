//  mockstore.go
package store

import (
	"context"
	"notification-service/internal/model"
	"time"
)

// Store type is used to perform unit tests.
// It does nothing in this example but you can change it to perform expected actions during tests.
type MockStore struct {
}

// This is to make sure MockStore implements all of the Store interface functions.
// var _ Store = (*MockStore)(nil)

func (ms *MockStore) GetNotification(ctx context.Context, notifications []*model.Notification, UserID string) error {
	return nil
}

func (ms *MockStore) UpdateLastLogin(ctx context.Context, docID string, lastLoginTimestamp time.Time) error {
	return nil
}
