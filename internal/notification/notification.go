package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"notification-service/internal/form"
	"notification-service/internal/model"
	"notification-service/internal/store"
	"notification-service/internal/util"
	"os"
	"strconv"
	"time"

	"github.com/FreedomCentral/central/audit"
	"github.com/FreedomCentral/central/queue"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

//audit messages
const (
	NotificationMessage = "GET_USER_NOTIFICATIONS"
)

func View(c *gin.Context, db store.Store, queue queue.Queue) ([]*form.NotificationResponseData, error) {
	var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	req := form.NotificationRequest{}
	if userID, ok := c.GetQuery("userid"); ok {
		req.UserID = userID
	}
	if notificationType, ok := c.GetQuery("type"); ok {
		tempType, err := strconv.Atoi(notificationType)
		if err != nil {
			return nil, &util.BadRequest{ErrMessage: err.Error()}
		}
		req.NotificationType = tempType
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, &util.BadRequest{ErrMessage: err.Error()}
	}

	var notifications []model.Notification

	err := db.GetNotification(c, &notifications, req.UserID)
	if err != nil {
		return nil, &util.InternalServer{ErrMessage: err.Error()}
	}
	//sample redis
	/*
		_, err := cache.Get(login.Username)
		if err != nil && err.Error() == "redigo: nil returned" {
			return nil, &util.UnAuthorized{ErrMessage: errUserNotFound.Error()}
		} else if err != nil {
			return nil, &util.InternalServer{ErrMessage: err.Error()}
		}
	*/
	//sample use of audit
	go audit.New(queue).Write(req.UserID, NotificationMessage)

	var response []*form.NotificationResponseData
	
	for _, notification := range notifications {
		fmt.Println(notification)
		fetchUserInfo(notification.DomainID)
		response = append(response, &form.NotificationResponseData{
			Type:           notification.Type,
			Message:        "need to frame message based on notification type and user",
			NotificationID: notification.DocumentID,
			PostLink:       "need to frame postLink based on user deatils",
		})
		
	}
	// fmt.Println(len(notifications))
	return response, nil
}

/* func getUserDetails(UserID string)  */
func fetchUserInfo(userID int) (string, int, error) {

	UserProfileBaseURL := os.Getenv("USER_PROFILE_API_BASE_URL")
	url := UserProfileBaseURL + string(userID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", 0, err
	}

	userAcountServiceKey := os.Getenv("USER_ACCOUNT_SERVICE_KEY")
	req.Header.Add("apikey", userAcountServiceKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	var r interface {
	}
	getContentResponseErr := json.Unmarshal(body, &r)
	if getContentResponseErr != nil {
		return "", 0, err
	}
	fmt.Println(r)
	return "", 0, err
}
