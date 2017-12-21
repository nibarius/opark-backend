package opark_backend

import (
	"encoding/base64"
	"net/http"
	"strings"
	"context"
	"google.golang.org/appengine"
	"encoding/json"
	"google.golang.org/appengine/log"
)

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
//https://gist.github.com/elithrar/9146306
func basicAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := appengine.NewContext(r)
		tokenInfo, err := verifyToken(ctx, pair[0])

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), contextKeyIdToken, tokenInfo.AccountId)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

func requireHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme != "https" && !appengine.IsDevAppServer() {
			body := makeErrorResponse("Please use HTTPS", http.StatusBadRequest)
			http.Error(w, body, http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func requirePostOrDelete(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodDelete {
			body := makeErrorResponse("Only HTTP POST and DELETE supported", http.StatusMethodNotAllowed)
			w.Header().Set("Allow", http.MethodPost+", "+http.MethodDelete)
			http.Error(w, body, http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func requirePushToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var reqBody = &waitListRequestBody{}
			json.NewDecoder(r.Body).Decode(reqBody)
			pushToken := reqBody.PushToken
			if pushToken == "" {
				body := makeErrorResponse("Parameter '"+contextKeyPushToken+"' missing.", http.StatusBadRequest)
				http.Error(w, body, http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), contextKeyPushToken, pushToken)
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	})
}

func recoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if appengine.IsDevAppServer() {
					panic(err)
				} else {
					log.Errorf(appengine.NewContext(r), "panic: %v", err)
					statusCode := http.StatusInternalServerError
					body := makeErrorResponse(http.StatusText(statusCode), statusCode)
					http.Error(w, body, statusCode)
				}
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func logAndIgnoreErrorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf(appengine.NewContext(r), "panic: %v", err)
			}
		}()
		h.ServeHTTP(w, r)
	})
}