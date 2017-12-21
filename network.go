package opark_backend

import (
	"google.golang.org/appengine/urlfetch"
	"io"
	"encoding/json"
	"bytes"
	"errors"
	"io/ioutil"
	"context"
	"net/http"
)

type requestParamaters struct {
	ctx         context.Context
	url         string
	method      string
	accessToken string
	body        interface{}
	caller      string
}

func makeHTTPRequest(params requestParamaters) string {
	client := urlfetch.Client(params.ctx)

	//log.Debugf(params.ctx, "Making %s request to: %s", params.method, params.url)

	var reqBody io.Reader = nil
	if params.body != nil {
		var jsonStr, err = json.Marshal(params.body)
		handleError(err, "Failed to create JSON for "+params.caller)
		reqBody = bytes.NewBuffer(jsonStr)
		//log.Debugf(params.ctx, "request body:\n %s", jsonStr)
	}

	req, err := http.NewRequest(params.method, params.url, reqBody)
	handleError(err, "Failed to create a new "+params.method+" request for "+params.caller)

	req.Header.Set("Authorization", "Bearer "+params.accessToken)
	if params.body != nil {
		req.Header.Set("Content-Type", "application/json; UTF-8")
	}

	resp, err := client.Do(req)
	handleError(err, "Failed to make a http request for "+params.caller)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	handleError(err, "Failed to read data from the http response for "+params.caller)

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status + "\n" + string(responseBody))
		handleError(err, "Unexpected status code returned when making http request for "+params.caller)
	}

	return string(responseBody)
}

// Makes a simple http get request to an url that responds with json.
// The result is decoded and inserted into the target parameter.
func getJson(context context.Context, url string, target interface{}) error {
	client := urlfetch.Client(context)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
