package server

import "net/http"

func Create(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusAccepted)
}
