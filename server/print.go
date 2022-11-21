package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/gorilla/mux"
)

func Print(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusAccepted)
}

func PrintSingle(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	ret, status, err := wlService.GetWorklogsByID(&model.Work{}, mux.Vars(req)["id"])
	resp.WriteHeader(status)

	if err != nil {
		helpers.LogError(fmt.Sprintf("failed to find work. %s", err.Error()), "print")
		return
	} else if len(ret) == 0 {
		return
	}
	err = json.NewEncoder(resp).Encode(ret[0])
	if err != nil {
		helpers.LogError("failed to encode work", "print")
		return
	}
}
