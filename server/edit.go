package server

import "net/http"

func Edit(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusAccepted)
}
