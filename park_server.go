package opark_backend

import (
	"golang.org/x/net/context"
)

type ParkServer struct {
	context context.Context
}



type parkServerResponse struct {
	Free int `json:"free"`
	Used []struct {
		Regno string `json:"regno"`
		Start string `json:"start"`
		User  string `json:"user"`
	} `json:"used"`
}

func newParkServer(context context.Context) *ParkServer {
	ret := new(ParkServer)
	ret.context = context
	return ret
}

func (p *ParkServer) checkStatus() parkServerResponse {
	response := parkServerResponse{}
	err := getJson(p.context, parkServerStatusUrl, &response)
	if err != nil {
		panic("Failed to get park status: " + err.Error())
	}
	return response
}
