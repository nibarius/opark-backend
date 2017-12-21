package opark_backend

import (
	"errors"
	"golang.org/x/net/context"
)

type tokenInfo struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Iss           string `json:"iss"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	AccountId     string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}



func verifyToken(context context.Context, token string) (tokenInfo, error) {
	var info = tokenInfo{}
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
	err := getJson(context, url, &info)
	if err != nil {
		return tokenInfo{}, err
	} else if !clientIdIsValid(info.Aud) {
		return tokenInfo{}, errors.New("invalid client id")
	}

	return info, nil
}

func clientIdIsValid(aud string) bool {
	for _, clientId := range appEngineClientId {
		if clientId == aud {
			return true
		}
	}
	return false
}
