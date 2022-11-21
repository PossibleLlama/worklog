package server

import "net/http"

func Edit(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusAccepted)
}
