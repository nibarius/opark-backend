package opark_backend

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"net/http"
	"time"
	"encoding/json"
)

const (
	firestoreScope             = "https://www.googleapis.com/auth/datastore"
	firestoreWaitListEndpoint  = "https://firestore.googleapis.com/v1beta1/projects/opera-park/databases/(default)/documents/waitList"
	firestoreBatchWiteEndpoint = "https://firestore.googleapis.com/v1beta1/projects/opera-park/databases/(default)/documents:commit"
)

type FirestoreHandle struct {
	context     context.Context
	client      *firestore.Client
	accessToken string
}

type FirestoreWaitList struct {
	Documents []struct {
		Name string `json:"name"`
		Fields struct {
			PushToken struct {
				StringValue string `json:"stringValue"`
			} `json:"pushToken"`
		} `json:"fields"`
		CreateTime time.Time `json:"createTime"`
		UpdateTime time.Time `json:"updateTime"`
	} `json:"documents"`
}

type FirestoreUpdateWaitListRequestBody struct {
	Fields struct {
		PushToken struct {
			StringValue string `json:"stringValue"`
		} `json:"pushToken"`
	} `json:"fields"`
}

func NewFirestoreHandle(context context.Context) *FirestoreHandle {
	f := new(FirestoreHandle)
	f.context = context
	f.initializeFirestore()
	return f
}

func (f *FirestoreHandle) initializeFirestore() {
	var client *firestore.Client
	var err error
	client, err = firestore.NewClient(f.context, "opera-park")

	if err != nil {
		log.Errorf(f.context, "Failed to create a new firestore client: %v", err)
		panic("Failed to create a new firestore client: " + err.Error())
	}
	f.accessToken = f.getAccessToken()
	f.client = client
}

func (f *FirestoreHandle) close() {
	f.client.Close()
}

func (f *FirestoreHandle) getAccessToken() string {
	token, _, err := appengine.AccessToken(f.context, firestoreScope)
	handleError(err, "Failed to get Firestore access token")
	return token
}

// Get all entries in the wait list
func (f *FirestoreHandle) getWaitList() *FirestoreWaitList {

	params := requestParamaters{
		ctx:         f.context,
		url:         firestoreWaitListEndpoint,
		method:      http.MethodGet,
		accessToken: f.accessToken,
		body:        nil,
		caller:      "getWaitList",
	}
	//api documentation: https://cloud.google.com/firestore/docs/reference/rest/v1beta1/projects.databases.documents/list
	response := makeHTTPRequest(params)
	target := new(FirestoreWaitList)
	err := json.Unmarshal([]byte(response), target)
	handleError(err, "Failed to parse Firestore response in getWaitList()")

	return target
}

// Add the given user to the wait list (or update it if it already exist)
func (f *FirestoreHandle) addToWaitList(id, token string) string {
	body := new(FirestoreUpdateWaitListRequestBody)
	body.Fields.PushToken.StringValue = token

	params := requestParamaters{
		ctx:         f.context,
		url:         firestoreWaitListEndpoint + "/" + id,
		method:      http.MethodPatch,
		body:        body,
		accessToken: f.accessToken,
		caller:      "addToWaitList",
	}
	// api documentation https://cloud.google.com/firestore/docs/reference/rest/v1beta1/projects.databases.documents/patch
	return makeHTTPRequest(params)
}

// Removes the given user from the wait list
func (f *FirestoreHandle) removeFromWaitList(id string) string {
	params := requestParamaters{
		ctx:         f.context,
		url:         firestoreWaitListEndpoint + "/" + id,
		method:      http.MethodDelete,
		body:        nil,
		accessToken: f.accessToken,
		caller:      "removeFromWaitList",
	}
	// api documentation https://cloud.google.com/firestore/docs/reference/rest/v1beta1/projects.databases.documents/delete
	return makeHTTPRequest(params)
}

type BatchDeleteBody struct {
	Writes []DeleteOperation `json:"writes"`
}

type DeleteOperation struct {
	Delete string `json:"delete"`
}

// Remove all entries in the wait list
func (f *FirestoreHandle) clearWaitList() {
	list := f.getWaitList()
	body := new(BatchDeleteBody)
	for _, document := range list.Documents {
		body.Writes = append(body.Writes, DeleteOperation{document.Name})
	}

	params := requestParamaters{
		ctx:         f.context,
		url:         firestoreBatchWiteEndpoint,
		method:      http.MethodPost,
		body:        body,
		accessToken: f.accessToken,
		caller:      "clearWaitList",
	}
	// api documentation https://cloud.google.com/firestore/docs/reference/rest/v1beta1/projects.databases.documents/commit
	makeHTTPRequest(params)
}
