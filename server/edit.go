package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/gorilla/mux"
)

func Edit(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")
	var body model.Work

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		helpers.LogError("error decoding body into work: "+err.Error(), "edit")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	newWl := model.NewWork(
		body.Title,
		body.Description,
		body.Author,
		body.Duration,
		body.Tags,
		body.When)
	newWl.ID = mux.Vars(req)["id"]

	new, status, err := wlService.EditWorklog(newWl.ID, newWl)
	resp.WriteHeader(status)
	if err != nil {
		helpers.LogError(fmt.Sprintf("failed to create work. %s", err.Error()), "create")
	}
	err = json.NewEncoder(resp).Encode(new)
	if err != nil {
		helpers.LogError("failed to encode work", "create")
	}
}
