package form

type NotificationRequest struct {
	UserID           string `json:"userId" validate:"required,min=1,max=30"`
	NotificationType int    `json:"notificationType,omitempty" validate:"number,omitempty"`
}

type NotificationResponseData struct {
	Type           int    `json:"type"`
	Message        string `json:"message"`
	PostLink       string `json:"postLink"`
	NotificationID string `json:"notificationId"`
}
