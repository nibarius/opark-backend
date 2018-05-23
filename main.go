package opark_backend

import (
	"fmt"
	"net/http"
	"encoding/json"
	"google.golang.org/appengine"
	"github.com/justinas/alice"
)

const (
	contextKeyIdToken   = "idToken"
	contextKeyPushToken = "pushToken"
)
const (
	RESULT_FAILURE = "failure"
	RESULT_SUCCESS = "success"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/api/v1/waitList", alice.New(recoverHandler, requireHTTPS, basicAuth, requirePostOrDelete, requirePushToken).ThenFunc(waitListHandler))
	http.Handle("/worker/v1/waitList", alice.New(logAndIgnoreErrorsHandler).ThenFunc(workerWaitListHandler))
	http.Handle("/privacy", alice.New(recoverHandler, requireHTTPS).ThenFunc(privacyHandler))
	http.Handle("/doc/api", alice.New(recoverHandler, requireHTTPS).ThenFunc(apiHandler))
	http.Handle("/doc/usage_statistics", alice.New(recoverHandler, requireHTTPS).ThenFunc(usageStatisticsHandler))
}

var fallbackErrorBody = `{"status":` + string(http.StatusInternalServerError) +
	`,"message":"` + http.StatusText(http.StatusInternalServerError) + `","result": "failure"}`

type responseBody struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
}

type waitListRequestBody struct {
	PushToken string `json:"pushToken"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://play.google.com/store/apps/details?id=se.barsk.park", http.StatusMovedPermanently)
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	renderDocument(w, PrivacyStatement)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	renderDocument(w, Api)
}

func usageStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	renderDocument(w, UsageStatistics)
}

func workerWaitListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	f := NewFirestoreHandle(ctx)
	waitList := newWaitList(ctx, f)
	waitList.runQueuedTask()
	w.WriteHeader(http.StatusNoContent)
}

func waitListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	f := NewFirestoreHandle(ctx)
	defer f.close()
	waitList := newWaitList(ctx, f)
	uid := r.Context().Value(contextKeyIdToken).(string)

	switch r.Method {
	case http.MethodPost:
		pushToken := r.Context().Value(contextKeyPushToken).(string)
		waitList.add(uid, pushToken)
		fmt.Fprintf(w, makeSuccessResponse())
	case http.MethodDelete:
		waitList.remove(uid)
		w.WriteHeader(http.StatusNoContent)
	}
}

func makeErrorResponse(message string, status int) string {
	body := responseBody{Status: status, Message: message, Result: RESULT_FAILURE}
	jsonBody, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return fallbackErrorBody
	}
	return string(jsonBody)
}

func makeSuccessResponse() string {
	body := responseBody{Status: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: RESULT_SUCCESS}
	jsonBody, err := json.MarshalIndent(body, "", "  ")
	handleError(err, "Failed to generate a success body")
	return string(jsonBody)
}

func handleError(err error, msg string) {
	if err != nil {
		message := fmt.Sprintf("%s: %v", msg, err)
		panic(message)
	}
}
