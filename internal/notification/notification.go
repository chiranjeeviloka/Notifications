package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"notification-service/constants"
	"notification-service/internal/form"
	"notification-service/internal/model"
	"notification-service/internal/store"
	"notification-service/internal/util"
	"strconv"
	"strings"
	"time"

	"github.com/FreedomCentral/central/audit"
	"github.com/FreedomCentral/central/queue"
	"github.com/FreedomCentral/central/secret"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func View(c *gin.Context, db store.Store, queue queue.Queue) ([]*form.NotificationResponseData, error) {
	var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	//	logger.Info("Notification view started")
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
		//	logger.Error("Notification view validation errors", zap.Error(err))
		return nil, &util.BadRequest{ErrMessage: err.Error()}
	}

	var notifications []model.Notification

	err := db.GetNotification(c, &notifications, req.UserID)
	if err != nil {
		//	logger.Error("Notification view db errors", zap.Error(err))
		return nil, &util.InternalServer{ErrMessage: err.Error()}
	}

	//sample use of audit
	go audit.New(queue).Write(req.UserID, constants.NOTIFICATION_MESSAGE)

	var response []*form.NotificationResponseData

	for _, notification := range notifications {

		username, status, err := fetchUserInfo(notification.DomainID)
		username, status, err = "test", 1, nil
		//logger.Info("Notification view username,status,error", username, status, err)
		if err == nil && status > 0 {
			response = append(response, &form.NotificationResponseData{
				Type:           notification.Type,
				Message:        username + constants.REQUEST_FOLLOW,
				NotificationID: notification.DocumentID,
				PostLink:       strings.Replace(constants.USER_PROFILE_API_BASE_URL, "##USER_NAME##", strconv.Itoa(notification.DomainID), 1),
			})
		}

	}
	return response, nil
}

/* func getUserDetails(UserID string)  */
func fetchUserInfo(userID int) (string, int, error) {
	sec, err := secret.Open(constants.SERVICE_NAME, secret.UseYAMLPlainText)
	if err != nil {
		//logger.Infof("failed to open config for %v and its error %s", constants.SERVICE_NAME, err)
		return "", 0, err
	}
	UserProfileBaseURL, err := sec.Get("USER_PROFILE_API_BASE_URL")
	if err != nil {
		//logger.Infof("failed to load the user profile api base url , error %s", err)
		return "", 0, err
	}

	url := UserProfileBaseURL + strconv.Itoa(userID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		//logger.Infof("failed to create http new request while fetching the user info  , error %s", err)
		return "", 0, err
	}

	userAcountServiceKey, err := sec.Get("USER_ACCOUNT_SERVICE_KEY")
	if err != nil {
		//logger.Infof("failed to load the user account service key , error %s", err)
		return "", 0, err
	}

	req.Header.Add("apikey", userAcountServiceKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		//logger.Infof("failed to get the user details  , error %s", err)
		return "", 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//logger.Infof("failed to process the user details  , error %s", err)
		return "", 0, err
	}
	var userResp form.ServiceResponse
	getContentResponseErr := json.Unmarshal(body, &userResp)
	if getContentResponseErr != nil {
		//logger.Infof("failed to unmarshal user response  , error %s", err)
		return "", 0, err
	}
	fmt.Println(userResp)
	return "", 0, nil
}
