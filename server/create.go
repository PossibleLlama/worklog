package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
)

func Create(resp http.ResponseWriter, req *http.Request) {
	var body model.Work

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		helpers.LogError("error decoding body into work: "+err.Error(), "create - server")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	// #nosec CWE-703 -- From my understanding, IO errors can occur, which are potentially an issue during file IO.
	// I haven't seen a similar example of harm to a network IO so will ignore for now.
	defer req.Body.Close()

	wl := model.NewWork(
		body.Title,
		body.Description,
		body.Author,
		body.Duration,
		body.Tags,
		body.When)

	status, err := wlService.CreateWorklog(wl)
	resp.WriteHeader(status)
	if err != nil {
		helpers.LogError(fmt.Sprintf("failed to create work. %s", err.Error()), "create")
	}
	err = json.NewEncoder(resp).Encode(wl)
	if err != nil {
		helpers.LogError("failed to encode work", "create")
	}
}
