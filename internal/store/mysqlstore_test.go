package store_test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetNotification(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		expectedError error
	}{
		{
			name:          "test with dumy userid",
			userID:        "123",
			expectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
		})
	}
}
