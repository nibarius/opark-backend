package opark_backend

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"net/http"
	"encoding/json"
	"encoding/base64"
)

const (
	fcmSendEndpoint  = "https://fcm.googleapis.com/v1/projects/opera-park/messages:send"
	fcmScope         = "https://www.googleapis.com/auth/firebase.messaging"
	notificationType = "not_full"
)

type fcm struct {
	context     context.Context
	accessToken string
}

type fcmPushRequestBody struct {
	ValidateOnly bool `json:"validate_only"`
	Message struct {
		Token string `json:"token"`
		Data struct {
			Ttl      string `json:"ttl"`
			Type     string `json:"type"`
			TypeData string `json:"type_data"`
		} `json:"data"`
	} `json:"message"`
}

type NotFullTypeData struct {
	Free string `json:"free"` //todo: better with int?
}

func newFcm(context context.Context) *fcm {
	var ret = new(fcm)
	ret.context = context
	ret.accessToken = ret.getFcmAccessToken()

	return ret
}

func newPushRequestBody(token string, free string) *fcmPushRequestBody {
	var ret = new(fcmPushRequestBody)
	ret.ValidateOnly = false
	ret.Message.Token = token
	ret.Message.Data.Ttl = "7200" // 2 hours
	ret.Message.Data.Type = notificationType

	data, err := json.Marshal(NotFullTypeData{free})
	handleError(err, "Unable to create JSON for the NotFullTypeData")
	ret.Message.Data.TypeData = base64.StdEncoding.EncodeToString([]byte(data))
	return ret
}

func (f *fcm) sendFcmPush(pushData *fcmPushRequestBody) string {
	params := requestParamaters{
		ctx:         f.context,
		url:         fcmSendEndpoint,
		method:      http.MethodPost,
		accessToken: f.accessToken,
		body:        pushData,
		caller:      "sendFcmPush",
	}
	return makeHTTPRequest(params)
}

func (f *fcm) getFcmAccessToken() string {
	token, _, err := appengine.AccessToken(f.context, fcmScope)
	handleError(err, "Failed to get FCM access token")
	return token
}
