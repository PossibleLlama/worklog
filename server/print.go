package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/gorilla/mux"
)

func Print(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json; charset=utf-8")

	var startDate, endDate time.Time
	var err error
	startDateString := req.URL.Query().Get("startDate")
	endDateString := req.URL.Query().Get("endDate")

	if startDateString == "" {
		startDate = time.Unix(int64(0), 0)
	} else {
		startDate, err = helpers.GetStringAsDateTime(startDateString)
		if err != nil {
			startDate = time.Now()
		}
	}
	if endDateString == "" {
		// Unix second for the year 3k
		endDate = time.Unix(int64(32503680000), 0)
	} else {
		endDate, err = helpers.GetStringAsDateTime(endDateString)
		if err != nil {
			endDate = time.Now()
		}
	}

	tags := []string{}
	for _, t := range strings.Split(req.URL.Query().Get("tags"), ",") {
		tags = append(tags, helpers.Sanitize(strings.TrimSpace(t)))
	}

	filter := model.Work{
		Title:       req.URL.Query().Get("title"),
		Description: req.URL.Query().Get("description"),
		Author:      req.URL.Query().Get("author"),
		Tags:        tags,
	}

	helpers.LogDebug(fmt.Sprintf("Getting worklogs from %v, to %v, with filter %+v", startDate, endDate, filter), "print")

	ret, status, err := wlService.GetWorklogsBetween(startDate, endDate, &filter)
	resp.WriteHeader(status)

	if err != nil {
		helpers.LogError(fmt.Sprintf("failed to find work. %s", err.Error()), "print")
		return
	}
	err = json.NewEncoder(resp).Encode(ret)
	if err != nil {
		helpers.LogError("failed to encode work", "print")
		return
	}
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
