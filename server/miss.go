package server

import (
	"net/http"
)

func NotFound(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusNotFound)
}

func InvalidMethod(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusMethodNotAllowed)
}

func CORS(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}
