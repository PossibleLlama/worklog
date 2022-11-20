package server

import "net/http"

func Print(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusAccepted)
}
